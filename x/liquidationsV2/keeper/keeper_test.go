package keeper_test

import (
	chain "github.com/comdex-official/comdex/app"
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	liquidationKeeper "github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	marketKeeper "github.com/comdex-official/comdex/x/market/keeper"
	vaultKeeper "github.com/comdex-official/comdex/x/vault/keeper"
	vaultTypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

type KeeperTestSuite struct {
	suite.Suite

	app               *chain.App
	ctx               sdk.Context
	vaultKeeper       vaultKeeper.Keeper
	assetKeeper       assetKeeper.Keeper
	liquidationKeeper liquidationKeeper.Keeper
	marketKeeper      marketKeeper.Keeper
	querier           liquidationKeeper.QueryServer
	vaultQuerier      vaultKeeper.QueryServer
	msgServer         types.MsgServer
	vaultMsgServer    vaultTypes.MsgServer
	lendKeeper        lendkeeper.Keeper
	lendQuerier       lendkeeper.QueryServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chain.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.vaultKeeper = s.app.VaultKeeper
	s.liquidationKeeper = s.app.NewliqKeeper
	s.assetKeeper = s.app.AssetKeeper
	s.querier = liquidationKeeper.QueryServer{Keeper: s.liquidationKeeper}
	s.msgServer = liquidationKeeper.NewMsgServerImpl(s.liquidationKeeper)
	s.vaultMsgServer = vaultKeeper.NewMsgServer(s.vaultKeeper)
	s.vaultQuerier = vaultKeeper.QueryServer{Keeper: s.vaultKeeper}
	s.marketKeeper = s.app.MarketKeeper
	s.lendKeeper = s.app.LendKeeper
	s.lendQuerier = lendkeeper.QueryServer{Keeper: s.lendKeeper}
}
