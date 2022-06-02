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

func (q *queryServer) QueryAllVaults(c context.Context, req *types.QueryAllVaultsRequest) (*types.QueryAllVaultsResponse, error) {

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	vaults := q.GetVaults(ctx)

	return &types.QueryAllVaultsResponse{
		Vault: vaults,
	}, nil
}

func (q *queryServer) QueryAllVaultsByProduct(c context.Context, req *types.QueryAllVaultsByProductRequest) (*types.QueryAllVaultsByProductResponse, error) {

	var (
		ctx           = sdk.UnwrapSDKContext(c)
		productVaults []types.Vault
	)
	vaults := q.GetVaults(ctx)
	for _, data := range vaults {
		if data.AppMappingId == req.AppId {
			productVaults = append(productVaults, data)
		}
	}

	return &types.QueryAllVaultsByProductResponse{
		Vault: productVaults,
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

	return &types.QueryVaultResponse{
		Vault: vault,
	}, nil
}
func (q *queryServer) QueryVaultInfo(c context.Context, req *types.QueryVaultInfoRequest) (*types.QueryVaultInfoResponse, error) {
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

	collateralizationRatio, err := q.CalculateCollaterlizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, vault.AmountOut)
	if err != nil {
		return nil, err
	}
	return &types.QueryVaultInfoResponse{
		VaultsInfo: types.VaultInfo{
			Id:                     req.Id,
			PairID:                 vault.ExtendedPairVaultID,
			Owner:                  vault.Owner,
			Collateral:             vault.AmountIn,
			Debt:                   vault.AmountOut,
			CollateralizationRatio: collateralizationRatio,
		},
	}, nil
}

func (q *queryServer) QueryVaultInfoByOwner(c context.Context, req *types.QueryVaultInfoByOwnerRequest) (*types.QueryVaultInfoByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx        = sdk.UnwrapSDKContext(c)
		vaultsIds  []string
		vaultsInfo []types.VaultInfo
	)

	userVaultAssetData, found := q.GetUserVaultExtendedPairMapping(ctx, req.Owner)
	if !found {
		return nil, status.Errorf(codes.NotFound, "data does not exists for user addesss %s", req.Owner)
	}
	for _, data := range userVaultAssetData.UserVaultApp {
		for _, inData := range data.UserExtendedPairVault {
			vaultsIds = append(vaultsIds, inData.VaultId)
		}
	}

	for _, id := range vaultsIds {
		vault, found := q.GetVault(ctx, id)
		if !found {
			return nil, status.Errorf(codes.NotFound, "vault does not exist for id %d", vault.Id)
		}

		collateralizationRatio, err := q.CalculateCollaterlizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, vault.AmountOut)
		if err != nil {
			return nil, err
		}
		vaults := types.VaultInfo{
			Id:                     vault.Id,
			PairID:                 vault.ExtendedPairVaultID,
			Owner:                  vault.Owner,
			Collateral:             vault.AmountIn,
			Debt:                   vault.AmountOut,
			CollateralizationRatio: collateralizationRatio,
		}
		vaultsInfo = append(vaultsInfo, vaults)

	}

	return &types.QueryVaultInfoByOwnerResponse{
		VaultsInfo: vaultsInfo,
	}, nil
}

