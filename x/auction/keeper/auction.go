package keeper

import (
	"time"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) checkStatusOfNetFeesCollectedAndStartAuction(ctx sdk.Context, appId, assetId uint64, assetToAuction collectortypes.AssetIdToAuctionLookupTable) (status uint64) {
	assetsCollectorDataUnderAppId, found := k.GetAppidToAssetCollectorMapping(ctx, appId)
	if !found {
		return
	}
	//assetCollector has netFeesCollected
	for _, assetCollector := range assetsCollectorDataUnderAppId.AssetCollector {
		if assetCollector.AssetId == assetId {
			//collectorLookupTable has surplusThreshhold for all assets
			collectorLookupTable, found := k.GetCollectorLookupTable(ctx, appId)
			if !found {
				return auctiontypes.NoAuction
			}
			for _, collector := range collectorLookupTable.AssetrateInfo {
				if collector.CollectorAssetId == assetId {
					// can extract both surplusThreshhold and netFeesCollected
					if assetCollector.Collector.NetFeesCollected.LTE(sdk.NewIntFromUint64(collector.DebtThreshold-collector.LotSize)) && assetToAuction.IsDebtAuction {
						//TODO START DEBT AUCTION .  LOTSIZE AS MINTED FOR SECONDARY ASSET and ACCEPT Collector assetid from user
						outflowAsset, found1 := k.asset.GetAsset(ctx, collector.SecondaryAssetId)
						inflowAsset, found2 := k.asset.GetAsset(ctx, collector.CollectorAssetId)
						if !found1 || !found2 {
							return auctiontypes.NoAuction
						}
						outflowToken := sdk.NewCoin(outflowAsset.Denom, sdk.NewInt(int64(collector.LotSize)))
						outflowTokenPrice, found3 := k.market.GetPriceForAsset(ctx, assetId)
						inflowTokenPrice, found4 := k.market.GetPriceForAsset(ctx, collector.SecondaryAssetId)
						if !found3 || !found4 {
							return auctiontypes.NoAuction
						}
						inflowTokenAmount := outflowToken.Amount.Mul(sdk.NewInt(int64(outflowTokenPrice)).Quo(sdk.NewInt(int64(inflowTokenPrice))))
						inflowToken := sdk.NewCoin(inflowAsset.Denom, inflowTokenAmount)
						//Mint the tokens when collector module sends tokens to user
						k.StartDebtAuction(ctx, outflowToken, inflowToken, appId, assetId)
						return auctiontypes.StartedDebtAuction
					} else if assetCollector.Collector.NetFeesCollected.GTE(sdk.NewIntFromUint64(collector.SurplusThreshold+collector.LotSize)) && assetToAuction.IsDebtAuction {
						//TODO START SURPLUS AUCTION .  WITH COLLECTOR ASSET ID AS token given to user of lot size and secondary asset received from user and burnt and bid factor
						outflowAsset, found1 := k.asset.GetAsset(ctx, collector.CollectorAssetId)
						inflowAsset, found2 := k.asset.GetAsset(ctx, collector.SecondaryAssetId)
						if !found1 || !found2 {
							return auctiontypes.NoAuction
						}
						outflowToken := sdk.NewCoin(outflowAsset.Denom, sdk.NewInt(int64(collector.LotSize)))
						outflowTokenPrice, found3 := k.market.GetPriceForAsset(ctx, assetId)
						inflowTokenPrice, found4 := k.market.GetPriceForAsset(ctx, collector.SecondaryAssetId)
						if !found3 || !found4 {
							return auctiontypes.NoAuction
						}
						inflowTokenAmount := outflowToken.Amount.Mul(sdk.NewInt(int64(outflowTokenPrice)).Quo(sdk.NewInt(int64(inflowTokenPrice))))
						inflowToken := sdk.NewCoin(inflowAsset.Denom, inflowTokenAmount)
						//Transfer balance from collector module to auction module
						k.bank.SendCoinsFromModuleToModule(ctx, collectortypes.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(outflowToken))
						k.StartSurplusAuction(ctx, outflowToken, inflowToken, *collector.BidFactor, appId, assetId)
						// TODO store netfeesaccumulated
						return auctiontypes.StartedSurplusAuction
					} else {
						return auctiontypes.NoAuction
					}
				}
			}
		}
	}
	return auctiontypes.NoAuction
}

