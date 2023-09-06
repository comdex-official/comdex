package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/comdex-official/comdex/x/common/types"
    "github.com/comdex-official/comdex/x/common/keeper"
    keepertest "github.com/comdex-official/comdex/testutil/keeper"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.CommonKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
