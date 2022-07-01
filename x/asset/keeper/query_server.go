package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/asset/types"
)

var (
	_ types.QueryServer = (*queryServer)(nil)
)

type queryServer struct {
	Keeper
}

func NewQueryServiceServer(k Keeper) types.QueryServer {
	return &queryServer{
		Keeper: k,
	}
}

func (q *queryServer) QueryAssets(c context.Context, req *types.QueryAssetsRequest) (*types.QueryAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.Asset
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.AssetKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Asset
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate {
				items = append(items, item)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAssetsResponse{
		Assets:     items,
		Pagination: pagination,
	}, nil
}

func (q *queryServer) QueryAsset(c context.Context, req *types.QueryAssetRequest) (*types.QueryAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetAsset(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", req.Id)
	}

	return &types.QueryAssetResponse{
		Asset: item,
	}, nil
}

func (q *queryServer) QueryPairs(c context.Context, req *types.QueryPairsRequest) (*types.QueryPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		pairsInfo []types.PairInfo
		ctx       = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.PairKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var pair types.Pair
			if err := q.cdc.Unmarshal(value, &pair); err != nil {
				return false, err
			}

			assetIn, found := q.GetAsset(ctx, pair.AssetIn)
			if !found {
				return false, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetIn)
			}

			assetOut, found := q.GetAsset(ctx, pair.AssetOut)
			if !found {
				return false, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetOut)
			}

			pairInfo := types.PairInfo{
				Id:       pair.Id,
				AssetIn:  pair.AssetIn,
				DenomIn:  assetIn.Denom,
				AssetOut: pair.AssetOut,
				DenomOut: assetOut.Denom,
			}

			if accumulate {
				pairsInfo = append(pairsInfo, pairInfo)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPairsResponse{
		PairsInfo:  pairsInfo,
		Pagination: pagination,
	}, nil
}

func (q *queryServer) QueryPair(c context.Context, req *types.QueryPairRequest) (*types.QueryPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	pair, found := q.GetPair(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pair does not exist for id %d", req.Id)
	}

	assetIn, found := q.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetIn)
	}

	assetOut, found := q.GetAsset(ctx, pair.AssetOut)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetOut)
	}

	pairInfo := types.PairInfo{
		Id:       pair.Id,
		AssetIn:  pair.AssetIn,
		DenomIn:  assetIn.Denom,
		AssetOut: pair.AssetOut,
		DenomOut: assetOut.Denom,
	}

	return &types.QueryPairResponse{
		PairInfo: pairInfo,
	}, nil
}

func (q *queryServer) QueryApps(c context.Context, _ *types.QueryAppsRequest) (*types.QueryAppsResponse, error) {
	var (
		ctx         = sdk.UnwrapSDKContext(c)
		apps, found = q.GetApps(ctx)
	)
	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for id")
	}

	return &types.QueryAppsResponse{
		Apps: apps,
	}, nil
}

func (q *queryServer) QueryApp(c context.Context, req *types.QueryAppRequest) (*types.QueryAppResponse, error) {
	var (
		ctx        = sdk.UnwrapSDKContext(c)
		app, found = q.GetApp(ctx, req.Id)
	)
	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for id %d", app.Id)
	}

	return &types.QueryAppResponse{
		App: app,
	}, nil
}

func (q *queryServer) QueryExtendedPairVault(c context.Context, req *types.QueryExtendedPairVaultRequest) (*types.QueryExtendedPairVaultResponse, error) {
	var (
		ctx         = sdk.UnwrapSDKContext(c)
		pair, found = q.GetPairsVault(ctx, req.Id)
	)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pair does not exist for id %d", pair.Id)
	}

	return &types.QueryExtendedPairVaultResponse{
		PairVault: pair,
	}, nil
}

func (q *queryServer) QueryAllExtendedPairVaults(c context.Context, _ *types.QueryAllExtendedPairVaultsRequest) (*types.QueryAllExtendedPairVaultsResponse, error) {
	var (
		ctx               = sdk.UnwrapSDKContext(c)
		pairVaults, found = q.GetPairsVaults(ctx)
	)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Extended pairs does not exist")
	}

	return &types.QueryAllExtendedPairVaultsResponse{
		PairVault: pairVaults,
	}, nil
}

func (q *queryServer) QueryAllExtendedPairVaultsByApp(c context.Context, req *types.QueryAllExtendedPairVaultsByAppRequest) (*types.QueryAllExtendedPairVaultsByAppResponse, error) {
	var (
		ctx               = sdk.UnwrapSDKContext(c)
		pairVaults, found = q.GetPairsVaults(ctx)
	)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Extended pairs does not exist")
	}
	var pairVaultsData []types.ExtendedPairVault
	for _, data := range pairVaults {
		if data.AppId == req.AppId {
			pairVaultsData = append(pairVaultsData, data)
		}
	}

	return &types.QueryAllExtendedPairVaultsByAppResponse{
		ExtendedPair: pairVaultsData,
	}, nil
}

func (q *queryServer) QueryAllExtendedPairStableVaultsIdByApp(c context.Context, req *types.QueryAllExtendedPairStableVaultsIdByAppRequest) (*types.QueryAllExtendedPairStableVaultsIdByAppResponse, error) {
	var (
		ctx               = sdk.UnwrapSDKContext(c)
		pairVaults, found = q.GetPairsVaults(ctx)
		pairVault         []uint64
	)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Extended pairs does not exist")
	}
	for _, data := range pairVaults {
		if (data.AppId == req.AppId) && (data.IsStableMintVault) {
			pairVault = append(pairVault, data.Id)
		}
	}

	return &types.QueryAllExtendedPairStableVaultsIdByAppResponse{
		ExtendedPairsId: pairVault,
	}, nil
}

func (q *queryServer) QueryGovTokenByApp(c context.Context, req *types.QueryGovTokenByAppRequest) (*types.QueryGovTokenByAppResponse, error) {
	var (
		ctx     = sdk.UnwrapSDKContext(c)
		assetID uint64
	)
	appData, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, types.AppIdsDoesntExist
	}
	for _, data := range appData.GenesisToken {
		if data.IsGovToken {
			assetID = data.AssetId
		}
	}

	return &types.QueryGovTokenByAppResponse{
		GovAssetId: assetID,
	}, nil
}

func (q *queryServer) QueryAllExtendedPairStableVaultsByApp(c context.Context, req *types.QueryAllExtendedPairStableVaultsByAppRequest) (*types.QueryAllExtendedPairStableVaultsByAppResponse, error) {
	var (
		ctx               = sdk.UnwrapSDKContext(c)
		pairVaults, found = q.GetPairsVaults(ctx)
	)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Extended pairs does not exist")
	}
	var pairVaultsData []types.ExtendedPairVault
	for _, data := range pairVaults {
		if data.AppId == req.AppId && data.IsStableMintVault {
			pairVaultsData = append(pairVaultsData, data)
		}
	}

	return &types.QueryAllExtendedPairStableVaultsByAppResponse{
		ExtendedPair: pairVaultsData,
	}, nil
}
