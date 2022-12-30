package keeper

import (
	"strconv"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"

	utils "github.com/comdex-official/comdex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
)

func (k Keeper) LendDutchActivator(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error {
	// basic check for differentiating between lockedVault of borrow and lockedVault of vault
	// when auction is not in progress
	if lockedVault.Kind != nil && !lockedVault.IsAuctionInProgress {
		extendedPair, found := k.lend.GetLendPair(ctx, lockedVault.ExtendedPairId)
		if !found {
			return auctiontypes.ErrorInvalidPair
		}
		assetIn, _ := k.asset.GetAsset(ctx, extendedPair.AssetIn)   // collateral
		assetOut, _ := k.asset.GetAsset(ctx, extendedPair.AssetOut) // debt

		AssetInPrice, err := k.market.CalcAssetPrice(ctx, assetIn.Id, sdk.OneInt())
		if err != nil {
			return auctiontypes.ErrorPrices
		}
		AssetOutPrice, err := k.market.CalcAssetPrice(ctx, assetOut.Id, sdk.OneInt())
		if err != nil {
			return auctiontypes.ErrorPrices
		}
		outflowToken := sdk.NewCoin(assetIn.Denom, lockedVault.CollateralToBeAuctioned.Quo(AssetInPrice).TruncateInt())  // Asset being auctioned
		inflowToken := sdk.NewCoin(assetOut.Denom, lockedVault.CollateralToBeAuctioned.Quo(AssetOutPrice).TruncateInt()) // Asset to be recovered

		AssetRatesStats, found := k.lend.GetAssetRatesParams(ctx, extendedPair.AssetIn)
		if !found {
			ctx.Logger().Error(auctiontypes.ErrorAssetRates.Error(), lockedVault.LockedVaultId)
			return auctiontypes.ErrorAssetRates
		}
		liquidationPenalty := AssetRatesStats.LiquidationPenalty
		// from here the lend dutch auction is started
		err1 := k.StartLendDutchAuction(ctx, outflowToken, inflowToken, lockedVault, assetOut.Id, assetIn.Id, liquidationPenalty)
		if err1 != nil {
			ctx.Logger().Error(auctiontypes.ErrorInStartDutchAuction.Error(), lockedVault.LockedVaultId)
			return auctiontypes.ErrorInStartDutchAuction
		}
	}

	return nil
}

func (k Keeper) StartLendDutchAuction(
	ctx sdk.Context,
	outFlowToken sdk.Coin,
	inFlowToken sdk.Coin,
	lockedVault liquidationtypes.LockedVault,
	assetInID, assetOutID uint64, // debt, collateral
	liquidationPenalty sdk.Dec,
) error {
	// If oracle Price required for the assetOut
	twaInData, found := k.market.GetTwa(ctx, assetInID)
	if !found || !twaInData.IsPriceActive {
		return auctiontypes.ErrorPrices
	}
	inFlowTokenPrice := twaInData.Twa

	auctionParams, found := k.lend.GetAddAuctionParamsData(ctx, lockedVault.AppId)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}

	twaData, found := k.market.GetTwa(ctx, assetOutID)
	if !found || !twaData.IsPriceActive {
		return auctiontypes.ErrorPrices
	}

	outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(sdk.NewIntFromUint64(twaData.Twa), auctionParams.Buffer)
	outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
	borrowOwner, err := sdk.AccAddressFromBech32(lockedVault.Owner)
	if err != nil {
		return err
	}

	auction := auctiontypes.DutchAuction{
		OutflowTokenInitAmount:    outFlowToken,
		OutflowTokenCurrentAmount: outFlowToken,
		InflowTokenTargetAmount:   inFlowToken,
		InflowTokenCurrentAmount:  sdk.NewCoin(inFlowToken.Denom, sdk.NewIntFromUint64(0)),
		OutflowTokenInitialPrice:  outFlowTokenInitialPrice,
		OutflowTokenCurrentPrice:  outFlowTokenInitialPrice,
		OutflowTokenEndPrice:      outFlowTokenEndPrice,
		InflowTokenCurrentPrice:   sdk.NewDecFromInt(sdk.NewIntFromUint64(inFlowTokenPrice)),
		StartTime:                 ctx.BlockTime(),
		EndTime:                   ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		AuctionStatus:             auctiontypes.AuctionStartNoBids,
		BiddingIds:                []*auctiontypes.BidOwnerMapping{},
		AuctionMappingId:          auctionParams.DutchId,
		AppId:                     lockedVault.AppId,
		AssetInId:                 assetInID,  // debt
		AssetOutId:                assetOutID, // collateral
		LockedVaultId:             lockedVault.LockedVaultId,
		VaultOwner:                borrowOwner,
		LiquidationPenalty:        liquidationPenalty,
	}

	auction.AuctionId = k.GetLendAuctionID(ctx) + 1
	k.SetLendAuctionID(ctx, auction.AuctionId)
	err = k.SetDutchLendAuction(ctx, auction)
	if err != nil {
		return err
	}
	err = k.liquidation.SetFlagIsAuctionInProgress(ctx, lockedVault.AppId, lockedVault.LockedVaultId, true)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			auctiontypes.EventTypeLendDutchNewAuction,
			sdk.NewAttribute(auctiontypes.DataAppID, strconv.FormatUint(auction.AppId, 10)),
			sdk.NewAttribute(auctiontypes.AttributeKeyOwner, auction.VaultOwner.String()),
			sdk.NewAttribute(auctiontypes.AttributeKeyCollateral, auction.OutflowTokenInitAmount.String()),
			sdk.NewAttribute(auctiontypes.AttributeKeyDebt, auction.InflowTokenTargetAmount.String()),
			sdk.NewAttribute(auctiontypes.AttributeKeyStartTime, auction.StartTime.String()),
			sdk.NewAttribute(auctiontypes.AttributeKeyEndTime, auction.EndTime.String()),
		),
	)

	return nil
}

