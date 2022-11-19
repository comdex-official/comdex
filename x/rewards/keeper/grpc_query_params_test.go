package keeper_test

import (
	"testing"

	"github.com/petrichormoney/petri/app"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/petrichormoney/petri/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	petriApp := app.Setup(false)
	ctx := petriApp.BaseApp.NewContext(false, tmproto.Header{})

	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	petriApp.Rewardskeeper.SetParams(ctx, params)

	response, err := petriApp.Rewardskeeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
