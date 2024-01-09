package testutil

import (
	"fmt"
	"github.com/comdex-official/comdex/app/params"
	"testing"

	"github.com/stretchr/testify/suite"

	chain "github.com/comdex-official/comdex/app"
	liquidityKeeper "github.com/comdex-official/comdex/x/liquidity/keeper"
	"github.com/comdex-official/comdex/x/liquidity/types"

	"github.com/comdex-official/comdex/testutil/network"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	networkI "github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"

	dbm "github.com/cometbft/cometbft-db"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
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

func NewAppConstructor() networkI.AppConstructor {
	return func(val networkI.ValidatorI) servertypes.Application {
		return chain.New(
			val.GetCtx().Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.GetCtx().Config.RootDir, 0,
			params.MakeEncodingConfig(), simtestutil.EmptyAppOptions{}, chain.GetWasmEnabledProposals(), chain.EmptyWasmOpts,
			baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
		)
	}
}

func (s *LiquidityIntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	if testing.Short() {
		s.T().Skip("skipping test in unit-tests mode.")
	}

	s.app = chain.Setup(s.T(), false)
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
