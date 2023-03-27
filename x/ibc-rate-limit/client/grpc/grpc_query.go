package grpc 

// THIS FILE IS GENERATED CODE, DO NOT EDIT
// SOURCE AT `proto/comdex/ibc-rate-limit/v1beta1/query.yml`

import (
	context "context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/comdex-official/comdex/x/ibc-rate-limit/client"
	"github.com/comdex-official/comdex/x/ibc-rate-limit/types"
)

type Querier struct {
	Q client.Querier
}

var _ types.QueryServer = Querier{}

func (q Querier) Params(grpcCtx context.Context,
	req *types.ParamsRequest,
) (*types.ParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(grpcCtx)
	return q.Q.Params(ctx, *req)
}