func (q *queryServer) QueryAllVaultsByAppAndExtendedPair(c context.Context, req *types.QueryAllVaultsByAppAndExtendedPairRequest) (*types.QueryAllVaultsByAppAndExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx       = sdk.UnwrapSDKContext(c)
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

	for _, data := range vaultData {
		if data.AppMappingId == req.AppId && data.ExtendedPairVaultID == req.ExtendedPairId {
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
	var (
		ctx     = sdk.UnwrapSDKContext(c)
		vaultId = ""
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

	userVaultAssetData, found := q.GetUserVaultExtendedPairMapping(ctx, req.Owner)
	if !found {
		return nil, status.Errorf(codes.NotFound, "data does not exists for user addesss %s", req.Owner)
	}

	for _, data := range userVaultAssetData.UserVaultApp {
		if data.AppMappingId == req.ProductId {
			for _, inData := range data.UserExtendedPairVault {
				if inData.ExtendedPairId == req.ExtendedPairId {
					vaultId = inData.VaultId
				}
			}
		}
	}

	return &types.QueryVaultOfOwnerByExtendedPairResponse{
		Vault_Id: vaultId,
	}, nil
}

func (q *queryServer) QueryVaultByProduct(c context.Context, req *types.QueryVaultByProductRequest) (*types.QueryVaultByProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx       = sdk.UnwrapSDKContext(c)
		vaultsIds []string
	)

	_, found := q.GetApp(ctx, req.ProductId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.ProductId)
	}
	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.ProductId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "extended pair vault data no found for product id %d", req.ProductId)
	}

	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		vaultsIds = append(vaultsIds, data.VaultIds...)
	}

	return &types.QueryVaultByProductResponse{
		VaultIds: vaultsIds,
	}, nil
}

func (q *queryServer) QueryAllVaultByOwner(c context.Context, req *types.QueryAllVaultByOwnerRequest) (*types.QueryAllVaultByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx       = sdk.UnwrapSDKContext(c)
		vaultsIds []string
	)

	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Address is not correct")
	}

	userVaultAssetData, found := q.GetUserVaultExtendedPairMapping(ctx, req.Owner)
	if !found {
		return nil, status.Errorf(codes.NotFound, "data does not exists for user addesss %s", req.Owner)
	}
	for _, data := range userVaultAssetData.UserVaultApp {
		for _, inData := range data.UserExtendedPairVault {
			vaultsIds = append(vaultsIds, inData.VaultId)
		}
	}

	return &types.QueryAllVaultByOwnerResponse{
		VaultIds: vaultsIds,
	}, nil
}

func (q *queryServer) QueryTokenMintedAllProductsByPair(c context.Context, req *types.QueryTokenMintedAllProductsByPairRequest) (*types.QueryTokenMintedAllProductsByPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx         = sdk.UnwrapSDKContext(c)
		tokenMinted = sdk.ZeroInt()
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

	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		if data.ExtendedPairId == req.ExtendedPairId {
			tokenMinted = data.TokenMintedAmount
		}
	}

	return &types.QueryTokenMintedAllProductsByPairResponse{
		TokenMinted: tokenMinted,
	}, nil
}

func (q *queryServer) QueryTokenMintedAllProducts(c context.Context, req *types.QueryTokenMintedAllProductsRequest) (*types.QueryTokenMintedAllProductsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx         = sdk.UnwrapSDKContext(c)
		tokenMinted = sdk.ZeroInt()
	)
	_, found := q.GetApp(ctx, req.ProductId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.ProductId)
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.ProductId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for product id %d", req.ProductId)
	}

	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		tokenMinted = tokenMinted.Add(data.TokenMintedAmount)
	}

	return &types.QueryTokenMintedAllProductsResponse{
		TokenMinted: tokenMinted,
	}, nil
}

func (q *queryServer) QueryVaultCountByProduct(c context.Context, req *types.QueryVaultCountByProductRequest) (*types.QueryVaultCountByProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx          = sdk.UnwrapSDKContext(c)
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
	var (
		ctx          = sdk.UnwrapSDKContext(c)
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
	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		if data.ExtendedPairId == req.ExtendedPairId {
			count = uint64(len(data.VaultIds))
		}
	}

	return &types.QueryVaultCountByProductAndPairResponse{
		VaultCount: count,
	}, nil
}

func (q *queryServer) QueryTotalValueLockedByProductExtendedPair(c context.Context, req *types.QueryTotalValueLockedByProductExtendedPairRequest) (*types.QueryTotalValueLockedByProductExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx         = sdk.UnwrapSDKContext(c)
		valueLocked = sdk.ZeroInt()
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
	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		if data.ExtendedPairId == req.ExtendedPairId {
			valueLocked = data.CollateralLockedAmount
		}
	}

	return &types.QueryTotalValueLockedByProductExtendedPairResponse{
		ValueLocked: &valueLocked,
	}, nil
}

