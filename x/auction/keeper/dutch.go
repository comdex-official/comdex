package keeper

import (
	"fmt"
	"strconv"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	vaulttypes "github.com/comdex-official/comdex/x/vault/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/comdex-official/comdex/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
)

func (k Keeper) DutchActivator(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error {
	if lockedVault.Kind == nil {
		extendedPair, found := k.asset.GetPairsVault(ctx, lockedVault.ExtendedPairId)
		if !found {
			return fmt.Errorf("pair not found for extended pair id %d", lockedVault.ExtendedPairId)
		}
		pair, _ := k.asset.GetPair(ctx, extendedPair.PairId)

		assetIn, _ := k.asset.GetAsset(ctx, pair.AssetIn) // collateral(cmdx)

		assetOut, _ := k.asset.GetAsset(ctx, pair.AssetOut) // debt(cmst)

		outflowToken := sdk.NewCoin(assetIn.Denom, lockedVault.AmountIn) // cmdx
		inflowToken := sdk.NewCoin(assetOut.Denom, sdkmath.ZeroInt())    // cmst

		liquidationPenalty := extendedPair.LiquidationPenalty

		err1 := k.StartDutchAuction(ctx, outflowToken, inflowToken, lockedVault.AppId, assetOut.Id, assetIn.Id, lockedVault.LockedVaultId, lockedVault.Owner, liquidationPenalty)
		if err1 != nil {
			return fmt.Errorf("error in start dutch auction for locked vault id %d", lockedVault.LockedVaultId)
		}
	}
	return nil
}

func (k Keeper) StartDutchAuction(
	ctx sdk.Context,
	outFlowToken sdk.Coin, // cmdx
	inFlowToken sdk.Coin, // cmst
	appID uint64,
	assetInID uint64, // cmst
	assetOutID uint64, // cmdx
	lockedVaultID uint64,
	lockedVaultOwner string,
	liquidationPenalty sdkmath.LegacyDec,
) error {
	var (
		inFlowTokenPrice  uint64
		outFlowTokenPrice uint64
		found             bool
	)

	lockedVault, found := k.liquidation.GetLockedVault(ctx, appID, lockedVaultID)
	if !found {
		return auctiontypes.ErrorInvalidLockedVault
	}
	k.vault.DeleteUserVaultExtendedPairMapping(ctx, lockedVault.Owner, appID, lockedVault.ExtendedPairId)

	extendedPairVault := lockedVault.ExtendedPairId

	ExtendedPairVault, found := k.asset.GetPairsVault(ctx, extendedPairVault)
	if !found {
		return auctiontypes.ErrorInvalidExtendedPairVault
	}
	if ExtendedPairVault.AssetOutOraclePrice {
		// If oracle Price required for the assetOut
		twaData, found := k.market.GetTwa(ctx, assetInID)
		if !found || !twaData.IsPriceActive {
			return auctiontypes.ErrorPrices
		}
		inFlowTokenPrice = twaData.Twa
	} else {
		// If oracle Price is not required for the assetOut
		inFlowTokenPrice = ExtendedPairVault.AssetOutPrice
	}

	auctionParams, found := k.GetAuctionParams(ctx, appID)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}
	if outFlowToken.Amount.GT(sdkmath.ZeroInt()) {
		err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(outFlowToken))
		if err != nil {
			return err
		}
	}
	twaData, found := k.market.GetTwa(ctx, assetOutID)
	if !found || !twaData.IsPriceActive {
		return auctiontypes.ErrorPrices
	}
	outFlowTokenPrice = twaData.Twa
	// set target amount for debt
	inFlowTokenTargetAmount := lockedVault.AmountOut
	mulfactor := sdkmath.LegacyNewDecFromInt(inFlowTokenTargetAmount).Mul(liquidationPenalty)
	inFlowTokenTargetAmount = inFlowTokenTargetAmount.Add(mulfactor.TruncateInt()).Add(lockedVault.InterestAccumulated)
	inFlowTokenTarget := sdk.NewCoin(inFlowToken.Denom, inFlowTokenTargetAmount)
	// These prices are in uusd
	outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(sdkmath.NewIntFromUint64(outFlowTokenPrice), auctionParams.Buffer)
	outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
	vaultOwner, err := sdk.AccAddressFromBech32(lockedVaultOwner)
	if err != nil {
		return err
	}
	timeNow := ctx.BlockTime()
	inFlowTokenCurrentAmount := sdk.NewCoin(inFlowToken.Denom, sdkmath.NewIntFromUint64(0))
	auction := auctiontypes.DutchAuction{
		OutflowTokenInitAmount:    outFlowToken,
		OutflowTokenCurrentAmount: outFlowToken,
		InflowTokenTargetAmount:   inFlowTokenTarget,
		InflowTokenCurrentAmount:  inFlowTokenCurrentAmount,
		OutflowTokenInitialPrice:  outFlowTokenInitialPrice,
		OutflowTokenCurrentPrice:  outFlowTokenInitialPrice,
		OutflowTokenEndPrice:      outFlowTokenEndPrice,
		InflowTokenCurrentPrice:   sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(inFlowTokenPrice)),
		StartTime:                 timeNow,
		EndTime:                   timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		AuctionStatus:             auctiontypes.AuctionStartNoBids,
		BiddingIds:                []*auctiontypes.BidOwnerMapping{},
		AuctionMappingId:          auctionParams.DutchId,
		AppId:                     appID,
		AssetInId:                 assetInID,  // cmst
		AssetOutId:                assetOutID, // cmdx
		LockedVaultId:             lockedVaultID,
		VaultOwner:                vaultOwner,
		LiquidationPenalty:        liquidationPenalty,
	}
	auction.AuctionId = k.GetAuctionID(ctx) + 1
	k.SetAuctionID(ctx, auction.AuctionId)
	err = k.SetDutchAuction(ctx, auction)
	if err != nil {
		return err
	}
	err = k.liquidation.SetFlagIsAuctionInProgress(ctx, appID, lockedVaultID, true)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			auctiontypes.EventTypeDutchNewAuction,
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

