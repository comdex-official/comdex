package keeper

import (
	"context"
	"fmt"
	"github.com/comdex-official/comdex/x/esm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ types.QueryServer = (*queryServer)(nil)
)

type queryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServer{
		Keeper: k,
	}
}

func (q queryServer) QueryESMTriggerParams(c context.Context, req *types.QueryESMTriggerParamsRequest) (*types.QueryESMTriggerParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetESMTriggerParams(ctx, req.Id)
	fmt.Println("item.....", item)
	if !found {
		return &types.QueryESMTriggerParamsResponse{}, nil
	}

	return &types.QueryESMTriggerParamsResponse{
		EsmTriggerParams: item,
	}, nil
}
