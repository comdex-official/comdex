package lend_test

import (
	"testing"

	keepertest "github.com/comdex-official/comdex/testutil/keeper"
	"github.com/comdex-official/comdex/x/lend"
	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:	types.DefaultParams(),
		
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.LendKeeper(t)
	lend.InitGenesis(ctx, *k, genesisState)
	got := lend.ExportGenesis(ctx, *k)
	require.NotNil(t, got)
	// this line is used by starport scaffolding # genesis/test/assert
}
