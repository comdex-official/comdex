package keeper

import (
	"fmt"
	"time"

	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) DutchActivator(ctx sdk.Context) error {
	lockedVaults := k.GetLockedVaults(ctx)
	if len(lockedVaults) == 0 {
		return auctiontypes.ErrorInvalidLockedVault
	}
	for _, lockedVault := range lockedVaults {
		extendedPair, found := k.GetPairsVault(ctx, lockedVault.ExtendedPairId)
		if !found {
			return auctiontypes.ErrorInvalidPair
		}
		pair, found := k.GetPair(ctx, extendedPair.PairId)
		if !found {
			return auctiontypes.ErrorInvalidPair
		}
		assetIn, found := k.GetAsset(ctx, pair.AssetIn)
		if !found {
			return auctiontypes.ErrorAssetNotFound
		}

		assetOut, found := k.GetAsset(ctx, pair.AssetOut)
		if !found {
			return auctiontypes.ErrorAssetNotFound
		}
		assetInPrice, found := k.GetPriceForAsset(ctx, assetIn.Id)
		if !found {
			return auctiontypes.ErrorPrices
		}
		//assetInPrice is the collateral price
		////Here collateral to be auctioned is received in ucollateral*uusd so inorder to get back amount we divide with uusd of assetIn
		outflowToken := sdk.NewCoin(assetIn.Denom, lockedVault.CollateralToBeAuctioned.Quo(sdk.NewDecFromInt(sdk.NewIntFromUint64(assetInPrice))).TruncateInt())
		inflowToken := sdk.NewCoin(assetOut.Denom, sdk.ZeroInt())

		extendedPairId := lockedVault.ExtendedPairId
		ExtendedPairVault, found := k.GetPairsVault(ctx, extendedPairId)
		if !found {
			return auctiontypes.ErrorInvalidExtendedPairVault
		}
		liquidationPenalty := ExtendedPairVault.LiquidationPenalty
		if !lockedVault.IsAuctionInProgress {
			err1 := k.StartDutchAuction(ctx, outflowToken, inflowToken, lockedVault.AppMappingId, assetOut.Id, assetIn.Id, lockedVault.LockedVaultId, lockedVault.Owner, liquidationPenalty)
			if err1 != nil {
				return err1
			}
		}


	}
	return nil
}