func (k Keeper) CreateSurplusAndDebtAuctions(ctx sdk.Context) error {
	appIds, found := k.GetApps(ctx)
	if !found {
		return assettypes.AppIdsDoesntExist
	}
	for _, appId := range appIds {
		//check if auction status for an asset is false
		auctionLookupTable, found := k.GetCollectorAuctionLookupTable(ctx, appId.Id)
		if !found {
			continue
		}
		for _, assetToAuction := range auctionLookupTable.AssetIdToAuctionLookup {
			if assetToAuction.IsSurplusAuction || assetToAuction.IsDebtAuction {
				if !assetToAuction.IsAuctionActive {
					status := k.checkStatusOfNetFeesCollectedAndStartAuction(ctx, appId.Id, assetToAuction.AssetId, assetToAuction)
					if status == auctiontypes.StartedDebtAuction {
						assetToAuction.IsAuctionActive = true
					} else if status == auctiontypes.StartedSurplusAuction {
						assetToAuction.IsAuctionActive = true
					}
					err := k.SetCollectorAuctionLookupTable(ctx, auctionLookupTable)
					if err == nil {
						continue
					}
				}
			}
		}
	}
	return nil
}

func (k Keeper) makeFalseForFlags(ctx sdk.Context, appId, assetId uint64) {

	auctionLookupTable, found := k.GetCollectorAuctionLookupTable(ctx, appId)
	if !found {
		return
	}
	for _, assetToAuction := range auctionLookupTable.AssetIdToAuctionLookup {
		if assetToAuction.AssetId == assetId {
			assetToAuction.IsAuctionActive = false
			err := k.SetCollectorAuctionLookupTable(ctx, auctionLookupTable)
			if err == nil {
				break
			}
		}
	}
}

//func (k Keeper) CreateNewAuctions(ctx sdk.Context) {
//	locked_vaults := k.GetLockedVaults(ctx)
//	for _, locked_vault := range locked_vaults {
//		pair, found := k.GetPair(ctx, locked_vault.PairId)
//		if !found {
//			continue
//		}
//		assetIn, found := k.GetAsset(ctx, pair.AssetIn)
//		if !found {
//			continue
//		}
//
//		assetOut, found := k.GetAsset(ctx, pair.AssetOut)
//		if !found {
//			continue
//		}
//		collateralizationRatio, err := k.CalculateCollaterlizationRatio(ctx, locked_vault.AmountIn, assetIn, locked_vault.AmountOut, assetOut)
//		if err != nil {
//			continue
//		}
//		//if sdk.Dec.LT(collateralizationRatio, pair.LiquidationRatio) && !locked_vault.IsAuctionInProgress {
//		//	k.StartCollateralAuction(ctx, locked_vault, pair, assetIn, assetOut)
//		//}
//	}
//}

func (k Keeper) CloseAuctions(ctx sdk.Context) {
	collateral_auctions := k.GetSurplusAuctions(ctx)
	for _, collateral_auction := range collateral_auctions {
		if ctx.BlockTime().After(collateral_auction.EndTime) {
			k.CloseSurplusAuction(ctx, collateral_auction)
		}
	}
}

func (k Keeper) CloseDebtAuctions(ctx sdk.Context) {
	debtAuctions := k.GetDebtAuctions(ctx)
	for _, debtAuction := range debtAuctions {
		if ctx.BlockTime().After(debtAuction.EndTime) {
			k.CloseDebtAuction(ctx, debtAuction)
		}
	}
}

