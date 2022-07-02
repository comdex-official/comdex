package lend_test

import (
	"github.com/comdex-official/comdex/app"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"

	"github.com/comdex-official/comdex/x/lend"
	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {

	comdexApp := app.Setup(false)
	ctx := comdexApp.BaseApp.NewContext(false, tmproto.Header{})

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	lend.InitGenesis(ctx, comdexApp.LendKeeper, genesisState)
	got := lend.ExportGenesis(ctx, comdexApp.LendKeeper)
	require.NotNil(t, got)
	// this line is used by starport scaffolding # genesis/test/assert
}
