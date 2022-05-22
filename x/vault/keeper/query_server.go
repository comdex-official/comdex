package keeper

import (
	"context"
	"strconv"

	// "github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/cosmos/cosmos-sdk/types/query"
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

// func (q *queryServer) QueryAllVaults(c context.Context, req *types.QueryAllVaultsRequest) (*types.QueryAllVaultsResponse, error) {
// 	if req == nil {
// 		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
// 	}

// 	var (
// 		items []types.VaultInfo
// 		ctx   = sdk.UnwrapSDKContext(c)
// 	)

// 	pagination, err := query.FilteredPaginate(
// 		prefix.NewStore(q.Store(ctx), types.VaultKeyPrefix),
// 		req.Pagination,
// 		func(_, value []byte, accumulate bool) (bool, error) {
// 			var item types.Vault
// 			if err := q.cdc.Unmarshal(value, &item); err != nil {
// 				return false, err
// 			}

// 			pair, found := q.GetPair(ctx, item.ExtendedPairVaultID)
// 			if !found {
// 				return false, status.Errorf(codes.NotFound, "pair does not exist for id %d", item.ExtendedPairVaultID)
// 			}

// 			assetIn, found := q.GetAsset(ctx, pair.AssetIn)
// 			if !found {
// 				return false, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetIn)
// 			}

// 			assetOut, found := q.GetAsset(ctx, pair.AssetOut)
// 			if !found {
// 				return false, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetOut)
// 			}

// 			collateralizationRatio, err := q.CalculateCollaterlizationRatio(ctx,pair.Id, item.AmountIn, item.AmountOut)
// 			if err != nil {
// 				return false, err
// 			}

// 			vaultInfo := types.VaultInfo{
// 				Id:                     item.Id,
// 				PairID:                 item.ExtendedPairVaultID,
// 				Owner:                  item.Owner,
// 				Collateral:             sdk.NewCoin(assetIn.Denom, item.AmountIn),
// 				Debt:                   sdk.NewCoin(assetOut.Denom, item.AmountOut),
// 				CollateralizationRatio: collateralizationRatio,
// 			}

// 			if accumulate {
// 				items = append(items, vaultInfo)
// 			}

// 			return true, nil
// 		},
// 	)

// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &types.QueryAllVaultsResponse{
// 		VaultsInfo: items,
// 		Pagination: pagination,
// 	}, nil
// }

func (q *queryServer) QueryVault(c context.Context, req *types.QueryVaultRequest) (*types.QueryVaultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.VaultInfo
		ctx   = sdk.UnwrapSDKContext(c)
	)
			var item types.Vault

			pair, found := q.GetPair(ctx, item.ExtendedPairVaultID)
			if !found {
				return nil, status.Errorf(codes.NotFound, "pair does not exist for id %d", item.ExtendedPairVaultID)
			}

			assetIn, found := q.GetAsset(ctx, pair.AssetIn)
			if !found {
				return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetIn)
			}

			assetOut, found := q.GetAsset(ctx, pair.AssetOut)
			if !found {
				return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetOut)
			}

			collateralizationRatio, err := q.CalculateCollaterlizationRatio(ctx,pair.Id, item.AmountIn, item.AmountOut)
			if err != nil {
				return nil, err
			}
			newItemId, err := strconv.ParseUint(item.Id, 10, 64)
			vaultInfo := types.VaultInfo{
				Id:                     newItemId,
				PairID:                 item.ExtendedPairVaultID,
				Owner:                  item.Owner,
				Collateral:             sdk.NewCoin(assetIn.Denom, item.AmountIn),
				Debt:                   sdk.NewCoin(assetOut.Denom, item.AmountOut),
				CollateralizationRatio: collateralizationRatio,
			}
			items = append(items, vaultInfo)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryVaultResponse{
		VaultInfo: items,
	}, nil
}
// func (q *queryServer) QueryVault(c context.Context, req *types.QueryVaultRequest) (*types.QueryVaultResponse, error) {
// 	if req == nil {
// 		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
// 	}

// 	var (
// 		ctx = sdk.UnwrapSDKContext(c)
// 	)

// 	vault, found := q.GetVault(ctx, req.Id)
// 	if !found {
// 		return nil, status.Errorf(codes.NotFound, "vault does not exist for id %d", req.Id)
// 	}

// 	pair, found := q.GetPair(ctx, vault.ExtendedPairVaultID)
// 	if !found {
// 		return nil, status.Errorf(codes.NotFound, "pair does not exist for id %d", vault.ExtendedPairVaultID)
// 	}

// 	assetIn, found := q.GetAsset(ctx, pair.AssetIn)
// 	if !found {
// 		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetIn)
// 	}

