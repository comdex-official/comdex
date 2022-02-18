package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/vault/types"
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

func (q *queryServer) QueryAllVaults(c context.Context, req *types.QueryAllVaultsRequest) (*types.QueryAllVaultsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.VaultInfo
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.VaultKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Vault
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

			vaultInfo := types.VaultInfo{
				Id:                     item.ID,
				PairID:                 item.PairID,
				Owner:                  item.Owner,
				Collateral:             sdk.NewCoin(assetIn.Denom, item.AmountIn),
				Debt:                   sdk.NewCoin(assetOut.Denom, item.AmountOut),
				CollateralizationRatio: collateralizationRatio,
			}

			if accumulate {
				items = append(items, vaultInfo)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllVaultsResponse{
		VaultsInfo: items,
		Pagination: pagination,
	}, nil
}

func (q *queryServer) QueryVaults(c context.Context, req *types.QueryVaultsRequest) (*types.QueryVaultsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.VaultInfo
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.VaultForAddressKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Vault
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

			vaultInfo := types.VaultInfo{
				Id:                     item.ID,
				PairID:                 item.PairID,
				Owner:                  item.Owner,
				Collateral:             sdk.NewCoin(assetIn.Denom, item.AmountIn),
				Debt:                   sdk.NewCoin(assetOut.Denom, item.AmountOut),
				CollateralizationRatio: collateralizationRatio,
			}

			if accumulate {
				items = append(items, vaultInfo)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryVaultsResponse{
		VaultsInfo: items,
		Pagination: pagination,
	}, nil
}

func (q *queryServer) QueryVault(c context.Context, req *types.QueryVaultRequest) (*types.QueryVaultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	vault, found := q.GetVault(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "vault does not exist for id %d", req.Id)
	}

	pair, found := q.GetPair(ctx, vault.PairID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pair does not exist for id %d", vault.PairID)
	}

	assetIn, found := q.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetIn)
	}

	assetOut, found := q.GetAsset(ctx, pair.AssetOut)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetOut)
	}

	collateralizationRatio, err := q.CalculateCollaterlizationRatio(ctx, vault.AmountIn, assetIn, vault.AmountOut, assetOut)
	if err != nil {
		return nil, err
	}
	return &types.QueryVaultResponse{
		VaultInfo: types.VaultInfo{
			Id:                     vault.ID,
			PairID:                 vault.PairID,
			Owner:                  vault.Owner,
			Collateral:             sdk.NewCoin(assetIn.Denom, vault.AmountIn),
			Debt:                   sdk.NewCoin(assetOut.Denom, vault.AmountOut),
			CollateralizationRatio: collateralizationRatio,
		},
	}, nil
}
