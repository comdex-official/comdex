package keeper_test

import (
	"fmt"
	"github.com/comdex-official/comdex/app/wasm/bindings"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const advanceSeconds = 21601

func (s *KeeperTestSuite) WasmSetCollectorLookupTableAndAuctionControl() {
	//userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	collectorKeeper, ctx := &s.collectorKeeper, &s.ctx

	for index, tc := range []struct {
		name string
		msg  bindings.MsgSetCollectorLookupTable
	}{
		{"Wasm Add MsgSetCollectorLookupTable AppID 1 CollectorAssetID 2",
			bindings.MsgSetCollectorLookupTable{
				AppID:            1,
				CollectorAssetID: 2,
				SecondaryAssetID: 3,
				SurplusThreshold: 10000000,
				DebtThreshold:    5000000,
				LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
				LotSize:          2000000,
				BidFactor:        sdk.MustNewDecFromStr("0.01"),
				DebtLotSize:      2000000,
			},
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &tc.msg)
			s.Require().NoError(err)
			result, found := collectorKeeper.GetCollectorLookupTable(*ctx, tc.msg.AppID)
			s.Require().True(found)
			s.Require().Equal(result.AssetRateInfo[index].AppId, tc.msg.AppID)
			s.Require().Equal(result.AssetRateInfo[index].CollectorAssetId, tc.msg.CollectorAssetID)
			s.Require().Equal(result.AssetRateInfo[index].SecondaryAssetId, tc.msg.SecondaryAssetID)
			s.Require().Equal(result.AssetRateInfo[index].SurplusThreshold, tc.msg.SurplusThreshold)
			s.Require().Equal(result.AssetRateInfo[index].DebtThreshold, tc.msg.DebtThreshold)
			s.Require().Equal(result.AssetRateInfo[index].LockerSavingRate, tc.msg.LockerSavingRate)
			s.Require().Equal(result.AssetRateInfo[index].LotSize, tc.msg.LotSize)
			s.Require().Equal(result.AssetRateInfo[index].BidFactor, tc.msg.BidFactor)
			s.Require().Equal(result.AssetRateInfo[index].DebtLotSize, tc.msg.DebtLotSize)
		})
	}
	//s.AddAuctionParams()
	for index, tc := range []struct {
		name string
		msg  bindings.MsgSetAuctionMappingForApp
	}{
		{
			"Wasm Add Auction Control AppID 1 AssetID 2",
			bindings.MsgSetAuctionMappingForApp{
				AppID:                1,
				AssetIDs:             []uint64{2},
				IsSurplusAuctions:    []bool{true},
				IsDebtAuctions:       []bool{true},
				AssetOutOraclePrices: []bool{false},
				AssetOutPrices:       []uint64{1000000},
			},
		},
	} {
		s.Run(tc.name, func() {
			err := collectorKeeper.WasmSetAuctionMappingForApp(*ctx, &tc.msg)
			s.Require().NoError(err)
			result1, found := collectorKeeper.GetAuctionMappingForApp(*ctx, tc.msg.AppID)
			s.Require().True(found)
			s.Require().Equal(result1.AssetIdToAuctionLookup[index].AssetId, tc.msg.AssetIDs[0])
			s.Require().Equal(result1.AssetIdToAuctionLookup[index].IsSurplusAuction, tc.msg.IsSurplusAuctions[0])
			s.Require().Equal(result1.AssetIdToAuctionLookup[index].IsDebtAuction, tc.msg.IsDebtAuctions[0])
			s.Require().Equal(result1.AssetIdToAuctionLookup[index].IsAuctionActive, false)
			s.Require().Equal(result1.AssetIdToAuctionLookup[index].AssetOutOraclePrice, tc.msg.AssetOutOraclePrices[0])
			s.Require().Equal(result1.AssetIdToAuctionLookup[index].AssetOutPrice, tc.msg.AssetOutPrices[0])
		})
	}

}