func (k Keeper) StartDutchAuction(
	ctx sdk.Context,
	outFlowToken sdk.Coin,
	inFlowToken sdk.Coin,
	appId uint64,
	assetOutId, assetInId uint64,
	lockedVaultId uint64,
	lockedVaultOwner string,
	liquidationPenalty sdk.Dec,
) error {
	var (
		inFlowTokenPrice  uint64
		outFlowTokenPrice uint64
		// found1            bool
		found2 bool
	)

	lockedVault, found := k.GetLockedVault(ctx, lockedVaultId)
	if !found {
		return auctiontypes.ErrorInvalidLockedVault
	}

	var extendedPairVault = lockedVault.ExtendedPairId

	ExtendedPairVault, found := k.GetPairsVault(ctx, extendedPairVault)
	if !found {
		return auctiontypes.ErrorInvalidExtendedPairVault
	}
	if ExtendedPairVault.AssetOutOraclePrice {
		//If oracle Price required for the assetOut
		inFlowTokenPrice, found = k.GetPriceForAsset(ctx, assetOutId)
	} else {
		//If oracle Price is not required for the assetOut
		inFlowTokenPrice = ExtendedPairVault.AssetOutPrice

	}

	err := k.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(outFlowToken))
	if err != nil {
		return err
	}
	auctionParams, found := k.GetAuctionParams(ctx, appId)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}
	//need to get real price instead of hard coding
	//calculate target amount of cmst to collect
	if auctiontypes.TestFlag != 1 {
		// inFlowTokenPrice, found1 = k.GetPriceForAsset(ctx, assetInId)
		// if !found1 {
		// 	return auctiontypes.ErrorPrices
		// }
		outFlowTokenPrice, found2 = k.GetPriceForAsset(ctx, assetOutId)
		if !found2 {
			return auctiontypes.ErrorPrices
		}
	} else {
		outFlowTokenPrice = uint64(2)
		inFlowTokenPrice = uint64(10)
	}
	//set target amount for debt
	inFlowTokenTargetAmount := k.getInflowTokenTargetAmount(outFlowToken.Amount, sdk.NewIntFromUint64(inFlowTokenPrice), sdk.NewIntFromUint64(outFlowTokenPrice))
	inFlowTokenTarget := sdk.NewCoin(inFlowToken.Denom, inFlowTokenTargetAmount)
	//These prices are in uusd
	outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(sdk.NewIntFromUint64(outFlowTokenPrice), auctionParams.Buffer)
	outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
	vaultOwner, err := sdk.AccAddressFromBech32(lockedVaultOwner)
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
		AppId:                     appId,
		AssetInId:                 assetInId,
		AssetOutId:                assetOutId,
		LockedVaultId:             lockedVaultId,
		VaultOwner:                vaultOwner,
		LiquidationPenalty:        liquidationPenalty,
		IsLockedVaultAmountInZero: false,
	}
	auction.AuctionId = k.GetAuctionID(ctx) + 1
	k.SetAuctionID(ctx, auction.AuctionId)
	err = k.SetDutchAuction(ctx, auction)
	if err != nil {
		return err
	}
	err = k.SetFlagIsAuctionInProgress(ctx, lockedVaultId, true)
	if err != nil {
		return err
	}
	isZero, err := k.DecreaseLockedVaultAmountIn(ctx, lockedVaultId, outFlowToken.Amount)
	if err != nil {
		return err
	}
	if isZero {
		auction.IsLockedVaultAmountInZero = true
	}
	err = k.SetDutchAuction(ctx, auction)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) PlaceDutchAuctionBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, max sdk.Dec) error {
	auction, err := k.GetDutchAuction(ctx, appId, auctionMappingId, auctionId)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction id %d not found", auctionId)
	}
	max = k.GetUUSDFromUSD(ctx, max)
	if bid.Denom != auction.OutflowTokenCurrentAmount.Denom {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bid denom %s not found", bid.Denom)
	}

	//Here OutflowToken current price is in uusd and max is in uusd
	if max.LT(auction.OutflowTokenCurrentPrice.Ceil()) {
		return auctiontypes.ErrorInvalidDutchPrice
	}

	// slice tells amount of collateral user should be given
	//using ceil as we need extract more from users
	outFlowTokenCurrentPrice := auction.OutflowTokenCurrentPrice.Ceil().TruncateInt()
	inFlowTokenCurrentPrice := auction.InflowTokenCurrentPrice.Ceil().TruncateInt()
	slice := sdk.MinInt(bid.Amount, auction.OutflowTokenCurrentAmount.Amount)
	//amount in usd to be given to user
	//owe is in usd * 10**-12
	owe := slice.Mul(outFlowTokenCurrentPrice)
	//required target cmst to raise in usd * 10**-12
	//here we are multiplying each ucmdx with uusd so cmdx tokens price will be calculated amount * 10**-12
	tab := auction.InflowTokenTargetAmount.Amount.Mul(inFlowTokenCurrentPrice).Sub(auction.InflowTokenCurrentAmount.Amount.Mul(inFlowTokenCurrentPrice))

	// here price ratio is taken so no problem whether price is usd or uusd
	inFlowTokenAmount := slice.ToDec().Mul(outFlowTokenCurrentPrice.ToDec()).Quo(inFlowTokenCurrentPrice.ToDec()).Ceil().TruncateInt()
	inFlowTokenCoin := sdk.NewCoin(auction.InflowTokenTargetAmount.Denom, inFlowTokenAmount)

	//here we are getting min dust for asset to be left in uusd(usd*10**-6)
	lockedVault, found := k.GetLockedVault(ctx, auction.LockedVaultId)
	if !found {
		return auctiontypes.ErrorInvalidLockedVault
	}
	ExtendedPairVault, found := k.GetPairsVault(ctx, lockedVault.ExtendedPairId)
	if !found {
		return auctiontypes.ErrorInvalidExtendedPairVault
	}

	//dust is in usd * 10*-6 (uusd)
	dust := sdk.NewIntFromUint64(ExtendedPairVault.MinUsdValueLeft)
	//here subtracting current amount and slice to get amount left in auction and also converting it to usd * 10**-12
	amountLeftInPUSD := auction.OutflowTokenCurrentAmount.Amount.ToDec().Sub(slice.ToDec()).Mul(outFlowTokenCurrentPrice.ToDec())
	//convert amountLeft to uusd from pusd(10**-12) so we can compare dust and amountLeft in UUSD . this happens by converting ucmdx to cmdx
	amountLeft := amountLeftInPUSD.Quo(sdk.MustNewDecFromStr("1000000")).TruncateInt()

	//check if bid in usd*10**-12 is greater than required target cmst in usd*10**-12
	//if user wants to buy more than target cmst then user should be sold only required cmst amount
	//so we need to divide tab by outflow token current price and we get outflowtoken amount to be sold to user
	//if user is not buying more than required cmst then we check for dust
	//here tab is divided by price again so price ration is only considered , we are not using owe again in this function so no problem
	//As tab is the amount calculated from difference of target and current inflow token we will be using same as inflow token
	if owe.GT(tab) && !auction.IsLockedVaultAmountInZero {
		slice = tab.Quo(auction.OutflowTokenCurrentPrice.Ceil().TruncateInt())
		inFlowTokenCoin.Amount = auction.InflowTokenTargetAmount.Amount.Sub(auction.InflowTokenCurrentAmount.Amount)
	} else if amountLeft.LT(dust) && amountLeft.GT(sdk.ZeroInt()) {
		//(outflowtokenavailableamount-slice) in uusd < chost in uusd
		//so ask user to increase his bid or decrease do that more than dust is left .
		coll := auction.OutflowTokenCurrentAmount.Amount.Uint64()
		dust := dust.Uint64()
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "either bid all the amount %d (UTOKEN) or bid amount by leaving dust greater than %d UUSD", coll, dust)
	}

	outFlowTokenCoin := sdk.NewCoin(auction.OutflowTokenInitAmount.Denom, slice)

	err = k.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(inFlowTokenCoin))
	if err != nil {
		return err
	}
	err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidder, sdk.NewCoins(outFlowTokenCoin))
	if err != nil {
		return err
	}
	//create user bidding
	biddingId, err := k.CreateNewDutchBid(ctx, appId, auctionMappingId, auctionId, bidder, inFlowTokenCoin, outFlowTokenCoin)
	if err != nil {
		return err
	}
	var bidIdOwner = &auctiontypes.BidOwnerMapping{BidId: biddingId, BidOwner: bidder.String()}
	auction.BiddingIds = append(auction.BiddingIds, bidIdOwner)
	if auction.AuctionStatus == auctiontypes.AuctionStartNoBids {
		auction.AuctionStatus = auctiontypes.AuctionGoingOn
	}

	//calculate inflow amount and outflow amount if  user  transaction successfull
	auction.OutflowTokenCurrentAmount = auction.OutflowTokenCurrentAmount.Sub(outFlowTokenCoin)
	auction.InflowTokenCurrentAmount = auction.InflowTokenCurrentAmount.Add(inFlowTokenCoin)

	//collateral not over but target cmst reached then send remaining collateral to owner
	//if inflow token current amount >= InflowTokenTargetAmount
	if auction.InflowTokenCurrentAmount.IsGTE(auction.InflowTokenTargetAmount) && !auction.IsLockedVaultAmountInZero {
		//send left overcollateral to vault owner as target cmst reached and also

		total := auction.OutflowTokenCurrentAmount
		err := k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, vaulttypes.ModuleName, sdk.NewCoins(total))
		if err != nil {
			return err
		}
		err = k.IncreaseLockedVaultAmountIn(ctx, auction.LockedVaultId, total.Amount)
		if err != nil {
			return err
		}
		err = k.SetDutchAuction(ctx, auction)
		if err != nil {
			return err
		}
		//remove dutch auction

		err = k.CloseDutchAuction(ctx, auction)
		if err != nil {
			return err
		}
	} else if auction.OutflowTokenCurrentAmount.Amount.IsZero() { //entire collateral sold out

		err = k.SetDutchAuction(ctx, auction)
		if err != nil {
			return err
		}
		//remove dutch auction

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

func (k Keeper) CreateNewDutchBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, outFlowTokenCoin sdk.Coin, inFlowTokenCoin sdk.Coin) (biddingId uint64, err error) {
	bidding := auctiontypes.DutchBiddings{
		BiddingId:          k.GetUserBiddingID(ctx) + 1,
		AuctionId:          auctionId,
		AuctionStatus:      auctiontypes.ActiveAuctionStatus,
		Bidder:             bidder.String(),
		OutflowTokenAmount: outFlowTokenCoin,
		InflowTokenAmount:  inFlowTokenCoin,
		BiddingTimestamp:   ctx.BlockTime(),
		BiddingStatus:      auctiontypes.SuccessBiddingStatus,
		AppId:              appId,
		AuctionMappingId:   auctionMappingId,
	}
	k.SetUserBiddingID(ctx, bidding.BiddingId)
	err = k.SetDutchUserBidding(ctx, bidding)
	if err != nil {
		return biddingId, err
	}
	return bidding.BiddingId, nil
}

