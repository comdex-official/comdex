package rewards_test

import (
	"testing"

	"github.com/petrichormoney/petri/app"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/petrichormoney/petri/x/rewards"
	"github.com/petrichormoney/petri/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	petriApp := app.Setup(false)
	ctx := petriApp.BaseApp.NewContext(false, tmproto.Header{})

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	rewards.InitGenesis(ctx, petriApp.Rewardskeeper, &genesisState)
	got := rewards.ExportGenesis(ctx, petriApp.Rewardskeeper)
	require.NotNil(t, got)
}
