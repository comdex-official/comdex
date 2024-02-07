package keeper_test

import (
	"testing"

	"github.com/comdex-official/comdex/app"

	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	comdexApp := app.Setup(t, false)
	ctx := comdexApp.BaseApp.NewContext(false)

	params := types.DefaultParams()

	comdexApp.Rewardskeeper.SetParams(ctx, params)

	require.EqualValues(t, params, comdexApp.Rewardskeeper.GetParams(ctx))
}
