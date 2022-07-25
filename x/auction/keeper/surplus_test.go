package keeper_test

import (
	"github.com/comdex-official/comdex/x/auction"
	"time"

	"github.com/comdex-official/comdex/app/wasm/bindings"
	auctionKeeper "github.com/comdex-official/comdex/x/auction/keeper"
	auctionTypes "github.com/comdex-official/comdex/x/auction/types"
	tokenmintKeeper1 "github.com/comdex-official/comdex/x/tokenmint/keeper"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) WasmUpdateCollectorLookupTable(surplusThreshhold uint64, debtThreshhold uint64, lotsize uint64, debtlotsize uint64) {
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx

	msg1 := bindings.MsgUpdateCollectorLookupTable{
		AppID:            1,
		AssetID:          2,
		SurplusThreshold: surplusThreshhold,
		DebtThreshold:    debtThreshhold,
		LSR:              sdk.MustNewDecFromStr("0.001"),
		LotSize:          lotsize,
		BidFactor:        sdk.MustNewDecFromStr("0.01"),
		DebtLotSize:      debtlotsize,
	}
	err := collectorKeeper.WasmUpdateCollectorLookupTable(*ctx, &msg1)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestSurplusActivatorBetweenThreshholdAndLotsize() {
	s.AddAppAsset()
	s.AddPairAndExtendedPairVault1()
	s.AddAuctionParams()
	s.WasmSetCollectorLookupTableAndAuctionControl()
	s.WasmUpdateCollectorLookupTable(19500, 1000, 501, 300)
	s.LiquidateVaults1()

	k, ctx := &s.keeper, &s.ctx

	auction.BeginBlocker(*ctx, s.keeper)
	//err := k.SurplusActivator(*ctx)
	//s.Require().NoError(err)

	appId := uint64(1)
	auctionMappingId := uint64(1)
	auctionId := uint64(1)

	_, err := k.GetSurplusAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().Error(err)
}
func (s *KeeperTestSuite) TestSurplusActivator() {

	s.AddAppAsset()
	s.AddPairAndExtendedPairVault1()
	s.AddAuctionParams()
	s.WasmSetCollectorLookupTableAndAuctionControl()
	s.WasmUpdateCollectorLookupTable(1000, 19800, 100, 300)
	s.LiquidateVaults1()

	k, collectorKeeper, ctx := &s.keeper, &s.collectorKeeper, &s.ctx

	auction.BeginBlocker(*ctx, s.keeper)
	//s.Require().NoError(err)

	appId := uint64(1)
	auctionMappingId := uint64(1)
	auctionId := uint64(1)

	surplusAuction, err := k.GetSurplusAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

	collectorLookUp, found := collectorKeeper.GetCollectorLookupTable(*ctx, 1)
	s.Require().True(found)

	netFees, found := k.GetNetFeeCollectedData(*ctx, uint64(1))
	s.Require().True(found)

	s.Require().Equal(surplusAuction.AppId, appId)
	s.Require().Equal(surplusAuction.AuctionId, auctionId)
	s.Require().Equal(surplusAuction.AuctionMappingId, auctionMappingId)
	s.Require().Equal(surplusAuction.ActiveBiddingId, uint64(0))
	s.Require().Equal(surplusAuction.AuctionStatus, auctionTypes.AuctionStartNoBids)
	s.Require().Equal(surplusAuction.AssetInId, collectorLookUp.AssetRateInfo[0].SecondaryAssetId)
	s.Require().Equal(surplusAuction.AssetOutId, collectorLookUp.AssetRateInfo[0].CollectorAssetId)
	s.Require().Equal(surplusAuction.BidFactor, collectorLookUp.AssetRateInfo[0].BidFactor)
	s.Require().Equal(surplusAuction.SellToken.Amount.Uint64(), collectorLookUp.AssetRateInfo[0].LotSize)
	s.Require().Equal(surplusAuction.BuyToken.Amount.Uint64(), uint64(0))
	s.Require().Equal(surplusAuction.Bid.Amount.Uint64(), uint64(0))
	s.Require().True(netFees.AssetIdToFeeCollected[0].NetFeesCollected.GTE(sdk.NewIntFromUint64(collectorLookUp.AssetRateInfo[0].SurplusThreshold + collectorLookUp.AssetRateInfo[0].LotSize)))

	//Test restart surplus auction
	s.advanceseconds(301)
	auction.BeginBlocker(*ctx, s.keeper)
	s.Require().NoError(err)
	surplusAuction1, err := k.GetSurplusAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

	s.Require().Equal(surplusAuction1.AppId, appId)
	s.Require().Equal(surplusAuction1.AuctionId, auctionId)
	s.Require().Equal(surplusAuction1.AuctionMappingId, auctionMappingId)
	s.Require().Equal(surplusAuction1.ActiveBiddingId, uint64(0))
	s.Require().Equal(surplusAuction1.AuctionStatus, auctionTypes.AuctionStartNoBids)
	s.Require().Equal(surplusAuction.AssetInId, collectorLookUp.AssetRateInfo[0].SecondaryAssetId)
	s.Require().Equal(surplusAuction.AssetOutId, collectorLookUp.AssetRateInfo[0].CollectorAssetId)
	s.Require().Equal(surplusAuction1.BidFactor, collectorLookUp.AssetRateInfo[0].BidFactor)
	s.Require().Equal(surplusAuction.SellToken.Amount.Uint64(), collectorLookUp.AssetRateInfo[0].LotSize)
	s.Require().Equal(surplusAuction.BuyToken.Amount.Uint64(), uint64(0))
	s.Require().Equal(surplusAuction.Bid.Amount.Uint64(), uint64(0))
	s.Require().True(netFees.AssetIdToFeeCollected[0].NetFeesCollected.GTE(sdk.NewIntFromUint64(collectorLookUp.AssetRateInfo[0].SurplusThreshold + collectorLookUp.AssetRateInfo[0].LotSize)))
	s.Require().Equal(ctx.BlockTime().Add(time.Second*time.Duration(300)), surplusAuction1.EndTime)
}

func (s *KeeperTestSuite) TestSurplusBid() {
	userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	userAddress2 := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"
	s.TestSurplusActivator()
	k, ctx := &s.keeper, &s.ctx
	appID := uint64(1)
	auctionMappingID := uint64(1)
	auctionID := uint64(1)

	for _, tc := range []struct {
		name            string
		msg             auctionTypes.MsgPlaceSurplusBidRequest
		bidID           uint64
		isErrorExpected bool
	}{
		{
			"TestSurplusBid : Less bid AppID 1 Asset 2 0uharbor",
			auctionTypes.MsgPlaceSurplusBidRequest{
				AuctionId:        1,
				Bidder:           userAddress1,
				Amount:           ParseCoin("0uharbor"),
				AppId:            appID,
				AuctionMappingId: auctionMappingID,
			},
			1,
			true,
		},
		{
			"TestSurplusBid : More bid AppID 1 Asset 2 1uharbor",
			auctionTypes.MsgPlaceSurplusBidRequest{
				AuctionId:        1,
				Bidder:           userAddress1,
				Amount:           ParseCoin("1000uharbor"),
				AppId:            appID,
				AuctionMappingId: auctionMappingID,
			},
			1,
			false,
		},
		{
			"TestSurplusBid : less than previous bid AppID 1 Asset 2 0uharbor",
			auctionTypes.MsgPlaceSurplusBidRequest{
				AuctionId:        1,
				Bidder:           userAddress2,
				Amount:           ParseCoin("0uharbor"),
				AppId:            appID,
				AuctionMappingId: auctionMappingID,
			},
			1,
			true,
		},
		{
			"TestSurplusBid : more than previous bid AppID 1 Asset 2 1001uharbor",
			auctionTypes.MsgPlaceSurplusBidRequest{
				AuctionId:        1,
				Bidder:           userAddress2,
				Amount:           ParseCoin("1001uharbor"),
				AppId:            appID,
				AuctionMappingId: auctionMappingID,
			},
			1,
			true,
		},
		{
			"TestSurplusBid : more than previous bid AppID 1 Asset 2 1010uharbor",
			auctionTypes.MsgPlaceSurplusBidRequest{
				AuctionId:        1,
				Bidder:           userAddress2,
				Amount:           ParseCoin("1010uharbor"),
				AppId:            appID,
				AuctionMappingId: auctionMappingID,
			},
			2,
			false,
		},
	} {
		s.Run(tc.name, func() {
			server := auctionKeeper.NewMsgServiceServer(*k)
			beforeAuction, err := k.GetSurplusAuction(*ctx, appID, auctionMappingID, auctionID)
			s.Require().NoError(err)
			beforeHarborBalance, err := s.getBalance(tc.msg.Bidder, "uharbor")
			s.Require().NoError(err)
			beforeCmstBalance, err := s.getBalance(tc.msg.Bidder, "ucmst")
			s.Require().NoError(err)
			previousUserAddress := ""
			bidToken := sdk.NewCoin("zero", sdk.NewIntFromUint64(10))
			beforeHarborBalance2 := sdk.NewCoin("zero", sdk.NewIntFromUint64(10))
			if tc.bidID != uint64(1) {
				previousUserAddress = beforeAuction.Bidder.String()
				beforeHarborBalance2, err = s.getBalance(previousUserAddress, "uharbor")
				s.Require().NoError(err)
				userBid3, err := k.GetSurplusUserBidding(*ctx, previousUserAddress, appID, tc.bidID-uint64(1))
				s.Require().NoError(err)
				bidToken = userBid3.Bid
			}

			//place bid
			_, err = server.MsgPlaceSurplusBid(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.isErrorExpected {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				if tc.bidID != uint64(1) {
					afterHarborBalance2, err := s.getBalance(previousUserAddress, "uharbor")
					s.Require().NoError(err)

					s.Require().Equal(beforeHarborBalance2.Add(bidToken), afterHarborBalance2)
				}

				afterHarborBalance, err := s.getBalance(tc.msg.Bidder, "uharbor")
				s.Require().NoError(err)
				afterCmstBalance, err := s.getBalance(tc.msg.Bidder, "ucmst")
				s.Require().NoError(err)

				afterAuction, err := k.GetSurplusAuction(*ctx, appID, auctionMappingID, auctionID)
				s.Require().NoError(err)

				userBid, err := k.GetSurplusUserBidding(*ctx, tc.msg.Bidder, appID, tc.bidID)
				s.Require().NoError(err)
				bid := tc.msg.Amount

				s.Require().Equal(afterAuction.ActiveBiddingId, tc.bidID)
				s.Require().Equal(afterAuction.Bidder.String(), tc.msg.Bidder)
				s.Require().Equal(afterAuction.BiddingIds[tc.bidID-uint64(1)].BidId, tc.bidID)
				s.Require().Equal(afterAuction.BiddingIds[tc.bidID-uint64(1)].BidOwner, tc.msg.Bidder)
				if tc.bidID != uint64(1) {
					s.Require().True(afterAuction.Bid.Amount.GTE(beforeAuction.Bid.Amount.ToDec().Mul(sdk.MustNewDecFromStr("1").Sub(beforeAuction.BidFactor)).TruncateInt()))
				}
				s.Require().Equal(beforeCmstBalance, afterCmstBalance)
				s.Require().Equal(beforeHarborBalance.Sub(bid), afterHarborBalance)
				s.Require().Equal(userBid.BiddingId, tc.bidID)
				s.Require().Equal(userBid.AppId, appID)
				s.Require().Equal(userBid.AuctionId, auctionID)
				s.Require().Equal(userBid.BiddingStatus, auctionTypes.PlacedBiddingStatus)
				s.Require().Equal(userBid.AuctionStatus, auctionTypes.ActiveAuctionStatus)
				s.Require().Equal(userBid.Bidder, tc.msg.Bidder)
				s.Require().Equal(userBid.AuctionMappingId, auctionMappingID)
				s.Require().Equal(userBid.AuctionedCollateral, afterAuction.SellToken)
				s.Require().Equal(userBid.Bid, bid)
			}
		})
	}

}

func (s *KeeperTestSuite) TestCloseSurplusAuction() {
	winnerAddress := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"
	s.TestSurplusBid()
	k, ctx := &s.keeper, &s.ctx
	appID := uint64(1)
	auctionMappingID := uint64(1)
	auctionID := uint64(1)

	tokenmintKeeper, ctx := &s.tokenmintKeeper, &s.ctx
	server := tokenmintKeeper1.NewMsgServer(*tokenmintKeeper)
	msg1 := tokenminttypes.MsgMintNewTokensRequest{
		From:    winnerAddress,
		AppId:   1,
		AssetId: 3,
	}
	_, err := server.MsgMintNewTokens(sdk.WrapSDKContext(*ctx), &msg1)
	s.Require().NoError(err)

	for _, tc := range []struct {
		name            string
		seconds         uint64
		isErrorExpected bool
	}{
		{
			name:            "TestCloseSurplusAuction : less seconds than auction duration ",
			seconds:         150,
			isErrorExpected: true,
		},
		{
			name:            "TestCloseSurplusAuction : equal seconds to auction duration ",
			seconds:         150,
			isErrorExpected: true,
		},
		{
			name:            "TestCloseSurplusAuction : more seconds than auction duration ",
			seconds:         30,
			isErrorExpected: false,
		},
	} {
		s.Run(tc.name, func() {

			beforeCmstBalance, err := s.getBalance(winnerAddress, "ucmst")
			s.Require().NoError(err)

			surplusAuction, err := k.GetSurplusAuction(*ctx, appID, auctionMappingID, auctionID)
			s.Require().NoError(err)

			s.advanceseconds(int64(tc.seconds))
			auction.BeginBlocker(*ctx, s.keeper)
			s.Require().NoError(err)

			afterCmstBalance, err := s.getBalance(winnerAddress, "ucmst")
			//s.Require().NoError(err)
			//s.Require().Equal(beforeHarborBalance.Add(auction.ExpectedMintedToken), afterHarborBalance)
			if tc.isErrorExpected {
				s.Require().NotEqual(beforeCmstBalance.Add(surplusAuction.SellToken), afterCmstBalance)
			} else {

				s.Require().Equal(beforeCmstBalance.Add(surplusAuction.SellToken), afterCmstBalance)
			}
		})
	}

}
