package keeper_test

import (
	collectorKeeper "github.com/comdex-official/comdex/x/collector/keeper"
	rewardsKeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

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
	querier       rewardsKeeper.QueryServer
	msgServer     rewardstypes.MsgServer
	collector     collectorKeeper.Keeper
	rewardsKeeper rewardsKeeper.Keeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chain.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.rewardsKeeper = s.app.Rewardskeeper
	s.assetKeeper = s.app.AssetKeeper
	s.querier = rewardsKeeper.QueryServer{Keeper: s.rewardsKeeper}
	s.msgServer = rewardsKeeper.NewMsgServerImpl(s.rewardsKeeper)
	s.collector = s.app.CollectorKeeper
	s.rewardsKeeper = s.app.Rewardskeeper
}

func (s *KeeperTestSuite) fundAddr(addr string, amt sdk.Coin) {
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, lockerTypes.ModuleName, sdk.NewCoins(amt))
	s.Require().NoError(err)
	addr1, err := sdk.AccAddressFromBech32(addr)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, lockerTypes.ModuleName, addr1, sdk.NewCoins(amt))
	s.Require().NoError(err)
}
