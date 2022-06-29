package keeper_test

import (
	"fmt"
	"github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) PrintDutchAuction(auction types.DutchAuction) {
	fmt.Println("")
	fmt.Println("printing dutch auction : ")
	fmt.Println(auction.AuctionId)
	fmt.Println(auction.OutflowTokenInitAmount)
	fmt.Println(auction.OutflowTokenCurrentAmount)
	fmt.Println(auction.InflowTokenTargetAmount)
	fmt.Println(auction.InflowTokenCurrentAmount) // make it 0
	fmt.Println(auction.OutflowTokenInitialPrice)
	fmt.Println(auction.OutflowTokenCurrentPrice)
	fmt.Println(auction.OutflowTokenEndPrice)
	fmt.Println(auction.InflowTokenCurrentPrice)
	fmt.Println(auction.EndTime)
	fmt.Println(auction.AuctionStatus)
	fmt.Println(auction.StartTime)
	fmt.Println(auction.BiddingIds)
	fmt.Println(auction.AuctionMappingId)
	fmt.Println(auction.AppId)
	fmt.Println(auction.AssetInId)
	fmt.Println(auction.AssetOutId)
	fmt.Println(auction.LockedVaultId)
	fmt.Println(auction.VaultOwner)
	fmt.Println(auction.LiquidationPenalty)

	fmt.Println("")
}

func (s *KeeperTestSuite) PrintDutchBid(bid types.DutchBiddings) {
	fmt.Println("printing dutch bid : ")
	fmt.Println(bid.BiddingId)
	fmt.Println(bid.AuctionId)
	fmt.Println(bid.AuctionStatus)
	fmt.Println(bid.OutflowTokenAmount)
	fmt.Println(bid.InflowTokenAmount)
	fmt.Println(bid.Bidder)
	fmt.Println(bid.BiddingTimestamp)
	fmt.Println(bid.BiddingStatus)
	fmt.Println(bid.AuctionMappingId)
	fmt.Println(bid.AppId)
}

func (s *KeeperTestSuite) TestCreateDutchAuction() {
	k, ctx := &s.keeper, &s.ctx
	outflowToken := ParseCoin("250denom1")
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, sdk.NewCoins(outflowToken))
	s.Require().NoError(err)

	dutchId := uint64(3)
	inflowToken := ParseCoin("100denom2")
	appId := uint64(1)
	assetInId := uint64(2)
	assetOutId := uint64(3)
	//outFlowTokenAddress, err := sdk.AccAddressFromBech32("cosmos1944ddjhecz5hfm8yp2496zu9u45vl5s7qc7vdc")
	//s.Require().NoError(err)
	//inFlowTokenAddress, err := sdk.AccAddressFromBech32("cosmos1hfml4tzwlc3mvynsg6vtgywyx00wfkhrtpkx6t")
	//s.Require().NoError(err)
	lockedVaultOwner, err := sdk.AccAddressFromBech32("cosmos1hfml4tzwlc3mvynsg6vtgywyx00wfkhrtpkx6t")
	s.Require().NoError(err)
	liquidationPenalty := sdk.MustNewDecFromStr("0.15")
	lockedVaultId := uint64(1)
	//start auction
	err = k.StartDutchAuction(*ctx, outflowToken, inflowToken, appId, assetInId, assetOutId, lockedVaultId, lockedVaultOwner.String(), liquidationPenalty)
	s.Require().NoError(err)

	auction, err := k.GetDutchAuction(*ctx, appId, dutchId, 1)
	s.Require().NoError(err)
	s.PrintDutchAuction(auction)
	s.Require().Equal(uint64(1), auction.AuctionId)
	err = k.DeleteDutchAuction(*ctx, auction)
	s.Require().NoError(err)
	d := k.GetDutchAuctions(*ctx, 1)
	fmt.Println(d)
}

func getPriceFromLinearDecreaseFunction(top sdk.Dec, tau, dur sdk.Int) sdk.Int {
	result1 := tau.Sub(dur)
	result2 := top.MulInt(result1)
	result3 := result2.Quo(tau.ToDec())
	return result3.TruncateInt()
}
func (s *KeeperTestSuite) TestLinearPriceFunction() {
	top := sdk.MustNewDecFromStr("100").Mul(sdk.MustNewDecFromStr("1"))
	tau := sdk.NewIntFromUint64(60)
	dur := sdk.NewInt(0)
	for n := 0; n <= 60; n++ {
		//fmt.Println("top tau dur seconds")
		//fmt.Println(top, tau, dur)
		//fmt.Println("price")
		//price := getPriceFromLinearDecreaseFunction(top, tau, dur)
		//fmt.Println(price)
		dur = dur.Add(sdk.NewInt(1))

	}
}