func (k Keeper) PlaceDutchAuctionBid(ctx sdk.Context, appID, auctionMappingID, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin) error {
	if bid.Amount.Equal(sdkmath.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "bid amount can't be Zero")
	}
	auction, err := k.GetDutchAuction(ctx, appID, auctionMappingID, auctionID)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "auction id %d not found", auctionID)
	}
	if bid.Denom != auction.OutflowTokenCurrentAmount.Denom {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "bid denom %s not found", bid.Denom)
	}
	if bid.Amount.GT(auction.OutflowTokenCurrentAmount.Amount) {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "bid amount can't be greater than collateral available")
	}

	// slice tells amount of collateral user should be given
	// using ceil as we need extract more from users
	outFlowTokenCurrentPrice := auction.OutflowTokenCurrentPrice // cmdx
	inFlowTokenCurrentPrice := auction.InflowTokenCurrentPrice   // cmst

	slice := bid.Amount // cmdx

	a := auction.InflowTokenTargetAmount.Amount
	b := auction.InflowTokenCurrentAmount.Amount
	tab := a.Sub(b) // leftover cmst

	// owe is $cmdx to be given to user

	owe, inFlowTokenAmount, err := k.vault.GetAmountOfOtherToken(ctx, auction.AssetOutId, outFlowTokenCurrentPrice, slice, auction.AssetInId, inFlowTokenCurrentPrice)
	if err != nil {
		return err
	}
	if inFlowTokenAmount.LTE(sdkmath.ZeroInt()) {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "Calculated Auction Amount is Zero")
	}

	TargetReachedFlag := false
	if inFlowTokenAmount.GT(tab) {
		TargetReachedFlag = true
		inFlowTokenAmount = tab // with precision

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
	ExtendedPairVault, found := k.asset.GetPairsVault(ctx, lockedVault.ExtendedPairId)
	if !found {
		return auctiontypes.ErrorInvalidExtendedPairVault
	}

	// dust is in usd * 10*-6 (uusd)
	dust := sdkmath.NewIntFromUint64(ExtendedPairVault.MinUsdValueLeft)
	// here subtracting current amount and slice to get amount left in auction and also converting it to usd * 10**-12
	outLeft, err := k.CalcDollarValueForToken(ctx, auction.AssetOutId, outFlowTokenCurrentPrice, auction.OutflowTokenCurrentAmount.Amount)
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
	if amountLeftInPUSD.LT(sdkmath.LegacyNewDecFromInt(dust)) && !amountLeftInPUSD.Equal(sdkmath.LegacyZeroDec()) && !TargetReachedFlag {
		coll := auction.OutflowTokenCurrentAmount.Amount.Uint64()
		dust := dust.Uint64()
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "Either bid all the collateral amount %d (UTOKEN) or bid amount by leaving dust greater than %d USD", coll, dust)
	}

	// Dust check for debt
	if amountLeftInPUSDforDebt.LT(sdkmath.LegacyNewDecFromInt(dust)) && !amountLeftInPUSDforDebt.Equal(sdkmath.LegacyZeroDec()) && !amountLeftInPUSD.Equal(sdkmath.LegacyZeroDec()) {
		dust := dust.Uint64()
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "Minimum amount left to be recovered should not be less than %d ", dust)
	}

	outFlowTokenCoin := sdk.NewCoin(auction.OutflowTokenInitAmount.Denom, slice)

	if inFlowTokenCoin.Amount.GT(sdkmath.ZeroInt()) {
		err = k.bank.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(inFlowTokenCoin))
		if err != nil {
			return err
		}
	}
	if outFlowTokenCoin.Amount.GT(sdkmath.ZeroInt()) {
		err = k.bank.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidder, sdk.NewCoins(outFlowTokenCoin))
		if err != nil {
			return err
		}
	}

	biddingID, err := k.CreateNewDutchBid(ctx, appID, auctionMappingID, auctionID, bidder.String(), inFlowTokenCoin, outFlowTokenCoin)
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
		// send left overcollateral to vault owner as target cmst reached and also
		vaultHolder, err := sdk.AccAddressFromBech32(lockedVault.Owner)
		if err != nil {
			panic(err)
		}

		total := auction.OutflowTokenCurrentAmount
		if total.Amount.GT(sdkmath.ZeroInt()) {
			err := k.bank.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, vaultHolder, sdk.NewCoins(total))
			if err != nil {
				return err
			}
		}

		err = k.SetDutchAuction(ctx, auction)
		if err != nil {
			return err
		}

		// remove dutch auction

		err = k.CloseDutchAuction(ctx, auction)
		if err != nil {
			return err
		}
	} else if auction.OutflowTokenCurrentAmount.Amount.IsZero() && auction.InflowTokenCurrentAmount.IsLT(auction.InflowTokenTargetAmount) { // entire collateral sold out, but debt is left
		requiredAmount := auction.InflowTokenTargetAmount.Sub(auction.InflowTokenCurrentAmount)
		_, err := k.collector.GetAmountFromCollector(ctx, auction.AppId, auction.AssetInId, requiredAmount.Amount)
		if err != nil {
			return err
		}

		// storing protocol loss
		k.SetProtocolStatistics(ctx, auction.AppId, auction.AssetInId, requiredAmount.Amount)

		err = k.SetDutchAuction(ctx, auction)
		if err != nil {
			return err
		}

		// remove dutch auction
		err = k.CloseDutchAuction(ctx, auction)
		if err != nil {
			return err
		}
	} else {
		err = k.SetDutchAuction(ctx, auction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) CreateNewDutchBid(ctx sdk.Context, appID, auctionMappingID, auctionID uint64, bidder string, outFlowTokenCoin sdk.Coin, inFlowTokenCoin sdk.Coin) (biddingID uint64, err error) {
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
	err = k.SetDutchUserBidding(ctx, bidding)
	if err != nil {
		return biddingID, err
	}
	return bidding.BiddingId, nil
}

func (k Keeper) CloseDutchAuction(
	ctx sdk.Context,
	dutchAuction auctiontypes.DutchAuction,
) error { // delete dutch biddings
	if dutchAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		for _, biddingID := range dutchAuction.BiddingIds {
			bidding, err := k.GetDutchUserBidding(ctx, biddingID.BidOwner, dutchAuction.AppId, biddingID.BidId)
			if err != nil {
				return err
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			err = k.SetDutchUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.DeleteDutchUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.SetHistoryDutchUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
		}
	}

	lockedVault, found := k.liquidation.GetLockedVault(ctx, dutchAuction.AppId, dutchAuction.LockedVaultId)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}
	// calculate penalty
	penaltyCoin := sdk.NewCoin(dutchAuction.InflowTokenCurrentAmount.Denom, sdkmath.ZeroInt())
	// burn and send target CMST to collector
	burnToken := sdk.NewCoin(dutchAuction.InflowTokenCurrentAmount.Denom, sdkmath.ZeroInt())
	burnToken.Amount = lockedVault.AmountOut
	penaltyCoin.Amount = dutchAuction.InflowTokenTargetAmount.Amount.Sub(burnToken.Amount)

	// burning
	if burnToken.Amount.GT(sdkmath.ZeroInt()) {
		err := k.bank.BurnCoins(ctx, auctiontypes.ModuleName, sdk.NewCoins(burnToken))
		if err != nil {
			return err
		}
	}

	// send penalty
	if penaltyCoin.Amount.GT(sdkmath.ZeroInt()) {
		err := k.bank.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(penaltyCoin))
		if err != nil {
			return err
		}
	}

	err := k.UpdateProtocolData(ctx, dutchAuction.OutflowTokenInitAmount, burnToken, lockedVault.ExtendedPairId)
	if err != nil {
		return err
	}

	// call increase function in collector
	err = k.collector.SetNetFeeCollectedData(ctx, dutchAuction.AppId, dutchAuction.AssetInId, penaltyCoin.Amount)
	if err != nil {
		return err
	}

	// set sell of history in locked vault
	err = k.liquidation.CreateLockedVaultHistory(ctx, lockedVault)
	if err != nil {
		return err
	}

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

	err = k.SetDutchAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}
	err = k.DeleteDutchAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}
	err = k.SetHistoryDutchAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}
	k.liquidation.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)

	return nil
}

