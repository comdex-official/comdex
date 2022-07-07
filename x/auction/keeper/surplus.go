package keeper

import (
	"time"

	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) SurplusActivator(ctx sdk.Context) error {
	auctionMapData, found := k.GetAllAuctionMappingForApp(ctx)
	if !found {
		return nil
	}
	for _, data := range auctionMapData {
		for _, inData := range data.AssetIdToAuctionLookup {
			if inData.IsSurplusAuction && !inData.IsAuctionActive {
				err := k.CreateSurplusAuction(ctx, data.AppId, inData.AssetId)
				if err != nil {
					return err
				}
			} else if inData.IsSurplusAuction && inData.IsAuctionActive {
				err := k.SurplusAuctionClose(ctx, data.AppId)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (k Keeper) CreateSurplusAuction(ctx sdk.Context, appID, assetID uint64) error { //check if auction status for an asset is false
	status, err := k.checkStatusOfNetFeesCollectedAndStartSurplusAuction(ctx, appID, assetID)
	if err != nil {
		return err
	}

	auctionLookupTable, _ := k.GetAuctionMappingForApp(ctx, appID)

	for i, assetToAuction := range auctionLookupTable.AssetIdToAuctionLookup {
		if assetToAuction.AssetId == assetID && status == auctiontypes.StartedSurplusAuction {
			auctionLookupTable.AssetIdToAuctionLookup[i].IsAuctionActive = true
		}
	}
	err1 := k.SetAuctionMappingForApp(ctx, auctionLookupTable)
	if err1 != nil {
		return err1
	}
	return nil
}

func (k Keeper) checkStatusOfNetFeesCollectedAndStartSurplusAuction(ctx sdk.Context, appID, assetID uint64) (status uint64, err error) {
	assetsCollectorDataUnderAppID, found := k.GetCollectorLookupTable(ctx, appID)
	if !found {
		return
	}
	//traverse this to access appId , collector asset id , surplus threshold
	for _, collector := range assetsCollectorDataUnderAppID.AssetRateInfo {
		if collector.CollectorAssetId == assetID {
			//collectorLookupTable has surplusThreshold for all assets

			NetFeeCollectedData, found := k.GetNetFeeCollectedData(ctx, appID)

			if !found {
				return auctiontypes.NoAuction, nil
			}
			//traverse this to access appId , collector asset id , netfees collected
			for _, AssetIDToFeeCollected := range NetFeeCollectedData.AssetIdToFeeCollected {
				if AssetIDToFeeCollected.AssetId == assetID {
					if AssetIDToFeeCollected.NetFeesCollected.GTE(sdk.NewIntFromUint64(collector.SurplusThreshold + collector.LotSize)) {
						// START SURPLUS AUCTION .  WITH COLLECTOR ASSET ID AS token given to user of lot size and secondary asset as received from user and burnt , bid factor
						//calculate inflow token amount

						assetBuyID := collector.SecondaryAssetId
						assetSellID := collector.CollectorAssetId

						//net = 900 surplusThreshhold = 500 , lotsize = 100
						amount := sdk.NewIntFromUint64(collector.LotSize)

						status, sellToken, buyToken := k.getSurplusBuyTokenAmount(ctx, assetBuyID, assetSellID, amount)

						if status == auctiontypes.NoAuction {
							return auctiontypes.NoAuction, nil
						}
						//Transfer balance from collector module to auction module

						_, err := k.GetAmountFromCollector(ctx, appID, assetID, sellToken.Amount)
						if err != nil {
							return status, err
						}

						err = k.StartSurplusAuction(ctx, sellToken, buyToken, collector.BidFactor, appID, assetID, assetBuyID, assetSellID)
						if err != nil {
							return status, err
						}
						return auctiontypes.StartedSurplusAuction, nil
					}
				}
			}
		}
	}
	return auctiontypes.NoAuction, nil
}

func (k Keeper) getSurplusBuyTokenAmount(ctx sdk.Context, AssetBuyID, AssetSellID uint64, lotSize sdk.Int) (status uint64, sellToken, buyToken sdk.Coin) {
	emptyCoin := sdk.NewCoin("empty", sdk.NewIntFromUint64(1))
	sellingAsset, found1 := k.GetAsset(ctx, AssetSellID)
	buyingAsset, found2 := k.GetAsset(ctx, AssetBuyID)
	if !found1 || !found2 {
		return auctiontypes.NoAuction, emptyCoin, emptyCoin
	}

	//outflow token will be of lot size
	sellToken = sdk.NewCoin(sellingAsset.Denom, lotSize)
	buyToken = sdk.NewCoin(buyingAsset.Denom, sdk.ZeroInt())
	return 5, sellToken, buyToken
}

func (k Keeper) StartSurplusAuction(
	ctx sdk.Context,
	sellToken sdk.Coin,
	buyToken sdk.Coin,
	bidFactor sdk.Dec,
	appID, assetID uint64,
	assetInID, assetOutID uint64,
) error {
	auctionParams, found := k.GetAuctionParams(ctx, appID)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}
	auction := auctiontypes.SurplusAuction{
		SellToken:        sellToken,
		BuyToken:         buyToken,
		ActiveBiddingId:  0,
		Bidder:           nil,
		Bid:              buyToken,
		EndTime:          ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		BidEndTime:       ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		BidFactor:        bidFactor,
		BiddingIds:       []*auctiontypes.BidOwnerMapping{},
		AuctionStatus:    auctiontypes.AuctionStartNoBids,
		AppId:            appID,
		AssetId:          assetID,
		AuctionMappingId: auctionParams.SurplusId,
		AssetInId:        assetInID,
		AssetOutId:       assetOutID,
	}
	auction.AuctionId = k.GetAuctionID(ctx) + 1
	k.SetAuctionID(ctx, auction.AuctionId)
	err := k.SetSurplusAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) SurplusAuctionClose(ctx sdk.Context, appID uint64) error {
	surplusAuctions := k.GetSurplusAuctions(ctx, appID)
	for _, surplusAuction := range surplusAuctions {
		if ctx.BlockTime().After(surplusAuction.EndTime) || ctx.BlockTime().After(surplusAuction.BidEndTime) {
			if surplusAuction.AuctionStatus == auctiontypes.AuctionStartNoBids {
				err := k.RestartSurplus(ctx, appID, surplusAuction)
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
	appID uint64,
	surplusAuction auctiontypes.SurplusAuction,
) error {
	status, _, buyToken := k.getSurplusBuyTokenAmount(ctx, surplusAuction.AssetInId, surplusAuction.AssetOutId, surplusAuction.BuyToken.Amount)
	if status == auctiontypes.NoAuction {
		return nil
	}
	auctionParams, found := k.GetAuctionParams(ctx, appID)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}
	surplusAuction.BuyToken = buyToken
	surplusAuction.Bid = buyToken
	surplusAuction.EndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
	surplusAuction.BidEndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
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

		//burn tokens by sending bid tokens from auction to tokenMint module and then call burn function
		err = k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, tokenminttypes.ModuleName, sdk.NewCoins(highestBidReceived))
		if err != nil {
			return err
		}
		err = k.BurnTokensForApp(ctx, surplusAuction.AppId, surplusAuction.AssetInId, highestBidReceived.Amount)
		if err != nil {
			return err
		}

		for _, biddingID := range surplusAuction.BiddingIds {
			bidding, err := k.GetSurplusUserBidding(ctx, biddingID.BidOwner, surplusAuction.AppId, biddingID.BidId)
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
			return auctiontypes.ErrorUnableToSetNetFees
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

func (k Keeper) PlaceSurplusAuctionBid(ctx sdk.Context, appID, auctionMappingID, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin) error {
	auction, err := k.GetSurplusAuction(ctx, appID, auctionMappingID, auctionID)
	auctionParam, _ := k.GetAuctionParams(ctx, appID)

	if err != nil {
		return auctiontypes.ErrorInvalidSurplusAuctionID
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
		if bid.Amount.LTE(auction.Bid.Amount) {
			return auctiontypes.ErrorLowBidAmount
		}
	}
	err = k.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(bid))
	if err != nil {
		return err
	}
	biddingID, err := k.CreateNewSurplusBid(ctx, appID, auctionMappingID, auctionID, bidder.String(), bid)
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
	auction.ActiveBiddingId = biddingID
	var bidIDOwner = &auctiontypes.BidOwnerMapping{BidId: biddingID, BidOwner: bidder.String()}
	auction.BiddingIds = append(auction.BiddingIds, bidIDOwner)
	auction.Bidder = bidder
	auction.Bid = bid
	auction.BidEndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParam.BidDurationSeconds))
	if auction.BidEndTime.After(auction.EndTime) {
		auction.BidEndTime = auction.EndTime
	}
	err = k.SetSurplusAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) CreateNewSurplusBid(ctx sdk.Context, appID, auctionMappingID, auctionID uint64, bidder string, bid sdk.Coin) (biddingID uint64, err error) {
	auction, err := k.GetSurplusAuction(ctx, appID, auctionMappingID, auctionID)
	if err != nil {
		return biddingID, err
	}
	bidding := auctiontypes.SurplusBiddings{
		BiddingId:           k.GetUserBiddingID(ctx) + 1,
		AuctionId:           auctionID,
		AuctionStatus:       auctiontypes.ActiveAuctionStatus,
		AuctionedCollateral: auction.SellToken,
		Bidder:              bidder,
		Bid:                 bid,
		BiddingTimestamp:    ctx.BlockTime(),
		BiddingStatus:       auctiontypes.PlacedBiddingStatus,
		AppId:               appID,
		AuctionMappingId:    auctionMappingID,
	}
	k.SetUserBiddingID(ctx, bidding.BiddingId)
	err = k.SetSurplusUserBidding(ctx, bidding)
	if err != nil {
		return biddingID, err
	}
	return bidding.BiddingId, nil
}
