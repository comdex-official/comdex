package keeper_test

import (
	"fmt"
	"github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) PrintAuction(auction types.SurplusAuction) {
	fmt.Println("printing surplus auction : ")
	fmt.Println(auction.AuctionId)
	fmt.Println(auction.AppId)
	fmt.Println(auction.AssetId)
	fmt.Println(auction.BuyToken)
	fmt.Println(auction.SellToken)
	fmt.Println(auction.Bidder)
	fmt.Println(auction.Bid)
	fmt.Println(auction.ActiveBiddingId)
	fmt.Println(auction.EndTime)
	fmt.Println(auction.BidFactor)
	fmt.Println(auction.BiddingIds)
	fmt.Println(auction.AuctionStatus)
}
func (s *KeeperTestSuite) PrintSurplusBid(bid types.SurplusBiddings) {
	fmt.Println("printing surplus bid : ")
	fmt.Println(bid.BiddingId)
	fmt.Println(bid.AuctionId)
	fmt.Println(bid.AuctionStatus)
	fmt.Println(bid.AuctionedCollateral)
	fmt.Println(bid.Bidder)
	fmt.Println(bid.Bid)
	fmt.Println(bid.BiddingTimestamp)
	fmt.Println(bid.BiddingStatus)

}
func (s *KeeperTestSuite) TestCreateSurplusAuction() {
	k, ctx := &s.keeper, &s.ctx
	outflowToken := ParseCoin("10000denom1")
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, sdk.NewCoins(outflowToken))
	s.Require().NoError(err)

	inflowToken := ParseCoin("100denom2")
	bidFactor := sdk.MustNewDecFromStr("0.01")
	appId := uint64(1)
	assetId := uint64(2)
	assetInId := uint64(1)
	assetOutId := uint64(2)
	err = k.StartSurplusAuction(*ctx, outflowToken, inflowToken, bidFactor, appId, assetId, assetInId, assetOutId)
	s.Require().NoError(err)
	auction, err := k.GetSurplusAuction(*ctx, appId, 1, 1)
	fmt.Println(auction)
	s.Require().NoError(err)
	//s.PrintAuction(auction)
	s.Require().Equal(uint64(1), auction.AuctionId)

}

func (s *KeeperTestSuite) TestBidSurplusAuction() {
	k, ctx := &s.keeper, &s.ctx
	appId := uint64(1)
	//create auction
	s.TestCreateSurplusAuction()
	userTokens := ParseCoin("1000denom2")
	bid := ParseCoin("100denom2")
	bidder, err := sdk.AccAddressFromBech32("cosmos155hjlwufdfu4c3hycylzz74ag9anz7lkfurxwg")
	s.Require().NoError(err)

	//fund bidder
	s.fundAddr(bidder, sdk.NewCoins(userTokens))

	//place bid
	err = k.PlaceSurplusAuctionBid(*ctx, appId, 1, 1, bidder, bid)
	s.Require().NoError(err)
	auction, err := k.GetSurplusAuction(*ctx, appId, 1, 1)
	s.Require().NoError(err)
	s.Require().Equal(auction.ActiveBiddingId, uint64(1))
	s.Require().Equal(sdk.NewInt(900), k.GetBalance(*ctx, bidder, "denom2").Amount)
	s.Require().Equal(bidder, auction.Bidder)
	s.PrintAuction(auction)
	bidding, err := k.GetSurplusUserBidding(*ctx, bidder.String(), appId, 1)
	s.Require().NoError(err)
	s.PrintSurplusBid(bidding)

	//close auction by advancing time
	s.advanceseconds(advanceSeconds)
	err = k.SurplusAuctionClose(*ctx, appId)
	s.Require().NoError(err)

	//check if user got collateral
	s.Require().Equal(sdk.NewInt(10000), k.GetBalance(*ctx, bidder, "denom1").Amount)

	//check status of bid
	bidding, err = k.GetSurplusUserBidding(*ctx, bidder.String(), appId, 1)
	s.Require().NoError(err)
	s.PrintSurplusBid(bidding)
	//get closed auction
	_, err = k.GetSurplusAuction(*ctx, appId, 1, 1)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestBidsSurplusAuction() {
	k, ctx := &s.keeper, &s.ctx
	appId := uint64(1)
	surplusAuction := uint64(1)
	s.TestCreateSurplusAuction()
	userTokens := ParseCoin("1000denom2")
	fmt.Println("bid 1")

	//bid1
	bid := ParseCoin("100denom2")
	bidder, err := sdk.AccAddressFromBech32("cosmos155hjlwufdfu4c3hycylzz74ag9anz7lkfurxwg")
	s.Require().NoError(err)
	s.fundAddr(bidder, sdk.NewCoins(userTokens))
	//place bid1
	err = k.PlaceSurplusAuctionBid(*ctx, appId, surplusAuction, 1, bidder, bid)
	fmt.Println(err)
	s.Require().NoError(err)
	auction, err := k.GetSurplusAuction(*ctx, appId, surplusAuction, 1)
	s.Require().NoError(err)
	s.Require().Equal(auction.ActiveBiddingId, uint64(1))
	//s.Require().Equal(sdk.NewInt(900), k.GetBalance(*ctx, bidder, "denom2").Amount)
	s.PrintAuction(auction)

	//bid2
	bidder2, err := sdk.AccAddressFromBech32("cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v")
	s.Require().NoError(err)

	//fund bidder2
	s.fundAddr(bidder2, sdk.NewCoins(userTokens))
	fmt.Println("bid 2")
	bid2 := ParseCoin("120denom2")

	//place bid
	err = k.PlaceSurplusAuctionBid(*ctx, appId, surplusAuction, 1, bidder2, bid2)
	fmt.Println(err)
	s.Require().NoError(err)
	//place bid2
	auction, err = k.GetSurplusAuction(*ctx, appId, surplusAuction, 1)
	s.Require().NoError(err)
	s.Require().Equal(auction.ActiveBiddingId, uint64(2))
	s.PrintAuction(auction)
	fmt.Println(k.GetBalance(*ctx, bidder2, "denom2"))
	s.Require().Equal(sdk.NewInt(880), k.GetBalance(*ctx, bidder2, "denom2").Amount)
	s.Require().Equal(sdk.NewInt(1000), k.GetBalance(*ctx, bidder, "denom2").Amount)

	//close auction by advancing time
	s.advanceseconds(61)
	err = k.SurplusAuctionClose(*ctx, appId)
	s.Require().NoError(err)
	//check if user got collateral
	s.Require().Equal(sdk.NewInt(10000), k.GetBalance(*ctx, bidder2, "denom1").Amount)

	//print bid1 and bid2
	//bidding, found1 := k.GetSurplusUserBidding(*ctx, bidder.String(), appId, 1)
	//s.Require().True(found1)
	//s.PrintSurplusBid(bidding)
	//bidding2, found2 := k.GetSurplusUserBidding(*ctx, bidder2.String(), appId, 2)
	//s.Require().True(found2)
	//s.PrintSurplusBid(bidding2)
	////get closed auction
	//_, found = k.GetSurplusAuction(*ctx, appId, 1, 1)
	//s.Require().False(found)
}
