package collector_test

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"

	app "github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/collector"
	"github.com/comdex-official/comdex/x/collector/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {

	comdexApp := app.Setup(false)
	ctx := comdexApp.BaseApp.NewContext(false, tmproto.Header{})
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k := comdexApp.CollectorKeeper
	collector.InitGenesis(ctx, k, &genesisState)
	got := collector.ExportGenesis(ctx, k)
	require.NotNil(t, got)

}