func (k Keeper) PlaceLendDutchAuctionBid(ctx sdk.Context, appID, auctionMappingID, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin, max sdk.Dec) error {
	if bid.Amount.Equal(sdk.ZeroInt()) {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bid amount can't be Zero")
	}
	auction, err := k.GetDutchLendAuction(ctx, appID, auctionMappingID, auctionID)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction id %d not found", auctionID)
	}
	if bid.Denom != auction.OutflowTokenCurrentAmount.Denom {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bid denom %s not found", bid.Denom)
	}
	if bid.Amount.GT(auction.OutflowTokenCurrentAmount.Amount) {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bid amount can't be greater than collateral available")
	}

	max = k.GetUUSDFromUSD(ctx, max)

	// Here OutflowToken current price is in uusd and max is in uusd
	if max.LT(auction.OutflowTokenCurrentPrice.Ceil()) {
		return auctiontypes.ErrorInvalidDutchPrice
	}

	// slice tells amount of collateral user should be given
	// using ceil as we need extract more from users
	outFlowTokenCurrentPrice := auction.OutflowTokenCurrentPrice
	inFlowTokenCurrentPrice := auction.InflowTokenCurrentPrice

	slice := bid.Amount // cmdx

	a := auction.InflowTokenTargetAmount.Amount
	b := auction.InflowTokenCurrentAmount.Amount
	tab := a.Sub(b)
	// owe is $token to be given to user
	owe, inFlowTokenAmount, err := k.vault.GetAmountOfOtherToken(ctx, auction.AssetOutId, outFlowTokenCurrentPrice, slice, auction.AssetInId, inFlowTokenCurrentPrice)
	if err != nil {
		return err
	}
	if inFlowTokenAmount.LTE(sdk.ZeroInt()) {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "Calculated Auction Amount is Zero")
	}

	TargetReachedFlag := false
	if inFlowTokenAmount.GT(tab) {
		TargetReachedFlag = true
		inFlowTokenAmount = tab
		owe, slice, err = k.vault.GetAmountOfOtherToken(ctx, auction.AssetInId, inFlowTokenCurrentPrice, inFlowTokenAmount, auction.AssetOutId, outFlowTokenCurrentPrice)
		if err != nil {
			return err
		}
	}
	inFlowTokenCoin := sdk.NewCoin(auction.InflowTokenTargetAmount.Denom, inFlowTokenAmount)

	// required target cmst to raise in usd * 10**-12
	// here we are multiplying each ucmdx with uusd so cmdx tokens price will be calculated amount * 10**-12

	lockedVault, found := k.liquidation.GetLockedVault(ctx, appID, auction.LockedVaultId)
	if !found {
		return auctiontypes.ErrorInvalidLockedVault
	}
	ExtendedPairVault, found := k.lend.GetLendPair(ctx, lockedVault.ExtendedPairId)
	if !found {
		return auctiontypes.ErrorInvalidExtendedPairVault
	}
	assetStats, _ := k.lend.GetAssetRatesParams(ctx, ExtendedPairVault.AssetIn)
	assetOutPool, _ := k.lend.GetPool(ctx, ExtendedPairVault.AssetOutPoolID)
	// dust is in usd * 10*-6 (uusd)
	dust := sdk.NewIntFromUint64(ExtendedPairVault.MinUsdValueLeft)
	// here subtracting current amount and slice to get amount left in auction and also converting it to usd * 10**-12
	outLeft, err := k.CalcDollarValueForToken(ctx, auction.AssetInId, outFlowTokenCurrentPrice, auction.OutflowTokenCurrentAmount.Amount)
	if err != nil {
		return err
	}
	outLeftDebt, err := k.CalcDollarValueForToken(ctx, auction.AssetInId, inFlowTokenCurrentPrice, tab)
	if err != nil {
		return err
	}
	amountLeftInPUSD := outLeft.Sub(owe)
	amountLeftInPUSDforDebt := outLeftDebt.Sub(owe)
	// convert amountLeft to uusd from pusd(10**-12) so we can compare dust and amountLeft in UUSD . this happens by converting ucmdx to cmdx

	// check if bid in usd*10**-12 is greater than required target cmst in usd*10**-12
	// if user wants to buy more than target cmst then user should be sold only required cmst amount
	// so we need to divide tab by outflow token current price and we get outflowtoken amount to be sold to user
	// if user is not buying more than required cmst then we check for dust
	// here tab is divided by price again so price ration is only considered , we are not using owe again in this function so no problem
	// As tab is the amount calculated from difference of target and current inflow token we will be using same as inflow token

	// Dust check for collateral
	if amountLeftInPUSD.LT(sdk.NewDecFromInt(dust)) && !amountLeftInPUSD.Equal(sdk.ZeroDec()) && !TargetReachedFlag {
		coll := auction.OutflowTokenCurrentAmount.Amount.Uint64()
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "Either bid all the collateral amount %d (UTOKEN) or bid amount by leaving dust greater than %d USD", coll, dust)
	}

	// Dust check for debt
	if amountLeftInPUSDforDebt.LT(sdk.NewDecFromInt(dust)) && !amountLeftInPUSDforDebt.Equal(sdk.ZeroDec()) && !amountLeftInPUSD.Equal(sdk.ZeroDec()) {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "Minimum amount left to be recovered should not be less than %d ", dust)
	}

	outFlowTokenCoin := sdk.NewCoin(auction.OutflowTokenInitAmount.Denom, slice)

	err = k.bank.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(inFlowTokenCoin))
	if err != nil {
		return err
	}
	// sending inflow token back to the pool
	err = k.bank.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, assetOutPool.ModuleName, sdk.NewCoins(inFlowTokenCoin))
	if err != nil {
		return err
	}

	// calculating additional auction bonus to the bidder

	auctionBonus := slice.ToDec().Mul(assetStats.LiquidationBonus)
	totalAmountToBidder := sdk.NewCoin(auction.OutflowTokenInitAmount.Denom, slice.Add(auctionBonus.TruncateInt()))

	biddingID, err := k.CreateNewDutchLendBid(ctx, appID, auctionMappingID, auctionID, bidder.String(), inFlowTokenCoin, outFlowTokenCoin)
	if err != nil {
		return err
	}
	bidIDOwner := &auctiontypes.BidOwnerMapping{BidId: biddingID, BidOwner: bidder.String()}
	auction.BiddingIds = append(auction.BiddingIds, bidIDOwner)
	if auction.AuctionStatus == auctiontypes.AuctionStartNoBids {
		auction.AuctionStatus = auctiontypes.AuctionGoingOn
	}

	// calculate inflow amount and outflow amount if  user  transaction successful
	auction.OutflowTokenCurrentAmount = auction.OutflowTokenCurrentAmount.Sub(outFlowTokenCoin)
	auction.InflowTokenCurrentAmount = auction.InflowTokenCurrentAmount.Add(inFlowTokenCoin)

	// collateral not over but target cmst reached then send remaining collateral to owner
	// if inflow token current amount >= InflowTokenTargetAmount
	if auction.InflowTokenCurrentAmount.IsGTE(auction.InflowTokenTargetAmount) {
		total := auction.OutflowTokenCurrentAmount
		vaultHolder, err := sdk.AccAddressFromBech32(lockedVault.Owner)
		if err != nil {
			panic(err)
		}
		if total.Amount.GT(sdk.ZeroInt()) {
			err = k.bank.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, vaultHolder, sdk.NewCoins(total))
			if err != nil {
				return err
			}
		}
		err = k.bank.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidder, sdk.NewCoins(totalAmountToBidder))
		if err != nil {
			return err
		}
		err = k.SetDutchLendAuction(ctx, auction)
		if err != nil {
			return err
		}

		err = k.CloseDutchLendAuction(ctx, auction)
		if err != nil {
			return err
		}
	} else if auction.OutflowTokenCurrentAmount.Amount.IsZero() && auction.InflowTokenCurrentAmount.IsLT(auction.InflowTokenTargetAmount) { // entire collateral sold out
		// take requiredAmount from reserve-pool
		requiredAmount := auction.InflowTokenTargetAmount.Sub(auction.InflowTokenCurrentAmount)
		// get reserve balance if the requiredAmount is available in the reserves or not
		modBal := k.lend.ModuleBalance(ctx, lendtypes.ModuleName, requiredAmount.Denom)
		if modBal.LT(requiredAmount.Amount) {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "Reserve pool having insufficient balance for this bid")
		}

		// reduce the qty from reserve pool
		pairID := lockedVault.ExtendedPairId
		lendPair, _ := k.lend.GetLendPair(ctx, pairID)
		inFlowTokenAssetID := lendPair.AssetOut
		err = k.bank.SendCoinsFromModuleToModule(ctx, lendtypes.ModuleName, assetOutPool.ModuleName, sdk.NewCoins(requiredAmount))
		if err != nil {
			return err
		}

		err = k.lend.UpdateReserveBalances(ctx, inFlowTokenAssetID, lendtypes.ModuleName, requiredAmount, false)
		if err != nil {
			return err
		}
		err = k.bank.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidder, sdk.NewCoins(totalAmountToBidder))
		if err != nil {
			return err
		}

		err = k.SetDutchLendAuction(ctx, auction)
		if err != nil {
			return err
		}

		// remove dutch auction
		err = k.CloseDutchLendAuction(ctx, auction)
		if err != nil {
			return err
		}
	} else {
		err = k.bank.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidder, sdk.NewCoins(totalAmountToBidder))
		if err != nil {
			return err
		}

		err = k.SetDutchLendAuction(ctx, auction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) CreateNewDutchLendBid(ctx sdk.Context, appID, auctionMappingID, auctionID uint64, bidder string, outFlowTokenCoin sdk.Coin, inFlowTokenCoin sdk.Coin) (biddingID uint64, err error) {
	bidding := auctiontypes.DutchBiddings{
		BiddingId:          k.GetUserBiddingID(ctx) + 1,
		AuctionId:          auctionID,
		AuctionStatus:      auctiontypes.ActiveAuctionStatus,
		Bidder:             bidder,
		OutflowTokenAmount: outFlowTokenCoin,
		InflowTokenAmount:  inFlowTokenCoin,
		BiddingTimestamp:   ctx.BlockTime(),
		BiddingStatus:      auctiontypes.SuccessBiddingStatus,
		AppId:              appID,
		AuctionMappingId:   auctionMappingID,
	}
	k.SetUserBiddingID(ctx, bidding.BiddingId)
	err = k.SetDutchUserLendBidding(ctx, bidding)
	if err != nil {
		return biddingID, err
	}
	return bidding.BiddingId, nil
}