func (s *KeeperTestSuite) TestBidDutchAuction() {
	k, ctx := &s.keeper, &s.ctx
	appId := uint64(1)
	auctionMappingId := uint64(3)
	auctionId := uint64(1)
	max := sdk.MustNewDecFromStr("12")
	//create auction
	s.TestCreateDutchAuction()
	userTokens := ParseCoin("2500denom2")

	bid := ParseCoin("5denom1")
	bidder, err := sdk.AccAddressFromBech32("cosmos155hjlwufdfu4c3hycylzz74ag9anz7lkfurxwg")
	s.Require().NoError(err)

	//fund bidder
	s.fundAddr(bidder, sdk.NewCoins(userTokens))

	//place bid
	err = k.PlaceDutchAuctionBid(*ctx, appId, auctionMappingId, auctionId, bidder, bid, max)
	s.Require().NoError(err)
	auction, err := k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)
	//check if user balance reduced
	fmt.Println(k.GetBalance(*ctx, bidder, "denom2"))
	//s.Require().Equal(sdk.NewInt(880), k.GetBalance(*ctx, bidder, "denom2").Amount)
	s.PrintDutchAuction(auction)
	bidding, err := k.GetDutchUserBidding(*ctx, bidder.String(), appId, 1)
	s.Require().NoError(err)
	s.PrintDutchBid(bidding)

	//close auction by advancing time
	//s.advanceseconds(advanceSeconds)
	//k.CloseDutchAuction(*ctx, auction)

	fmt.Println(k.GetBalance(*ctx, bidder, "denom1"))
	fmt.Println(k.GetBalance(*ctx, bidder, "denom2"))
	//check if user got collateral
	s.Require().Equal(sdk.NewInt(10), k.GetBalance(*ctx, bidder, "denom1").Amount)

	//check status of bid
	bidding, err = k.GetHistoryDutchUserBidding(*ctx, bidder.String(), appId, 1)
	s.Require().NoError(err)
	s.PrintDutchBid(bidding)

	//get closed auction
	_, err = k.GetHistoryDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestBidsDutchAuction() {
	k, ctx := &s.keeper, &s.ctx
	appId := uint64(1)
	auctionMappingId := uint64(3)
	auctionId := uint64(1)
	max := sdk.MustNewDecFromStr("12")
	//create auction
	s.TestCreateDutchAuction()
	userTokens := ParseCoin("1000denom2")

	bid := ParseCoin("10denom1")
	bidder, err := sdk.AccAddressFromBech32("cosmos155hjlwufdfu4c3hycylzz74ag9anz7lkfurxwg")
	s.Require().NoError(err)

	//fund bidder
	s.fundAddr(bidder, sdk.NewCoins(userTokens))

	//place bid
	err = k.PlaceDutchAuctionBid(*ctx, appId, auctionMappingId, auctionId, bidder, bid, max)
	s.Require().NoError(err)
	auction, err := k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)
	//check if user balance reduced
	fmt.Println(k.GetBalance(*ctx, bidder, "denom2"))
	//s.Require().Equal(sdk.NewInt(880), k.GetBalance(*ctx, bidder, "denom2").Amount)
	s.PrintDutchAuction(auction)
	bidding, err := k.GetDutchUserBidding(*ctx, bidder.String(), appId, 1)
	s.Require().NoError(err)
	s.PrintDutchBid(bidding)

	//place bid
	err = k.PlaceDutchAuctionBid(*ctx, appId, auctionMappingId, auctionId, bidder, bid, max)
	s.Require().NoError(err)
	auction, err = k.GetDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)
	//check if user balance reduced
	fmt.Println(k.GetBalance(*ctx, bidder, "denom2"))
	//s.Require().Equal(sdk.NewInt(880), k.GetBalance(*ctx, bidder, "denom2").Amount)
	s.PrintDutchAuction(auction)
	bidding, err = k.GetDutchUserBidding(*ctx, bidder.String(), appId, 1)
	s.Require().NoError(err)
	s.PrintDutchBid(bidding)
	//close auction by advancing time
	s.advanceseconds(advanceSeconds)
	k.CloseDutchAuction(*ctx, auction)

	fmt.Println(k.GetBalance(*ctx, bidder, "denom1"))
	//check if user got collateral
	s.Require().Equal(sdk.NewInt(10), k.GetBalance(*ctx, bidder, "denom1").Amount)

	//check status of bid
	bidding, err = k.GetHistoryDutchUserBidding(*ctx, bidder.String(), appId, 1)
	s.Require().NoError(err)
	s.PrintDutchBid(bidding)

	//get closed auction
	_, err = k.GetHistoryDutchAuction(*ctx, appId, auctionMappingId, auctionId)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) Test() {
	_, ctx := &s.keeper, &s.ctx
	time1 := ctx.BlockTime()
	s.advanceseconds(2933222)
	fmt.Println("red")
	fmt.Println(ctx.BlockTime().Sub(time1).Seconds())
}
