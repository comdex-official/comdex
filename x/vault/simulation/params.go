package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

const (
	keyDistrVaultIdentifier = "DistrVaultIdentifier"
)

var (
	// TODO: remove hardcoded params
	// refer x/epochs/simulation/genesis.go
	vaultIdentifiers = []string{"day", "hour"}
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyDistrVaultIdentifier,
			func(r *rand.Rand) string {
				return fmt.Sprintf(`{"%s: %s"}`, ParamsDistrVaultIdentifier, GenParamsDistrVaultIdentifier(r))
			},
		),
	}
}

func GenParamsDistrVaultIdentifier(r *rand.Rand) string {
	return vaultIdentifiers[r.Intn(len(vaultIdentifiers))]
}