func (k Keeper) RestartDutchAuctions(ctx sdk.Context, appID uint64) error {
	auctionParams, found := k.GetAuctionParams(ctx, appID)
	if !found {
		return nil
	}
	dutchAuctions := k.GetDutchAuctions(ctx, appID)
	// SET current price of inflow token and outflow token
	for _, dutchAuction := range dutchAuctions {
		_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
			lockedVault, found := k.liquidation.GetLockedVault(ctx, appID, dutchAuction.LockedVaultId)
			if !found {
				return auctiontypes.ErrorInvalidLockedVault
			}

			ExtendedPairVault, found := k.asset.GetPairsVault(ctx, lockedVault.ExtendedPairId)
			if !found {
				return auctiontypes.ErrorInvalidExtendedPairVault
			}
			var inFlowTokenCurrentPrice uint64
			if ExtendedPairVault.AssetOutOraclePrice {
				// If oracle Price required for the assetOut
				twaData, found := k.market.GetTwa(ctx, dutchAuction.AssetInId)
				if !found || !twaData.IsPriceActive {
					return auctiontypes.ErrorPrices
				}
				inFlowTokenCurrentPrice = twaData.Twa
			} else {
				// If oracle Price is not required for the assetOut
				inFlowTokenCurrentPrice = ExtendedPairVault.AssetOutPrice
			}
			tnume := dutchAuction.OutflowTokenInitialPrice.Mul(sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(auctionParams.AuctionDurationSeconds)))
			tdeno := dutchAuction.OutflowTokenInitialPrice.Sub(dutchAuction.OutflowTokenEndPrice)
			ntau := tnume.Quo(tdeno)
			tau := sdkmath.NewInt(ntau.TruncateInt64())
			dur := ctx.BlockTime().Sub(dutchAuction.StartTime)
			seconds := sdkmath.NewInt(int64(dur.Seconds()))
			outFlowTokenCurrentPrice := k.getPriceFromLinearDecreaseFunction(dutchAuction.OutflowTokenInitialPrice, tau, seconds)
			dutchAuction.InflowTokenCurrentPrice = sdkmath.LegacyNewDec(int64(inFlowTokenCurrentPrice))
			dutchAuction.OutflowTokenCurrentPrice = outFlowTokenCurrentPrice
			err := k.SetDutchAuction(ctx, dutchAuction)
			if err != nil {
				return err
			}
			// check if auction need to be restarted
			if ctx.BlockTime().After(dutchAuction.EndTime) {
				esmStatus, found := k.esm.GetESMStatus(ctx, lockedVault.AppId)
				status := false
				if found {
					status = esmStatus.Status
				}

				if status {
					// check user mapping of if vault exists for user
					// if not create new vault of user with cmdx cmst
					// if exists append in existing
					// close auction func call
					inflowLeft := dutchAuction.InflowTokenTargetAmount.Amount.Sub(dutchAuction.InflowTokenCurrentAmount.Amount)
					penaltyAmt := dutchAuction.InflowTokenTargetAmount.Amount.Sub(lockedVault.AmountOut)
					flag := false
					if dutchAuction.InflowTokenCurrentAmount.Amount.GTE(lockedVault.AmountOut) {
						flag = true
					}
					penaltyCoin := sdk.NewCoin(dutchAuction.InflowTokenCurrentAmount.Denom, sdkmath.ZeroInt())
					// burn and send target CMST to collector
					burnToken := sdk.NewCoin(dutchAuction.InflowTokenCurrentAmount.Denom, sdkmath.ZeroInt())
					burnToken.Amount = lockedVault.AmountOut
					penaltyCoin.Amount = dutchAuction.InflowTokenCurrentAmount.Amount.Sub(burnToken.Amount)
					vaultID, userExists := k.vault.GetUserAppExtendedPairMappingData(ctx, dutchAuction.VaultOwner.String(), dutchAuction.AppId, lockedVault.ExtendedPairId)
					if !flag {
						if userExists {
							vaultData, _ := k.vault.GetVault(ctx, vaultID.VaultId)
							if dutchAuction.OutflowTokenCurrentAmount.Amount.GT(sdkmath.ZeroInt()) {
								err := k.bank.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, vaulttypes.ModuleName, sdk.NewCoins(dutchAuction.OutflowTokenCurrentAmount))
								if err != nil {
									return err
								}
							}
							// append to existing vault
							vaultData.AmountIn = vaultData.AmountIn.Add(dutchAuction.OutflowTokenCurrentAmount.Amount)
							vaultData.AmountOut = vaultData.AmountOut.Add(inflowLeft).Sub(penaltyAmt)
							k.vault.SetVault(ctx, vaultData)
						} else {
							if dutchAuction.OutflowTokenCurrentAmount.Amount.GT(sdkmath.ZeroInt()) {
								err1 := k.bank.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, vaulttypes.ModuleName, sdk.NewCoins(dutchAuction.OutflowTokenCurrentAmount))
								if err1 != nil {
									return err1
								}
							}
							// create new vault done
							err := k.vault.CreateNewVault(ctx, dutchAuction.VaultOwner.String(), lockedVault.AppId, lockedVault.ExtendedPairId, dutchAuction.OutflowTokenCurrentAmount.Amount, inflowLeft.Sub(penaltyAmt))
							if err != nil {
								return err
							}
							length := k.vault.GetLengthOfVault(ctx)
							k.vault.SetLengthOfVault(ctx, length+1)
						}
						burnToken.Amount = dutchAuction.InflowTokenCurrentAmount.Amount
					}
					// burning
					if burnToken.Amount.GT(sdkmath.ZeroInt()) {
						err := k.bank.BurnCoins(ctx, auctiontypes.ModuleName, sdk.NewCoins(burnToken))
						if err != nil {
							return err
						}
					}
					if flag {
						// send penalty
						if penaltyCoin.Amount.GT(sdkmath.ZeroInt()) {
							err := k.bank.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(penaltyCoin))
							if err != nil {
								return err
							}
						}

						err = k.collector.SetNetFeeCollectedData(ctx, dutchAuction.AppId, dutchAuction.AssetInId, penaltyCoin.Amount)
						if err != nil {
							return err
						}
						rateIn, found := k.esm.GetSnapshotOfPrices(ctx, appID, dutchAuction.AssetOutId)
						if !found {
							return esmtypes.ErrPriceNotFound
						}
						assetData, _ := k.asset.GetAsset(ctx, dutchAuction.AssetOutId)
						var coolOffData esmtypes.DataAfterCoolOff
						coolOffData.AppId = appID
						var itemc esmtypes.AssetToAmount
						itemc.AppId = appID
						itemc.AssetID = dutchAuction.AssetOutId
						itemc.Amount = dutchAuction.OutflowTokenCurrentAmount.Amount
						itemc.IsCollateral = true
						coolOffData.CollateralTotalAmount = k.esm.CalcDollarValueOfToken(ctx, rateIn, itemc.Amount, assetData.Decimals)
						coolOffData.DebtTotalAmount = sdkmath.LegacyZeroDec()
						itemc.Share = sdkmath.LegacyOneDec()
						err := k.bank.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, esmtypes.ModuleName, sdk.NewCoins(sdk.NewCoin(dutchAuction.OutflowTokenCurrentAmount.Denom, itemc.Amount)))
						if err != nil {
							return err
						}
						k.esm.SetAssetToAmount(ctx, itemc)
						k.esm.SetDataAfterCoolOff(ctx, coolOffData)
						// send that collateral to esm data for asset
					}

					err := k.UpdateProtocolData(ctx, dutchAuction.OutflowTokenInitAmount.Sub(dutchAuction.OutflowTokenCurrentAmount), burnToken, lockedVault.ExtendedPairId)
					if err != nil {
						return err
					}

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

					err = k.SetDutchAuction(ctx, dutchAuction)
					if err != nil {
						return err
					}
					err = k.DeleteDutchAuction(ctx, dutchAuction)
					if err != nil {
						return err
					}
					err = k.SetHistoryDutchAuction(ctx, dutchAuction)
					if err != nil {
						return err
					}
					k.liquidation.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
				} else {
					twaData, found := k.market.GetTwa(ctx, dutchAuction.AssetOutId)
					if !found || !twaData.IsPriceActive {
						return auctiontypes.ErrorPrices
					}
					OutFlowTokenCurrentPrice := twaData.Twa
					timeNow := ctx.BlockTime()
					dutchAuction.StartTime = timeNow
					dutchAuction.EndTime = timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
					outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(sdkmath.NewIntFromUint64(OutFlowTokenCurrentPrice), auctionParams.Buffer)
					outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
					dutchAuction.OutflowTokenInitialPrice = outFlowTokenInitialPrice
					dutchAuction.OutflowTokenEndPrice = outFlowTokenEndPrice
					dutchAuction.OutflowTokenCurrentPrice = outFlowTokenInitialPrice
					err := k.SetDutchAuction(ctx, dutchAuction)
					if err != nil {
						return err
					}
				}
				// SET initial price fetched from market module and also end price , start time , end time
			}
			return nil
		})
	}
	return nil
}

func (k Keeper) UpdateProtocolData(ctx sdk.Context, collateralToken, burnToken sdk.Coin, extPairID uint64) error {
	ExtendedPairVault, found2 := k.asset.GetPairsVault(ctx, extPairID)
	if !found2 {
		return auctiontypes.ErrorInvalidExtendedPairVault
	}

	appExtendedPairVaultData, found3 := k.vault.GetAppExtendedPairVaultMappingData(ctx, ExtendedPairVault.AppId, ExtendedPairVault.Id)
	if !found3 {
		return sdkerrors.ErrNotFound
	}

	k.vault.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, burnToken.Amount, false)
	k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId, collateralToken.Amount, false)
	return nil
}

func (k Keeper) RestartDutch(ctx sdk.Context, appID uint64) error {
	err := k.RestartDutchAuctions(ctx, appID)
	if err != nil {
		return err
	}
	return nil
}