func (k Keeper) CloseDutchAuction(
	ctx sdk.Context,
	dutchAuction auctiontypes.DutchAuction,
) error {

	//delete dutch biddings
	if dutchAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		for _, biddingId := range dutchAuction.BiddingIds {
			bidding, err := k.GetDutchUserBidding(ctx, biddingId.BidOwner, dutchAuction.AppId, biddingId.BidId)
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

	lockedVault, found := k.GetLockedVault(ctx, dutchAuction.LockedVaultId)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}

	// burn and send target CMST to collector
	burnToken := sdk.NewCoin(dutchAuction.InflowTokenCurrentAmount.Denom, sdk.ZeroInt())
	//doing burn amount  = inflowtokencurrentamount / (1 + liq_penalty)
	burnToken.Amount = burnToken.Amount.Add(k.getBurnAmount(dutchAuction.InflowTokenCurrentAmount.Amount, dutchAuction.LiquidationPenalty))
	//calculate penalty
	penaltyAmount := dutchAuction.InflowTokenCurrentAmount.Amount.Sub(burnToken.Amount)
	//if amountInZero is true
	//if burnAmount is greater than amount out
	//add burnAmount-amountout out to penalty
	//make burn amount = amountout

	//if burnAmount is less than amount out
	// get amountout - burnamount from collector
	// make burnamount = amountout
	if dutchAuction.IsLockedVaultAmountInZero {
		if burnToken.Amount.GT(lockedVault.AmountOut) {

			penaltyAmount = penaltyAmount.Add(burnToken.Amount.Sub(lockedVault.AmountOut))
			burnToken.Amount = lockedVault.AmountOut
		} else if burnToken.Amount.LT(lockedVault.AmountOut) {

			//Transfer balance from collector module to auction module
			requiredAmount := lockedVault.AmountOut.Sub(burnToken.Amount)

			_, err := k.GetAmountFromCollector(ctx, dutchAuction.AppId, dutchAuction.AssetInId, requiredAmount)
			if err != nil {

				return err
			}

			//storing protocol loss
			k.SetProtocolStatistics(ctx, dutchAuction.AppId, dutchAuction.AssetInId, requiredAmount)
			burnToken.Amount = lockedVault.AmountOut
		}
	}
	fmt.Println("before burn tokens in auction")
	fmt.Println(k.GetModuleAccountBalance(ctx, auctiontypes.ModuleName, "ucmst"))
	fmt.Println("BURN TOKEN 1 : ", burnToken)
	//burn the burn tokens
	fmt.Println("penalty",penaltyAmount)
	burnToken.Amount = k.GetModuleAccountBalance(ctx, auctiontypes.ModuleName, "ucmst").Sub(penaltyAmount)
	err := k.BurnCoins(ctx, auctiontypes.ModuleName, burnToken)
	if err != nil {
		return err
	}
	fmt.Println("after burn tokens in auction")
	fmt.Println(k.GetModuleAccountBalance(ctx, auctiontypes.ModuleName, "ucmst"))
	ExtendedPairVault, found := k.GetPairsVault(ctx, lockedVault.ExtendedPairId)
	if !found {
		return auctiontypes.ErrorInvalidExtendedPairVault
	}

	appExtendedPairVaultData, found := k.GetAppExtendedPairVaultMapping(ctx, ExtendedPairVault.AppMappingId)
	if !found {
		return sdkerrors.ErrNotFound
	}
	fmt.Println("BURN TOKEN 1 : ", burnToken)
	fmt.Println("appExtendedPairVaultData : ", appExtendedPairVaultData)
	fmt.Println("ExtendedPairVault")
	k.UpdateTokenMintedAmountLockerMapping(ctx, appExtendedPairVaultData, ExtendedPairVault.Id, burnToken.Amount, false)

	//send penalty
	err = k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(burnToken.Denom, penaltyAmount)))
	if err != nil {
		return err
	}
	//call increase function in collector
	err = k.SetNetFeeCollectedData(ctx, dutchAuction.AppId, dutchAuction.AssetInId, penaltyAmount)
	if err != nil {
		return err
	}
	lockedVault.AmountOut = lockedVault.AmountOut.Sub(burnToken.Amount)
	lockedVault.UpdatedAmountOut = lockedVault.UpdatedAmountOut.Sub(burnToken.Amount)

	//set sell of history in locked vault
	outFlowToken := dutchAuction.OutflowTokenInitAmount.Sub(dutchAuction.OutflowTokenCurrentAmount)
	sellOfHistory := outFlowToken.String() + dutchAuction.InflowTokenCurrentAmount.String()
	lockedVault.SellOffHistory = append(lockedVault.SellOffHistory, sellOfHistory)

	k.SetLockedVault(ctx, lockedVault)

	dutchAuction.AuctionStatus = auctiontypes.AuctionEnded

	fmt.Println("outflow token --------",outFlowToken)
	fmt.Println("burn token -----------",burnToken)
	k.UpdateCollateralLockedAmountLockerMapping(ctx, appExtendedPairVaultData, ExtendedPairVault.Id, outFlowToken.Amount, false)
	//update locked vault
	err = k.SetFlagIsAuctionComplete(ctx, dutchAuction.LockedVaultId, true)
	if err != nil {
		return err
	}

	err = k.SetFlagIsAuctionInProgress(ctx, dutchAuction.LockedVaultId, false)
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

	return nil
}



