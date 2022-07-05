package rewards_test

import (
	"github.com/comdex-official/comdex/app"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"

	"github.com/comdex-official/comdex/x/rewards"
	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {

	comdexApp := app.Setup(false)
	ctx := comdexApp.BaseApp.NewContext(false, tmproto.Header{})

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	rewards.InitGenesis(ctx, comdexApp.Rewardskeeper, genesisState)
	got := rewards.ExportGenesis(ctx, comdexApp.Rewardskeeper)
	require.NotNil(t, got)

	// this line is used by starport scaffolding # genesis/test/assert
}
