package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/poolapi/types"
)

type queryServer struct {
	Keeper
}

func NewQueryServiceServer(k Keeper) types.QueryServiceServer {
	return &queryServer{
		Keeper: k,
	}
}

func (q *queryServer) IndividualPoolLiquidity(c context.Context, req *types.QueryIndividualPoolLiquidityRequest) (*types.QueryIndividualPoolLiquidityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	liquidity, found := q.PoolLiquidity(ctx, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pool does not exist for id %d", req.PoolId)
	}

	return &types.QueryIndividualPoolLiquidityResponse{
		PoolLiquidity: liquidity,
	}, nil
}

func (q *queryServer) PoolsLiquidity(c context.Context, req *types.QueryTotalLiquidityRequest) (*types.QueryTotalLiquidityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	total_liquidity, found := q.TotalLiquidity(ctx)
	if !found {
		return nil, status.Errorf(codes.NotFound, "one or multiple pools or denoms does not exist ")
	}

	return &types.QueryTotalLiquidityResponse{
		TotalLiquidity: total_liquidity,
	}, nil
}

func (q *queryServer) TotalCollateral(c context.Context, req *types.QueryTotalCollateralRequest) (*types.QueryTotalCollateralResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	total_collateral, _ := q.GetTotalCollateral(c)

	return &types.QueryTotalCollateralResponse{
		TotalCollateral: total_collateral,
	}, nil
}

func (q *queryServer) PoolAPR(c context.Context, req *types.QueryPoolAPRRequest) (*types.QueryPoolAPRResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	pool_apr, _ := q.GetAPR(c)

	return &types.QueryPoolAPRResponse{
		Apr: uint64(pool_apr),
	}, nil
}
