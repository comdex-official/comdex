package bandoracle_test

import (
	"testing"

	keepertest "github.com/comdex-official/comdex/testutil/keeper"
	"github.com/comdex-official/comdex/x/bandoracle"
	"github.com/comdex-official/comdex/x/bandoracle/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BandoracleKeeper(t)
	bandoracle.InitGenesis(ctx, *k, genesisState)
	got := bandoracle.ExportGenesis(ctx, *k)
	require.NotNil(t, got)


	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
