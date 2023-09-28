package keeper_test

import (
	"testing"
	"time"

	rewardsKeeper "github.com/comdex-official/comdex/x/rewards/keeper"

	collectorKeeper "github.com/comdex-official/comdex/x/collector/keeper"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/comdex-official/comdex/app"
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
	lockerKeeper "github.com/comdex-official/comdex/x/locker/keeper"
	lockerTypes "github.com/comdex-official/comdex/x/locker/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app           *chain.App
	ctx           sdk.Context
	assetKeeper   assetKeeper.Keeper
	lockerKeeper  lockerKeeper.Keeper
	querier       lockerKeeper.QueryServer
	msgServer     lockerTypes.MsgServer
	collector     collectorKeeper.Keeper
	rewardsKeeper rewardsKeeper.Keeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chain.Setup(s.T(), false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.lockerKeeper = s.app.LockerKeeper
	s.assetKeeper = s.app.AssetKeeper
	s.querier = lockerKeeper.QueryServer{Keeper: s.lockerKeeper}
	s.msgServer = lockerKeeper.NewMsgServer(s.lockerKeeper)
	s.collector = s.app.CollectorKeeper
	s.rewardsKeeper = s.app.Rewardskeeper
}

//
//// Below are just shortcuts to frequently-used functions.
//func (s *KeeperTestSuite) getBalances(addr sdk.AccAddress) sdk.Coins {
//	return s.app.bankKeeper.GetAllBalances(s.ctx, addr)
//}
//
//func (s *KeeperTestSuite) getBalance(addr sdk.AccAddress, denom string) sdk.Coin {
//	return s.app.bankKeeper.GetBalance(s.ctx, addr, denom)
//}
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
	err := s.app.BankKeeper.MintCoins(s.ctx, lockerTypes.ModuleName, sdk.NewCoins(amt))
	s.Require().NoError(err)
	addr1, err := sdk.AccAddressFromBech32(addr)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, lockerTypes.ModuleName, addr1, sdk.NewCoins(amt))
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
