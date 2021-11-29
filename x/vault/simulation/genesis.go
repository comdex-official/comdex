package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Simulation parameter constants
const (
	ParamsDistrVaultIdentifier = "distr_epoch_identifier"
)

// RandomizedGenState generates a random GenesisState for gov
func RandomizedGenState(simState *module.SimulationState) {
	// Parameter for how often rewards get distributed
	var distrEpochIdentifier string
	simState.AppParams.GetOrGenerate(
		simState.Cdc, ParamsDistrVaultIdentifier, &distrEpochIdentifier, simState.Rand,
		func(r *rand.Rand) { distrEpochIdentifier = GenParamsDistrVaultIdentifier(r) },
	)

	vaultGenesis := types.GenesisState{
		Vaults: []types.Vault{
			{
				ID:        1,
				PairID:    1,
				Owner:     "comdex11111111111111",
				AmountIn:  sdk.Int(sdk.NewInt(100)),
				AmountOut: sdk.Int(sdk.NewInt(66)),
			},
		},
	}

	bz, err := json.MarshalIndent(&vaultGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated incentives parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&vaultGenesis)
}
