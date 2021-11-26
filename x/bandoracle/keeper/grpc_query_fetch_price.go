package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/bandoracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FetchPriceResult returns the FetchPrice result by RequestId
func (k Keeper) FetchPriceResult(c context.Context, req *types.QueryFetchPriceRequest) (*types.QueryFetchPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	result, err := k.GetFetchPriceResult(ctx, types.OracleRequestID(req.RequestId))
	if err != nil {
		return nil, err
	}
	return &types.QueryFetchPriceResponse{Result: &result}, nil
}

// LastFetchPriceId returns the last FetchPrice request Id
func (k Keeper) LastFetchPriceId(c context.Context, req *types.QueryLastFetchPriceIdRequest) (*types.QueryLastFetchPriceIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	id := k.GetLastFetchPriceID(ctx)
	return &types.QueryLastFetchPriceIdResponse{RequestId: id}, nil
}