func (k Keeper) RestartDutchAuctions(ctx sdk.Context) {
	dutchAuctions := k.GetDutchAuctions(ctx)
	auctionParams := k.GetParams(ctx)
	// SET current price of inflow token and outflow token
	for _, dutchAuction := range dutchAuctions {
		//TODO get price of dai from market module and set it in InflowTokenCurrentPrice of auction

		inFlowTokenCurrentPrice := sdk.MustNewDecFromStr("1")
		dutchAuction.InflowTokenCurrentPrice = inFlowTokenCurrentPrice
		tau := sdk.NewInt(int64(auctionParams.AuctionDurationSeconds))
		dur := ctx.BlockTime().Sub(dutchAuction.StartTime)
		seconds := sdk.NewInt(int64(dur.Seconds()))
		outFlowTokenCurrentPrice := k.getPriceFromLinearDecreaseFunction(dutchAuction.OutflowTokenInitialPrice, tau, seconds)
		dutchAuction.OutflowTokenCurrentPrice = outFlowTokenCurrentPrice

		//check if auction need to be restarted
		if ctx.BlockTime().After(dutchAuction.EndTime) || outFlowTokenCurrentPrice.LT(dutchAuction.OutflowTokenEndPrice) {
			//SET initial price fetched from market module and also end price , start time , end time
			//TODO get price of outflowtoken from market module
			outFlowTokenCurrentPrice := sdk.MustNewDecFromStr("10")
			timeNow := ctx.BlockTime()
			dutchAuction.StartTime = timeNow
			dutchAuction.EndTime = timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
			outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(outFlowTokenCurrentPrice, auctionParams.Buffer)
			outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
			dutchAuction.OutflowTokenInitialPrice = outFlowTokenInitialPrice
			dutchAuction.OutflowTokenEndPrice = outFlowTokenEndPrice
			dutchAuction.OutflowTokenCurrentPrice = outFlowTokenInitialPrice
		}
		k.SetDutchAuction(ctx, dutchAuction)
	}
}

func (k Keeper) StartSurplusAuction(
	ctx sdk.Context,
	outflowToken sdk.Coin,
	inflowToken sdk.Coin,
	bidFactor sdk.Dec,
	appId, assetId uint64,
) error {

	auctionParams := k.GetParams(ctx)

	auction := auctiontypes.SurplusAuction{
		OutflowToken:    outflowToken,
		InflowToken:     inflowToken,
		ActiveBiddingId: 0,
		Bidder:          nil,
		Bid:             inflowToken,
		EndTime:         ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		BidFactor:       bidFactor,
		BiddingIds:      []uint64{},
		AuctionStatus:   auctiontypes.AuctionStartNoBids,
		AppId:           appId,
		AssetId:         assetId,
	}
	auction.Id = k.GetSurplusAuctionID(ctx) + 1
	k.SetSurplusAuctionID(ctx, auction.Id)
	k.SetSurplusAuction(ctx, auction)
	return nil
}

func (k Keeper) StartDebtAuction(
	ctx sdk.Context,
	auctionToken sdk.Coin,
	expectedUserToken sdk.Coin,
	appId, assetId uint64,
) error {

	auctionParams := k.GetParams(ctx)
	auction := auctiontypes.DebtAuction{
		AuctionedToken:      auctionToken,
		ExpectedMintedToken: auctionToken,
		ExpectedUserToken:   expectedUserToken,
		ActiveBiddingId:     0,
		Bidder:              nil,
		EndTime:             ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		CurrentBidAmount:    sdk.NewCoin(auctionToken.Denom, sdk.NewInt(0)),
		AuctionStatus:       auctiontypes.AuctionStartNoBids,
		AppId:               appId,
		AssetId:             assetId,
	}
	auction.AuctionId = k.GetDebtAuctionID(ctx) + 1
	k.SetDebtAuctionID(ctx, auction.AuctionId)
	k.SetDebtAuction(ctx, auction)
	return nil
}

