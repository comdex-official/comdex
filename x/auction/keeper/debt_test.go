package keeper_test

import (
	"time"

	"github.com/comdex-official/comdex/x/auction"

	"github.com/comdex-official/comdex/app/wasm/bindings"
	auctionKeeper "github.com/comdex-official/comdex/x/auction/keeper"
	auctionTypes "github.com/comdex-official/comdex/x/auction/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	tokenmintKeeper1 "github.com/comdex-official/comdex/x/tokenmint/keeper"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const advanceSeconds = 21601

func (s *KeeperTestSuite) WasmSetCollectorLookupTableAndAuctionControlForSurplus() {
	// userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx

	for _, tc := range []struct {
		name string
		msg  bindings.MsgSetCollectorLookupTable
	}{
		{
			"Wasm Add MsgSetCollectorLookupTable AppID 1 CollectorAssetID 2",
			bindings.MsgSetCollectorLookupTable{
				AppID:            1,
				CollectorAssetID: 2,
				SecondaryAssetID: 3,
				SurplusThreshold: sdk.NewInt(10000000),
				DebtThreshold:    sdk.NewInt(5000000),
				LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
				LotSize:          sdk.NewInt(200000),
				BidFactor:        sdk.MustNewDecFromStr("0.01"),
				DebtLotSize:      sdk.NewInt(2000000),
			},
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &tc.msg)
			s.Require().NoError(err)
			result, found := collectorKeeper.GetCollectorLookupTable(*ctx, tc.msg.AppID, tc.msg.CollectorAssetID)
			s.Require().True(found)
			s.Require().Equal(result.AppId, tc.msg.AppID)
			s.Require().Equal(result.CollectorAssetId, tc.msg.CollectorAssetID)
			s.Require().Equal(result.SecondaryAssetId, tc.msg.SecondaryAssetID)
			s.Require().Equal(result.SurplusThreshold, tc.msg.SurplusThreshold)
			s.Require().Equal(result.DebtThreshold, tc.msg.DebtThreshold)
			s.Require().Equal(result.LockerSavingRate, tc.msg.LockerSavingRate)
			s.Require().Equal(result.LotSize, tc.msg.LotSize)
			s.Require().Equal(result.BidFactor, tc.msg.BidFactor)
			s.Require().Equal(result.DebtLotSize, tc.msg.DebtLotSize)
		})
	}
	// s.AddAuctionParams()
	for _, tc := range []struct {
		name string
		msg  bindings.MsgSetAuctionMappingForApp
	}{
		{
			"Wasm Add Auction Control AppID 1 AssetID 2",
			bindings.MsgSetAuctionMappingForApp{
				AppID:                1,
				AssetIDs:             uint64(2),
				IsSurplusAuctions:    bool(true),
				IsDebtAuctions:       bool(false),
				IsDistributor:        bool(false),
				AssetOutOraclePrices: bool(false),
				AssetOutPrices:       uint64(1000000),
			},
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.WasmSetAuctionMappingForApp(*ctx, &tc.msg)
			s.Require().NoError(err)
			result1, found := collectorKeeper.GetAuctionMappingForApp(*ctx, tc.msg.AppID, tc.msg.AssetIDs)
			s.Require().True(found)
			s.Require().Equal(result1.AssetId, tc.msg.AssetIDs)
			s.Require().Equal(result1.IsSurplusAuction, tc.msg.IsSurplusAuctions)
			s.Require().Equal(result1.IsDebtAuction, tc.msg.IsDebtAuctions)
			s.Require().Equal(result1.IsDistributor, tc.msg.IsDistributor)
			s.Require().Equal(result1.IsAuctionActive, false)
			s.Require().Equal(result1.AssetOutOraclePrice, tc.msg.AssetOutOraclePrices)
			s.Require().Equal(result1.AssetOutPrice, tc.msg.AssetOutPrices)
		})
	}
}

