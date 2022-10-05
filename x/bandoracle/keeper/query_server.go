package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/bandoracle/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) FetchPriceResult(c context.Context, req *types.QueryFetchPriceRequest) (*types.QueryFetchPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	result, err := k.GetFetchPriceResult(ctx, types.OracleRequestID(req.RequestId))
	if err != nil {
		return nil, err
	}
	return &types.QueryFetchPriceResponse{Result: &result}, nil
}

func (k Keeper) LastFetchPriceID(c context.Context, req *types.QueryLastFetchPriceIdRequest) (*types.QueryLastFetchPriceIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	id := k.GetLastFetchPriceID(ctx)
	return &types.QueryLastFetchPriceIdResponse{RequestId: id}, nil
}