func (k Keeper) StartDutchAuction(
	ctx sdk.Context,
	outFlowToken sdk.Coin,
	inFlowToken sdk.Coin,
	outFlowTokenAddress string,
	inFlowTokenAddress string,
) error {

	outFlowTokenAddress1, err := sdk.AccAddressFromBech32(outFlowTokenAddress)
	inFlowTokenAddress1, err := sdk.AccAddressFromBech32(inFlowTokenAddress)
	if err != nil {
		return auctiontypes.ErrorInvalidAddress
	}
	outFlowTokenBalance := k.bank.GetBalance(ctx, outFlowTokenAddress1, outFlowToken.Denom)
	if outFlowTokenBalance.Amount.LT(outFlowToken.Amount) {
		return auctiontypes.ErrorInvalidOutFlowTokenBalance
	}
	auctionParams := k.GetParams(ctx)
	//TODO need to get real price instead of hard coding
	outFlowTokenPrice := sdk.MustNewDecFromStr("100")
	inFlowTokenPrice := sdk.MustNewDecFromStr("1")
	outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(outFlowTokenPrice, auctionParams.Buffer)
	outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
	timeNow := ctx.BlockTime()
	auction := auctiontypes.DutchAuction{
		OutflowTokenInitAmount:    outFlowToken,
		OutflowTokenCurrentAmount: outFlowToken,
		InflowTokenTargetAmount:   inFlowToken,
		InflowTokenCurrentAmount:  inFlowToken,
		OutflowTokenInitialPrice:  outFlowTokenInitialPrice,
		OutflowTokenCurrentPrice:  outFlowTokenInitialPrice,
		OutflowTokenEndPrice:      outFlowTokenEndPrice,
		InflowTokenCurrentPrice:   inFlowTokenPrice,
		StartTime:                 timeNow,
		EndTime:                   timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		AuctionStatus:             auctiontypes.AuctionStartNoBids,
		InflowTokenAddress:        outFlowTokenAddress1,
		OutflowTokenAddress:       inFlowTokenAddress1,
	}
	auction.AuctionId = k.GetDutchAuctionID(ctx) + 1
	k.SetDebtAuctionID(ctx, auction.AuctionId)
	k.SetDutchAuction(ctx, auction)
	return nil
}

func (k Keeper) CloseSurplusAuction(
	ctx sdk.Context,
	surplusAuction auctiontypes.SurplusAuction,
) error {

	if surplusAuction.Bidder != nil && surplusAuction.Bid.Amount.GTE(surplusAuction.Bid.Amount) {

		highestBidReceived := surplusAuction.Bid

		err := k.bank.SendCoinsFromModuleToAccount(ctx, collectortypes.ModuleName, surplusAuction.Bidder, sdk.NewCoins(highestBidReceived))
		if err != nil {
			return err
		}
		bidding, _ := k.GetSurplusBidding(ctx, surplusAuction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.SuccessBiddingStatus
		k.SetSurplusBidding(ctx, bidding)
		//TODO decide how to burn tokens
		k.bank.BurnCoins(ctx, collectortypes.ModuleName, sdk.NewCoins(surplusAuction.Bid))

		for _, biddingId := range surplusAuction.BiddingIds {
			bidding, found := k.GetSurplusBidding(ctx, biddingId)
			if !found {
				continue
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			k.SetSurplusBidding(ctx, bidding)
		}
	}
	k.makeFalseForFlags(ctx, surplusAuction.AppId, surplusAuction.AssetId)
	k.DeleteSurplusAuction(ctx, surplusAuction.Id)
	return nil
}

func (k Keeper) CloseDebtAuction(
	ctx sdk.Context,
	debtAuction auctiontypes.DebtAuction,
) error {

	if debtAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		k.bank.MintCoins(ctx, collectortypes.ModuleName, sdk.NewCoins(debtAuction.CurrentBidAmount))
		err := k.bank.SendCoinsFromModuleToAccount(ctx, collectortypes.ModuleName, debtAuction.Bidder, sdk.NewCoins(debtAuction.CurrentBidAmount))
		if err != nil {
			return err
		}
		bidding, _ := k.GetDebtBidding(ctx, debtAuction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.SuccessBiddingStatus
		k.SetDebtBidding(ctx, bidding)

		//for _, biddingId := range debtAuction.BiddingIds {
		//	bidding, found := k.GetBidding(ctx, biddingId)
		//	if !found {
		//		continue
		//	}
		//	bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
		//	k.SetBidding(ctx, bidding)
		//}
	}
	k.bank.SendCoinsFromModuleToModule(ctx, collectortypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(debtAuction.ExpectedUserToken))
	k.makeFalseForFlags(ctx, debtAuction.AppId, debtAuction.AssetId)
	k.DeleteDebtAuction(ctx, debtAuction.AuctionId)
	return nil
}

func (k Keeper) CloseDutchAuction(
	ctx sdk.Context,
	dutchtAuction auctiontypes.DutchAuction,
) error {

	if dutchtAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		//for _, biddingId := range debtAuction.BiddingIds {
		//	bidding, found := k.GetBidding(ctx, biddingId)
		//	if !found {
		//		continue
		//	}
		//	bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
		//	k.SetBidding(ctx, bidding)
		//}
	}
	k.DeleteDutchAuction(ctx, dutchtAuction.AuctionId)
	return nil
}

func (k Keeper) CreateNewSurplusBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) (biddingId uint64, err error) {
	auction, found := k.GetSurplusAuction(ctx, auctionId)
	if !found {
		return 0, auctiontypes.ErrorInvalidSurplusAuctionId
	}
	bidding := auctiontypes.Biddings{
		Id:                  k.GetSurplusBiddingID(ctx) + 1,
		AuctionId:           auctionId,
		AuctionStatus:       auctiontypes.ActiveAuctionStatus,
		AuctionedCollateral: auction.OutflowToken,
		Bidder:              bidder.String(),
		Bid:                 bid,
		BiddingTimestamp:    ctx.BlockTime(),
		BiddingStatus:       auctiontypes.PlacedBiddingStatus,
	}
	k.SetSurplusBiddingID(ctx, bidding.Id)
	k.SetSurplusBidding(ctx, bidding)

	userBiddings, found := k.GetSurplusUserBiddings(ctx, bidder.String())
	if !found {
		userBiddings = auctiontypes.UserBiddings{
			Id:         k.GetSurplusUserBiddingID(ctx) + 1,
			Bidder:     bidder.String(),
			BiddingIds: []uint64{},
		}
		k.SetSurplusUserBiddingID(ctx, userBiddings.Id)
	}
	userBiddings.BiddingIds = append(userBiddings.BiddingIds, bidding.Id)
	k.SetSurplusUserBidding(ctx, userBiddings)
	return bidding.Id, nil
}

