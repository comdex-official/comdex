package keeper

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
)

func (k Keeper) LendDutchActivator(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error {
	if lockedVault.Kind != nil {
		if !lockedVault.IsAuctionInProgress {
			extendedPair, found := k.GetLendPair(ctx, lockedVault.ExtendedPairId)
			if !found {
				return auctiontypes.ErrorInvalidPair
			}
			assetIn, _ := k.GetAsset(ctx, extendedPair.AssetIn)

			assetOut, _ := k.GetAsset(ctx, extendedPair.AssetOut)

			assetInPrice, found := k.GetPriceForAsset(ctx, assetIn.Id)
			if !found {
				ctx.Logger().Error(auctiontypes.ErrorPrices.Error(), lockedVault.LockedVaultId)
				return nil
			}
			assetOutPrice, found := k.GetPriceForAsset(ctx, assetOut.Id)
			if !found {
				ctx.Logger().Error(auctiontypes.ErrorPrices.Error(), lockedVault.LockedVaultId)
				return nil
			}
			//assetInPrice is the collateral price
			////Here collateral to be auctioned is received in ucollateral*uusd so inorder to get back amount we divide with uusd of assetIn
			AssetInPrice := sdk.NewDecFromInt(sdk.NewIntFromUint64(assetInPrice))
			if AssetInPrice.Equal(sdk.ZeroDec()) {
				ctx.Logger().Error(auctiontypes.ErrorPrices.Error(), lockedVault.LockedVaultId)
				return nil
			}
			AssetOutPrice := sdk.NewDecFromInt(sdk.NewIntFromUint64(assetOutPrice))
			if AssetOutPrice.Equal(sdk.ZeroDec()) {
				ctx.Logger().Error(auctiontypes.ErrorPrices.Error(), lockedVault.LockedVaultId)
				return nil
			}
			outflowToken := sdk.NewCoin(assetIn.Denom, lockedVault.CollateralToBeAuctioned.Quo(AssetInPrice).TruncateInt())
			inflowToken := sdk.NewCoin(assetOut.Denom, lockedVault.CollateralToBeAuctioned.Quo(AssetOutPrice).TruncateInt())

			AssetRatesStats, found := k.GetAssetRatesParams(ctx, extendedPair.AssetIn)
			if !found {
				ctx.Logger().Error(auctiontypes.ErrorAssetRates.Error(), lockedVault.LockedVaultId)
				return nil
			}
			liquidationPenalty := AssetRatesStats.LiquidationPenalty

			err1 := k.StartLendDutchAuction(ctx, outflowToken, inflowToken, lockedVault.AppId, assetOut.Id, assetIn.Id, lockedVault.LockedVaultId, lockedVault.Owner, liquidationPenalty)
			if err1 != nil {
				ctx.Logger().Error(auctiontypes.ErrorInStartDutchAuction.Error(), lockedVault.LockedVaultId)
				return nil
			}
		}
	}

	return nil
}

func (k Keeper) StartLendDutchAuction(
	ctx sdk.Context,
	outFlowToken sdk.Coin,
	inFlowToken sdk.Coin,
	appID uint64,
	assetInID, assetOutID uint64,
	lockedVaultID uint64,
	lockedVaultOwner string,
	liquidationPenalty sdk.Dec,
) error {
	var (
		inFlowTokenPrice  uint64
		outFlowTokenPrice uint64
		found             bool
	)

	lockedVault, found := k.GetLockedVault(ctx, appID, lockedVaultID)
	if !found {
		return auctiontypes.ErrorInvalidLockedVault
	}

	// If oracle Price required for the assetOut
	inFlowTokenPrice, found = k.GetPriceForAsset(ctx, assetInID)
	if !found {
		return auctiontypes.ErrorPrices
	}

	auctionParams, found := k.GetAddAuctionParamsData(ctx, appID)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}
	BorrowMetaData := lockedVault.GetBorrowMetaData()
	LendPos, _ := k.GetLend(ctx, BorrowMetaData.LendingId)
	pool, _ := k.GetPool(ctx, LendPos.PoolID)
	assetStat, _ := k.GetAssetRatesParams(ctx, LendPos.AssetID)
	cAsset, _ := k.GetAsset(ctx, assetStat.CAssetID)
	if outFlowToken.Amount.GT(sdk.ZeroInt()) {
		err := k.SendCoinsFromModuleToModule(ctx, pool.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(outFlowToken))
		if err != nil {
			return err
		}
		err = k.BurnCoins(ctx, pool.ModuleName, sdk.NewCoin(cAsset.Denom, outFlowToken.Amount))
		if err != nil {
			return err
		}
	}
	outFlowTokenPrice, found = k.GetPriceForAsset(ctx, assetOutID)
	if !found {
		return auctiontypes.ErrorPrices
	}
	// set target amount for debt
	inFlowTokenTargetAmount := inFlowToken.Amount
	mulfactor := inFlowTokenTargetAmount.ToDec().Mul(liquidationPenalty)
	inFlowTokenTargetAmount = inFlowTokenTargetAmount.Add(mulfactor.TruncateInt()).Add(lockedVault.InterestAccumulated)
	inFlowTokenTarget := sdk.NewCoin(inFlowToken.Denom, inFlowTokenTargetAmount)
	// These prices are in uusd
	outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(sdk.NewIntFromUint64(outFlowTokenPrice), auctionParams.Buffer)
	outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
	borrowOwner, err := sdk.AccAddressFromBech32(lockedVaultOwner)
	if err != nil {
		return err
	}
	timeNow := ctx.BlockTime()
	inFlowTokenCurrentAmount := sdk.NewCoin(inFlowToken.Denom, sdk.NewIntFromUint64(0))
	auction := auctiontypes.DutchAuction{
		OutflowTokenInitAmount:    outFlowToken,
		OutflowTokenCurrentAmount: outFlowToken,
		InflowTokenTargetAmount:   inFlowTokenTarget,
		InflowTokenCurrentAmount:  inFlowTokenCurrentAmount,
		OutflowTokenInitialPrice:  outFlowTokenInitialPrice,
		OutflowTokenCurrentPrice:  outFlowTokenInitialPrice,
		OutflowTokenEndPrice:      outFlowTokenEndPrice,
		InflowTokenCurrentPrice:   sdk.NewDecFromInt(sdk.NewIntFromUint64(inFlowTokenPrice)),
		StartTime:                 timeNow,
		EndTime:                   timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		AuctionStatus:             auctiontypes.AuctionStartNoBids,
		BiddingIds:                []*auctiontypes.BidOwnerMapping{},
		AuctionMappingId:          auctionParams.DutchId,
		AppId:                     appID,
		AssetInId:                 assetInID,
		AssetOutId:                assetOutID,
		LockedVaultId:             lockedVaultID,
		VaultOwner:                borrowOwner,
		LiquidationPenalty:        liquidationPenalty,
	}

	auction.AuctionId = k.GetLendAuctionID(ctx) + 1
	k.SetLendAuctionID(ctx, auction.AuctionId)
	err = k.SetDutchLendAuction(ctx, auction)
	if err != nil {
		return err
	}
	err = k.SetFlagIsAuctionInProgress(ctx, appID, lockedVaultID, true)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) PlaceLendDutchAuctionBid(ctx sdk.Context, appID, auctionMappingID, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin, max sdk.Dec) error {
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
	outFlowTokenCurrentPrice := auction.OutflowTokenCurrentPrice.Ceil().TruncateInt()
	inFlowTokenCurrentPrice := auction.InflowTokenCurrentPrice.Ceil().TruncateInt()

	slice := bid.Amount // cmdx

	a := auction.InflowTokenTargetAmount.Amount
	b := auction.InflowTokenCurrentAmount.Amount
	tab := a.Sub(b)
	// owe is $token to be given to user
	owe := slice.Mul(outFlowTokenCurrentPrice)
	inFlowTokenAmount := owe.Quo(inFlowTokenCurrentPrice)
	TargetReachedFlag := false
	if inFlowTokenAmount.GT(tab) {
		TargetReachedFlag = true
		inFlowTokenAmount = tab
		owe = inFlowTokenAmount.Mul(inFlowTokenCurrentPrice)
		slice = owe.Quo(outFlowTokenCurrentPrice)
		owe = slice.Mul(outFlowTokenCurrentPrice)
	}
	inFlowTokenCoin := sdk.NewCoin(auction.InflowTokenTargetAmount.Denom, inFlowTokenAmount)

	// required target cmst to raise in usd * 10**-12
	// here we are multiplying each ucmdx with uusd so cmdx tokens price will be calculated amount * 10**-12

	lockedVault, found := k.GetLockedVault(ctx, appID, auction.LockedVaultId)
	if !found {
		return auctiontypes.ErrorInvalidLockedVault
	}
	ExtendedPairVault, found := k.GetLendPair(ctx, lockedVault.ExtendedPairId)
	if !found {
		return auctiontypes.ErrorInvalidExtendedPairVault
	}
	assetStats, _ := k.GetAssetRatesParams(ctx, ExtendedPairVault.AssetIn)
	assetOutPool, _ := k.GetPool(ctx, ExtendedPairVault.AssetOutPoolID)
	// dust is in usd * 10*-6 (uusd)
	dust := sdk.NewIntFromUint64(ExtendedPairVault.MinUsdValueLeft)
	// here subtracting current amount and slice to get amount left in auction and also converting it to usd * 10**-12
	outLeft := auction.OutflowTokenCurrentAmount.Amount.Mul(outFlowTokenCurrentPrice)
	amountLeftInPUSD := outLeft.Sub(owe)
	// convert amountLeft to uusd from pusd(10**-12) so we can compare dust and amountLeft in UUSD . this happens by converting ucmdx to cmdx

	// check if bid in usd*10**-12 is greater than required target cmst in usd*10**-12
	// if user wants to buy more than target cmst then user should be sold only required cmst amount
	// so we need to divide tab by outflow token current price and we get outflowtoken amount to be sold to user
	// if user is not buying more than required cmst then we check for dust
	// here tab is divided by price again so price ration is only considered , we are not using owe again in this function so no problem
	// As tab is the amount calculated from difference of target and current inflow token we will be using same as inflow token

	if amountLeftInPUSD.LT(dust) && !amountLeftInPUSD.Equal(sdk.ZeroInt()) && !TargetReachedFlag {
		coll := auction.OutflowTokenCurrentAmount.Amount.Uint64()
		dust := dust.Uint64()
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "either bid all the amount %d (UTOKEN) or bid amount by leaving dust greater than %d PUSD", coll, dust)
	}

	outFlowTokenCoin := sdk.NewCoin(auction.OutflowTokenInitAmount.Denom, slice)

	err = k.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(inFlowTokenCoin))
	if err != nil {
		return err
	}
	// sending inflow token back to the pool
	err = k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, assetOutPool.ModuleName, sdk.NewCoins(inFlowTokenCoin))
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
		_ = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, sdk.AccAddress(lockedVault.Owner), sdk.NewCoins(total))
		err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidder, sdk.NewCoins(totalAmountToBidder))
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
		modBal := k.ModuleBalance(ctx, lendtypes.ModuleName, requiredAmount.Denom)
		if modBal.LT(requiredAmount.Amount) {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "Reserve pool having insufficient balance for this bid")
		}

		// reduce the qty from reserve pool
		pairID := lockedVault.ExtendedPairId
		lendPair, _ := k.GetLendPair(ctx, pairID)
		inFlowTokenAssetID := lendPair.AssetOut

		err = k.UpdateReserveBalances(ctx, inFlowTokenAssetID, lendtypes.ModuleName, requiredAmount, false)
		if err != nil {
			return err
		}
		err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidder, sdk.NewCoins(totalAmountToBidder))
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
	} else if auction.OutflowTokenCurrentAmount.Amount.IsZero() { // entire collateral sold out
		err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidder, sdk.NewCoins(totalAmountToBidder))
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

	lockedVault, found := k.GetLockedVault(ctx, dutchAuction.AppId, dutchAuction.LockedVaultId)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}

	lockedVault.AmountIn = lockedVault.AmountIn.Sub(dutchAuction.OutflowTokenInitAmount.Amount)
	// set sell of history in locked vault
	err := k.CreateLockedVaultHistory(ctx, lockedVault)
	if err != nil {
		return err
	}
	borrowMetaData := lockedVault.GetBorrowMetaData()
	lendPos, _ := k.GetLend(ctx, borrowMetaData.LendingId)
	lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(dutchAuction.OutflowTokenInitAmount.Amount)
	lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(dutchAuction.OutflowTokenInitAmount.Amount)
	k.SetLend(ctx, lendPos)

	lockedVault.AmountOut = lockedVault.AmountOut.Sub(dutchAuction.InflowTokenTargetAmount.Amount)
	lockedVault.UpdatedAmountOut = lockedVault.UpdatedAmountOut.Sub(dutchAuction.InflowTokenTargetAmount.Amount)
	if lockedVault.AmountOut.LTE(sdk.ZeroInt()) {
		lockedVault.AmountOut = sdk.ZeroInt()
	}
	k.SetLockedVault(ctx, lockedVault)
	dutchAuction.AuctionStatus = auctiontypes.AuctionEnded

	// update locked vault
	err = k.SetFlagIsAuctionComplete(ctx, dutchAuction.AppId, dutchAuction.LockedVaultId, true)
	if err != nil {
		return err
	}

	err = k.SetFlagIsAuctionInProgress(ctx, dutchAuction.AppId, dutchAuction.LockedVaultId, false)
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
	err = k.UnLiquidateLockedBorrows(ctx, lockedVault.AppId, lockedVault.LockedVaultId, dutchAuction)
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
	dutchAuctions := k.GetDutchLendAuctions(ctx, appID)
	auctionParams, found := k.GetAddAuctionParamsData(ctx, appID)
	if !found {
		return nil
	}
	// SET current price of inflow token and outflow token
	for _, dutchAuction := range dutchAuctions {
		inFlowTokenCurrentPrice, found := k.GetPriceForAsset(ctx, dutchAuction.AssetInId)
		if !found {
			return auctiontypes.ErrorPrices
		}

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
			OutFlowTokenCurrentPrice, found := k.GetPriceForAsset(ctx, dutchAuction.AssetOutId)
			if !found {
				return auctiontypes.ErrorPrices
			}
			timeNow := ctx.BlockTime()
			dutchAuction.StartTime = timeNow
			dutchAuction.EndTime = timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
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
			// outFlowTokenCurrentPrice := sdk.NewIntFromUint64(10)
		}
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
