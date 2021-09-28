package keeper_test

import (
	"github.com/comdex-official/comdex/app"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	types1 "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/vault/keeper"
	types2 "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	"math/big"
	"testing"
)

type VaultTestSuite struct {
	suite.Suite
	keeper keeper.Keeper
	assetKeeper assetkeeper.Keeper
	app    app.TestApp
	ctx    sdk.Context
}

func (suite *VaultTestSuite) SetupTest() {
	testApp := app.NewTestApp()
	k := testApp.GetVaultKeeper()
	ctx := testApp.NewContext(true, tmproto.Header{Height: 1, Time: tmtime.Now()})
	suite.app = testApp
	suite.keeper = k
	suite.ctx = ctx

	return
}
func TestVaultTestSuite(t *testing.T) {
	suite.Run(t, new(VaultTestSuite))
}

func (suite *VaultTestSuite) TestVault_SetGet() {
	vault := types2.Vault{ID: types.DefaultIndex, PairID: 1, Owner: "abc", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(100)}
	suite.keeper.SetVault(suite.ctx, vault)

	t, found := suite.keeper.GetVault(suite.ctx, types.DefaultIndex)
	suite.True(found)
	suite.Equal(vault, t)
	_, found = suite.keeper.GetVault(suite.ctx, 100)
	suite.False(found)
	suite.keeper.DeleteVault(suite.ctx, 100)
	_, found = suite.keeper.GetVault(suite.ctx, 100)
	suite.False(found)
}

func (suite *VaultTestSuite) TestID_SetGet() {
	suite.keeper.SetID(suite.ctx, types.DefaultIndex)
	id := suite.keeper.GetID(suite.ctx)
	suite.Equal(types.DefaultIndex, id)
}

func (suite *VaultTestSuite) TestGetVaults() {
	vaults := suite.keeper.GetVaults(suite.ctx)
	if vaults == nil{
		print(suite)
	}
}

func (suite *VaultTestSuite) TestSetVaultForAddressByPair(){
	suite.keeper.SetVaultForAddressByPair(suite.ctx,sdk.AccAddress("comdex1yples84d8avjl"), 1,1)
	suite.keeper.HasVaultForAddressByPair(suite.ctx,sdk.AccAddress("comdex1yples84d8avjl"), 1)
	suite.keeper.DeleteVaultForAddressByPair(suite.ctx,sdk.AccAddress("comdex1yples84d8avjl"), 1)
}

func (suite *VaultTestSuite) TestVerifyCollaterlizationRatio(){
	assetin := types1.Asset{
		Id:       1,
		Name:     "comdex",
		Denom:    "cmdx",
		Decimals: 100,
	}
	err := suite.keeper.VerifyCollaterlizationRatio(suite.ctx, sdk.NewInt(100),assetin, sdk.NewInt(100), assetin, sdk.NewDecFromBigInt(big.NewInt(150)))
	suite.Error(err)
}