func (s *KeeperTestSuite) WasmSetCollectorLookupTableAndAuctionControlForDebt() {
	// userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx

	for _, tc := range []struct {
		name string
		msg  bindings.MsgSetCollectorLookupTable
	}{
		{
			"Wasm Add MsgSetCollectorLookupTable AppID 1 CollectorAssetID 2",
			bindings.MsgSetCollectorLookupTable{
				AppID:            1,
				CollectorAssetID: 2,
				SecondaryAssetID: 3,
				SurplusThreshold: sdk.NewInt(10000000),
				DebtThreshold:    sdk.NewInt(5000000),
				LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
				LotSize:          sdk.NewInt(200000),
				BidFactor:        sdk.MustNewDecFromStr("0.01"),
				DebtLotSize:      sdk.NewInt(2000000),
			},
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &tc.msg)
			s.Require().NoError(err)
			result, found := collectorKeeper.GetCollectorLookupTable(*ctx, tc.msg.AppID, tc.msg.CollectorAssetID)
			s.Require().True(found)
			s.Require().Equal(result.AppId, tc.msg.AppID)
			s.Require().Equal(result.CollectorAssetId, tc.msg.CollectorAssetID)
			s.Require().Equal(result.SecondaryAssetId, tc.msg.SecondaryAssetID)
			s.Require().Equal(result.SurplusThreshold, tc.msg.SurplusThreshold)
			s.Require().Equal(result.DebtThreshold, tc.msg.DebtThreshold)
			s.Require().Equal(result.LockerSavingRate, tc.msg.LockerSavingRate)
			s.Require().Equal(result.LotSize, tc.msg.LotSize)
			s.Require().Equal(result.BidFactor, tc.msg.BidFactor)
			s.Require().Equal(result.DebtLotSize, tc.msg.DebtLotSize)
		})
	}
	// s.AddAuctionParams()
	for _, tc := range []struct {
		name string
		msg  bindings.MsgSetAuctionMappingForApp
	}{
		{
			"Wasm Add Auction Control AppID 1 AssetID 2",
			bindings.MsgSetAuctionMappingForApp{
				AppID:                1,
				AssetIDs:             uint64(2),
				IsSurplusAuctions:    bool(false),
				IsDebtAuctions:       bool(true),
				IsDistributor:        bool(false),
				AssetOutOraclePrices: bool(false),
				AssetOutPrices:       uint64(1000000),
			},
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.WasmSetAuctionMappingForApp(*ctx, &tc.msg)
			s.Require().NoError(err)
			result1, found := collectorKeeper.GetAuctionMappingForApp(*ctx, tc.msg.AppID, tc.msg.AssetIDs)
			s.Require().True(found)
			s.Require().Equal(result1.AssetId, tc.msg.AssetIDs)
			s.Require().Equal(result1.IsSurplusAuction, tc.msg.IsSurplusAuctions)
			s.Require().Equal(result1.IsDebtAuction, tc.msg.IsDebtAuctions)
			s.Require().Equal(result1.IsDistributor, tc.msg.IsDistributor)
			s.Require().Equal(result1.IsAuctionActive, false)
			s.Require().Equal(result1.AssetOutOraclePrice, tc.msg.AssetOutOraclePrices)
			s.Require().Equal(result1.AssetOutPrice, tc.msg.AssetOutPrices)
		})
	}
}