func (k Keeper) CreateNewDebtBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) (biddingId uint64, err error) {
	bidding := auctiontypes.Biddings{
		Id:               k.GetDebtBiddingID(ctx) + 1,
		AuctionId:        auctionId,
		AuctionStatus:    auctiontypes.ActiveAuctionStatus,
		Bidder:           bidder.String(),
		Bid:              bid,
		BiddingTimestamp: ctx.BlockTime(),
		BiddingStatus:    auctiontypes.PlacedBiddingStatus,
	}
	k.SetDebtBiddingID(ctx, bidding.Id)
	k.SetDebtBidding(ctx, bidding)

	userBiddings, found := k.GetDebtUserBiddings(ctx, bidder.String())
	if !found {
		userBiddings = auctiontypes.UserBiddings{
			Id:         k.GetDebtUserBiddingID(ctx) + 1,
			Bidder:     bidder.String(),
			BiddingIds: []uint64{},
		}
		k.SetDebtUserBiddingID(ctx, userBiddings.Id)
	}
	userBiddings.BiddingIds = append(userBiddings.BiddingIds, bidding.Id)
	k.SetDebtUserBidding(ctx, userBiddings)
	return bidding.Id, nil
}

func (k Keeper) CreateNewDutchBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, outFlowTokenCoin sdk.Coin, inFlowTokenCoin sdk.Coin) (biddingId uint64) {
	bidding := auctiontypes.DutchBiddings{
		BiddingId:          k.GetDutchBiddingID(ctx) + 1,
		AuctionId:          auctionId,
		AuctionStatus:      auctiontypes.ActiveAuctionStatus,
		Bidder:             bidder.String(),
		OutflowTokenAmount: outFlowTokenCoin,
		InflowTokenAmount:  inFlowTokenCoin,
		BiddingTimestamp:   ctx.BlockTime(),
		BiddingStatus:      auctiontypes.SuccessBiddingStatus,
	}
	k.SetDutchBiddingID(ctx, bidding.BiddingId)
	k.SetDutchBidding(ctx, bidding)

	userDutchBiddings, found := k.GetDutchUserBiddings(ctx, bidder.String())
	if !found {
		userDutchBiddings = auctiontypes.UserBiddings{
			Id:         k.GetDutchUserBiddingID(ctx) + 1,
			Bidder:     bidder.String(),
			BiddingIds: []uint64{},
		}
		k.SetDutchUserBiddingID(ctx, userDutchBiddings.Id)
	}
	userDutchBiddings.BiddingIds = append(userDutchBiddings.BiddingIds, bidding.BiddingId)
	k.SetDutchUserBidding(ctx, userDutchBiddings)
	return bidding.BiddingId
}

