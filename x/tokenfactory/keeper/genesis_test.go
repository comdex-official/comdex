package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/comdex-official/comdex/x/tokenfactory/types"
)

func (suite *KeeperTestSuite) TestGenesis() {
	genesisState := types.GenesisState{
		FactoryDenoms: []types.GenesisDenom{
			{
				Denom: "factory/cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t/bitcoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				},
			},
			{
				Denom: "factory/cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t/diff-admin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				},
			},
			{
				Denom: "factory/cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t/litecoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				},
			},
		},
	}

	suite.SetupTestForInitGenesis()
	app := suite.App

	// Test both with bank denom metadata set, and not set.
	for i, denom := range genesisState.FactoryDenoms {
		// hacky, sets bank metadata to exist if i != 0, to cover both cases.
		if i != 0 {
			app.BankKeeper.SetDenomMetaData(suite.Ctx, banktypes.Metadata{Base: denom.GetDenom()})
		}
	}

	if err := app.TokenFactoryKeeper.SetParams(suite.Ctx, types.Params{DenomCreationFee: sdk.Coins{sdk.NewInt64Coin("stake", 100)}}); err != nil {
		panic(err)
	}
	app.TokenFactoryKeeper.InitGenesis(suite.Ctx, genesisState)

	// check that the module account is now initialized
	tokenfactoryModuleAccount := app.AccountKeeper.GetAccount(suite.Ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	suite.Require().NotNil(tokenfactoryModuleAccount)

	exportedGenesis := app.TokenFactoryKeeper.ExportGenesis(suite.Ctx)
	suite.Require().NotNil(exportedGenesis)
	suite.Require().Equal(genesisState, *exportedGenesis)
}
