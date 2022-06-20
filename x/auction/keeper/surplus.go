package keeper

import (
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) SurplusActivator(ctx sdk.Context) error {

	auctionMapData, found := k.GetAllAuctionMappingForApp(ctx)
	if !found{
		return nil
	}
	for _, data := range auctionMapData{
		for _, indata := range data.AssetIdToAuctionLookup{
			if indata.IsSurplusAuction && !indata.IsAuctionActive{
				err:= k.CreateSurplusAuction(ctx, data.AppId, indata.AssetId )
				if err !=nil{
					return err
				}
			}else if indata.IsSurplusAuction && indata.IsAuctionActive{
				err := k.SurplusAuctionClose(ctx, data.AppId)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (k Keeper) CreateSurplusAuction(ctx sdk.Context, app_id, asset_id uint64) error {

		//check if auction status for an asset is false

		status, err := k.checkStatusOfNetFeesCollectedAndStartSurplusAuction(ctx, app_id, asset_id)
		if err != nil {
			return err
		}

		auctionLookupTable, _ := k.GetAuctionMappingForApp(ctx, app_id)

		for i, assetToAuction := range auctionLookupTable.AssetIdToAuctionLookup {
			if assetToAuction.AssetId == asset_id && status == auctiontypes.StartedSurplusAuction{

				auctionLookupTable.AssetIdToAuctionLookup[i].IsAuctionActive = true
			}
		}
		err1 := k.SetAuctionMappingForApp(ctx, auctionLookupTable)
		if err1 != nil {
			return err1
		}
	return nil
}

func (k Keeper) checkStatusOfNetFeesCollectedAndStartSurplusAuction(ctx sdk.Context, appId, assetId uint64) (status uint64, err error) {
	assetsCollectorDataUnderAppId, found := k.GetCollectorLookupTable(ctx, appId)
	if !found {
		return
	}
	//traverse this to access appId , collector asset id , surplus threshhold
	for _, collector := range assetsCollectorDataUnderAppId.AssetRateInfo {

		if collector.CollectorAssetId == assetId {
			//collectorLookupTable has surplusThreshhold for all assets

			NetFeeCollectedData, found := k.GetNetFeeCollectedData(ctx, appId)

			if !found {

				return auctiontypes.NoAuction, nil
			}
			//traverse this to access appId , collector asset id , netfees collected
			for _, AssetIdToFeeCollected := range NetFeeCollectedData.AssetIdToFeeCollected {

				if AssetIdToFeeCollected.AssetId == assetId {

					 if AssetIdToFeeCollected.NetFeesCollected.GTE(sdk.NewIntFromUint64(collector.SurplusThreshold+collector.LotSize)){
						// START SURPLUS AUCTION .  WITH COLLECTOR ASSET ID AS token given to user of lot size and secondary asset as received from user and burnt , bid factor
						//calculate inflow token amount

						assetBuyId := collector.SecondaryAssetId
						assetSellId := collector.CollectorAssetId

						//net = 900 surplusThreshhold = 500 , lotsize = 100
						amount := sdk.NewIntFromUint64(collector.LotSize)

						status, sellToken, buyToken := k.getSurplusBuyTokenAmount(ctx, appId, assetBuyId, assetSellId, amount)

						if status == auctiontypes.NoAuction {
							return auctiontypes.NoAuction, nil
						}
						//Transfer balance from collector module to auction module

						_, err := k.GetAmountFromCollector(ctx, appId, assetId, sellToken.Amount)
						if err != nil {

							return status, err
						}

						err = k.startSurplusAuction(ctx, sellToken, buyToken, collector.BidFactor, appId, assetId, assetBuyId, assetSellId)
						if err != nil {
							return status, err
						}
						return auctiontypes.StartedSurplusAuction, nil
					} else {

						return auctiontypes.NoAuction, nil
					}
				}
			}
		}
	}
	return auctiontypes.NoAuction, nil
}

func (k Keeper) getSurplusBuyTokenAmount(ctx sdk.Context, appId, AssetBuyId, AssetSellId uint64, lotSize sdk.Int) (status uint64, sellToken, buyToken sdk.Coin) {
	emptyCoin := sdk.NewCoin("empty", sdk.NewIntFromUint64(1))
	sellingAsset, found1 := k.GetAsset(ctx, AssetSellId)
	buyingAsset, found2 := k.GetAsset(ctx, AssetBuyId)
	if !found1 || !found2 {
		return auctiontypes.NoAuction, emptyCoin, emptyCoin
	}


	//outflow token will be of lot size
	sellToken = sdk.NewCoin(sellingAsset.Denom, lotSize)
	buyToken = sdk.NewCoin(buyingAsset.Denom, sdk.ZeroInt())
	return 5, sellToken, buyToken
}

func (k Keeper) startSurplusAuction(
	ctx sdk.Context,
	sellToken sdk.Coin,
	buyToken sdk.Coin,
	bidFactor sdk.Dec,
	appId, assetId uint64,
	assetInId, assetOutId uint64,
) error {

	auctionParams, found := k.GetAuctionParams(ctx, appId)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}
	auction := auctiontypes.SurplusAuction{
		SellToken:     sellToken,
		BuyToken:      buyToken,
		ActiveBiddingId:  0,
		Bidder:           nil,
		Bid:              buyToken,
		EndTime:          ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		BidFactor:        bidFactor,
		BiddingIds:       []*auctiontypes.BidOwnerMapping{},
		AuctionStatus:    auctiontypes.AuctionStartNoBids,
		AppId:            appId,
		AssetId:          assetId,
		AuctionMappingId: auctionParams.SurplusId,
		AssetInId:        assetInId,
		AssetOutId:       assetOutId,
	}
	auction.AuctionId = k.GetAuctionID(ctx) + 1
	k.SetAuctionID(ctx, auction.AuctionId)
	err := k.SetSurplusAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) SurplusAuctionClose(ctx sdk.Context, appId uint64) error {
	surplusAuctions := k.GetSurplusAuctions(ctx, appId)
	for _, surplusAuction := range surplusAuctions {
		if ctx.BlockTime().After(surplusAuction.EndTime) || ctx.BlockTime().After(surplusAuction.BidEndTime) {

			if surplusAuction.AuctionStatus == auctiontypes.AuctionStartNoBids {

				err := k.RestartSurplus(ctx, appId, surplusAuction)
				if err != nil {
					return err
				}
			} else {

				err := k.closeSurplusAuction(ctx, surplusAuction)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (k Keeper) RestartSurplus(
	ctx sdk.Context,
	appId uint64,
	surplusAuction auctiontypes.SurplusAuction,
) error {
	status, _, buyToken := k.getSurplusBuyTokenAmount(ctx, appId, surplusAuction.AssetInId, surplusAuction.AssetOutId, surplusAuction.BuyToken.Amount)
	if status == auctiontypes.NoAuction {
		return nil
	}
	auctionParams, found := k.GetAuctionParams(ctx, appId)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}
	surplusAuction.BuyToken = buyToken
	surplusAuction.Bid = buyToken
	surplusAuction.EndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
	err := k.SetSurplusAuction(ctx, surplusAuction)
	if err != nil {
		return err
	}
	return nil
}


func (k Keeper) closeSurplusAuction(
	ctx sdk.Context,
	surplusAuction auctiontypes.SurplusAuction,
) error {

	if surplusAuction.Bidder != nil {

		highestBidReceived := surplusAuction.Bid

		err := k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, surplusAuction.Bidder, sdk.NewCoins(surplusAuction.SellToken))
		if err != nil {

			return err
		}

		bidding, err := k.GetSurplusUserBidding(ctx, surplusAuction.Bidder.String(), surplusAuction.AppId, surplusAuction.ActiveBiddingId)
		if err != nil {

			return err
		}
		bidding.BiddingStatus = auctiontypes.SuccessBiddingStatus
		err = k.SetSurplusUserBidding(ctx, bidding)
		if err != nil {
			return err
		}

		if auctiontypes.TestFlag == 1 {

			err = k.BurnCoins(ctx, auctiontypes.ModuleName, highestBidReceived)
			if err != nil {
				return auctiontypes.ErrorInvalidBurn
			}
		} else {

			//burn tokens by sending bid tokens from auction to tokenmint module and then call burn function
			err = k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, tokenminttypes.ModuleName, sdk.NewCoins(highestBidReceived))
			if err != nil {
				return err
			}
			err = k.BurnTokensForApp(ctx, surplusAuction.AppId, surplusAuction.AssetInId, highestBidReceived.Amount)
			if err != nil {

				return err
			}

		}

		for _, biddingId := range surplusAuction.BiddingIds {
			bidding, err := k.GetSurplusUserBidding(ctx, biddingId.BidOwner, surplusAuction.AppId, biddingId.BidId)
			if err != nil {
				continue
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			err = k.SetSurplusUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.DeleteSurplusUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.SetHistorySurplusUserBidding(ctx, bidding)
			if err != nil {
				return err
			}

		}
	} else {
		err1 := k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(surplusAuction.SellToken))
		if err1 != nil {
			return err1
		}
		err2 := k.SetNetFeeCollectedData(ctx, surplusAuction.AppId, surplusAuction.AssetOutId, surplusAuction.SellToken.Amount)
		if err2 != nil {
			return auctiontypes.ErrorUnableToSetNetfees
		}
	}
	err := k.makeFalseForFlags(ctx, surplusAuction.AppId, surplusAuction.AssetId)
	if err != nil {
		return auctiontypes.ErrorUnableToMakeFlagsFalse
	}
	err = k.DeleteSurplusAuction(ctx, surplusAuction)
	if err != nil {
		return err
	}
	surplusAuction.AuctionStatus = auctiontypes.AuctionEnded
	err = k.SetHistorySurplusAuction(ctx, surplusAuction)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) PlaceSurplusAuctionBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) error {
	auction, err := k.GetSurplusAuction(ctx, appId, auctionMappingId, auctionId)
	auctionParam, _ := k.GetAuctionParams(ctx,appId)

	if err != nil {
		return auctiontypes.ErrorInvalidSurplusAuctionId
	}
	if bid.Denom != auction.BuyToken.Denom {
		return auctiontypes.ErrorInvalidBiddingDenom
	}

	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		change := auction.BidFactor.MulInt(auction.Bid.Amount).Ceil().TruncateInt()
		minBidAmount := auction.Bid.Amount.Add(change)
		if bid.Amount.LT(minBidAmount) {
			return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bid should be greater than or equal to %d ", minBidAmount)
		}
	} else {
		if bid.Amount.LT(auction.Bid.Amount) {
			return auctiontypes.ErrorLowBidAmount
		}
	}
	err = k.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(bid))
	if err != nil {
		return err
	}
	biddingId, err := k.CreateNewSurplusBid(ctx, appId, auctionMappingId, auctionId, bidder, bid)
	if err != nil {
		return err
	}
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		// auction.Bidder as previous bidder . refund previous bidder
		err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, auction.Bidder, sdk.NewCoins(auction.Bid))
		if err != nil {
			return err
		}
		bidding, err := k.GetSurplusUserBidding(ctx, auction.Bidder.String(), auction.AppId, auction.ActiveBiddingId)
		if err != nil {
			return err
		}
		bidding.BiddingStatus = auctiontypes.RejectedBiddingStatus
		err = k.SetSurplusUserBidding(ctx, bidding)
		if err != nil {
			return err
		}
	} else {
		auction.AuctionStatus = auctiontypes.AuctionGoingOn
	}
	auction.ActiveBiddingId = biddingId
	var bidIdOwner = &auctiontypes.BidOwnerMapping{BidId: biddingId, BidOwner: bidder.String()}
	auction.BiddingIds = append(auction.BiddingIds, bidIdOwner)
	auction.Bidder = bidder
	auction.Bid = bid
	auction.BidEndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParam.BidDurationSeconds))
	if auction.BidEndTime.After((auction.EndTime)){
		auction.BidEndTime = auction.EndTime
	}
	err = k.SetSurplusAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) CreateNewSurplusBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) (biddingId uint64, err error) {
	auction, err := k.GetSurplusAuction(ctx, appId, auctionMappingId, auctionId)
	if err != nil {
		return biddingId, err
	}
	bidding := auctiontypes.SurplusBiddings{
		BiddingId:           k.GetUserBiddingID(ctx) + 1,
		AuctionId:           auctionId,
		AuctionStatus:       auctiontypes.ActiveAuctionStatus,
		AuctionedCollateral: auction.SellToken,
		Bidder:              bidder.String(),
		Bid:                 bid,
		BiddingTimestamp:    ctx.BlockTime(),
		BiddingStatus:       auctiontypes.PlacedBiddingStatus,
		AppId:               appId,
		AuctionMappingId:    auctionMappingId,
	}
	k.SetUserBiddingID(ctx, bidding.BiddingId)
	err = k.SetSurplusUserBidding(ctx, bidding)
	if err != nil {
		return biddingId, err
	}
	return bidding.BiddingId, nil
}