func (k Keeper) PlaceSurplusBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) error {
	auction, found := k.GetSurplusAuction(ctx, auctionId)
	if !found {
		return auctiontypes.ErrorInvalidSurplusAuctionId
	}
	if bid.Denom != auction.InflowToken.Denom {
		return auctiontypes.ErrorInvalidBiddingDenom
	}
	//Test this multiplication
	minBidCoin := sdk.NewCoin(bid.Denom, auction.BidFactor.MulInt(auction.Bid.Amount).Ceil().TruncateInt())
	if bid.Amount.LT(minBidCoin.Amount) {
		return auctiontypes.ErrorLowBidAmount
	}

	err := k.SendCoinsFromAccountToModule(ctx, bidder, collectortypes.ModuleName, sdk.NewCoins(bid))
	if err != nil {
		return err
	}
	biddingId, err := k.CreateNewSurplusBid(ctx, auctionId, bidder, bid)
	if err != nil {
		return err
	}
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		// auction.Bidder as previous bidder
		err = k.bank.SendCoinsFromModuleToAccount(ctx, collectortypes.ModuleName, auction.Bidder, sdk.NewCoins(auction.Bid))
		if err != nil {
			return err
		}
		bidding, _ := k.GetSurplusBidding(ctx, auction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.RejectedBiddingStatus
		k.SetSurplusBidding(ctx, bidding)
	} else {
		auction.AuctionStatus = auctiontypes.AuctionGoingOn
	}
	auction.ActiveBiddingId = biddingId
	auction.BiddingIds = append(auction.BiddingIds, biddingId)
	auction.Bidder = bidder
	auction.Bid = bid
	k.SetSurplusAuction(ctx, auction)
	return nil
}

func (k Keeper) PlaceDebtBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, expectedUserToken sdk.Coin) error {
	auction, found := k.GetDebtAuction(ctx, auctionId)
	if !found {
		return auctiontypes.ErrorInvalidDebtAuctionId
	}
	if expectedUserToken.Denom != auction.ExpectedUserToken.Denom {
		return auctiontypes.ErrorInvalidDebtUserExpectedDenom
	}
	if expectedUserToken.Amount.Equal(auction.ExpectedUserToken.Amount) == false {
		return auctiontypes.ErrorDebtExpectedUserAmount
	}
	if bid.Denom != auction.ExpectedMintedToken.Denom {
		return auctiontypes.ErrorInvalidDebtMintedDenom
	}
	if bid.Amount.GT(auction.ExpectedMintedToken.Amount) {
		return auctiontypes.ErrorDebtMoreBidAmount
	}

	err := k.SendCoinsFromAccountToModule(ctx, bidder, collectortypes.ModuleName, sdk.NewCoins(expectedUserToken))
	if err != nil {
		return err
	}
	biddingId, err := k.CreateNewDebtBid(ctx, auctionId, bidder, bid)
	if err != nil {
		return err
	}
	//If auction gets bid from second time onwards
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		err = k.bank.SendCoinsFromModuleToAccount(ctx, collectortypes.ModuleName, auction.Bidder, sdk.NewCoins(auction.ExpectedUserToken))
		if err != nil {
			return err
		}
		bidding, _ := k.GetDebtBidding(ctx, auction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.RejectedBiddingStatus
		k.SetDebtBidding(ctx, bidding)
	} else {
		auction.AuctionStatus = auctiontypes.AuctionGoingOn
	}
	auction.ActiveBiddingId = biddingId
	auction.Bidder = bidder
	auction.CurrentBidAmount = bid
	//decreasing expected minted token for next bid
	expectedMintedToken := sdk.NewDecFromInt(auction.ExpectedMintedToken.Amount)
	decreaseAmount := expectedMintedToken.Mul(auctiontypes.DefaultDebtMintTokenDecreasePercentage)
	expectedMintedToken = expectedMintedToken.Sub(decreaseAmount).Ceil() // As of now ceiling is done
	auction.ExpectedMintedToken = sdk.NewCoin(auction.ExpectedMintedToken.Denom, expectedMintedToken.TruncateInt())
	k.SetDebtAuction(ctx, auction)
	return nil
}