func (s *KeeperTestSuite) TestDebtActivatorBetweenThreshholdAndLotsize() {
	s.AddAppAsset()
	s.AddPairAndExtendedPairVault1()
	s.AddAuctionParams()
	s.WasmSetCollectorLookupTableAndAuctionControlForDebt()
	s.WasmUpdateCollectorLookupTable(sdk.NewInt(30000), sdk.NewInt(20500), sdk.NewInt(800), sdk.NewInt(501))

	k, ctx := &s.keeper, &s.ctx

	auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)
	// s.Require().NoError(err)

	appId := uint64(1)
	auctionMappingId := uint64(2)
	auctionId := uint64(1)

	_, err := k.GetDebtAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestDebtActivator() {
	// userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	// addr1, err := sdk.AccAddressFromBech32(userAddress1)
	// s.Require().NoError(err)
	s.AddAppAsset()
	s.AddPairAndExtendedPairVault1()
	s.AddAuctionParams()
	s.WasmSetCollectorLookupTableAndAuctionControlForDebt()

	k, collectorKeeper, ctx := &s.keeper, &s.collectorKeeper, &s.ctx

	auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)

	appId := uint64(1)
	auctionMappingId := uint64(2)
	auctionId := uint64(1)

	collectorLookUp, found := collectorKeeper.GetCollectorLookupTable(*ctx, 1, 2)
	s.Require().True(found)

	auctionMapData, auctionMappingFound := collectorKeeper.GetAuctionMappingForApp(*ctx, appId, collectorLookUp.CollectorAssetId)
	s.Require().True(auctionMappingFound)
	klswData := esmtypes.KillSwitchParams{
		AppId:         1,
		BreakerEnable: false,
	}
	err := collectorKeeper.SetNetFeeCollectedData(*ctx, uint64(1), 2, sdk.NewIntFromUint64(4700000))
	s.Require().NoError(err)
	err1 := k.DebtActivator(*ctx, auctionMapData, klswData, false)
	s.Require().NoError(err1)
	netFees, found := collectorKeeper.GetNetFeeCollectedData(*ctx, uint64(1), 2)
	s.Require().True(found)

	debtAuction, err := k.GetDebtAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

	s.Require().Equal(debtAuction.AppId, appId)
	s.Require().Equal(debtAuction.AuctionId, auctionId)
	s.Require().Equal(debtAuction.AuctionMappingId, auctionMappingId)
	s.Require().Equal(debtAuction.ActiveBiddingId, uint64(0))
	s.Require().Equal(debtAuction.AuctionStatus, auctionTypes.AuctionStartNoBids)
	s.Require().Equal(debtAuction.AssetInId, collectorLookUp.CollectorAssetId)
	s.Require().Equal(debtAuction.AssetOutId, collectorLookUp.SecondaryAssetId)
	s.Require().Equal(debtAuction.BidFactor, collectorLookUp.BidFactor)
	s.Require().Equal(debtAuction.ExpectedUserToken.Amount, collectorLookUp.LotSize)
	s.Require().Equal(debtAuction.AuctionedToken.Amount, collectorLookUp.DebtLotSize)
	s.Require().Equal(debtAuction.ExpectedMintedToken.Amount, collectorLookUp.DebtLotSize)
	s.Require().True(netFees.NetFeesCollected.LTE(collectorLookUp.DebtThreshold.Sub(collectorLookUp.LotSize)))

	// Test restart debt auction
	s.advanceseconds(301)
	auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)
	s.Require().NoError(err)
	debtAuction1, err := k.GetDebtAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)

	s.Require().Equal(debtAuction1.AppId, appId)
	s.Require().Equal(debtAuction1.AuctionId, auctionId)
	s.Require().Equal(debtAuction1.AuctionMappingId, auctionMappingId)
	s.Require().Equal(debtAuction1.ActiveBiddingId, uint64(0))
	s.Require().Equal(debtAuction1.AuctionStatus, auctionTypes.AuctionStartNoBids)
	s.Require().Equal(debtAuction1.AssetInId, collectorLookUp.CollectorAssetId)
	s.Require().Equal(debtAuction1.AssetOutId, collectorLookUp.SecondaryAssetId)
	s.Require().Equal(debtAuction1.BidFactor, collectorLookUp.BidFactor)
	s.Require().Equal(debtAuction1.ExpectedUserToken.Amount, collectorLookUp.LotSize)
	s.Require().Equal(debtAuction1.AuctionedToken.Amount, collectorLookUp.DebtLotSize)
	s.Require().Equal(debtAuction1.ExpectedMintedToken.Amount, collectorLookUp.DebtLotSize)
	s.Require().True(netFees.NetFeesCollected.LTE(collectorLookUp.DebtThreshold.Sub(collectorLookUp.LotSize)))
	s.Require().Equal(ctx.BlockTime().Add(time.Second*time.Duration(300)), debtAuction1.EndTime)
}

