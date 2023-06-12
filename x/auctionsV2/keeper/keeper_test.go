package keeper_test

import (
	chain "github.com/comdex-official/comdex/app"
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/auctionsV2/keeper"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	collectKeeper "github.com/comdex-official/comdex/x/collector/keeper"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	liquidationKeeper "github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	liquidationTypes "github.com/comdex-official/comdex/x/liquidationsV2/types"
	marketKeeper "github.com/comdex-official/comdex/x/market/keeper"
	tokenmintKeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
	vaultKeeper "github.com/comdex-official/comdex/x/vault/keeper"
	vaultTypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

type KeeperTestSuite struct {
	suite.Suite

	app                  *chain.App
	ctx                  sdk.Context
	vaultKeeper          vaultKeeper.Keeper
	assetKeeper          assetKeeper.Keeper
	liquidationKeeper    liquidationKeeper.Keeper
	tokenmintKeeper      tokenmintKeeper.Keeper
	marketKeeper         marketKeeper.Keeper
	collectorKeeper      collectKeeper.Keeper
	liquidationQuerier   liquidationKeeper.QueryServer
	vaultQuerier         vaultKeeper.QueryServer
	liquidationMsgServer liquidationTypes.MsgServer
	vaultMsgServer       vaultTypes.MsgServer
	keeper               keeper.Keeper
	auctionMsgServer     auctionsV2types.MsgServer
	lendKeeper           lendkeeper.Keeper
	lendQuerier          lendkeeper.QueryServer
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
	s.collectorKeeper = s.app.CollectorKeeper
	s.liquidationQuerier = liquidationKeeper.QueryServer{Keeper: s.liquidationKeeper}
	s.liquidationMsgServer = liquidationKeeper.NewMsgServerImpl(s.liquidationKeeper)
	s.vaultMsgServer = vaultKeeper.NewMsgServer(s.vaultKeeper)
	s.vaultQuerier = vaultKeeper.QueryServer{Keeper: s.vaultKeeper}
	s.marketKeeper = s.app.MarketKeeper
	s.keeper = s.app.NewaucKeeper
	s.auctionMsgServer = keeper.NewMsgServerImpl(s.keeper)
	s.tokenmintKeeper = s.app.TokenmintKeeper
	s.lendKeeper = s.app.LendKeeper
	s.lendQuerier = lendkeeper.QueryServer{Keeper: s.lendKeeper}
}

func (s *KeeperTestSuite) getBalance(addr string, denom string) (coin sdk.Coin, err error) {
	addr1, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return coin, err
	}
	return s.app.BankKeeper.GetBalance(s.ctx, addr1, denom), nil
}

func (s *KeeperTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coin) {
	amt1 := sdk.NewCoins(amt)
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, liquidationTypes.ModuleName, amt1)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, liquidationTypes.ModuleName, addr, amt1)
	s.Require().NoError(err)
}
