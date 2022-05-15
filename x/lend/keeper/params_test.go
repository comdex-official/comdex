package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/comdex-official/comdex/testutil/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.LendKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
