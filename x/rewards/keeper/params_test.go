package keeper_test

import (
	"testing"

	"github.com/petrichormoney/petri/app"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/petrichormoney/petri/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	petriApp := app.Setup(false)
	ctx := petriApp.BaseApp.NewContext(false, tmproto.Header{})

	params := types.DefaultParams()

	petriApp.Rewardskeeper.SetParams(ctx, params)

	require.EqualValues(t, params, petriApp.Rewardskeeper.GetParams(ctx))
}