func (k Keeper) RestartDutchAuctions(ctx sdk.Context, appId uint64) error {
	dutchAuctions := k.GetDutchAuctions(ctx, appId)
	auctionParams, found := k.GetAuctionParams(ctx, appId)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}
	// SET current price of inflow token and outflow token
	for _, dutchAuction := range dutchAuctions {
		lockedVault, found := k.GetLockedVault(ctx, dutchAuction.LockedVaultId)
		if !found {
			return auctiontypes.ErrorInvalidLockedVault
		}

		ExtendedPairVault, found := k.GetPairsVault(ctx, lockedVault.ExtendedPairId)
		if !found {
			return auctiontypes.ErrorInvalidExtendedPairVault
		}
		var inFlowTokenCurrentPrice uint64
		if ExtendedPairVault.AssetOutOraclePrice {
			//If oracle Price required for the assetOut
			inFlowTokenCurrentPrice, found = k.GetPriceForAsset(ctx, dutchAuction.AssetInId)
			if !found {
				return auctiontypes.ErrorPrices
			}
		} else {
			//If oracle Price is not required for the assetOut
			inFlowTokenCurrentPrice = ExtendedPairVault.AssetOutPrice

		}
		//inFlowTokenCurrentPrice := sdk.MustNewDecFromStr("1")
		dutchAuction.InflowTokenCurrentPrice = sdk.NewDec(int64(inFlowTokenCurrentPrice))
		tau := sdk.NewInt(int64(auctionParams.AuctionDurationSeconds))
		dur := ctx.BlockTime().Sub(dutchAuction.StartTime)
		seconds := sdk.NewInt(int64(dur.Seconds()))
		outFlowTokenCurrentPrice := k.getPriceFromLinearDecreaseFunction(dutchAuction.OutflowTokenInitialPrice, tau, seconds)
		dutchAuction.OutflowTokenCurrentPrice = outFlowTokenCurrentPrice

		//check if auction need to be restarted
		if ctx.BlockTime().After(dutchAuction.EndTime) || outFlowTokenCurrentPrice.LT(dutchAuction.OutflowTokenEndPrice) {
			//SET initial price fetched from market module and also end price , start time , end time
			//outFlowTokenCurrentPrice := sdk.NewIntFromUint64(10)
			outFlowTokenCurrentPrice, found := k.GetPriceForAsset(ctx, dutchAuction.AssetOutId)
			if !found {
				return auctiontypes.ErrorPrices
			}
			timeNow := ctx.BlockTime()
			dutchAuction.StartTime = timeNow
			dutchAuction.EndTime = timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
			outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(sdk.NewIntFromUint64(outFlowTokenCurrentPrice), auctionParams.Buffer)
			outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
			dutchAuction.OutflowTokenInitialPrice = outFlowTokenInitialPrice
			dutchAuction.OutflowTokenEndPrice = outFlowTokenEndPrice
			dutchAuction.OutflowTokenCurrentPrice = outFlowTokenInitialPrice
		}
		err := k.SetDutchAuction(ctx, dutchAuction)
		if err != nil {
			return err
		}
	}
	return nil
}