// 	assetOut, found := q.GetAsset(ctx, pair.AssetOut)
// 	if !found {
// 		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", pair.AssetOut)
// 	}

// 	collateralizationRatio, err := q.CalculateCollaterlizationRatio(ctx,pair.Id, item.AmountIn, item.AmountOut)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &types.QueryVaultResponse{
// 		VaultInfo: types.VaultInfo{
// 			Id:                     vault.Id,
// 			PairID:                 vault.ExtendedPairVaultID,
// 			Owner:                  vault.Owner,
// 			Collateral:             sdk.NewCoin(assetIn.Denom, vault.AmountIn),
// 			Debt:                   sdk.NewCoin(assetOut.Denom, vault.AmountOut),
// 			CollateralizationRatio: collateralizationRatio,
// 		},
// 	}, nil
// }

func (q *queryServer) QueryVaultOfOwnerByPair(c context.Context, req *types.QueryVaultOfOwnerByPairRequest) (*types.QueryVaultOfOwnerByPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		vault_id = ""
	)

			_, found := q.GetApp(ctx, req.ProductId)
			if !found {
				return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.ProductId)
			}

			_, err := sdk.AccAddressFromBech32(req.Owner)
 			if err != nil {
  				return nil, status.Errorf(codes.NotFound, "Address is not correct")
 			}

			_, nfound := q.GetPairsVault(ctx, req.ExtendedPairId)
			if !nfound {
				return nil, status.Errorf(codes.NotFound, "extended pair does not exist for id %d", req.ExtendedPairId)
			}

			userVaultAssetData,found := q.GetUserVaultExtendedPairMapping(ctx, req.Owner)
			if !found {
				return nil, status.Errorf(codes.NotFound, "data does not exists for user addesss %s", req.Owner)
			}
			
			for _, data := range userVaultAssetData.UserVaultApp{
				if data.AppMappingId == req.ProductId{
					for _, inData := range data.UserExtendedPairVault{
						if inData.ExtendedPairId == req.ExtendedPairId {
							vault_id = inData.VaultId
						}
					}
				}
			}

	return &types.QueryVaultOfOwnerByPairResponse{
		Vault_Id: vault_id,
	}, nil
}

func (q *queryServer) QueryVaultByProduct(c context.Context, req *types.QueryVaultByProductRequest) (*types.QueryVaultByProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		vaultsIds []string
	)

			_, found := q.GetApp(ctx, req.ProductId)
			if !found {
				return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.ProductId)
			}
			appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.ProductId)
			if !found{
				return nil, status.Errorf(codes.NotFound, "extended pair vault data no found for product id %d", req.ProductId)
			}

			
			for _, data := range appExtendedPairVaultData.ExtendedPairVaults{
				vaultsIds = append(vaultsIds, data.VaultIds...)
			}

	return &types.QueryVaultByProductResponse{
		Vault_Ids: vaultsIds,
	}, nil
}

func (q *queryServer) QueryAllVaultByProducts(c context.Context, req *types.QueryAllVaultByProductsRequest) (*types.QueryAllVaultByProductsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
		var( 
			ctx   = sdk.UnwrapSDKContext(c)
			vaultsIds []string
		)

		_, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			 return nil, status.Errorf(codes.NotFound, "Address is not correct")
		}

		userVaultAssetData,found := q.GetUserVaultExtendedPairMapping(ctx, req.Owner)
		if !found {
			return nil, status.Errorf(codes.NotFound, "data does not exists for user addesss %s", req.Owner)
		}
			for _, data := range userVaultAssetData.UserVaultApp{
				for _, inData := range data.UserExtendedPairVault{
					vaultsIds = append(vaultsIds, inData.VaultId)
				}
			}

	return &types.QueryAllVaultByProductsResponse{
		Vault_Ids: vaultsIds,
	}, nil
}

func (q *queryServer) QueryTokenMintedAllProductsByPair(c context.Context, req *types.QueryTokenMintedAllProductsByPairRequest) (*types.QueryTokenMintedAllProductsByPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		token_minted = sdk.ZeroInt()
	)
		_, found := q.GetApp(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.ProductId)
		}

		_, nfound := q.GetPairsVault(ctx, req.ExtendedPairId)
		if !nfound {
			return nil, status.Errorf(codes.NotFound, "extended pair does not exist for id %d", req.ExtendedPairId)
		}
		appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for product id %d", req.ProductId)
		}
			
			for _, data := range appExtendedPairVaultData.ExtendedPairVaults{
				if data.ExtendedPairId == req.ExtendedPairId {
					token_minted = *data.TokenMintedAmount 
				}
			}

	return &types.QueryTokenMintedAllProductsByPairResponse{
		Token_Minted: &token_minted,
	}, nil
}

