package tokenmint_test

import (
	"testing"

	"github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/tokenmint"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/comdex-official/comdex/x/tokenmint/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	comdexApp := app.Setup(false)
	ctx := comdexApp.BaseApp.NewContext(false, tmproto.Header{})

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	tokenmint.InitGenesis(ctx, comdexApp.TokenmintKeeper, &genesisState)
	got := tokenmint.ExportGenesis(ctx, comdexApp.TokenmintKeeper)
	require.NotNil(t, got)
}
