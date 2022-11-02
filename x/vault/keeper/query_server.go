package keeper

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/vault/types"
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

func (q QueryServer) QueryAllVaults(c context.Context, req *types.QueryAllVaultsRequest) (*types.QueryAllVaultsResponse, error) {
	var (
		items []types.Vault
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

			if accumulate {
				items = append(items, item)
			}

			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllVaultsResponse{
		Vault:      items,
		Pagination: pagination,
	}, nil
}

func (q QueryServer) QueryAllVaultsByApp(c context.Context, req *types.QueryAllVaultsByAppRequest) (*types.QueryAllVaultsByAppResponse, error) {
	var (
		ctx   = sdk.UnwrapSDKContext(c)
		items []types.Vault
	)
	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.VaultKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Vault
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate {
				if item.AppId == req.AppId {
					items = append(items, item)
				}
			}

			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllVaultsByAppResponse{
		Vault:      items,
		Pagination: pagination,
	}, nil
}

func (q QueryServer) QueryVault(c context.Context, req *types.QueryVaultRequest) (*types.QueryVaultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	vault, found := q.GetVault(ctx, req.Id)
	if !found {
		return &types.QueryVaultResponse{}, nil
	}

	return &types.QueryVaultResponse{
		Vault: vault,
	}, nil
}

func (q QueryServer) QueryVaultInfoByVaultID(c context.Context, req *types.QueryVaultInfoByVaultIDRequest) (*types.QueryVaultInfoByVaultIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	vault, found := q.GetVault(ctx, req.Id)
	if !found {
		return &types.QueryVaultInfoByVaultIDResponse{}, nil
	}

	collateralizationRatio, err := q.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, vault.AmountOut)
	if err != nil {
		return nil, err
	}
	pairVaults, _ := q.asset.GetPairsVault(ctx, vault.ExtendedPairVaultID)
	pairID, _ := q.asset.GetPair(ctx, pairVaults.PairId)
	assetIn, _ := q.asset.GetAsset(ctx, pairID.AssetIn)
	assetOut, _ := q.asset.GetAsset(ctx, pairID.AssetOut)
	return &types.QueryVaultInfoByVaultIDResponse{
		VaultsInfo: types.VaultInfo{
			Id:                     req.Id,
			ExtendedPairID:         vault.ExtendedPairVaultID,
			Owner:                  vault.Owner,
			Collateral:             vault.AmountIn,
			Debt:                   vault.AmountOut,
			CollateralizationRatio: collateralizationRatio,
			ExtendedPairName:       pairVaults.PairName,
			InterestRate:           pairVaults.StabilityFee,
			AssetInDenom:           assetIn.Denom,
			AssetOutDenom:          assetOut.Denom,
			MinCr:                  pairVaults.MinCr,
		},
	}, nil
}

func (q QueryServer) QueryVaultInfoOfOwnerByApp(c context.Context, req *types.QueryVaultInfoOfOwnerByAppRequest) (*types.QueryVaultInfoOfOwnerByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx        = sdk.UnwrapSDKContext(c)
		vaultsIds  []uint64
		vaultsInfo []types.VaultInfo
	)
	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Address is not correct")
	}

	userVaultAssetData, _ := q.GetUserAppMappingData(ctx, req.Owner, req.AppId)

	for _, data := range userVaultAssetData {
		vaultsIds = append(vaultsIds, data.VaultId)
	}
	count := len(vaultsIds)
	for _, id := range vaultsIds {
		vault, found := q.GetVault(ctx, id)
		if !found {
			count--
			continue
		}

		collateralizationRatio, err := q.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, vault.AmountOut)
		if err != nil {
			return nil, err
		}
		pairVaults, _ := q.asset.GetPairsVault(ctx, vault.ExtendedPairVaultID)
		pairID, _ := q.asset.GetPair(ctx, pairVaults.PairId)
		assetIn, _ := q.asset.GetAsset(ctx, pairID.AssetIn)
		assetOut, _ := q.asset.GetAsset(ctx, pairID.AssetOut)

		vaults := types.VaultInfo{
			Id:                     vault.Id,
			ExtendedPairID:         vault.ExtendedPairVaultID,
			Owner:                  vault.Owner,
			Collateral:             vault.AmountIn,
			Debt:                   vault.AmountOut,
			CollateralizationRatio: collateralizationRatio,
			ExtendedPairName:       pairVaults.PairName,
			InterestRate:           pairVaults.StabilityFee,
			AssetInDenom:           assetIn.Denom,
			AssetOutDenom:          assetOut.Denom,
			MinCr:                  pairVaults.MinCr,
		}
		vaultsInfo = append(vaultsInfo, vaults)
	}
	if count == 0 {
		return &types.QueryVaultInfoOfOwnerByAppResponse{}, nil
	}

	return &types.QueryVaultInfoOfOwnerByAppResponse{
		VaultsInfo: vaultsInfo,
	}, nil
}

