package simulation

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Simulation parameter constants
const (
	ParamsDistrVaultIdentifier = "distr_epoch_identifier"
)

// RandomizedGenState generates a random GenesisState for gov
func RandomizedGenState(simState *module.SimulationState) {

	// vaultGenesis := types.GenesisState{
	// 	// Vaults: []types.Vault{
	// 	// 	{
	// 	// 		ID:        1,
	// 	// 		PairID:    1,
	// 	// 		Owner:     "comdex1wpd3scc4ulcw3nq3s5q79escym6cccgzm76rxq",
	// 	// 		AmountIn:  sdk.Int(sdk.NewInt(100)),
	// 	// 		AmountOut: sdk.Int(sdk.NewInt(66)),
	// 	// 	},
	// 	// },
	// }

	// bz, err := json.MarshalIndent(&vaultGenesis, "", " ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Selected randomly generated vault parameters:\n%s\n", bz)
	// simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&vaultGenesis)
}
