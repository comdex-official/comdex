package tokenmint_test

import (
	"testing"

	"github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/tokenmint"

	"github.com/comdex-official/comdex/x/tokenmint/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	comdexApp := app.Setup(t, false)
	ctx := comdexApp.BaseApp.NewContext(false)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	tokenmint.InitGenesis(ctx, comdexApp.TokenmintKeeper, &genesisState)
	got := tokenmint.ExportGenesis(ctx, comdexApp.TokenmintKeeper)
	require.NotNil(t, got)
}
