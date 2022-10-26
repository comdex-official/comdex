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

var _ types.QueryServer = QueryServer{}

type QueryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &QueryServer{
		Keeper: k,
	}
}

func (q QueryServer) QueryAssets(c context.Context, req *types.QueryAssetsRequest) (*types.QueryAssetsResponse, error) {
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

func (q QueryServer) QueryAsset(c context.Context, req *types.QueryAssetRequest) (*types.QueryAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	item, found := q.GetAsset(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", req.Id)
	}

	return &types.QueryAssetResponse{
		Asset: item,
	}, nil
}

func (q QueryServer) QueryPairs(c context.Context, req *types.QueryPairsRequest) (*types.QueryPairsResponse, error) {
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

func (q QueryServer) QueryPair(c context.Context, req *types.QueryPairRequest) (*types.QueryPairResponse, error) {
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

func (q QueryServer) QueryApps(c context.Context, req *types.QueryAppsRequest) (*types.QueryAppsResponse, error) {
	var (
		items []types.AppData
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.NewAppKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.AppData
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

	return &types.QueryAppsResponse{
		Apps:       items,
		Pagination: pagination,
	}, nil
}

func (q QueryServer) QueryApp(c context.Context, req *types.QueryAppRequest) (*types.QueryAppResponse, error) {
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

func (q QueryServer) QueryExtendedPairVault(c context.Context, req *types.QueryExtendedPairVaultRequest) (*types.QueryExtendedPairVaultResponse, error) {
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

func (q QueryServer) QueryAllExtendedPairVaults(c context.Context, req *types.QueryAllExtendedPairVaultsRequest) (*types.QueryAllExtendedPairVaultsResponse, error) {
	var (
		items []types.ExtendedPairVault
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.PairsVaultKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.ExtendedPairVault
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

	return &types.QueryAllExtendedPairVaultsResponse{
		PairVault:  items,
		Pagination: pagination,
	}, nil
}

func (q QueryServer) QueryAllExtendedPairVaultsByApp(c context.Context, req *types.QueryAllExtendedPairVaultsByAppRequest) (*types.QueryAllExtendedPairVaultsByAppResponse, error) {
	var (
		items []types.ExtendedPairVault
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.PairsVaultKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.ExtendedPairVault
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate && item.AppId == req.AppId {
				items = append(items, item)
			}

			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllExtendedPairVaultsByAppResponse{
		ExtendedPair: items,
		Pagination:   pagination,
	}, nil
}

func (q QueryServer) QueryAllExtendedPairStableVaultsIDByApp(c context.Context, req *types.QueryAllExtendedPairStableVaultsIDByAppRequest) (*types.QueryAllExtendedPairStableVaultsIDByAppResponse, error) {
	var (
		items []uint64
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.PairsVaultKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.ExtendedPairVault
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate && (item.AppId == req.AppId) && (item.IsStableMintVault) {
				items = append(items, item.Id)
			}

			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllExtendedPairStableVaultsIDByAppResponse{
		ExtendedPairsId: items,
		Pagination:      pagination,
	}, nil
}

func (q QueryServer) QueryGovTokenByApp(c context.Context, req *types.QueryGovTokenByAppRequest) (*types.QueryGovTokenByAppResponse, error) {
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

func (q QueryServer) QueryAllExtendedPairStableVaultsByApp(c context.Context, req *types.QueryAllExtendedPairStableVaultsByAppRequest) (*types.QueryAllExtendedPairStableVaultsByAppResponse, error) {
	var (
		items []types.ExtendedPairVault
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.PairsVaultKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.ExtendedPairVault
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate && (item.AppId == req.AppId) && (item.IsStableMintVault) {
				items = append(items, item)
			}

			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllExtendedPairStableVaultsByAppResponse{
		ExtendedPair: items,
		Pagination:   pagination,
	}, nil
}