func (k Keeper) PlaceDutchBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, max sdk.Dec) error {
	auction, found := k.GetDutchAuction(ctx, auctionId)
	if !found {
		return auctiontypes.ErrorInvalidDutchAuctionId
	}

	if bid.Denom != auction.OutflowTokenCurrentAmount.Denom {
		return auctiontypes.ErrorInvalidDutchUserbidDenom
	}

	if max.LT(auction.OutflowTokenCurrentPrice) {
		return auctiontypes.ErrorInvalidDutchPrice
	}
	// slice tells amount of collateral user should be given
	auctionParams := k.GetParams(ctx)
	outFlowTokenCurrentPrice := sdk.NewIntFromBigInt(auction.OutflowTokenCurrentPrice.BigInt())
	inFlowTokenCurrentPrice := sdk.NewIntFromBigInt(auction.InflowTokenCurrentPrice.BigInt())
	slice := sdk.MinInt(bid.Amount, auction.OutflowTokenCurrentAmount.Amount)
	owe := slice.Mul(outFlowTokenCurrentPrice)
	tab := auction.InflowTokenCurrentAmount.Amount.Mul(inFlowTokenCurrentPrice)
	//check if bid is greater than target dai
	if owe.GT(tab) {
		slice = tab.Quo(sdk.NewIntFromBigInt(auction.OutflowTokenCurrentPrice.BigInt()))
	} else if auction.OutflowTokenCurrentAmount.Amount.Sub(slice).Mul(outFlowTokenCurrentPrice).LT(sdk.NewIntFromBigInt(auctionParams.Chost.BigInt())) {
		//(outflowtokenavailableamount-slice) in usd < chost in usd
		//see if user has balance to buy whole collateral
		user_balance_usd := k.bank.GetBalance(ctx, bidder, bid.Denom).Amount.Mul(outFlowTokenCurrentPrice)
		collateral_available_usd := auction.OutflowTokenCurrentAmount.Amount.Mul(outFlowTokenCurrentPrice)
		if user_balance_usd.LT(collateral_available_usd) {
			return auctiontypes.ErrorDutchinsufficientUserBalance
		}
		slice = auction.OutflowTokenCurrentAmount.Amount
	}

	inFlowTokenToCharge := slice.Mul(outFlowTokenCurrentPrice).Quo(inFlowTokenCurrentPrice)
	inFlowTokenCoin := sdk.NewCoin(auction.InflowTokenTargetAmount.Denom, inFlowTokenToCharge)
	outFlowTokenCoin := sdk.NewCoin(auction.OutflowTokenInitAmount.Denom, slice)
	err := k.SendCoinsFromAccountToModule(ctx, bidder, collectortypes.ModuleName, sdk.NewCoins(inFlowTokenCoin))
	if err != nil {
		return err
	}
	err = k.SendCoinsFromModuleToAccount(ctx, collectortypes.ModuleName, bidder, sdk.NewCoins(outFlowTokenCoin))
	if err != nil {
		//refund tokens to user as he is not getting collateral
		k.SendCoinsFromModuleToAccount(ctx, collectortypes.ModuleName, bidder, sdk.NewCoins(inFlowTokenCoin))
		return err
	}

	k.CreateNewDutchBid(ctx, auctionId, bidder, inFlowTokenCoin, outFlowTokenCoin)
	auction.OutflowTokenCurrentAmount = auction.OutflowTokenCurrentAmount.Sub(outFlowTokenCoin)
	auction.InflowTokenCurrentAmount = auction.InflowTokenCurrentAmount.Add(inFlowTokenCoin)
	if auction.InflowTokenTargetAmount.IsEqual(inFlowTokenCoin) {
		//TODO return collateral to vault as target dai reached
		//remove dutch auction
		k.CloseDutchAuction(ctx, auction)
		return nil
	}
	if auction.OutflowTokenCurrentAmount.Amount.IsZero() {
		k.CloseDutchAuction(ctx, auction)
		//TODO send target dai to the module specified
		return nil
	}
	k.SetDutchAuction(ctx, auction)
	return nil
}
