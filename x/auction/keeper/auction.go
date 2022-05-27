package keeper

import (
	"fmt"
	"time"

	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) checkStatusOfNetFeesCollectedAndStartAuction(ctx sdk.Context, appId, assetId uint64, assetToAuction collectortypes.AssetIdToAuctionLookupTable) (status uint64) {
	assetsCollectorDataUnderAppId, found := k.GetCollectorLookupTable(ctx, appId)
	if !found {
		return
	}
	//traverse this to access appId , collector asset id , surplus threshhold , debt threshhold
	for _, collector := range assetsCollectorDataUnderAppId.AssetrateInfo {
		if collector.CollectorAssetId == assetId {
			//collectorLookupTable has surplusThreshhold for all assets

			NetFeeCollectedData, found := k.GetNetFeeCollectedData(ctx, appId)
			if !found {
				return auctiontypes.NoAuction
			}
			//traverse this to access appId , collector asset id , netfees collected
			for _, AssetIdToFeeCollected := range NetFeeCollectedData.AssetIdToFeeCollected {
				if AssetIdToFeeCollected.AssetId == assetId {
					// if netfees <= debt threshhold -lotsize the start debt auction with lot size and debt auction is allowed true
					if AssetIdToFeeCollected.NetFeesCollected.LTE(sdk.NewIntFromUint64(collector.DebtThreshold-collector.LotSize)) && assetToAuction.IsDebtAuction {
						// START DEBT AUCTION .  LOTSIZE AS MINTED FOR SECONDARY ASSET and ACCEPT Collector assetid from user
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
						err := k.StartDebtAuction(ctx, outflowToken, inflowToken, appId, assetId)
						if err != nil {
							break
						}
						return auctiontypes.StartedDebtAuction
						// if netfees >= surplus threshhold+lotsize the start surplus auction with lot size and surplus auction is allowed true
					} else if AssetIdToFeeCollected.NetFeesCollected.GTE(sdk.NewIntFromUint64(collector.SurplusThreshold+collector.LotSize)) && assetToAuction.IsDebtAuction {
						// START SURPLUS AUCTION .  WITH COLLECTOR ASSET ID AS token given to user of lot size and secondary asset as received from user and burnt , bid factor
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
						_, err := k.collector.GetAmountFromCollector(ctx, appId, assetId, outflowToken.Amount)
						if err != nil {
							break
						}
						//k.bank.SendCoinsFromModuleToModule(ctx, collectortypes.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(outflowToken))
						err = k.StartSurplusAuction(ctx, outflowToken, inflowToken, *collector.BidFactor, appId, assetId)
						if err != nil {
							break
						}
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
					status := k.checkStatusOfNetFeesCollectedAndStartAuction(ctx, appId.Id, assetToAuction.AssetId, *assetToAuction)
					if status == auctiontypes.StartedDebtAuction {
						assetToAuction.IsAuctionActive = true
					} else if status == auctiontypes.StartedSurplusAuction {
						assetToAuction.IsAuctionActive = true
					} else {
						continue
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

func (k Keeper) CloseAndRestartSurplusAuctions(ctx sdk.Context) error {
	appIds, found := k.GetApps(ctx)
	if !found {
		return assettypes.AppIdsDoesntExist
	}
	for _, appId := range appIds {
		k.CloseSurplusAuctions(ctx, appId.Id)
		k.CloseDebtAuctions(ctx, appId.Id)
		k.RestartDutchAuctions(ctx, appId.Id)
		k.CreateSurplusAndDebtAuctions(ctx)
		k.CreateNewDutchAuctions(ctx)
	}
	return nil
}
func (k Keeper) CreateNewDutchAuctions(ctx sdk.Context) {
	lockedVaults := k.GetLockedVaults(ctx)
	for _, lockedVault := range lockedVaults {
		pair, found := k.GetPair(ctx, lockedVault.ExtendedPairId)
		if !found {
			continue
		}
		assetIn, found := k.GetAsset(ctx, pair.AssetIn)
		if !found {
			continue
		}

		assetOut, found := k.GetAsset(ctx, pair.AssetOut)
		if !found {
			continue
		}

		outflowToken := sdk.NewCoin(assetIn.Denom, lockedVault.CollateralToBeAuctioned.TruncateInt())
		inflowToken := sdk.NewCoin(assetOut.Denom, sdk.ZeroInt())

		if !lockedVault.IsAuctionInProgress {
			k.StartDutchAuction(ctx, outflowToken, inflowToken, lockedVault.AppMappingId, assetOut.Id, assetIn.Id, lockedVault.LockedVaultId, lockedVault.Owner)
		}
	}
}

func (k Keeper) CloseSurplusAuctions(ctx sdk.Context, appId uint64) {
	surplus_auctions := k.GetSurplusAuctions(ctx, appId)
	for _, surplus_auction := range surplus_auctions {
		if ctx.BlockTime().After(surplus_auction.EndTime) {
			err := k.CloseSurplusAuction(ctx, surplus_auction)
			if err == nil {
				continue
			}
		}
	}
}

// get all app ids and call RestartDutchAuctions with app id
func (k Keeper) CloseDebtAuctions(ctx sdk.Context, appId uint64) {
	debtAuctions := k.GetDebtAuctions(ctx, appId)
	for _, debtAuction := range debtAuctions {
		if ctx.BlockTime().After(debtAuction.EndTime) {
			err := k.CloseDebtAuction(ctx, debtAuction)
			if err == nil {
				continue
			}
		}
	}
}

//get all app ids and call RestartDutchAuctions with app id
func (k Keeper) RestartDutchAuctions(ctx sdk.Context, appId uint64) {
	dutchAuctions := k.GetDutchAuctions(ctx, appId)
	auctionParams := k.GetParams(ctx)
	// SET current price of inflow token and outflow token
	for _, dutchAuction := range dutchAuctions {
		inFlowTokenCurrentPrice, found := k.market.GetPriceForAsset(ctx, dutchAuction.AssetInId)
		if !found {
			fmt.Println("not able fetch price from oracle")
			return
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
			outFlowTokenCurrentPrice, found := k.market.GetPriceForAsset(ctx, dutchAuction.AssetOutId)
			if !found {
				fmt.Println("not able fetch price from oracle")
				return
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
		OutflowToken:     outflowToken,
		InflowToken:      inflowToken,
		ActiveBiddingId:  0,
		Bidder:           nil,
		Bid:              inflowToken,
		EndTime:          ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		BidFactor:        bidFactor,
		BiddingIds:       []*auctiontypes.BidOwnerMapping{},
		AuctionStatus:    auctiontypes.AuctionStartNoBids,
		AppId:            appId,
		AssetId:          assetId,
		AuctionMappingId: auctionParams.SurplusId,
	}
	auction.AuctionId = k.GetAuctionID(ctx) + 1
	k.SetAuctionID(ctx, auction.AuctionId)
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
		BiddingIds:          []*auctiontypes.BidOwnerMapping{},
		AuctionMappingId:    auctionParams.DebtId,
	}
	fmt.Println(auctionParams.DebtId)
	auction.AuctionId = k.GetAuctionID(ctx) + 1
	k.SetAuctionID(ctx, auction.AuctionId)
	k.SetDebtAuction(ctx, auction)
	return nil
}

func (k Keeper) StartDutchAuction(
	ctx sdk.Context,
	outFlowToken sdk.Coin,
	inFlowToken sdk.Coin,
	appId uint64,
	assetInId, assetOutId uint64,
	lockedVaultId uint64,
	lockedVaultOwner string,
) error {

	auctionParams := k.GetParams(ctx)
	//need to get real price instead of hard coding
	//calculate target amount of cmst to collect
	inFlowTokenPrice, found1 := k.market.GetPriceForAsset(ctx, assetInId)
	outFlowTokenPrice, found2 := k.market.GetPriceForAsset(ctx, assetOutId)

	if !(found1 && found2) {
		return auctiontypes.ErrorInvalidBidId
	}
	//outFlowTokenPrice := sdk.MustNewDecFromStr("100")
	//inFlowTokenPrice := sdk.MustNewDecFromStr("1")

	outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(sdk.NewIntFromUint64(outFlowTokenPrice), auctionParams.Buffer)
	outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
	vaultOwner, err := sdk.AccAddressFromBech32(lockedVaultOwner)
	if err != nil {
		return err
	}
	timeNow := ctx.BlockTime()
	auction := auctiontypes.DutchAuction{
		OutflowTokenInitAmount:    outFlowToken,
		OutflowTokenCurrentAmount: outFlowToken,
		InflowTokenTargetAmount:   inFlowToken,
		InflowTokenCurrentAmount:  inFlowToken,
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
	}

	auction.AuctionId = k.GetAuctionID(ctx) + 1
	k.SetAuctionID(ctx, auction.AuctionId)
	k.SetDutchAuction(ctx, auction)
	k.SetFlagIsAuctionInProgress(ctx, lockedVaultId, true)
	return nil
}

func (k Keeper) CloseSurplusAuction(
	ctx sdk.Context,
	surplusAuction auctiontypes.SurplusAuction,
) error {

	if surplusAuction.Bidder != nil {

		highestBidReceived := surplusAuction.Bid

		err := k.bank.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, surplusAuction.Bidder, sdk.NewCoins(surplusAuction.OutflowToken))
		if err != nil {
			return err
		}
		bidding, found := k.GetSurplusUserBidding(ctx, surplusAuction.Bidder.String(), surplusAuction.AppId, surplusAuction.ActiveBiddingId)
		if !found {
			return auctiontypes.ErrorInvalidBidId
		}
		bidding.BiddingStatus = auctiontypes.SuccessBiddingStatus
		k.SetSurplusUserBidding(ctx, bidding)
		//burn tokens by sending bid tokens from auction to tokenmint module and then call burn function
		err = k.bank.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, tokenminttypes.ModuleName, sdk.NewCoins(highestBidReceived))
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = k.tokenmint.BurnTokensForApp(ctx, surplusAuction.AppId, surplusAuction.AssetId, highestBidReceived.Amount)
		if err != nil {
			fmt.Println(err)
			return err
		}
		//following 4 lines used for testing purpose
		//err = k.bank.BurnCoins(ctx, auctiontypes.ModuleName, sdk.NewCoins(highestBidReceived))
		//if err != nil {
		//	return auctiontypes.ErrorInvalidBurn
		//}

		for _, biddingId := range surplusAuction.BiddingIds {
			bidding, found := k.GetSurplusUserBidding(ctx, biddingId.BidOwner, surplusAuction.AppId, biddingId.BidId)
			if !found {
				continue
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			k.SetSurplusUserBidding(ctx, bidding)
			k.DeleteSurplusUserBidding(ctx, bidding)
			k.SetHistorySurplusUserBidding(ctx, bidding)
		}
	}
	k.makeFalseForFlags(ctx, surplusAuction.AppId, surplusAuction.AssetId)
	k.DeleteSurplusAuction(ctx, surplusAuction)
	k.SetHistorySurplusAuction(ctx, surplusAuction)
	//store auctions and user bidding in history after they are deleted
	return nil
}

func (k Keeper) CloseDebtAuction(
	ctx sdk.Context,
	debtAuction auctiontypes.DebtAuction,
) error {

	if debtAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		// ask token mint to mint new tokens for bidder address
		err := k.tokenmint.MintNewTokensForApp(ctx, debtAuction.AppId, debtAuction.AssetId, debtAuction.Bidder.String(), debtAuction.CurrentBidAmount.Amount)
		if err != nil {
			return err
		}
		//following 6 lines used for testing purpose
		//err := k.bank.MintCoins(ctx, auctiontypes.ModuleName, sdk.NewCoins(debtAuction.CurrentBidAmount))
		//err = k.bank.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, debtAuction.Bidder, sdk.NewCoins(debtAuction.CurrentBidAmount))
		//if err != nil {
		//	return err
		//}

		bidding, found := k.GetDebtUserBidding(ctx, debtAuction.Bidder.String(), debtAuction.AppId, debtAuction.ActiveBiddingId)
		if !found {
			return auctiontypes.ErrorInvalidBidId
		}
		bidding.BiddingStatus = auctiontypes.SuccessBiddingStatus
		k.SetDebtUserBidding(ctx, bidding)
		for _, biddingId := range debtAuction.BiddingIds {
			bidding, found := k.GetDebtUserBidding(ctx, biddingId.BidOwner, debtAuction.AppId, biddingId.BidId)
			if !found {
				continue
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			k.SetDebtUserBidding(ctx, bidding)
			k.DeleteDebtUserBidding(ctx, bidding)
			k.SetHistoryDebtUserBidding(ctx, bidding)
		}
	}
	//send to collector module the amount collected in debt auction
	err := k.bank.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(debtAuction.ExpectedUserToken))
	if err != nil {
		return err
	}
	err = k.SetNetFeeCollectedData(ctx, debtAuction.AuctionId, debtAuction.AssetId, debtAuction.ExpectedUserToken.Amount)
	if err != nil {
		return err
	}
	k.makeFalseForFlags(ctx, debtAuction.AppId, debtAuction.AssetId)
	k.DeleteDebtAuction(ctx, debtAuction)
	k.SetHistoryDebtAuction(ctx, debtAuction)
	return nil
}

func (k Keeper) CloseDutchAuction(
	ctx sdk.Context,
	dutchAuction auctiontypes.DutchAuction,
) error {

	if dutchAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		for _, biddingId := range dutchAuction.BiddingIds {
			bidding, found := k.GetDutchUserBidding(ctx, biddingId.BidOwner, dutchAuction.AppId, biddingId.BidId)
			if !found {
				continue
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			k.SetDutchUserBidding(ctx, bidding)
			k.DeleteDutchUserBidding(ctx, bidding)
			k.SetHistoryDutchUserBidding(ctx, bidding)
		}
	}
	k.SetFlagIsAuctionComplete(ctx, dutchAuction.LockedVaultId, true)
	k.SetFlagIsAuctionInProgress(ctx, dutchAuction.LockedVaultId, false)
	k.DeleteDutchAuction(ctx, dutchAuction)
	k.SetHistoryDutchAuction(ctx, dutchAuction)
	return nil
}

func (k Keeper) CreateNewSurplusBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) (biddingId uint64, err error) {
	auction, found := k.GetSurplusAuction(ctx, appId, auctionMappingId, auctionId)
	if !found {
		return 0, auctiontypes.ErrorInvalidSurplusAuctionId
	}
	bidding := auctiontypes.SurplusBiddings{
		BiddingId:           k.GetUserBiddingID(ctx) + 1,
		AuctionId:           auctionId,
		AuctionStatus:       auctiontypes.ActiveAuctionStatus,
		AuctionedCollateral: auction.OutflowToken,
		Bidder:              bidder.String(),
		Bid:                 bid,
		BiddingTimestamp:    ctx.BlockTime(),
		BiddingStatus:       auctiontypes.PlacedBiddingStatus,
		AppId:               appId,
		AuctionMappingId:    auctionMappingId,
	}
	k.SetUserBiddingID(ctx, bidding.BiddingId)
	k.SetSurplusUserBidding(ctx, bidding)
	return bidding.BiddingId, nil
}

func (k Keeper) CreateNewDebtBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, expectedUserToken sdk.Coin) (biddingId uint64, err error) {
	bidding := auctiontypes.DebtBiddings{
		BiddingId:        k.GetUserBiddingID(ctx) + 1,
		AuctionId:        auctionId,
		AuctionStatus:    auctiontypes.ActiveAuctionStatus,
		Bidder:           bidder.String(),
		Bid:              bid,
		BiddingTimestamp: ctx.BlockTime(),
		BiddingStatus:    auctiontypes.PlacedBiddingStatus,
		AppId:            appId,
		AuctionMappingId: auctionMappingId,
		OutflowTokens:    expectedUserToken,
	}
	fmt.Println(auctionMappingId)
	k.SetUserBiddingID(ctx, bidding.BiddingId)

	k.SetDebtUserBidding(ctx, bidding)

	return bidding.BiddingId, nil
}

func (k Keeper) CreateNewDutchBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, outFlowTokenCoin sdk.Coin, inFlowTokenCoin sdk.Coin) (biddingId uint64) {
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
	k.SetDutchUserBidding(ctx, bidding)
	return bidding.BiddingId
}

func (k Keeper) PlaceSurplusBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) error {
	auction, found := k.GetSurplusAuction(ctx, appId, auctionMappingId, auctionId)
	if !found {
		return auctiontypes.ErrorInvalidSurplusAuctionId
	}
	if bid.Denom != auction.InflowToken.Denom {
		return auctiontypes.ErrorInvalidBiddingDenom
	}
	//Test this multiplication check if new bid greater than previous bid by bid factor
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		minBidCoin := sdk.NewCoin(bid.Denom, auction.BidFactor.MulInt(auction.Bid.Amount).Ceil().TruncateInt())
		if bid.Amount.LT(minBidCoin.Amount) {
			return auctiontypes.ErrorLowBidAmount
		}
	} else {
		if bid.Amount.LT(auction.Bid.Amount) {
			return auctiontypes.ErrorLowBidAmount
		}
	}
	err := k.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(bid))
	if err != nil {
		return err
	}
	fmt.Println(bidder)
	biddingId, err := k.CreateNewSurplusBid(ctx, auctionId, auctionMappingId, auctionId, bidder, bid)
	if err != nil {
		return err
	}
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		// auction.Bidder as previous bidder
		err = k.bank.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, auction.Bidder, sdk.NewCoins(auction.Bid))
		if err != nil {
			return err
		}
		bidding, _ := k.GetSurplusUserBidding(ctx, auction.Bidder.String(), auction.AppId, auction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.RejectedBiddingStatus
		k.SetSurplusUserBidding(ctx, bidding)
	} else {
		auction.AuctionStatus = auctiontypes.AuctionGoingOn
	}
	auction.ActiveBiddingId = biddingId
	var bidIdOwner = &auctiontypes.BidOwnerMapping{biddingId, bidder.String()}
	auction.BiddingIds = append(auction.BiddingIds, bidIdOwner)
	auction.Bidder = bidder
	auction.Bid = bid
	k.SetSurplusAuction(ctx, auction)
	return nil
}

