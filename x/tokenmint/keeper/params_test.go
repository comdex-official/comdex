package keeper_test

import (
	"testing"

	testkeeper "github.com/comdex-official/comdex/testutil/keeper"
	"github.com/comdex-official/comdex/x/tokenmint/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.TokenmintKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
