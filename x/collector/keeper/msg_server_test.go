package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/comdex-official/comdex/testutil/keeper"
	"github.com/comdex-official/comdex/x/collector/keeper"
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.CollectorKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
