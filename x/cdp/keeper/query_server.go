package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/cdp/types"
)

var (
	_ types.QueryServiceServer = (*queryServer)(nil)
)

type queryServer struct {
	Keeper
}

func NewQueryServiceServer(k Keeper) types.QueryServiceServer {
	return &queryServer{
		Keeper: k,
	}
}

func (q *queryServer) QueryCDPs(c context.Context, req *types.QueryCDPsRequest) (*types.QueryCDPsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.CDPInfo
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.CDPKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.CDP
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			pair, found := q.GetPair(ctx, item.PairID)
			if !found {
				return false, status.Errorf(codes.NotFound, "pair does not exist for id %d", item.PairID)
			}

			assetIn, found := q.GetAsset(ctx, pair.AssetIn)
			if !found {
				return false, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetIn)
			}

			assetOut, found := q.GetAsset(ctx, pair.AssetOut)
			if !found {
				return false, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetOut)
			}

			collateralizationRatio, err := q.CalculateCollaterlizationRatio(ctx, item.AmountIn, assetIn, item.AmountOut, assetOut)
			if err != nil {
				return false, err
			}

			cdpInfo := types.CDPInfo{
				Id:                    item.ID,
				PairID:                item.PairID,
				Owner:                 item.Owner,
				Collateral:            sdk.NewCoin(assetIn.Denom, item.AmountIn),
				Debt:                  sdk.NewCoin(assetOut.Denom, item.AmountOut),
				CollaterlizationRatio: collateralizationRatio,
			}

			if accumulate {
				items = append(items, cdpInfo)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCDPsResponse{
		CDPsInfo:   items,
		Pagination: pagination,
	}, nil
}

func (q *queryServer) QueryCDP(c context.Context, req *types.QueryCDPRequest) (*types.QueryCDPResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	cdp, found := q.GetCDP(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "cdp does not exist for id %d", req.Id)
	}

	pair, found := q.GetPair(ctx, cdp.PairID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pair does not exist for id %d", cdp.PairID)
	}

	assetIn, found := q.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetIn)
	}

	assetOut, found := q.GetAsset(ctx, pair.AssetOut)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetOut)
	}

	collateralizationRatio, err := q.CalculateCollaterlizationRatio(ctx, cdp.AmountIn, assetIn, cdp.AmountOut, assetOut)
	if err != nil {
		return nil, err
	}
	return &types.QueryCDPResponse{
		CDPInfo: types.CDPInfo{
			Id:                    cdp.ID,
			PairID:                cdp.PairID,
			Owner:                 cdp.Owner,
			Collateral:            sdk.NewCoin(assetIn.Denom, cdp.AmountIn),
			Debt:                  sdk.NewCoin(assetOut.Denom, cdp.AmountOut),
			CollaterlizationRatio: collateralizationRatio,
		},
	}, nil
}
