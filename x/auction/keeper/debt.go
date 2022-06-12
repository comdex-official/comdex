package keeper

import (
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) DebtActivator(ctx sdk.Context) error {

	auctionMapData, found := k.GetAllAuctionMappingForApp(ctx)
	if !found{
		return nil
	}
	for _, data := range auctionMapData{
		for _, indata := range data.AssetIdToAuctionLookup{
			if indata.IsDebtAuction && !indata.IsAuctionActive{
				err:= k.CreateDebtAuction(ctx, data.AppId, indata.AssetId )
				if err !=nil{
					return err
				}
			}else if indata.IsDebtAuction && indata.IsAuctionActive{
				err := k.DebtAuctionClose(ctx, data.AppId)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (k Keeper) CreateDebtAuction(ctx sdk.Context, app_id, asset_id uint64) error {

	//check if auction status for an asset is false

	status, err := k.checkStatusOfNetFeesCollectedAndStartDebtAuction(ctx, app_id, asset_id)
	if err != nil {
		return err
	}

	auctionLookupTable, _ := k.GetAuctionMappingForApp(ctx, app_id)

	for i, assetToAuction := range auctionLookupTable.AssetIdToAuctionLookup {
		if assetToAuction.AssetId == asset_id && status == auctiontypes.StartedDebtAuction{

			auctionLookupTable.AssetIdToAuctionLookup[i].IsAuctionActive = true
		}
	}
	err1 := k.SetAuctionMappingForApp(ctx, auctionLookupTable)
	if err1 != nil {
		return err1
	}
return nil
}

func (k Keeper) checkStatusOfNetFeesCollectedAndStartDebtAuction(ctx sdk.Context, appId, assetId uint64) (status uint64, err error) {
	assetsCollectorDataUnderAppId, found := k.GetCollectorLookupTable(ctx, appId)
	if !found {
		return
	}
	//traverse this to access appId , collector asset id , surplus threshhold , debt threshhold
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

					if AssetIdToFeeCollected.NetFeesCollected.LTE(sdk.NewIntFromUint64(collector.DebtThreshold-collector.LotSize)) {
						// START DEBT AUCTION .  LOTSIZE AS MINTED FOR SECONDARY ASSET and ACCEPT Collector assetid from user
						//calculate inflow token amount
						assetInId := collector.CollectorAssetId
						assetOutId := collector.SecondaryAssetId
						//net = 200 debtThreshhold = 500 , lotsize = 100
						amount := sdk.NewIntFromUint64(collector.DebtThreshold).Sub(AssetIdToFeeCollected.NetFeesCollected)

						status, outflowToken, inflowToken := k.getDebtSellTokenAmount(ctx, appId, assetInId, assetOutId, amount)
						if status == auctiontypes.NoAuction {
							return auctiontypes.NoAuction, nil
						}

						//Mint the tokens when collector module sends tokens to user
						err := k.startDebtAuction(ctx, outflowToken, inflowToken, collector.BidFactor, appId, assetId, assetInId, assetOutId)
						if err != nil {
							break
						}
						return auctiontypes.StartedDebtAuction, nil
						// if netfees >= surplus threshhold+lotsize the start surplus auction with lot size and surplus auction is allowed true
					} else {

						return auctiontypes.NoAuction, nil
					}
				}
			}
		}
	}
	return auctiontypes.NoAuction, nil
}


