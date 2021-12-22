package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/comdex-official/comdex/x/asset/types"
)

// Simulation parameter constants
const (
	AssetAdmin = "asset_admin"
)

func GenAssetAdmin(r *rand.Rand) string {
	types.SetAccountAddressPrefixes()
	assetadmin := "comdex1pwu5sjk2lje94cwhnh0dhr0fery2c75w4y09t4"
	return assetadmin
	//replace with random account generating func
}

// RandomizedGenState generates a random GenesisState for gov
func RandomizedGenState(simState *module.SimulationState) {
	// Parameter for how often rewards get distributed
	var AssetAdmin string
	simState.AppParams.GetOrGenerate(
		simState.Cdc, AssetAdmin, &AssetAdmin, simState.Rand,
		func(r *rand.Rand) { AssetAdmin = GenAssetAdmin(r) },
	)

	assetGenesis := types.GenesisState{
		Params: types.Params{
			Admin: AssetAdmin,
		},
	}

	bz, err := json.MarshalIndent(&assetGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated asset parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&assetGenesis)
}