func (q QueryServer) QueryAllVaultsByAppAndExtendedPair(c context.Context, req *types.QueryAllVaultsByAppAndExtendedPairRequest) (*types.QueryAllVaultsByAppAndExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx       = sdk.UnwrapSDKContext(c)
		vaultList []types.Vault
	)

	appExtendedPairData, found := q.GetAppExtendedPairVaultMappingData(ctx, req.AppId, req.ExtendedPairId)
	if !found {
		return nil, types.ErrorAppExtendedPairDataDoesNotExists
	}
	vaultIDs := appExtendedPairData.VaultIds

	for _, data := range vaultIDs {
		vaultData, _ := q.GetVault(ctx, data)
		vaultList = append(vaultList, vaultData)
	}

	return &types.QueryAllVaultsByAppAndExtendedPairResponse{
		Vault: vaultList,
	}, nil
}

func (q QueryServer) QueryVaultIDOfOwnerByExtendedPairAndApp(c context.Context, req *types.QueryVaultIDOfOwnerByExtendedPairAndAppRequest) (*types.QueryVaultIDOfOwnerByExtendedPairAndAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	userVault, found := q.GetUserAppExtendedPairMappingData(ctx, req.Owner, req.AppId, req.ExtendedPairId)
	if !found {
		return &types.QueryVaultIDOfOwnerByExtendedPairAndAppResponse{}, nil
	}

	return &types.QueryVaultIDOfOwnerByExtendedPairAndAppResponse{
		Vault_Id: userVault.VaultId,
	}, nil
}

func (q QueryServer) QueryVaultIdsByAppInAllExtendedPairs(c context.Context, req *types.QueryVaultIdsByAppInAllExtendedPairsRequest) (*types.QueryVaultIdsByAppInAllExtendedPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx       = sdk.UnwrapSDKContext(c)
		vaultsIds []uint64
	)

	appExtendedPairVaultData, found := q.GetAppMappingData(ctx, req.AppId)
	if !found {
		return &types.QueryVaultIdsByAppInAllExtendedPairsResponse{}, nil
	}

	for _, data := range appExtendedPairVaultData {
		vaultsIds = append(vaultsIds, data.VaultIds...)
	}

	return &types.QueryVaultIdsByAppInAllExtendedPairsResponse{
		VaultIds: vaultsIds,
	}, nil
}

func (q QueryServer) QueryAllVaultIdsByAnOwner(c context.Context, req *types.QueryAllVaultIdsByAnOwnerRequest) (*types.QueryAllVaultIdsByAnOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx       = sdk.UnwrapSDKContext(c)
		vaultsIds []uint64
	)

	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Address is not correct")
	}

	userVaultAssetData := q.GetUserMappingData(ctx, req.Owner)

	for _, data := range userVaultAssetData {
		vaultsIds = append(vaultsIds, data.VaultId)
	}

	return &types.QueryAllVaultIdsByAnOwnerResponse{
		VaultIds: vaultsIds,
	}, nil
}

func (q QueryServer) QueryTokenMintedByAppAndExtendedPair(c context.Context, req *types.QueryTokenMintedByAppAndExtendedPairRequest) (*types.QueryTokenMintedByAppAndExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMappingData(ctx, req.AppId, req.ExtendedPairId)
	if !found {
		return &types.QueryTokenMintedByAppAndExtendedPairResponse{}, nil
	}

	return &types.QueryTokenMintedByAppAndExtendedPairResponse{
		TokenMinted: appExtendedPairVaultData.TokenMintedAmount,
	}, nil
}

func (q QueryServer) QueryTokenMintedAssetWiseByApp(c context.Context, req *types.QueryTokenMintedAssetWiseByAppRequest) (*types.QueryTokenMintedAssetWiseByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx        = sdk.UnwrapSDKContext(c)
		mintedData []types.MintedDataMap
	)

	appExtendedPairVaultData, found := q.GetAppMappingData(ctx, req.AppId)
	if !found {
		return &types.QueryTokenMintedAssetWiseByAppResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData {
		extPairVault, _ := q.asset.GetPairsVault(ctx, data.ExtendedPairId)
		pairID, _ := q.asset.GetPair(ctx, extPairVault.PairId)

		var minted types.MintedDataMap

		denom, found := q.asset.GetAsset(ctx, pairID.AssetOut)
		if !found {
			return nil, types.ErrorAssetDoesNotExist
		}
		minted.AssetDenom = denom.Denom
		minted.MintedAmount = data.TokenMintedAmount

		mintedData = append(mintedData, minted)
	}

	return &types.QueryTokenMintedAssetWiseByAppResponse{
		MintedData: mintedData,
	}, nil
}

