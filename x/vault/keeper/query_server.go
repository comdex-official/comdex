package keeper

import (
	"context"

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

func (q *queryServer) QueryAllVaults(c context.Context, req *types.QueryAllVaultsRequest) (*types.QueryAllVaultsResponse, error) {

	var (
		ctx   = sdk.UnwrapSDKContext(c)
	)
	vaults := q.GetVaults(ctx)
	
	return &types.QueryAllVaultsResponse{
		Vault: vaults,
	}, nil
} 

func (q *queryServer) QueryAllVaultsByProduct(c context.Context, req *types.QueryAllVaultsByProductRequest) (*types.QueryAllVaultsByProductResponse, error) {

	var (
		ctx   = sdk.UnwrapSDKContext(c)
		productvaults []types.Vault
	)
	vaults := q.GetVaults(ctx)
	for _, data := range vaults{
		if data.AppMappingId == req.AppId {
			productvaults = append(productvaults, data)
		}
	}
	
	return &types.QueryAllVaultsByProductResponse{
		Vault: productvaults,
	}, nil
}

func (q *queryServer) QueryVault(c context.Context, req *types.QueryVaultRequest) (*types.QueryVaultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx   = sdk.UnwrapSDKContext(c)
	)
			vault, found := q.GetVault(ctx, req.Id)
			if !found {
				return nil, status.Errorf(codes.NotFound, "vault does not exist for id %d", req.Id)
			}

	return &types.QueryVaultResponse{
		Vault: vault,
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


func (q *queryServer) QueryAllVaultsByAppAndExtendedPair(c context.Context, req *types.QueryAllVaultsByAppAndExtendedPairRequest) (*types.QueryAllVaultsByAppAndExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		vaultList []types.Vault
	)

			_, found := q.GetApp(ctx, req.AppId)
			if !found {
				return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
			}

			_, nfound := q.GetPairsVault(ctx, req.ExtendedPairId)
			if !nfound {
				return nil, status.Errorf(codes.NotFound, "extended pair does not exist for id %d", req.ExtendedPairId)
			}

			vaultData := q.GetVaults(ctx)
			
			for _, data := range vaultData{
				if data.AppMappingId == req.AppId && data.ExtendedPairVaultID == req.ExtendedPairId{
					vaultList = append(vaultList, data)
				}
			}

	return &types.QueryAllVaultsByAppAndExtendedPairResponse{
		Vault: vaultList,
	}, nil
}


func (q *queryServer) QueryVaultOfOwnerByExtendedPair(c context.Context, req *types.QueryVaultOfOwnerByExtendedPairRequest) (*types.QueryVaultOfOwnerByExtendedPairResponse, error) {
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

	return &types.QueryVaultOfOwnerByExtendedPairResponse{
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

func (q *queryServer) QueryAllVaultByOwner(c context.Context, req *types.QueryAllVaultByOwnerRequest) (*types.QueryAllVaultByOwnerResponse, error) {
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

	return &types.QueryAllVaultByOwnerResponse{
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
				token_minted = token_minted.Add(*data.TokenMintedAmount)
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
		ExtendedPairIds: pairIds,
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

func (q *queryServer) QueryExtendedPairVaultMappingByAppAndExtendedPairId(c context.Context, req *types.QueryExtendedPairVaultMappingByAppAndExtendedPairIdRequest) (*types.QueryExtendedPairVaultMappingByAppAndExtendedPairIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		extendedpairIdData types.ExtendedPairVaultMapping
	)
		_, found := q.GetApp(ctx, req.AppId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
		}

		appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for product id %d", req.AppId)
		}
		for _, data := range appExtendedPairVaultData.ExtendedPairVaults{
			if data.ExtendedPairId == req.ExtendedPairId{
				extendedpairIdData = *data
			}
		}
	

	return &types.QueryExtendedPairVaultMappingByAppAndExtendedPairIdResponse{
		ExtendedPairVaultMapping: &extendedpairIdData,
	}, nil
} 

func (q *queryServer) QueryExtendedPairVaultMappingByApp(c context.Context, req *types.QueryExtendedPairVaultMappingByAppRequest) (*types.QueryExtendedPairVaultMappingByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		pairIds []*types.ExtendedPairVaultMapping
	)
		_, found := q.GetApp(ctx, req.AppId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
		}

		appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for product id %d", req.AppId)
		}
		pairIds = append(pairIds, appExtendedPairVaultData.ExtendedPairVaults...)
	

	return &types.QueryExtendedPairVaultMappingByAppResponse{
		ExtendedPairVaultMapping: pairIds,
	}, nil
} 

func (q *queryServer) QueryExtendedPairVaultMappingByOwnerAndApp(c context.Context, req *types.QueryExtendedPairVaultMappingByOwnerAndAppRequest) (*types.QueryExtendedPairVaultMappingByOwnerAndAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		extendedPairVault []*types.ExtendedPairToVaultMapping
	)
		_, found := q.GetApp(ctx, req.AppId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
		}

		userVaultAssetData, found := q.GetUserVaultExtendedPairMapping(ctx, req.Owner)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for owner %d", req.Owner)
		}
		for _, data := range userVaultAssetData.UserVaultApp{
			if data.AppMappingId == req.AppId {
					extendedPairVault = append(extendedPairVault, data.UserExtendedPairVault...)
			}
		}
	

	return &types.QueryExtendedPairVaultMappingByOwnerAndAppResponse{
		ExtendedPairtoVaultMapping: extendedPairVault,
	}, nil
} 

func (q *queryServer) QueryExtendedPairVaultMappingByOwnerAndAppAndExtendedPairID(c context.Context, req *types.QueryExtendedPairVaultMappingByOwnerAndAppAndExtendedPairIDRequest) (*types.QueryExtendedPairVaultMappingByOwnerAndAppAndExtendedPairIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		vaultId string
	)
		_, found := q.GetApp(ctx, req.AppId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
		}

		userVaultAssetData, found := q.GetUserVaultExtendedPairMapping(ctx, req.Owner)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for owner %d", req.Owner)
		}
		for _, data := range userVaultAssetData.UserVaultApp{
			if data.AppMappingId == req.AppId {
				for _, inData := range data.UserExtendedPairVault{
					if inData.ExtendedPairId == req.ExtendedPair{
						vaultId = inData.VaultId
					}
				}
			}
		}
	

	return &types.QueryExtendedPairVaultMappingByOwnerAndAppAndExtendedPairIDResponse{
		VaultId: vaultId,
	}, nil
} 

func (q *queryServer) QueryTVLlockedByApp(c context.Context, req *types.QueryTVLlockedByAppRequest) (*types.QueryTVLlockedByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		tvlData []*types.TvlLockedDataMap
	)
		_, found := q.GetApp(ctx, req.AppId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
		}

		appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for product id %d", req.AppId)
		}
		for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
			extPairVault, _ := q.GetPairsVault(ctx, data.ExtendedPairId)
			pairId, _ := q.GetPair(ctx,extPairVault.PairId)

			var tvl types.TvlLockedDataMap

			tvl.AssetId = pairId.AssetIn
			tvl.CollateralLockedAmount = data.CollateralLockedAmount
			
			tvlData = append(tvlData, &tvl)
		}

	return &types.QueryTVLlockedByAppResponse{
		Tvldata: tvlData,
	}, nil
} 