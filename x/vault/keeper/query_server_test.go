package keeper_test

import (
	"github.com/comdex-official/comdex/app"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/vault/keeper"
	types2 "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	"math/big"
	"testing"
)
type QueryTestSuite struct {
	suite.Suite
	keeper keeper.Keeper
	assetKeeper assetkeeper.Keeper
	app    app.TestApp
	ctx    sdk.Context
	prq    query.PageRequest
}

func (suite *QueryTestSuite) SetupTest() {
	testApp := app.NewTestApp()
	k := testApp.GetVaultKeeper()
	ak := testApp.GetAssetKeeper()
	ctx := testApp.NewContext(true, tmproto.Header{Height: 1, Time: tmtime.Now()})
	suite.app = testApp
	suite.keeper = k
	suite.assetKeeper = ak
	suite.ctx = ctx

	return
}
func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}

func (suite *QueryTestSuite)TestQueryVaults() {
	queryserver := keeper.NewQueryServiceServer(suite.keeper)
	qvr := types2.QueryVaultsRequest{"abc",&suite.prq }
	_,err := queryserver.QueryVaults(sdk.WrapSDKContext(suite.ctx), nil )
	suite.Error(err)
	_,err = queryserver.QueryVaults(sdk.WrapSDKContext(suite.ctx), &qvr )
	suite.NoError(err)
	vault := types2.Vault{ID: 1, PairID: 1, Owner: "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(100)}
	suite.keeper.SetVault(suite.ctx, vault)
	_,err = queryserver.QueryVaults(sdk.WrapSDKContext(suite.ctx), &qvr )
	suite.Error(err)
	pair := assettypes.Pair{
		Id:               1,
		AssetIn:          1500000,
		AssetOut:         1000000,
		LiquidationRatio: sdk.NewDecFromBigInt(big.NewInt(150)),
	}
	suite.assetKeeper.SetPair(suite.ctx, pair)
	_,err = queryserver.QueryVaults(sdk.WrapSDKContext(suite.ctx), &qvr )
	suite.Error(err)

}

func (suite *QueryTestSuite)TestQueryVault() {
	queryserver := keeper.NewQueryServiceServer(suite.keeper)
	qvr := types2.QueryVaultRequest{1 }

	_,err := queryserver.QueryVault(sdk.WrapSDKContext(suite.ctx), nil )
	suite.Error(err)
	_,err = queryserver.QueryVault(sdk.WrapSDKContext(suite.ctx), &qvr )
	suite.Error(err)
	vault := types2.Vault{ID: types.DefaultIndex, PairID: 1, Owner: "abc", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(100)}
	suite.keeper.SetVault(suite.ctx, vault)
	_,err = queryserver.QueryVault(sdk.WrapSDKContext(suite.ctx), &qvr )
	suite.Error(err)
}