func (s *KeeperTestSuite) TestDebtActivator() {
	//userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	//addr1, err := sdk.AccAddressFromBech32(userAddress1)
	//s.Require().NoError(err)
	s.AddAppAsset()
	s.AddPairAndExtendedPairVault1()
	s.AddAuctionParams()
	s.WasmSetCollectorLookupTableAndAuctionControl()
	s.LiquidateVaults1()

	k, ctx := &s.keeper, &s.ctx

	err := k.DebtActivator(*ctx)
	s.Require().NoError(err)

	debtAuctions := k.GetDebtAuctions(*ctx, 1)
	for _, auction := range debtAuctions {
		fmt.Printf("%+v\n", auction)
	}
	//appId := uint64(1)
	//auctionMappingId := uint64(3)
	//auctionId := uint64(1)

}

//func (s *KeeperTestSuite) TestCreateDebtAuction() {
//	bidFactor := sdk.MustNewDecFromStr("0.01")
//	assetInId := uint64(1)
//	assetOutId := uint64(2)
//	k, ctx := &s.keeper, &s.ctx
//	auctionMapppingId := uint64(2)
//	outflowToken := ParseCoin("250denom1")
//	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, sdk.NewCoins(outflowToken))
//	s.Require().NoError(err)
//
//	inflowToken := ParseCoin("100denom2")
//	appId := uint64(1)
//	assetId := uint64(2)
//
//	//start auction
//	err = k.StartDebtAuction(*ctx, outflowToken, inflowToken, bidFactor, appId, assetId, assetInId, assetOutId)
//	s.Require().NoError(err)
//
//	auction, err := k.GetDebtAuction(*ctx, appId, auctionMapppingId, 1)
//	s.Require().NoError(err)
//	s.PrintDebtAuction(auction)
//	s.Require().Equal(uint64(1), auction.AuctionId)
//}
//
//func (s *KeeperTestSuite) TestBidDebtAuction() {
//	k, ctx := &s.keeper, &s.ctx
//	appId := uint64(1)
//	auctionMappingId := uint64(2)
//	auctionId := uint64(1)
//	//create auction
//	s.TestCreateDebtAuction()
//	userTokens := ParseCoin("1000denom2")
//
//	expectedUserToken := ParseCoin("100denom2")
//	bid := ParseCoin("245denom1")
//	bidder, err := sdk.AccAddressFromBech32("cosmos155hjlwufdfu4c3hycylzz74ag9anz7lkfurxwg")
//	s.Require().NoError(err)
//
//	//fund bidder
//	s.fundAddr(bidder, userTokens)
//
//	//place bid
//	err = k.PlaceDebtAuctionBid(*ctx, appId, auctionMappingId, auctionId, bidder, bid, expectedUserToken)
//	s.Require().NoError(err)
//
//	auction, err := k.GetDebtAuction(*ctx, appId, auctionMappingId, auctionId)
//	s.Require().NoError(err)
//	s.Require().Equal(auction.ActiveBiddingId, uint64(1))
//	//check if user balance reduced
//	s.Require().Equal(sdk.NewInt(900), k.GetBalance(*ctx, bidder, "denom2").Amount)
//	s.Require().Equal(bidder, auction.Bidder)
//	s.PrintDebtAuction(auction)
//
//	bidding, err := k.GetDebtUserBidding(*ctx, bidder.String(), appId, 1)
//	s.Require().NoError(err)
//	s.PrintDebtBid(bidding)
//
//	//close auction by advancing time
//	s.advanceseconds(advanceSeconds)
//	err = k.DebtAuctionClose(*ctx, appId)
//	s.Require().NoError(err)
//	fmt.Println(k.GetBalance(*ctx, bidder, "denom1").Amount)
//	//check if user got collateral
//	s.Require().Equal(sdk.NewInt(250), k.GetBalance(*ctx, bidder, "denom1").Amount)
//
//	//check status of bid
//	bidding, err = k.GetHistoryDebtUserBidding(*ctx, bidder.String(), appId, 1)
//	s.Require().NoError(err)
//	s.PrintDebtBid(bidding)
//
//	//get closed auction
//	_, err = k.GetHistoryDebtAuction(*ctx, appId, auctionMappingId, 1)
//	s.Require().NoError(err)
//}
//
//func (s *KeeperTestSuite) TestBidsDebtAuction() {
//	k, ctx := &s.keeper, &s.ctx
//	appId := uint64(1)
//	auctionMappingId := uint64(2)
//	auctionId := uint64(1)
//	//create auction
//	s.TestCreateDebtAuction()
//	userTokens := ParseCoin("1000denom2")
//
//	expectedUserToken := ParseCoin("100denom2")
//
//	//bid1
//	bid := ParseCoin("240denom1")
//	bidder, err := sdk.AccAddressFromBech32("cosmos155hjlwufdfu4c3hycylzz74ag9anz7lkfurxwg")
//	s.Require().NoError(err)
//
//	//fund bidder
//	s.fundAddr(bidder, userTokens)
//
//	//place bid1
//	err = k.PlaceDebtAuctionBid(*ctx, appId, auctionMappingId, auctionId, bidder, bid, expectedUserToken)
//	s.Require().NoError(err)
//	auction, err := k.GetDebtAuction(*ctx, appId, auctionMappingId, auctionId)
//	s.Require().NoError(err)
//	s.Require().Equal(auction.ActiveBiddingId, uint64(1))
//	//check if user balance reduced
//	s.Require().Equal(sdk.NewInt(900), k.GetBalance(*ctx, bidder, "denom2").Amount)
//	s.Require().Equal(bidder, auction.Bidder)
//	s.PrintDebtAuction(auction)
//	bidding, err := k.GetDebtUserBidding(*ctx, bidder.String(), appId, 1)
//	s.Require().NoError(err)
//	s.PrintDebtBid(bidding)
//
//	//bid2
//	bidder2, err := sdk.AccAddressFromBech32("cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v")
//	s.Require().NoError(err)
//
//	//fund bidder2
//	s.fundAddr(bidder2, userTokens)
//	bid = ParseCoin("235denom1")
//
//	//place bid
//	err = k.PlaceDebtAuctionBid(*ctx, appId, auctionMappingId, auctionId, bidder2, bid, expectedUserToken)
//	fmt.Println(err)
//	s.Require().NoError(err)
//	//place bid2
//	auction, err = k.GetDebtAuction(*ctx, appId, auctionMappingId, auctionId)
//	s.Require().NoError(err)
//	s.Require().Equal(auction.ActiveBiddingId, uint64(2))
//	s.PrintDebtAuction(auction)
//	s.Require().Equal(sdk.NewInt(900), k.GetBalance(*ctx, bidder2, "denom2").Amount)
//	s.Require().Equal(sdk.NewInt(1000), k.GetBalance(*ctx, bidder, "denom2").Amount)
//
//	//close auction by advancing time
//	s.advanceseconds(advanceSeconds)
//	err = k.DebtAuctionClose(*ctx, appId)
//	s.Require().NoError(err)
//	fmt.Println(k.GetBalance(*ctx, bidder2, "denom1").Amount)
//	//check if user got collateral
//	s.Require().Equal(sdk.NewInt(243), k.GetBalance(*ctx, bidder2, "denom1").Amount)
//
//	//check status of bid
//	bidding, err = k.GetHistoryDebtUserBidding(*ctx, bidder.String(), appId, 1)
//	s.Require().NoError(err)
//	s.PrintDebtBid(bidding)
//
//	//get closed auction
//	_, err = k.GetHistoryDebtAuction(*ctx, appId, auctionMappingId, auctionId)
//	s.Require().NoError(err)
//}

/* 1. auction closed but no bids then auction should restart
2. make sure code enter every if condition or every part of code
3. make sure to change numbers
4.test maths this where quo involved with smaller numbers and also places sdk.Int involved in quo
5.check chost block entered
6.advance time in dutch test so that collateral gets over but target cmst is still left
7.make fetch price dynamic in dutch as of now linear is hard coded
8.handle if auctiontype is empty in setters
// */