func (q *queryServer) QueryExtendedPairIDByProduct(c context.Context, req *types.QueryExtendedPairIDByProductRequest) (*types.QueryExtendedPairIDByProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx     = sdk.UnwrapSDKContext(c)
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
	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
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
	var (
		ctx = sdk.UnwrapSDKContext(c)
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
	var (
		ctx            = sdk.UnwrapSDKContext(c)
		stableMintData []*types.StableMintVault
	)
	stableMint := q.GetStableMintVaults(ctx)
	for _, data := range stableMint {
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
	var (
		ctx            = sdk.UnwrapSDKContext(c)
		stableMintData types.StableMintVault
	)
	stableMint := q.GetStableMintVaults(ctx)
	for _, data := range stableMint {
		if data.AppMappingId == req.AppId && data.ExtendedPairVaultID == req.ExtendedPairId {
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
	var (
		ctx                = sdk.UnwrapSDKContext(c)
		extendedPairIdData types.ExtendedPairVaultMapping
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
		if data.ExtendedPairId == req.ExtendedPairId {
			extendedPairIdData = *data
		}
	}

	return &types.QueryExtendedPairVaultMappingByAppAndExtendedPairIdResponse{
		ExtendedPairVaultMapping: &extendedPairIdData,
	}, nil
}

func (q *queryServer) QueryExtendedPairVaultMappingByApp(c context.Context, req *types.QueryExtendedPairVaultMappingByAppRequest) (*types.QueryExtendedPairVaultMappingByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx     = sdk.UnwrapSDKContext(c)
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
	var (
		ctx               = sdk.UnwrapSDKContext(c)
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
	for _, data := range userVaultAssetData.UserVaultApp {
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
	var (
		ctx     = sdk.UnwrapSDKContext(c)
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
	for _, data := range userVaultAssetData.UserVaultApp {
		if data.AppMappingId == req.AppId {
			for _, inData := range data.UserExtendedPairVault {
				if inData.ExtendedPairId == req.ExtendedPair {
					vaultId = inData.VaultId
				}
			}
		}
	}

	return &types.QueryExtendedPairVaultMappingByOwnerAndAppAndExtendedPairIDResponse{
		VaultId: vaultId,
	}, nil
}

func (q *queryServer) QueryTVLLockedByAppOfAllExtendedPairs(c context.Context, req *types.QueryTVLLockedByAppOfAllExtendedPairsRequest) (*types.QueryTVLLockedByAppOfAllExtendedPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx     = sdk.UnwrapSDKContext(c)
		tvlData []types.TvlLockedDataMap
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
		pairId, _ := q.GetPair(ctx, extPairVault.PairId)

		var tvl types.TvlLockedDataMap

		denom, found := q.GetAsset(ctx, pairId.AssetIn)
		if !found {
			return nil, types.ErrorAssetDoesNotExist
		}
		tvl.AssetDenom = denom.Denom
		tvl.CollateralLockedAmount = data.CollateralLockedAmount

		tvlData = append(tvlData, tvl)
	}

	return &types.QueryTVLLockedByAppOfAllExtendedPairsResponse{
		Tvldata: tvlData,
	}, nil
}

func (q *queryServer) QueryTotalTVLByApp(c context.Context, req *types.QueryTotalTVLByAppRequest) (*types.QueryTotalTVLByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx    = sdk.UnwrapSDKContext(c)
		locked uint64
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
		pairId, _ := q.GetPair(ctx, extPairVault.PairId)

		rate, _ := q.GetPriceForAsset(ctx, pairId.AssetIn)
		locked = locked + (rate * data.CollateralLockedAmount.Uint64())
	}

	return &types.QueryTotalTVLByAppResponse{
		CollateralLocked: locked,
	}, nil
}