func (k Keeper) getDebtSellTokenAmount(ctx sdk.Context, appId, AssetInId, AssetOutId uint64, lotSize sdk.Int) (status uint64, sellToken, buyToken sdk.Coin) {
	emptyCoin := sdk.NewCoin("empty", sdk.NewIntFromUint64(1))
	sellAsset, found1 := k.GetAsset(ctx, AssetOutId)
	buyAsset, found2 := k.GetAsset(ctx, AssetInId)
	if !found1 || !found2 {
		return auctiontypes.NoAuction, emptyCoin, emptyCoin
	}

	var buyTokenPrice uint64
	collectorAuction, found := k.GetAuctionMappingForApp(ctx, appId)
	if !found {
		return auctiontypes.NoAuction, emptyCoin, emptyCoin
	}
	for _, data := range collectorAuction.AssetIdToAuctionLookup {

		if data.AssetOutOraclePrice {
			//If oracle Price required for the assetOut
			buyTokenPrice, found = k.GetPriceForAsset(ctx, AssetInId)

		} else {
			//If oracle Price is not required for the assetOut
			buyTokenPrice = data.AssetOutPrice
		}
	}
	sellTokenPrice, found := k.GetPriceForAsset(ctx, AssetOutId)
	// inflow token will be of lot size
	buyToken = sdk.NewCoin(buyAsset.Denom, lotSize)
	outflowTokenAmount := buyToken.Amount.Mul(sdk.NewIntFromUint64(buyTokenPrice)).Quo(sdk.NewIntFromUint64(sellTokenPrice))
	sellToken = sdk.NewCoin(sellAsset.Denom, outflowTokenAmount)
	return 5, sellToken, buyToken
}

