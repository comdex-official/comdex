package keeper

import (
	"time"

	esmtypes "github.com/comdex-official/comdex/x/esm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
)

func (k Keeper) DebtActivator(ctx sdk.Context, data collectortypes.AppAssetIdToAuctionLookupTable, killSwitchParams esmtypes.KillSwitchParams, status bool) error {
	if data.IsDebtAuction && !data.IsAuctionActive && !killSwitchParams.BreakerEnable && !status {
		err := k.CreateDebtAuction(ctx, data.AppId, data.AssetId)
		if err != nil {
			return err
		}
	}
	if data.IsDebtAuction && data.IsAuctionActive {
		err := k.DebtAuctionClose(ctx, data.AppId, status)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) CreateDebtAuction(ctx sdk.Context, appID, assetID uint64) error { // check if auction status for an asset is false
	status, err := k.checkStatusOfNetFeesCollectedAndStartDebtAuction(ctx, appID, assetID)
	if err != nil {
		return err
	}

	auctionLookupTable, _ := k.GetAuctionMappingForApp(ctx, appID, assetID)

	if status == auctiontypes.StartedDebtAuction {
		auctionLookupTable.IsAuctionActive = true
		err1 := k.SetAuctionMappingForApp(ctx, auctionLookupTable)
		if err1 != nil {
			return err1
		}
	}

	return nil
}

//nolint:unparam
func (k Keeper) checkStatusOfNetFeesCollectedAndStartDebtAuction(ctx sdk.Context, appID, assetID uint64) (status uint64, err error) {
	collector, found := k.GetCollectorLookupTable(ctx, appID, assetID)
	if !found {
		return
	}
	// traverse this to access appId , collector asset id  , debt threshold

	NetFeeCollectedData, found := k.GetNetFeeCollectedData(ctx, appID, assetID)

	if !found {
		return auctiontypes.NoAuction, nil
	}
	// traverse this to access appId , collector asset id , netfees collected

	if NetFeeCollectedData.NetFeesCollected.LTE(sdk.NewIntFromUint64(collector.DebtThreshold - collector.LotSize)) {
		// START DEBT AUCTION .  LOTSIZE AS MINTED FOR SECONDARY ASSET and ACCEPT Collector assetid from user
		// calculate inflow token amount
		assetInID := collector.CollectorAssetId
		assetOutID := collector.SecondaryAssetId
		// net = 200 debtThreshold = 500 , lotsize = 100
		amount := sdk.NewIntFromUint64(collector.LotSize)

		status, outflowToken, inflowToken := k.getDebtSellTokenAmount(ctx, appID, assetInID, assetOutID, amount)
		if status == auctiontypes.NoAuction {
			return auctiontypes.NoAuction, nil
		}

		// Mint the tokens when collector module sends tokens to user
		err := k.StartDebtAuction(ctx, outflowToken, inflowToken, collector.BidFactor, appID, assetID, assetInID, assetOutID)
		if err != nil {
			return auctiontypes.NoAuction, nil
		}
		return auctiontypes.StartedDebtAuction, nil
	}

	return auctiontypes.NoAuction, nil
}

func (k Keeper) getDebtSellTokenAmount(ctx sdk.Context, appID, AssetInID, AssetOutID uint64, lotSize sdk.Int) (status uint64, sellToken, buyToken sdk.Coin) {
	emptyCoin := sdk.NewCoin("empty", sdk.NewIntFromUint64(1))
	sellAsset, found1 := k.GetAsset(ctx, AssetOutID)
	buyAsset, found2 := k.GetAsset(ctx, AssetInID)
	if !found1 || !found2 {
		return auctiontypes.NoAuction, emptyCoin, emptyCoin
	}
	var debtLot uint64
	collectorData, _ := k.GetCollectorLookupTable(ctx, appID, AssetInID)

	debtLot = collectorData.DebtLotSize

	buyToken = sdk.NewCoin(buyAsset.Denom, lotSize)
	sellToken = sdk.NewCoin(sellAsset.Denom, sdk.NewIntFromUint64(debtLot))
	return 5, sellToken, buyToken
}

func (k Keeper) StartDebtAuction(
	ctx sdk.Context,
	auctionToken sdk.Coin, // sell token
	expectedUserToken sdk.Coin, // buy token
	bidFactor sdk.Dec,
	appID, assetID uint64,
	assetInID, assetOutID uint64,
) error {
	auctionParams, found := k.GetAuctionParams(ctx, appID)
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
		BidEndTime:          ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		CurrentBidAmount:    sdk.NewCoin(auctionToken.Denom, sdk.NewInt(0)),
		AuctionStatus:       auctiontypes.AuctionStartNoBids,
		AppId:               appID,
		AssetId:             assetID,
		BiddingIds:          []*auctiontypes.BidOwnerMapping{},
		AuctionMappingId:    auctionParams.DebtId,
		BidFactor:           bidFactor,
		AssetInId:           assetInID,
		AssetOutId:          assetOutID,
	}
	auction.AuctionId = k.GetAuctionID(ctx) + 1
	k.SetAuctionID(ctx, auction.AuctionId)
	err := k.SetDebtAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) DebtAuctionClose(ctx sdk.Context, appID uint64, statusEsm bool) error {
	debtAuctions := k.GetDebtAuctions(ctx, appID)

	for _, debtAuction := range debtAuctions {
		if ctx.BlockTime().After(debtAuction.EndTime) || ctx.BlockTime().After(debtAuction.BidEndTime) || statusEsm {
			if (debtAuction.AuctionStatus == auctiontypes.AuctionStartNoBids) && !statusEsm {
				err := k.RestartDebt(ctx, appID, debtAuction)
				if err != nil {
					return err
				}
			} else {
				err := k.closeDebtAuction(ctx, debtAuction, statusEsm)
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
	appID uint64,
	debtAuction auctiontypes.DebtAuction,
) error {
	status, _, inflowToken := k.getDebtSellTokenAmount(ctx, appID, debtAuction.AssetInId, debtAuction.AssetOutId, debtAuction.ExpectedUserToken.Amount)
	if status == auctiontypes.NoAuction {
		return nil
	}
	auctionParams, found := k.GetAuctionParams(ctx, appID)
	if !found {
		return auctiontypes.ErrorInvalidAuctionParams
	}
	debtAuction.ExpectedUserToken = inflowToken
	debtAuction.EndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
	debtAuction.BidEndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
	err := k.SetDebtAuction(ctx, debtAuction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) closeDebtAuction(
	ctx sdk.Context,
	debtAuction auctiontypes.DebtAuction,
	statusEsm bool,
) error { // If there are bids
	if statusEsm && debtAuction.BiddingIds != nil {
		bidding, err := k.GetDebtUserBidding(ctx, debtAuction.Bidder.String(), debtAuction.AppId, debtAuction.ActiveBiddingId)
		if err != nil {
			return err
		}
		err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, sdk.AccAddress(bidding.Bidder), sdk.NewCoins(debtAuction.ExpectedUserToken))
		if err != nil {
			return err
		}
	} else if (debtAuction.AuctionStatus != auctiontypes.AuctionStartNoBids) && !statusEsm {
		err := k.MintNewTokensForApp(ctx, debtAuction.AppId, debtAuction.AssetOutId, debtAuction.Bidder.String(), debtAuction.CurrentBidAmount.Amount)
		if err != nil {
			return err
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
		for _, biddingID := range debtAuction.BiddingIds {
			bidding, err := k.GetDebtUserBidding(ctx, biddingID.BidOwner, debtAuction.AppId, biddingID.BidId)
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

		// send to collector module the amount collected in debt auction

		err = k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(debtAuction.ExpectedUserToken))

		if err != nil {
			return err
		}

		err = k.SetNetFeeCollectedData(ctx, debtAuction.AppId, debtAuction.AssetInId, debtAuction.ExpectedUserToken.Amount)
		if err != nil {
			return auctiontypes.ErrorUnableToSetNetFees
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

func (k Keeper) PlaceDebtAuctionBid(ctx sdk.Context, appID, auctionMappingID, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin, expectedUserToken sdk.Coin) error {
	auction, err := k.GetDebtAuction(ctx, appID, auctionMappingID, auctionID)
	auctionParam, _ := k.GetAuctionParams(ctx, appID)

	if err != nil {
		return auctiontypes.ErrorInvalidDebtAuctionID
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

	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		change := auction.BidFactor.MulInt(auction.ExpectedMintedToken.Amount).Ceil().TruncateInt()
		maxBidAmount := auction.ExpectedMintedToken.Amount.Sub(change)
		if bid.Amount.GT(maxBidAmount) {
			return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bid should be less than or equal to %d ", maxBidAmount.Uint64())
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

	biddingID, err := k.CreateNewDebtBid(ctx, appID, auctionMappingID, auctionID, bidder.String(), bid, expectedUserToken)
	if err != nil {
		return err
	}
	// If auction gets bid from second time onwards . refund previous bidder
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
	auction.ActiveBiddingId = biddingID
	bidIDOwner := &auctiontypes.BidOwnerMapping{BidId: biddingID, BidOwner: bidder.String()}
	auction.BiddingIds = append(auction.BiddingIds, bidIDOwner)
	auction.Bidder = bidder
	auction.CurrentBidAmount = bid
	auction.ExpectedMintedToken = bid
	auction.BidEndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParam.BidDurationSeconds))
	if auction.BidEndTime.After(auction.EndTime) {
		auction.BidEndTime = auction.EndTime
	}
	err = k.SetDebtAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) CreateNewDebtBid(ctx sdk.Context, appID, auctionMappingID, auctionID uint64, bidder string, bid sdk.Coin, expectedUserToken sdk.Coin) (biddingID uint64, err error) {
	bidding := auctiontypes.DebtBiddings{
		BiddingId:        k.GetUserBiddingID(ctx) + 1,
		AuctionId:        auctionID,
		AuctionStatus:    auctiontypes.ActiveAuctionStatus,
		Bidder:           bidder,
		Bid:              bid,
		BiddingTimestamp: ctx.BlockTime(),
		BiddingStatus:    auctiontypes.PlacedBiddingStatus,
		AppId:            appID,
		AuctionMappingId: auctionMappingID,
		OutflowTokens:    expectedUserToken,
	}

	k.SetUserBiddingID(ctx, bidding.BiddingId)

	err = k.SetDebtUserBidding(ctx, bidding)
	if err != nil {
		return biddingID, err
	}

	return bidding.BiddingId, nil
}