func (q QueryServer) QueryVaultCountByApp(c context.Context, req *types.QueryVaultCountByAppRequest) (*types.QueryVaultCountByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx   = sdk.UnwrapSDKContext(c)
		count uint64
	)
	_, found := q.asset.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppMappingData(ctx, req.AppId)
	if !found {
		return &types.QueryVaultCountByAppResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData {
		count += uint64(len(data.VaultIds))
	}

	return &types.QueryVaultCountByAppResponse{
		VaultCount: count,
	}, nil
}

func (q QueryServer) QueryVaultCountByAppAndExtendedPair(c context.Context, req *types.QueryVaultCountByAppAndExtendedPairRequest) (*types.QueryVaultCountByAppAndExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx   = sdk.UnwrapSDKContext(c)
		count uint64
	)
	_, found := q.asset.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMappingData(ctx, req.AppId, req.ExtendedPairId)
	if !found {
		return &types.QueryVaultCountByAppAndExtendedPairResponse{}, nil
	}

	count = uint64(len(appExtendedPairVaultData.VaultIds))

	return &types.QueryVaultCountByAppAndExtendedPairResponse{
		VaultCount: count,
	}, nil
}

func (q QueryServer) QueryTotalValueLockedByAppAndExtendedPair(c context.Context, req *types.QueryTotalValueLockedByAppAndExtendedPairRequest) (*types.QueryTotalValueLockedByAppAndExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx         = sdk.UnwrapSDKContext(c)
		valueLocked = sdk.ZeroInt()
	)
	_, found := q.asset.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}
	_, nfound := q.asset.GetPairsVault(ctx, req.ExtendedPairId)
	if !nfound {
		return &types.QueryTotalValueLockedByAppAndExtendedPairResponse{}, nil
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMappingData(ctx, req.AppId, req.ExtendedPairId)
	if !found {
		return &types.QueryTotalValueLockedByAppAndExtendedPairResponse{}, nil
	}

	valueLocked = appExtendedPairVaultData.CollateralLockedAmount

	return &types.QueryTotalValueLockedByAppAndExtendedPairResponse{
		ValueLocked: &valueLocked,
	}, nil
}

func (q QueryServer) QueryExtendedPairIDsByApp(c context.Context, req *types.QueryExtendedPairIDsByAppRequest) (*types.QueryExtendedPairIDsByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx     = sdk.UnwrapSDKContext(c)
		pairIDs []uint64
	)

	appExtendedPairVaultData, found := q.GetAppMappingData(ctx, req.AppId)
	if !found {
		return &types.QueryExtendedPairIDsByAppResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData {
		pairIDs = append(pairIDs, data.ExtendedPairId)
	}

	return &types.QueryExtendedPairIDsByAppResponse{
		ExtendedPairIds: pairIDs,
	}, nil
}

func (q QueryServer) QueryStableVaultByVaultID(c context.Context, req *types.QueryStableVaultByVaultIDRequest) (*types.QueryStableVaultByVaultIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	stableMintData, found := q.GetStableMintVault(ctx, req.StableVaultId)
	if !found {
		return &types.QueryStableVaultByVaultIDResponse{}, nil
	}

	return &types.QueryStableVaultByVaultIDResponse{
		StableMintVault: &stableMintData,
	}, nil
}

func (q QueryServer) QueryStableVaultByApp(c context.Context, req *types.QueryStableVaultByAppRequest) (*types.QueryStableVaultByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx            = sdk.UnwrapSDKContext(c)
		stableMintData []types.StableMintVault
	)
	stableMint := q.GetStableMintVaults(ctx)

	for _, data := range stableMint {
		if data.AppId == req.AppId {
			stableMintData = append(stableMintData, data)
		}
	}

	return &types.QueryStableVaultByAppResponse{
		StableMintVault: stableMintData,
	}, nil
}