func (k Keeper) startDebtAuction(
	ctx sdk.Context,
	auctionToken sdk.Coin, //sell token
	expectedUserToken sdk.Coin, // buy token
	bidFactor sdk.Dec,
	appId, assetId uint64,
	assetInId, assetOutId uint64,
) error {

	auctionParams, found := k.GetAuctionParams(ctx, appId)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}
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
		BidFactor:           bidFactor,
		AssetInId:           assetInId,
		AssetOutId:          assetOutId,
	}
	auction.AuctionId = k.GetAuctionID(ctx) + 1
	k.SetAuctionID(ctx, auction.AuctionId)
	err := k.SetDebtAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) DebtAuctionClose(ctx sdk.Context, appId uint64) error {
	debtAuctions := k.GetDebtAuctions(ctx, appId)

	for _, debtAuction := range debtAuctions {

		if ctx.BlockTime().After(debtAuction.EndTime) {

			if debtAuction.AuctionStatus == auctiontypes.AuctionStartNoBids {

				err := k.RestartDebt(ctx, appId, debtAuction)
				if err != nil {
					return err
				}
			} else {

				err := k.closeDebtAuction(ctx, debtAuction)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (k Keeper) RestartDebt(
	ctx sdk.Context,
	appId uint64,
	debtAuction auctiontypes.DebtAuction,
) error {
	status, _, inflowToken := k.getDebtSellTokenAmount(ctx, appId, debtAuction.AssetInId, debtAuction.AssetOutId, debtAuction.ExpectedUserToken.Amount)
	if status == auctiontypes.NoAuction {
		return nil
	}
	auctionParams, found := k.GetAuctionParams(ctx, appId)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}
	debtAuction.ExpectedUserToken = inflowToken
	debtAuction.EndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
	err := k.SetDebtAuction(ctx, debtAuction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) closeDebtAuction(
	ctx sdk.Context,
	debtAuction auctiontypes.DebtAuction,
) error {

	//If there are bids
	if debtAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {

		if auctiontypes.TestFlag == 1 {
			//following 6 lines used for testing purpose
			err := k.MintCoins(ctx, auctiontypes.ModuleName, debtAuction.CurrentBidAmount)
			err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, debtAuction.Bidder, sdk.NewCoins(debtAuction.CurrentBidAmount))
			if err != nil {
				return err
			}
		} else {
			//ask token mint to mint new tokens for bidder address

			err := k.MintNewTokensForApp(ctx, debtAuction.AppId, debtAuction.AssetOutId, debtAuction.Bidder.String(), debtAuction.CurrentBidAmount.Amount)
			if err != nil {

				return err
			}
		}

		bidding, err := k.GetDebtUserBidding(ctx, debtAuction.Bidder.String(), debtAuction.AppId, debtAuction.ActiveBiddingId)
		if err != nil {
			return err
		}
		bidding.BiddingStatus = auctiontypes.SuccessBiddingStatus
		err = k.SetDebtUserBidding(ctx, bidding)

		if err != nil {
			return err
		}
		for _, biddingId := range debtAuction.BiddingIds {
			bidding, err := k.GetDebtUserBidding(ctx, biddingId.BidOwner, debtAuction.AppId, biddingId.BidId)
			if err != nil {
				return err
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			err = k.SetDebtUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.DeleteDebtUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.SetHistoryDebtUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
		}

		//send to collector module the amount collected in debt auction

		err = k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(debtAuction.ExpectedUserToken))

		if err != nil {

			return err
		}

		err = k.SetNetFeeCollectedData(ctx, debtAuction.AuctionId, debtAuction.AssetInId, debtAuction.ExpectedUserToken.Amount)
		if err != nil {

			return auctiontypes.ErrorUnableToSetNetfees
		}

	}

	err := k.makeFalseForFlags(ctx, debtAuction.AppId, debtAuction.AssetId)
	if err != nil {
		return auctiontypes.ErrorUnableToMakeFlagsFalse
	}
	err = k.DeleteDebtAuction(ctx, debtAuction)
	if err != nil {
		return err
	}
	debtAuction.AuctionStatus = auctiontypes.AuctionEnded
	err = k.SetHistoryDebtAuction(ctx, debtAuction)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) PlaceDebtAuctionBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, expectedUserToken sdk.Coin) error {
	auction, err := k.GetDebtAuction(ctx, appId, auctionMappingId, auctionId)

	if err != nil {
		return auctiontypes.ErrorInvalidDebtAuctionId
	}
	if expectedUserToken.Denom != auction.ExpectedUserToken.Denom {
		return auctiontypes.ErrorInvalidDebtUserExpectedDenom
	}

	if !expectedUserToken.Amount.Equal(auction.ExpectedUserToken.Amount) {
		return auctiontypes.ErrorDebtExpectedUserAmount
	}
	if bid.Denom != auction.ExpectedMintedToken.Denom {
		return auctiontypes.ErrorInvalidDebtMintedDenom
	}

	//Test this multiplication check if new bid greater than previous bid by bid factor
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		change := auction.BidFactor.MulInt(auction.ExpectedMintedToken.Amount).Ceil().TruncateInt()
		maxBidAmount := auction.ExpectedMintedToken.Amount.Sub(change)
		if bid.Amount.GT(maxBidAmount) {
			sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bid should be less than or equal to %d ", maxBidAmount)
		}
	} else {
		if bid.Amount.GT(auction.AuctionedToken.Amount) {
			return auctiontypes.ErrorMaxBidAmount
		}
	}
	err = k.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(expectedUserToken))
	if err != nil {
		return err
	}

	biddingId, err := k.CreateNewDebtBid(ctx, appId, auctionMappingId, auctionId, bidder, bid, expectedUserToken)
	if err != nil {
		return err
	}
	//If auction gets bid from second time onwards . refund previous bidder
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, auction.Bidder, sdk.NewCoins(auction.ExpectedUserToken))
		if err != nil {
			return err
		}
		bidding, err := k.GetDebtUserBidding(ctx, auction.Bidder.String(), auction.AppId, auction.ActiveBiddingId)
		if err != nil {
			return err
		}
		bidding.BiddingStatus = auctiontypes.RejectedBiddingStatus

		err = k.SetDebtUserBidding(ctx, bidding)
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
	auction.CurrentBidAmount = bid
	auction.ExpectedMintedToken = bid
	err = k.SetDebtAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
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

	k.SetUserBiddingID(ctx, bidding.BiddingId)

	err = k.SetDebtUserBidding(ctx, bidding)
	if err != nil {
		return biddingId, err
	}

	return bidding.BiddingId, nil
}