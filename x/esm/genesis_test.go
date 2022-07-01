package esm_test

import (
	"testing"

	keepertest "github.com/comdex-official/comdex/testutil/keeper"
	"github.com/comdex-official/comdex/testutil/nullify"
	"github.com/comdex-official/comdex/x/esm"
	"github.com/comdex-official/comdex/x/esm/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:	types.DefaultParams(),
		
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.EsmKeeper(t)
	esm.InitGenesis(ctx, *k, genesisState)
	got := esm.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	

	// this line is used by starport scaffolding # genesis/test/assert
}
