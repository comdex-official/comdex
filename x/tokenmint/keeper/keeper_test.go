package keeper_test

import (
	"testing"
	"time"

	collectorTypes "github.com/petrichormoney/petri/x/collector/types"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	chain "github.com/petrichormoney/petri/app"
	assetKeeper "github.com/petrichormoney/petri/x/asset/keeper"
	auctionKeeper "github.com/petrichormoney/petri/x/auction/keeper"
	collectorKeeper "github.com/petrichormoney/petri/x/collector/keeper"
	tokenmintKeeper "github.com/petrichormoney/petri/x/tokenmint/keeper"
	tokenmintTypes "github.com/petrichormoney/petri/x/tokenmint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app             *chain.App
	ctx             sdk.Context
	assetKeeper     assetKeeper.Keeper
	collectorKeeper collectorKeeper.Keeper
	auctionKeeper   auctionKeeper.Keeper
	tokenmintKeeper tokenmintKeeper.Keeper
	querier         tokenmintKeeper.QueryServer
	msgServer       tokenmintTypes.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chain.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.collectorKeeper = s.app.CollectorKeeper
	s.assetKeeper = s.app.AssetKeeper
	s.auctionKeeper = s.app.AuctionKeeper
	s.tokenmintKeeper = s.app.TokenmintKeeper
	s.querier = tokenmintKeeper.QueryServer{Keeper: s.tokenmintKeeper}
	s.msgServer = tokenmintKeeper.NewMsgServer(s.tokenmintKeeper)
}

// // Below are just shortcuts to frequently-used functions.
//
//	func (s *KeeperTestSuite) getBalances(addr sdk.AccAddress) sdk.Coins {
//		return s.app.bankKeeper.GetAllBalances(s.ctx, addr)
//	}
func (s *KeeperTestSuite) getBalance(addr string, denom string) (coin sdk.Coin, err error) {
	addr1, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return coin, err
	}
	return s.app.BankKeeper.GetBalance(s.ctx, addr1, denom), nil
}

//
//func (s *KeeperTestSuite) sendCoins(fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) {
//	s.T().Helper()
//	err := s.app.bankKeeper.SendCoins(s.ctx, fromAddr, toAddr, amt)
//	s.Require().NoError(err)
//}
//
//func (s *KeeperTestSuite) nextBlock() {
//	liquidity.EndBlocker(s.ctx, s.keeper)
//	liquidity.BeginBlocker(s.ctx, s.keeper)
//}
//
//// Below are useful helpers to write test code easily.
//func (s *KeeperTestSuite) addr(addrNum int) sdk.AccAddress {
//	addr := make(sdk.AccAddress, 20)
//	binary.PutVarint(addr, int64(addrNum))
//	return addr
//}

func (s *KeeperTestSuite) fundAddr(addr string, amt sdk.Coin) {
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, collectorTypes.ModuleName, sdk.NewCoins(amt))
	s.Require().NoError(err)
	addr1, err := sdk.AccAddressFromBech32(addr)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, collectorTypes.ModuleName, addr1, sdk.NewCoins(amt))
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) advanceseconds(dur int64) {
	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Second * time.Duration(dur)))
}

// ParseCoins parses and returns sdk.Coins.
func ParseCoin(s string) sdk.Coin {
	coins, err := sdk.ParseCoinNormalized(s)
	if err != nil {
		panic(err)
	}
	return coins
}
