package rewards_test

import (
	"testing"

	"github.com/comdex-official/comdex/app"

	"github.com/comdex-official/comdex/x/rewards"
	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	comdexApp := app.Setup(t, false)
	ctx := comdexApp.BaseApp.NewContext(false)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	rewards.InitGenesis(ctx, comdexApp.Rewardskeeper, &genesisState)
	got := rewards.ExportGenesis(ctx, comdexApp.Rewardskeeper)
	require.NotNil(t, got)
}
