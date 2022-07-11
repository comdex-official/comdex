package testutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	chain "github.com/comdex-official/comdex/app"
	liquidityKeeper "github.com/comdex-official/comdex/x/liquidity/keeper"
	"github.com/comdex-official/comdex/x/liquidity/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

type LiquidityIntegrationTestSuite struct {
	suite.Suite

	app       *chain.App
	cfg       network.Config
	network   *network.Network
	val       *network.Validator
	ctx       sdk.Context
	msgServer types.MsgServer
}

func NewAppConstructor() network.AppConstructor {
	return func(val network.Validator) servertypes.Application {
		return chain.New(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
			chain.MakeEncodingConfig(), simapp.EmptyAppOptions{}, chain.GetWasmEnabledProposals(), chain.EmptyWasmOpts,
			baseapp.SetPruning(store.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)
	}
}

func (s *LiquidityIntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	if testing.Short() {
		s.T().Skip("skipping test in unit-tests mode.")
	}

	s.app = chain.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.msgServer = liquidityKeeper.NewMsgServerImpl(s.app.LiquidityKeeper)

	cfg := network.DefaultConfig()
	cfg.AppConstructor = NewAppConstructor()
	cfg.GenesisState = chain.ModuleBasics.DefaultGenesis(cfg.Codec)
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	s.val = s.network.Validators[0]

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.Create()

	fmt.Println("Setting up suit.....")
}

func (s *LiquidityIntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *LiquidityIntegrationTestSuite) Create() {
	err := s.network.WaitForNextBlock()
	s.Require().NoError(err)
}

func (s *LiquidityIntegrationTestSuite) TestQueryPairsCmd() {
	// val := s.network.Validators[0]
}