func (q *queryServer) QueryTokenMintedAllProducts(c context.Context, req *types.QueryTokenMintedAllProductsRequest) (*types.QueryTokenMintedAllProductsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		token_minted = sdk.ZeroInt()
	)
		_, found := q.GetApp(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.ProductId)
		}

		appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for product id %d", req.ProductId)
		}
			
			for _, data := range appExtendedPairVaultData.ExtendedPairVaults{
				token_minted = *data.TokenMintedAmount
			}

	return &types.QueryTokenMintedAllProductsResponse{
		Token_Minted: &token_minted,
	}, nil
}

func (q *queryServer) QueryVaultCountByProduct(c context.Context, req *types.QueryVaultCountByProductRequest) (*types.QueryVaultCountByProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		count uint64 = 0
	)
		_, found := q.GetApp(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.ProductId)
		}

		appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for product id %d", req.ProductId)
		}
		
		count = appExtendedPairVaultData.Counter

	return &types.QueryVaultCountByProductResponse{
		VaultCount: count,
	}, nil
}

func (q *queryServer) QueryVaultCountByProductAndPair(c context.Context, req *types.QueryVaultCountByProductAndPairRequest) (*types.QueryVaultCountByProductAndPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		count uint64 = 0
	)
		_, found := q.GetApp(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.ProductId)
		}

		appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for product id %d", req.ProductId)
		}
		for _, data := range appExtendedPairVaultData.ExtendedPairVaults{
			if data.ExtendedPairId == req.ExtendedPairId{
				count = uint64(len(data.VaultIds))
			}
		}
		count = appExtendedPairVaultData.Counter

	return &types.QueryVaultCountByProductAndPairResponse{
		VaultCount: count,
	}, nil
}

func (q *queryServer) QueryTotalValueLockedByProductExtendedPair(c context.Context, req *types.QueryTotalValueLockedByProductExtendedPairRequest) (*types.QueryTotalValueLockedByProductExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		value_locked = sdk.ZeroInt()
	)
		_, found := q.GetApp(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.ProductId)
		}
		_, nfound := q.GetPairsVault(ctx, req.ExtendedPairId)
		if !nfound {
			return nil, status.Errorf(codes.NotFound, "extended pair does not exist for id %d", req.ExtendedPairId)
		}

		appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for product id %d", req.ProductId)
		}
		for _, data := range appExtendedPairVaultData.ExtendedPairVaults{
			if data.ExtendedPairId == req.ExtendedPairId{
				value_locked = *data.CollateralLockedAmount
			}
		}
	

	return &types.QueryTotalValueLockedByProductExtendedPairResponse{
		ValueLocked: &value_locked,
	}, nil
} 


func (q *queryServer) QueryExtendedPairIDByProduct(c context.Context, req *types.QueryExtendedPairIDByProductRequest) (*types.QueryExtendedPairIDByProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		pairIds []uint64
	)
		_, found := q.GetApp(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.ProductId)
		}

		appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.ProductId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for product id %d", req.ProductId)
		}
		for _, data := range appExtendedPairVaultData.ExtendedPairVaults{
			pairIds = append(pairIds, data.ExtendedPairId)
		}
	

	return &types.QueryExtendedPairIDByProductResponse{
		PairId: pairIds,
	}, nil
} 

func (q *queryServer) QueryStableVaultInfo(c context.Context, req *types.QueryStableVaultInfoRequest) (*types.QueryStableVaultInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
	)
		stableMintData, found := q.GetStableMintVault(ctx, req.StableVaultId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "stable mint data not exist for id %d", req.StableVaultId)
		}
	

	return &types.QueryStableVaultInfoResponse{
		StableMintVault: &stableMintData,
	}, nil
}

func (q *queryServer) QueryAllStableVaults(c context.Context, req *types.QueryAllStableVaultsRequest) (*types.QueryAllStableVaultsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		stableMintData[] *types.StableMintVault
	)
		stableMint := q.GetStableMintVaults(ctx)
		for _ ,data := range stableMint{
			if data.AppMappingId == req.AppId {
				stableMintData = append(stableMintData, &data)
			}
		} 
	
	return &types.QueryAllStableVaultsResponse{
		StableMintVault: stableMintData,
	}, nil
}

func (q *queryServer) QueryStableVaultByProductExtendedPair(c context.Context, req *types.QueryStableVaultByProductExtendedPairRequest) (*types.QueryStableVaultByProductExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		stableMintData types.StableMintVault
	)
		stableMint := q.GetStableMintVaults(ctx)
		for _ ,data := range stableMint{
			if data.AppMappingId == req.AppId && data.ExtendedPairVaultID == req.ExtendedPairId{
				stableMintData = data
			}
		} 
	
	return &types.QueryStableVaultByProductExtendedPairResponse{
		StableMintVault: &stableMintData,
	}, nil
}