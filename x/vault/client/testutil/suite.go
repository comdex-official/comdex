package testutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	chain "github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/vault/client/cli"
	vaultKeeper "github.com/comdex-official/comdex/x/vault/keeper"
	"github.com/comdex-official/comdex/x/vault/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	// store "github.com/cosmos/cosmos-sdk/store/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"

	dbm "github.com/cometbft/cometbft-db"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

type VaultIntegrationTestSuite struct {
	suite.Suite

	app       *chain.App
	cfg       network.Config
	network   *network.Network
	val       *network.Validator
	ctx       sdk.Context
	msgServer types.MsgServer
}

func NewAppConstructor() network.AppConstructor {
	return func(val network.ValidatorI) servertypes.Application {
		return chain.New(
			val.GetCtx().Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.GetCtx().Config.RootDir, 0,
			chain.MakeEncodingConfig(), simtestutil.EmptyAppOptions{}, chain.GetWasmEnabledProposals(), chain.EmptyWasmOpts,
			baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
		)
	}
}

func (s *VaultIntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	if testing.Short() {
		s.T().Skip("skipping test in unit-tests mode.")
	}

	s.app = chain.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.msgServer = vaultKeeper.NewMsgServer(s.app.VaultKeeper)

	cfg := network.DefaultConfig()
	cfg.AppConstructor = NewAppConstructor()
	cfg.GenesisState = chain.ModuleBasics.DefaultGenesis(cfg.Codec)
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network, _ = network.New(s.T(), s.T().TempDir(), cfg)

	s.val = s.network.Validators[0]

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.Create()

	fmt.Println("Setting up suit.....")
}

func (s *VaultIntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *VaultIntegrationTestSuite) Create() {
	appID := s.CreateNewApp("appone")
	assetInID := s.CreateNewAsset("ASSETONE", "denom1", 2000000)
	assetOutID := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)
	pairID := s.CreateNewPair(assetInID, assetOutID)
	extendedVaultPairID := s.CreateNewExtendedVaultPair("CMDX-C", appID, pairID)

	_, _ = MsgCreate(s.val.ClientCtx, appID, extendedVaultPairID, sdk.NewInt(3), sdk.NewInt(2), s.val.Address.String())
	// s.Require().NoError(err)

	err := s.network.WaitForNextBlock()
	s.Require().NoError(err)
}

func (s *VaultIntegrationTestSuite) TestQueryPairsCmd() {
	val := s.network.Validators[0]

	for _, tc := range []struct {
		name        string
		args        []string
		expectedErr string
		postRun     func(resp types.QueryAllVaultsResponse)
	}{
		{
			"valid case",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			"",
			func(resp types.QueryAllVaultsResponse) {
				// WIP - vault created but not present in the client context ? IDK how it works.
				// s.Require().Len(resp.Vault, 1)
			},
		},
	} {
		s.Run(tc.name, func() {
			cmd := cli.QueryAllVaults()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				var resp types.QueryAllVaultsResponse
				s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resp), out.String())
				tc.postRun(resp)
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}
