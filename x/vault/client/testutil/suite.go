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
	"github.com/cosmos/cosmos-sdk/simapp"
	store "github.com/cosmos/cosmos-sdk/store/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tmcli "github.com/tendermint/tendermint/libs/cli"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
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
	return func(val network.Validator) servertypes.Application {
		return chain.New(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
			chain.MakeEncodingConfig(), simapp.EmptyAppOptions{}, chain.GetWasmEnabledProposals(), chain.EmptyWasmOpts,
			baseapp.SetPruning(store.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
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
	s.network = network.New(s.T(), cfg)

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
	appID := s.CreateNewApp("appOne")
	fmt.Println("app created....", appID)

	assetInID := s.CreateNewAsset("asset1", "denom1", 2000000)
	fmt.Println("asset1 created....", assetInID)

	assetOutID := s.CreateNewAsset("asset2", "denom2", 1000000)
	fmt.Println("asset2 created....", assetOutID)

	pairID := s.CreateNewPair(assetInID, assetOutID)
	fmt.Println("pair created....", pairID)

	extendedVaultPairID := s.CreateNewExtendedVaultPair("CMDX C", appID, pairID)
	fmt.Println("extendedVaultPair created....", extendedVaultPairID)

	// msg := types.MsgCreateRequest{
	// 	From:                s.val.Address.String(),
	// 	AppMappingId:        appID,
	// 	ExtendedPairVaultId: extendedVaultPairID,
	// 	AmountIn:            sdk.NewInt(300000000),
	// 	AmountOut:           sdk.NewInt(200000000),
	// }
	// s.fundAddr(s.val.Address, sdk.NewCoins(sdk.NewCoin("denom1", msg.AmountIn), sdk.NewCoin("denom2", msg.AmountOut)))

	// _, err := s.msgServer.MsgCreate(sdk.WrapSDKContext(s.ctx), &msg)
	// s.Require().NoError(err)

	_, err := MsgCreate(s.val.ClientCtx, appID, extendedVaultPairID, sdk.NewInt(3), sdk.NewInt(2), s.val.Address.String())
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	fmt.Println("all vaults.....", s.app.VaultKeeper.GetVaults(s.ctx))

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
			"happy case",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			"",
			func(resp types.QueryAllVaultsResponse) {
				// WIP - vault created but not present in the client context ? IDK how it works.
				fmt.Println("Response....", resp)
				s.Require().Len(resp.Vault, 1)
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