func (q QueryServer) QueryStableVaultByAppAndExtendedPair(c context.Context, req *types.QueryStableVaultByAppAndExtendedPairRequest) (*types.QueryStableVaultByAppAndExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx            = sdk.UnwrapSDKContext(c)
		stableMintData types.StableMintVault
	)
	stableMint := q.GetStableMintVaults(ctx)
	for _, data := range stableMint {
		if data.AppId == req.AppId && data.ExtendedPairVaultID == req.ExtendedPairId {
			stableMintData = data
		}
	}

	return &types.QueryStableVaultByAppAndExtendedPairResponse{
		StableMintVault: &stableMintData,
	}, nil
}

// QueryExtendedPairVaultMappingByAppAndExtendedPair to query vault by app and extended pair.
func (q QueryServer) QueryExtendedPairVaultMappingByAppAndExtendedPair(c context.Context, req *types.QueryExtendedPairVaultMappingByAppAndExtendedPairRequest) (*types.QueryExtendedPairVaultMappingByAppAndExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	_, found := q.asset.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMappingData(ctx, req.AppId, req.ExtendedPairId)
	if !found {
		return &types.QueryExtendedPairVaultMappingByAppAndExtendedPairResponse{}, nil
	}

	return &types.QueryExtendedPairVaultMappingByAppAndExtendedPairResponse{
		ExtendedPairVaultMapping: &appExtendedPairVaultData,
	}, nil
}

func (q QueryServer) QueryExtendedPairVaultMappingByApp(c context.Context, req *types.QueryExtendedPairVaultMappingByAppRequest) (*types.QueryExtendedPairVaultMappingByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	_, found := q.asset.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppMappingData(ctx, req.AppId)
	if !found {
		return &types.QueryExtendedPairVaultMappingByAppResponse{}, nil
	}

	return &types.QueryExtendedPairVaultMappingByAppResponse{
		ExtendedPairVaultMapping: appExtendedPairVaultData,
	}, nil
}

func (q QueryServer) QueryTVLByAppOfAllExtendedPairs(c context.Context, req *types.QueryTVLByAppOfAllExtendedPairsRequest) (*types.QueryTVLByAppOfAllExtendedPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx     = sdk.UnwrapSDKContext(c)
		tvlData []types.TvlLockedDataMap
	)

	appExtendedPairVaultData, found := q.GetAppMappingData(ctx, req.AppId)
	if !found {
		return &types.QueryTVLByAppOfAllExtendedPairsResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData {
		extPairVault, _ := q.asset.GetPairsVault(ctx, data.ExtendedPairId)
		pairID, _ := q.asset.GetPair(ctx, extPairVault.PairId)

		var tvl types.TvlLockedDataMap

		denom, found := q.asset.GetAsset(ctx, pairID.AssetIn)
		if !found {
			return nil, types.ErrorAssetDoesNotExist
		}
		tvl.AssetDenom = denom.Denom
		tvl.CollateralLockedAmount = data.CollateralLockedAmount

		tvlData = append(tvlData, tvl)
	}

	return &types.QueryTVLByAppOfAllExtendedPairsResponse{
		Tvldata: tvlData,
	}, nil
}

func (q QueryServer) QueryTVLByApp(c context.Context, req *types.QueryTVLByAppRequest) (*types.QueryTVLByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx    = sdk.UnwrapSDKContext(c)
		locked = sdk.ZeroDec()
	)
	_, found := q.asset.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppMappingData(ctx, req.AppId)
	if !found {
		return &types.QueryTVLByAppResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData {
		extPairVault, _ := q.asset.GetPairsVault(ctx, data.ExtendedPairId)
		pairID, _ := q.asset.GetPair(ctx, extPairVault.PairId)

		twaData, _ := q.oracle.CalcAssetPrice(ctx, pairID.AssetIn, data.CollateralLockedAmount)
		locked = twaData.Add(locked)
	}
	// locked = locked.Quo(sdk.NewInt(1000000))

	return &types.QueryTVLByAppResponse{
		CollateralLocked: locked.TruncateInt(),
	}, nil
}