func (k Keeper) CloseDutchLendAuction(
	ctx sdk.Context,
	dutchAuction auctiontypes.DutchAuction,
) error { // delete dutch biddings
	if dutchAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		for _, biddingID := range dutchAuction.BiddingIds {
			bidding, err := k.GetDutchLendUserBidding(ctx, biddingID.BidOwner, dutchAuction.AppId, biddingID.BidId)
			if err != nil {
				return err
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			err = k.SetDutchUserLendBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.DeleteDutchLendUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.SetHistoryDutchUserLendBidding(ctx, bidding)
			if err != nil {
				return err
			}
		}
	}

	lockedVault, found := k.liquidation.GetLockedVault(ctx, dutchAuction.AppId, dutchAuction.LockedVaultId)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}

	// set sell of history in locked vault
	err := k.liquidation.CreateLockedVaultHistory(ctx, lockedVault)
	if err != nil {
		return err
	}

	lockedVault.AmountOut = lockedVault.AmountOut.Sub(dutchAuction.InflowTokenTargetAmount.Amount)
	lockedVault.UpdatedAmountOut = lockedVault.UpdatedAmountOut.Sub(dutchAuction.InflowTokenTargetAmount.Amount)
	if lockedVault.AmountOut.LTE(sdk.ZeroInt()) {
		lockedVault.AmountOut = sdk.ZeroInt()
	}
	if lockedVault.UpdatedAmountOut.LTE(sdk.ZeroInt()) {
		lockedVault.UpdatedAmountOut = sdk.ZeroInt()
	}
	k.liquidation.SetLockedVault(ctx, lockedVault)
	dutchAuction.AuctionStatus = auctiontypes.AuctionEnded

	// update locked vault
	err = k.liquidation.SetFlagIsAuctionComplete(ctx, dutchAuction.AppId, dutchAuction.LockedVaultId, true)
	if err != nil {
		return err
	}

	err = k.liquidation.SetFlagIsAuctionInProgress(ctx, dutchAuction.AppId, dutchAuction.LockedVaultId, false)
	if err != nil {
		return err
	}

	err = k.SetDutchLendAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}
	err = k.SetHistoryDutchLendAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}
	err = k.liquidation.UnLiquidateLockedBorrows(ctx, lockedVault.AppId, lockedVault.LockedVaultId, dutchAuction)
	if err != nil {
		return err
	}
	err = k.DeleteDutchLendAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) RestartDutchLendAuctions(ctx sdk.Context, appID uint64) error {
	auctionParams, found := k.lend.GetAddAuctionParamsData(ctx, appID)
	if !found {
		return nil
	}
	dutchAuctions := k.GetDutchLendAuctions(ctx, appID)
	// SET current price of inflow token and outflow token
	for _, dutchAuction := range dutchAuctions {
		_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
			twaData, found := k.market.GetTwa(ctx, dutchAuction.AssetInId)
			if !found || !twaData.IsPriceActive {
				return auctiontypes.ErrorPrices
			}
			inFlowTokenCurrentPrice := twaData.Twa

			// inFlowTokenCurrentPrice := sdk.MustNewDecFromStr("1")
			// tau := sdk.NewInt(int64(auctionParams.AuctionDurationSeconds))
			tnume := dutchAuction.OutflowTokenInitialPrice.Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(auctionParams.AuctionDurationSeconds)))
			tdeno := dutchAuction.OutflowTokenInitialPrice.Sub(dutchAuction.OutflowTokenEndPrice)
			ntau := tnume.Quo(tdeno)
			tau := sdk.NewInt(ntau.TruncateInt64())
			dur := ctx.BlockTime().Sub(dutchAuction.StartTime)
			seconds := sdk.NewInt(int64(dur.Seconds()))
			outFlowTokenCurrentPrice := k.getPriceFromLinearDecreaseFunction(dutchAuction.OutflowTokenInitialPrice, tau, seconds)
			dutchAuction.InflowTokenCurrentPrice = sdk.NewDec(int64(inFlowTokenCurrentPrice))
			dutchAuction.OutflowTokenCurrentPrice = outFlowTokenCurrentPrice
			err := k.SetDutchLendAuction(ctx, dutchAuction)
			if err != nil {
				return err
			}
			// check if auction need to be restarted
			if ctx.BlockTime().After(dutchAuction.EndTime) {
				twaData, found := k.market.GetTwa(ctx, dutchAuction.AssetOutId)
				if !found || !twaData.IsPriceActive {
					return auctiontypes.ErrorPrices
				}
				OutFlowTokenCurrentPrice := twaData.Twa
				dutchAuction.StartTime = ctx.BlockTime()
				dutchAuction.EndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
				outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(sdk.NewIntFromUint64(OutFlowTokenCurrentPrice), auctionParams.Buffer)
				outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
				dutchAuction.OutflowTokenInitialPrice = outFlowTokenInitialPrice
				dutchAuction.OutflowTokenEndPrice = outFlowTokenEndPrice
				dutchAuction.OutflowTokenCurrentPrice = outFlowTokenInitialPrice
				err = k.SetDutchLendAuction(ctx, dutchAuction)
				if err != nil {
					return err
				}
				// SET initial price fetched from market module and also end price , start time , end time
			}
			return nil
		})
	}
	return nil
}

func (k Keeper) RestartLendDutch(ctx sdk.Context, appID uint64) error {
	err := k.RestartDutchLendAuctions(ctx, appID)
	if err != nil {
		return err
	}
	return nil
}
