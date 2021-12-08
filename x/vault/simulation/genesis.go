package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"

	"github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Simulation parameter constants
const (
	ParamsDistrVaultIdentifier = "distr_vault_identifier"
)

// RandomizedGenState generates a random GenesisState for gov
func RandomizedGenState(simState *module.SimulationState) {

	vaultGenesis := types.GenesisState{
		Vaults: []types.Vault{
			{
				ID:        1,
				PairID:    1,
				Owner:     "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v",
				AmountIn:  sdk.Int(sdk.NewInt(100000)),
				AmountOut: sdk.Int(sdk.NewInt(66666)),
			},
		},
	}

	bz, err := json.MarshalIndent(&vaultGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated vault parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&vaultGenesis)
}