func (s *KeeperTestSuite) TestDebtBid() {
	userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	userAddress2 := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"
	s.TestDebtActivator()
	k, ctx := &s.keeper, &s.ctx
	appID := uint64(1)
	auctionMappingID := uint64(2)
	auctionID := uint64(1)

	for _, tc := range []struct {
		name            string
		msg             auctionTypes.MsgPlaceDebtBidRequest
		bidID           uint64
		isErrorExpected bool
	}{
		{
			"TestDebtBid : Less ExpectedUserToken AppID 1 Asset 2 21000ucmst",
			auctionTypes.MsgPlaceDebtBidRequest{
				AuctionId:         1,
				Bidder:            userAddress1,
				Bid:               ParseCoin("2000000uharbor"),
				ExpectedUserToken: ParseCoin("21000ucmst"),
				AppId:             appID,
				AuctionMappingId:  auctionMappingID,
			},
			1,
			true,
		},
		{
			"TestDebtBid : More ExpectedUserToken AppID 1 Asset 2 200001ucmst",
			auctionTypes.MsgPlaceDebtBidRequest{
				AuctionId:         1,
				Bidder:            userAddress1,
				Bid:               ParseCoin("2000000uharbor"),
				ExpectedUserToken: ParseCoin("200001ucmst"),
				AppId:             appID,
				AuctionMappingId:  auctionMappingID,
			},
			1,
			true,
		},
		{
			"TestDebtBid : More ExpectedMintedToken AppID 1 Asset 2 2000001uharbor",
			auctionTypes.MsgPlaceDebtBidRequest{
				AuctionId:         1,
				Bidder:            userAddress1,
				Bid:               ParseCoin("2000001uharbor"),
				ExpectedUserToken: ParseCoin("200000ucmst"),
				AppId:             appID,
				AuctionMappingId:  auctionMappingID,
			},
			1,
			true,
		},
		{
			"TestDebtBid : Exact mint token AppID 1 Asset 2 21000 uharbor",
			auctionTypes.MsgPlaceDebtBidRequest{
				AuctionId:         1,
				Bidder:            userAddress1,
				Bid:               ParseCoin("2000000uharbor"),
				ExpectedUserToken: ParseCoin("200000ucmst"),
				AppId:             appID,
				AuctionMappingId:  auctionMappingID,
			},
			1,
			false,
		},
		{
			"TestDebtBid : minting more than bid factor AppID 1 Asset 2 21000 uharbor",
			auctionTypes.MsgPlaceDebtBidRequest{
				AuctionId:         1,
				Bidder:            userAddress2,
				Bid:               ParseCoin("1980001uharbor"),
				ExpectedUserToken: ParseCoin("200000ucmst"),
				AppId:             appID,
				AuctionMappingId:  auctionMappingID,
			},
			1,
			true,
		},
		{
			"TestDebtBid : minting less than bid factor AppID 1 Asset 2 21000 uharbor",
			auctionTypes.MsgPlaceDebtBidRequest{
				AuctionId:         1,
				Bidder:            userAddress2,
				Bid:               ParseCoin("1980000uharbor"),
				ExpectedUserToken: ParseCoin("200000ucmst"),
				AppId:             appID,
				AuctionMappingId:  auctionMappingID,
			},
			2,
			false,
		},
	} {
		s.Run(tc.name, func() {
			server := auctionKeeper.NewMsgServiceServer(*k)
			beforeAuction, err := k.GetDebtAuction(*ctx, appID, auctionMappingID, auctionID)
			s.Require().NoError(err)
			beforeHarborBalance, err := s.getBalance(tc.msg.Bidder, "uharbor")
			s.Require().NoError(err)
			beforeCmstBalance, err := s.getBalance(tc.msg.Bidder, "ucmst")
			s.Require().NoError(err)
			previousUserAddress := ""
			mintedToken := sdk.NewCoin("zero", sdk.NewIntFromUint64(10))
			beforeCmstBalance2 := sdk.NewCoin("zero", sdk.NewIntFromUint64(10))
			if tc.bidID != uint64(1) {
				previousUserAddress = beforeAuction.Bidder.String()
				beforeCmstBalance2, err = s.getBalance(previousUserAddress, "ucmst")
				s.Require().NoError(err)
				userBid3, err := k.GetDebtUserBidding(*ctx, previousUserAddress, appID, tc.bidID-uint64(1))
				s.Require().NoError(err)
				mintedToken = userBid3.OutflowTokens
			}

			// place bid
			_, err = server.MsgPlaceDebtBid(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.isErrorExpected {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				if tc.bidID != uint64(1) {
					afterCmstBalance2, err := s.getBalance(previousUserAddress, "ucmst")
					s.Require().NoError(err)

					s.Require().Equal(beforeCmstBalance2.Add(mintedToken), afterCmstBalance2)
				}

				afterHarborBalance, err := s.getBalance(tc.msg.Bidder, "uharbor")
				s.Require().NoError(err)
				afterCmstBalance, err := s.getBalance(tc.msg.Bidder, "ucmst")
				s.Require().NoError(err)

				afterAuction, err := k.GetDebtAuction(*ctx, appID, auctionMappingID, auctionID)
				s.Require().NoError(err)

				userBid, err := k.GetDebtUserBidding(*ctx, tc.msg.Bidder, appID, tc.bidID)
				s.Require().NoError(err)
				bid := tc.msg.Bid
				expectedUserToken := tc.msg.ExpectedUserToken
				s.Require().Equal(beforeAuction.ExpectedUserToken, afterAuction.ExpectedUserToken)
				s.Require().Equal(afterAuction.ExpectedMintedToken, bid)
				s.Require().Equal(afterAuction.ActiveBiddingId, tc.bidID)
				s.Require().Equal(afterAuction.Bidder.String(), tc.msg.Bidder)
				s.Require().Equal(afterAuction.BiddingIds[tc.bidID-uint64(1)].BidId, tc.bidID)
				s.Require().Equal(afterAuction.BiddingIds[tc.bidID-uint64(1)].BidOwner, tc.msg.Bidder)
				if tc.bidID != uint64(1) {
					s.Require().True(afterAuction.ExpectedMintedToken.Amount.LTE(sdk.NewDecFromInt(beforeAuction.ExpectedMintedToken.Amount).Mul(sdk.MustNewDecFromStr("1").Sub(beforeAuction.BidFactor)).TruncateInt()))
				}
				s.Require().Equal(beforeHarborBalance, afterHarborBalance)
				s.Require().Equal(beforeCmstBalance.Sub(expectedUserToken), afterCmstBalance)
				s.Require().Equal(userBid.BiddingId, tc.bidID)
				s.Require().Equal(userBid.AppId, appID)
				s.Require().Equal(userBid.AuctionId, auctionID)
				s.Require().Equal(userBid.BiddingStatus, auctionTypes.PlacedBiddingStatus)
				s.Require().Equal(userBid.AuctionStatus, auctionTypes.ActiveAuctionStatus)
				s.Require().Equal(userBid.Bidder, tc.msg.Bidder)
				s.Require().Equal(userBid.AuctionMappingId, auctionMappingID)
				s.Require().Equal(userBid.OutflowTokens, expectedUserToken)
				s.Require().Equal(userBid.Bid, bid)
			}
		})
	}
}

func (s *KeeperTestSuite) TestCloseDebtAuction() {
	winnerAddress := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"
	s.TestDebtBid()
	k, ctx := &s.keeper, &s.ctx
	appID := uint64(1)
	auctionMappingID := uint64(2)
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
			name:            "TestCloseDebtAuction : less seconds than auction duration ",
			seconds:         150,
			isErrorExpected: true,
		},
		{
			name:            "TestCloseDebtAuction : equal seconds to auction duration ",
			seconds:         150,
			isErrorExpected: true,
		},
		{
			name:            "TestCloseDebtAuction : more seconds than auction duration ",
			seconds:         30,
			isErrorExpected: false,
		},
	} {
		s.Run(tc.name, func() {
			beforeHarborBalance, err := s.getBalance(winnerAddress, "uharbor")
			s.Require().NoError(err)

			debtAuction, err := k.GetDebtAuction(*ctx, appID, auctionMappingID, auctionID)
			s.Require().NoError(err)

			s.advanceseconds(int64(tc.seconds))
			auction.BeginBlocker(*ctx, s.app.AuctionKeeper, s.app.AssetKeeper, s.app.CollectorKeeper, s.app.EsmKeeper)
			s.Require().NoError(err)

			afterHarborBalance, err := s.getBalance(winnerAddress, "uharbor")
			// s.Require().NoError(err)
			// s.Require().Equal(beforeHarborBalance.Add(auction.ExpectedMintedToken), afterHarborBalance)
			if tc.isErrorExpected {
				s.Require().NotEqual(beforeHarborBalance.Add(debtAuction.ExpectedMintedToken), afterHarborBalance)
			} else {
				s.Require().Equal(beforeHarborBalance.Add(debtAuction.ExpectedMintedToken), afterHarborBalance)
			}
		})
	}
}