func (q QueryServer) QueryUserMyPositionByApp(c context.Context, req *types.QueryUserMyPositionByAppRequest) (*types.QueryUserMyPositionByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx             = sdk.UnwrapSDKContext(c)
		vaultsIds       []uint64
		totalLocked     = sdk.ZeroDec()
		totalDue        = sdk.ZeroDec()
		availableBorrow = sdk.ZeroDec()
		averageCr       sdk.Dec
		totalCr         = sdk.ZeroDec()
	)

	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Address is not correct")
	}

	userVaultAssetData, found := q.GetUserAppMappingData(ctx, req.Owner, req.AppId)
	if !found {
		return &types.QueryUserMyPositionByAppResponse{}, nil
	}

	for _, data := range userVaultAssetData {
		vaultsIds = append(vaultsIds, data.VaultId)
	}
	count := len(vaultsIds)

	if count == 0 {
		return &types.QueryUserMyPositionByAppResponse{}, nil
	}

	for _, data := range vaultsIds {
		vault, found := q.GetVault(ctx, data)
		if !found {
			continue
		}

		extPairVault, _ := q.asset.GetPairsVault(ctx, vault.ExtendedPairVaultID)
		pairID, _ := q.asset.GetPair(ctx, extPairVault.PairId)
		assetOutData, found := q.asset.GetAsset(ctx, pairID.AssetOut)
		if !found {
			continue
		}

		assetInTotalPrice, _ := q.oracle.CalcAssetPrice(ctx, pairID.AssetIn, vault.AmountIn)
		var assetOutTotalPrice sdk.Dec
		totalLocked = assetInTotalPrice.Add(totalLocked)

		if extPairVault.AssetOutOraclePrice {
			assetOutTotalPrice, _ = q.oracle.CalcAssetPrice(ctx, pairID.AssetOut, vault.AmountOut)
		} else {
			assetOutTotalPrice = (sdk.NewDecFromInt(sdk.NewIntFromUint64(extPairVault.AssetOutPrice)).Mul(sdk.NewDecFromInt(vault.AmountOut))).Quo(sdk.NewDecFromInt(assetOutData.Decimals))
		}
		totalDue = assetOutTotalPrice.Add(totalDue)

		collaterlizationRatio, err := q.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, vault.AmountOut)
		if err != nil {
			return nil, err
		}

		totalCr = collaterlizationRatio.Add(totalCr)
		minCr := extPairVault.MinCr

		AmtIn := assetInTotalPrice
		AmtOut := assetOutTotalPrice

		av := AmtIn.Quo(minCr)
		av = av.Sub(AmtOut)

		availableBorrow = av.Quo(sdk.OneDec()).Add(availableBorrow)
	}

	// totalLocked = totalLocked.Quo(sdk.NewInt(1000000))
	// totalDue = totalDue.Quo(sdk.NewInt(1000000))
	// availableBorrow = availableBorrow.Quo(sdk.NewInt(1000000))
	t, _ := sdk.NewDecFromStr(strconv.Itoa(len(vaultsIds)))
	averageCr = totalCr.Quo(t)

	return &types.QueryUserMyPositionByAppResponse{
		CollateralLocked:  totalLocked.TruncateInt(),
		TotalDue:          totalDue.TruncateInt(),
		AvailableToBorrow: availableBorrow.TruncateInt(),
		AverageCrRatio:    averageCr,
	}, nil
}

func (q QueryServer) QueryUserExtendedPairTotalData(c context.Context, req *types.QueryUserExtendedPairTotalDataRequest) (*types.QueryUserExtendedPairTotalDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	userVaultAssetData := q.GetUserMappingData(ctx, req.Owner)

	return &types.QueryUserExtendedPairTotalDataResponse{
		UserTotalData: userVaultAssetData,
	}, nil
}

func (q QueryServer) QueryPairsLockedAndMintedStatisticByApp(c context.Context, req *types.QueryPairsLockedAndMintedStatisticByAppRequest) (*types.QueryPairsLockedAndMintedStatisticByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx            = sdk.UnwrapSDKContext(c)
		pairStatistics []types.PairStatisticData
	)

	appExtendedPairVaultData, found := q.GetAppMappingData(ctx, req.AppId)
	if !found {
		return &types.QueryPairsLockedAndMintedStatisticByAppResponse{}, nil
	}

	for _, data := range appExtendedPairVaultData {
		extPairVault, _ := q.asset.GetPairsVault(ctx, data.ExtendedPairId)
		pairID, _ := q.asset.GetPair(ctx, extPairVault.PairId)

		var statistics types.PairStatisticData
		inDenom, _ := q.asset.GetAsset(ctx, pairID.AssetIn)
		outDenom, _ := q.asset.GetAsset(ctx, pairID.AssetOut)

		statistics.AssetInDenom = inDenom.Denom
		statistics.AssetOutDenom = outDenom.Denom
		statistics.CollateralAmount = data.CollateralLockedAmount
		statistics.MintedAmount = data.TokenMintedAmount
		statistics.ExtendedPairVaultID = data.ExtendedPairId

		pairStatistics = append(pairStatistics, statistics)
	}

	return &types.QueryPairsLockedAndMintedStatisticByAppResponse{
		PairStatisticData: pairStatistics,
	}, nil
}
