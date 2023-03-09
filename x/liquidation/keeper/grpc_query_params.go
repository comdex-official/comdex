package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/liquidation/types"
)

func (k Keeper) Params(c context.Context, req *types.QueryLiquidationParamsRequest) (*types.QueryLiquidationParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryLiquidationParamsResponse{Params: k.GetParams(ctx)}, nil
}