func (k Keeper) PlaceDebtBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, expectedUserToken sdk.Coin) error {
	fmt.Println(auctionMappingId)
	auction, found := k.GetDebtAuction(ctx, appId, auctionMappingId, auctionId)

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

	err := k.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(expectedUserToken))
	if err != nil {
		return err
	}
	fmt.Println(auctionMappingId)
	biddingId, err := k.CreateNewDebtBid(ctx, appId, auctionMappingId, auctionId, bidder, bid, expectedUserToken)
	if err != nil {
		return err
	}
	//If auction gets bid from second time onwards
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		err = k.bank.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, auction.Bidder, sdk.NewCoins(auction.ExpectedUserToken))
		if err != nil {
			return err
		}
		bidding, _ := k.GetDebtUserBidding(ctx, auction.Bidder.String(), auction.AppId, auction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.RejectedBiddingStatus

		k.SetDebtUserBidding(ctx, bidding)
	} else {
		auction.AuctionStatus = auctiontypes.AuctionGoingOn
	}
	auction.ActiveBiddingId = biddingId
	var bidIdOwner = &auctiontypes.BidOwnerMapping{biddingId, bidder.String()}
	auction.BiddingIds = append(auction.BiddingIds, bidIdOwner)
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

func (k Keeper) PlaceDutchBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, max sdk.Dec) error {
	auction, found := k.GetDutchAuction(ctx, appId, auctionMappingId, auctionId)
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
	//check if bid is greater than target cmst
	if owe.GT(tab) {
		slice = tab.Quo(sdk.NewIntFromBigInt(auction.OutflowTokenCurrentPrice.BigInt()))
	} else if auction.OutflowTokenCurrentAmount.Amount.Sub(slice).Mul(outFlowTokenCurrentPrice).LT(sdk.NewIntFromBigInt(auctionParams.Chost.BigInt())) {
		//(outflowtokenavailableamount-slice) in usd < chost in usd
		//see if user has balance to buy whole collateral
		userBalanceUsd := k.bank.GetBalance(ctx, bidder, bid.Denom).Amount.Mul(outFlowTokenCurrentPrice)
		collateralAvailableUsd := auction.OutflowTokenCurrentAmount.Amount.Mul(outFlowTokenCurrentPrice)
		if userBalanceUsd.LT(collateralAvailableUsd) {
			return auctiontypes.ErrorDutchinsufficientUserBalance
		}
		slice = auction.OutflowTokenCurrentAmount.Amount
	}

	inFlowTokenToCharge := slice.Mul(outFlowTokenCurrentPrice).Quo(inFlowTokenCurrentPrice)
	inFlowTokenCoin := sdk.NewCoin(auction.InflowTokenTargetAmount.Denom, inFlowTokenToCharge)
	outFlowTokenCoin := sdk.NewCoin(auction.OutflowTokenInitAmount.Denom, slice)
	err := k.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(inFlowTokenCoin))
	if err != nil {
		return err
	}
	err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidder, sdk.NewCoins(outFlowTokenCoin))
	if err != nil {
		//refund tokens to user as he is not getting collateral
		err := k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidder, sdk.NewCoins(inFlowTokenCoin))
		return err
	}
	biddingId := k.CreateNewDutchBid(ctx, appId, auctionMappingId, auctionId, bidder, inFlowTokenCoin, outFlowTokenCoin)
	var bidIdOwner = &auctiontypes.BidOwnerMapping{biddingId, bidder.String()}
	auction.BiddingIds = append(auction.BiddingIds, bidIdOwner)
	auction.OutflowTokenCurrentAmount = auction.OutflowTokenCurrentAmount.Sub(outFlowTokenCoin)
	auction.InflowTokenCurrentAmount = auction.InflowTokenCurrentAmount.Add(inFlowTokenCoin)
	//collateral not over but target cmst reached
	if auction.InflowTokenTargetAmount.IsEqual(inFlowTokenCoin) {

		//return collateral to vault owner as target cmst reached and also
		total := auction.OutflowTokenCurrentAmount
		err := k.bank.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, auction.VaultOwner, sdk.NewCoins(total))
		if err != nil {
			return err
		}
		//burn and send target CMST to collector
		params := k.GetParams(ctx)
		inFlowAmount := auction.InflowTokenCurrentAmount
		burnAmount := k.getBurnAmountAmount(inFlowAmount.Amount, params.LiquidationPenaltyPercent)
		err = k.tokenmint.BurnTokensForApp(ctx, auction.AppId, auction.AssetInId, burnAmount)
		if err != nil {
			return err
		}
		penaltyAmount := inFlowAmount.Amount.Sub(burnAmount)
		err = k.bank.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(inFlowAmount.Denom, penaltyAmount)))
		if err != nil {
			return err
		}
		// call increase function in collector
		k.SetNetFeeCollectedData(ctx, auction.AppId, auction.AssetInId, penaltyAmount)
		//remove dutch auction
		k.CloseDutchAuction(ctx, auction)
		return nil
	}
	//entire collateral sold out
	if auction.OutflowTokenCurrentAmount.Amount.IsZero() {
		//burn and send target CMST to collector
		params := k.GetParams(ctx)
		inFlowAmount := auction.InflowTokenCurrentAmount
		burnAmount := k.getBurnAmountAmount(inFlowAmount.Amount, params.LiquidationPenaltyPercent)
		err = k.tokenmint.BurnTokensForApp(ctx, auction.AppId, auction.AssetInId, burnAmount)
		if err != nil {
			return err
		}
		penaltyAmount := inFlowAmount.Amount.Sub(burnAmount)
		err = k.bank.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(inFlowAmount.Denom, penaltyAmount)))
		if err != nil {
			return err
		}

		//call increase function in collector
		k.SetNetFeeCollectedData(ctx, auction.AppId, auction.AssetInId, penaltyAmount)

		//remove dutch auction
		k.CloseDutchAuction(ctx, auction)
		return nil
	}
	k.SetDutchAuction(ctx, auction)
	return nil
}
