package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/vault/types"
)

type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func NewQueryServer(k Keeper) types.QueryServer {
	return Querier{
		Keeper: k,
	}
}

func (q Querier) QueryAllVaults(c context.Context, req *types.QueryAllVaultsRequest) (*types.QueryAllVaultsResponse, error) {
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	vaults := q.GetVaults(ctx)

	return &types.QueryAllVaultsResponse{
		Vault: vaults,
	}, nil
}

func (q Querier) QueryAllVaultsByApp(c context.Context, req *types.QueryAllVaultsByAppRequest) (*types.QueryAllVaultsByAppResponse, error) {
	var (
		ctx       = sdk.UnwrapSDKContext(c)
		AppVaults []types.Vault
	)
	vaults := q.GetVaults(ctx)
	for _, data := range vaults {
		if data.AppMappingId == req.AppId {
			AppVaults = append(AppVaults, data)
		}
	}

	return &types.QueryAllVaultsByAppResponse{
		Vault: AppVaults,
	}, nil
}

func (q Querier) QueryVault(c context.Context, req *types.QueryVaultRequest) (*types.QueryVaultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	vault, found := q.GetVault(ctx, req.Id)
	if !found {
		return &types.QueryVaultResponse{}, nil
	}

	return &types.QueryVaultResponse{
		Vault: vault,
	}, nil
}
func (q Querier) QueryVaultInfoByVaultId(c context.Context, req *types.QueryVaultInfoByVaultIdRequest) (*types.QueryVaultInfoByVaultIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	vault, found := q.GetVault(ctx, req.Id)
	if !found {
		return &types.QueryVaultInfoByVaultIdResponse{}, nil
	}

	collateralizationRatio, err := q.CalculateCollaterlizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, vault.AmountOut)
	if err != nil {
		return nil, err
	}
	pairVaults, _ := q.GetPairsVault(ctx, vault.ExtendedPairVaultID)
	pairID, _ := q.GetPair(ctx, pairVaults.PairId)
	assetIn, _ := q.GetAsset(ctx, pairID.AssetIn)
	assetOut, _ := q.GetAsset(ctx, pairID.AssetOut)
	return &types.QueryVaultInfoByVaultIdResponse{
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

func (q Querier) QueryVaultInfoOfOwnerByApp(c context.Context, req *types.QueryVaultInfoOfOwnerByAppRequest) (*types.QueryVaultInfoOfOwnerByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	// nolint
	var (
		ctx        = sdk.UnwrapSDKContext(c)
		vaultsIds  []string
		vaultsInfo []types.VaultInfo
	)
	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Address is not correct")
	}

	userVaultAssetData, _ := q.GetUserVaultExtendedPairMapping(ctx, req.Owner)

	for _, data := range userVaultAssetData.UserVaultApp {
		if data.AppMappingId == req.AppId {
			for _, inData := range data.UserExtendedPairVault {
				vaultsIds = append(vaultsIds, inData.VaultId)
			}
		}
	}
	var count = len(vaultsIds)
	for _, id := range vaultsIds {
		vault, found := q.GetVault(ctx, id)
		if !found {
			count--
			continue
		}

		collateralizationRatio, err := q.CalculateCollaterlizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, vault.AmountOut)
		if err != nil {
			return nil, err
		}
		pairVaults, _ := q.GetPairsVault(ctx, vault.ExtendedPairVaultID)
		pairID, _ := q.GetPair(ctx, pairVaults.PairId)
		assetIn, _ := q.GetAsset(ctx, pairID.AssetIn)
		assetOut, _ := q.GetAsset(ctx, pairID.AssetOut)

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

func (q Querier) QueryAllVaultsByAppAndExtendedPair(c context.Context, req *types.QueryAllVaultsByAppAndExtendedPairRequest) (*types.QueryAllVaultsByAppAndExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx       = sdk.UnwrapSDKContext(c)
		vaultList []types.Vault
	)

	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	_, nfound := q.GetPairsVault(ctx, req.ExtendedPairId)
	if !nfound {
		return &types.QueryAllVaultsByAppAndExtendedPairResponse{}, nil
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

func (q Querier) QueryVaultIdOfOwnerByExtendedPairAndApp(c context.Context, req *types.QueryVaultIdOfOwnerByExtendedPairAndAppRequest) (*types.QueryVaultIdOfOwnerByExtendedPairAndAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx     = sdk.UnwrapSDKContext(c)
		vaultID = ""
	)

	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Address is not correct")
	}

	_, nfound := q.GetPairsVault(ctx, req.ExtendedPairId)
	if !nfound {
		return &types.QueryVaultIdOfOwnerByExtendedPairAndAppResponse{}, nil
	}

	vaultData := q.GetVaults(ctx)

	for _, data := range vaultData {
		if data.AppMappingId == req.AppId && data.ExtendedPairVaultID == req.ExtendedPairId && data.Owner == req.Owner {
			vaultID = data.Id
		}
	}

	return &types.QueryVaultIdOfOwnerByExtendedPairAndAppResponse{
		Vault_Id: vaultID,
	}, nil
}

func (q Querier) QueryVaultIdsByAppInAllExtendedPairs(c context.Context, req *types.QueryVaultIdsByAppInAllExtendedPairsRequest) (*types.QueryVaultIdsByAppInAllExtendedPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx       = sdk.UnwrapSDKContext(c)
		vaultsIds []string
	)

	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}
	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
	if !found {
		return &types.QueryVaultIdsByAppInAllExtendedPairsResponse{}, nil
	}

	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		vaultsIds = append(vaultsIds, data.VaultIds...)
	}

	return &types.QueryVaultIdsByAppInAllExtendedPairsResponse{
		VaultIds: vaultsIds,
	}, nil
}

func (q Querier) QueryAllVaultIdsByAnOwner(c context.Context, req *types.QueryAllVaultIdsByAnOwnerRequest) (*types.QueryAllVaultIdsByAnOwnerResponse, error) {
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

	userVaultAssetData, _ := q.GetUserVaultExtendedPairMapping(ctx, req.Owner)

	for _, data := range userVaultAssetData.UserVaultApp {
		for _, inData := range data.UserExtendedPairVault {
			vaultsIds = append(vaultsIds, inData.VaultId)
		}
	}

	return &types.QueryAllVaultIdsByAnOwnerResponse{
		VaultIds: vaultsIds,
	}, nil
}

func (q Querier) QueryTokenMintedByAppAndExtendedPair(c context.Context, req *types.QueryTokenMintedByAppAndExtendedPairRequest) (*types.QueryTokenMintedByAppAndExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx         = sdk.UnwrapSDKContext(c)
		tokenMinted = sdk.ZeroInt()
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	_, nfound := q.GetPairsVault(ctx, req.ExtendedPairId)
	if !nfound {
		return nil, status.Errorf(codes.NotFound, "extended pair does not exist for id %d", req.ExtendedPairId)
	}
	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Pair vault does not exist for App id %d", req.AppId)
	}

	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		if data.ExtendedPairId == req.ExtendedPairId {
			tokenMinted = data.TokenMintedAmount
		}
	}

	return &types.QueryTokenMintedByAppAndExtendedPairResponse{
		TokenMinted: tokenMinted,
	}, nil
}

func (q Querier) QueryTokenMintedAssetWiseByApp(c context.Context, req *types.QueryTokenMintedAssetWiseByAppRequest) (*types.QueryTokenMintedAssetWiseByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	// nolint
	var (
		ctx        = sdk.UnwrapSDKContext(c)
		mintedData []types.MintedDataMap
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
	if !found {
		return &types.QueryTokenMintedAssetWiseByAppResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		extPairVault, _ := q.GetPairsVault(ctx, data.ExtendedPairId)
		pairID, _ := q.GetPair(ctx, extPairVault.PairId)

		var minted types.MintedDataMap

		denom, found := q.GetAsset(ctx, pairID.AssetOut)
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

func (q Querier) QueryVaultCountByApp(c context.Context, req *types.QueryVaultCountByAppRequest) (*types.QueryVaultCountByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx   = sdk.UnwrapSDKContext(c)
		count uint64
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
	if !found {
		return &types.QueryVaultCountByAppResponse{}, nil
	}

	count = appExtendedPairVaultData.Counter

	return &types.QueryVaultCountByAppResponse{
		VaultCount: count,
	}, nil
}

func (q Querier) QueryVaultCountByAppAndExtendedPair(c context.Context, req *types.QueryVaultCountByAppAndExtendedPairRequest) (*types.QueryVaultCountByAppAndExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx          = sdk.UnwrapSDKContext(c)
		count uint64 = 0
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
	if !found {
		return &types.QueryVaultCountByAppAndExtendedPairResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		if data.ExtendedPairId == req.ExtendedPairId {
			count = uint64(len(data.VaultIds))
		}
	}

	return &types.QueryVaultCountByAppAndExtendedPairResponse{
		VaultCount: count,
	}, nil
}

func (q Querier) QueryTotalValueLockedByAppExtendedPair(c context.Context, req *types.QueryTotalValueLockedByAppExtendedPairRequest) (*types.QueryTotalValueLockedByAppExtendedPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx         = sdk.UnwrapSDKContext(c)
		valueLocked = sdk.ZeroInt()
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}
	_, nfound := q.GetPairsVault(ctx, req.ExtendedPairId)
	if !nfound {
		return &types.QueryTotalValueLockedByAppExtendedPairResponse{}, nil
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
	if !found {
		return &types.QueryTotalValueLockedByAppExtendedPairResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		if data.ExtendedPairId == req.ExtendedPairId {
			valueLocked = data.CollateralLockedAmount
		}
	}

	return &types.QueryTotalValueLockedByAppExtendedPairResponse{
		ValueLocked: &valueLocked,
	}, nil
}

func (q Querier) QueryExtendedPairIDsByApp(c context.Context, req *types.QueryExtendedPairIDsByAppRequest) (*types.QueryExtendedPairIDsByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	// nolint
	var (
		ctx     = sdk.UnwrapSDKContext(c)
		pairIDs []uint64
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
	if !found {
		return &types.QueryExtendedPairIDsByAppResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		pairIDs = append(pairIDs, data.ExtendedPairId)
	}

	return &types.QueryExtendedPairIDsByAppResponse{
		ExtendedPairIds: pairIDs,
	}, nil
}

func (q Querier) QueryStableVaultByVaultId(c context.Context, req *types.QueryStableVaultByVaultIdRequest) (*types.QueryStableVaultByVaultIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	stableMintData, found := q.GetStableMintVault(ctx, req.StableVaultId)
	if !found {
		return &types.QueryStableVaultByVaultIdResponse{}, nil
	}

	return &types.QueryStableVaultByVaultIdResponse{
		StableMintVault: &stableMintData,
	}, nil
}

func (q Querier) QueryStableVaultByApp(c context.Context, req *types.QueryStableVaultByAppRequest) (*types.QueryStableVaultByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx            = sdk.UnwrapSDKContext(c)
		stableMintData []types.StableMintVault
	)
	stableMint := q.GetStableMintVaults(ctx)

	for _, data := range stableMint {
		if data.AppMappingId == req.AppId {
			stableMintData = append(stableMintData, data)
		}
	}

	return &types.QueryStableVaultByAppResponse{
		StableMintVault: stableMintData,
	}, nil
}

func (q Querier) QueryStableVaultByAppExtendedPair(c context.Context, req *types.QueryStableVaultByAppExtendedPairRequest) (*types.QueryStableVaultByAppExtendedPairResponse, error) {
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

	return &types.QueryStableVaultByAppExtendedPairResponse{
		StableMintVault: &stableMintData,
	}, nil
}

// nolint
func (q Querier) QueryExtendedPairVaultMappingByAppAndExtendedPairId(c context.Context, req *types.QueryExtendedPairVaultMappingByAppAndExtendedPairIdRequest) (*types.QueryExtendedPairVaultMappingByAppAndExtendedPairIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx                = sdk.UnwrapSDKContext(c)
		extendedPairIDData types.ExtendedPairVaultMapping
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
	if !found {
		return &types.QueryExtendedPairVaultMappingByAppAndExtendedPairIdResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		if data.ExtendedPairId == req.ExtendedPairId {
			extendedPairIDData = *data
		}
	}

	return &types.QueryExtendedPairVaultMappingByAppAndExtendedPairIdResponse{
		ExtendedPairVaultMapping: &extendedPairIDData,
	}, nil
}

func (q Querier) QueryExtendedPairVaultMappingByApp(c context.Context, req *types.QueryExtendedPairVaultMappingByAppRequest) (*types.QueryExtendedPairVaultMappingByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx     = sdk.UnwrapSDKContext(c)
		pairIDs []*types.ExtendedPairVaultMapping
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
	if !found {
		return &types.QueryExtendedPairVaultMappingByAppResponse{}, nil
	}
	pairIDs = append(pairIDs, appExtendedPairVaultData.ExtendedPairVaults...)

	return &types.QueryExtendedPairVaultMappingByAppResponse{
		ExtendedPairVaultMapping: pairIDs,
	}, nil
}

func (q Querier) QueryTVLLockedByAppOfAllExtendedPairs(c context.Context, req *types.QueryTVLLockedByAppOfAllExtendedPairsRequest) (*types.QueryTVLLockedByAppOfAllExtendedPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	// nolint
	var (
		ctx     = sdk.UnwrapSDKContext(c)
		tvlData []types.TvlLockedDataMap
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
	if !found {
		return &types.QueryTVLLockedByAppOfAllExtendedPairsResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		extPairVault, _ := q.GetPairsVault(ctx, data.ExtendedPairId)
		pairID, _ := q.GetPair(ctx, extPairVault.PairId)

		var tvl types.TvlLockedDataMap

		denom, found := q.GetAsset(ctx, pairID.AssetIn)
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

func (q Querier) QueryTotalTVLByApp(c context.Context, req *types.QueryTotalTVLByAppRequest) (*types.QueryTotalTVLByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx    = sdk.UnwrapSDKContext(c)
		locked = sdk.ZeroInt()
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "App does not exist for id %d", req.AppId)
	}

	appExtendedPairVaultData, found := q.GetAppExtendedPairVaultMapping(ctx, req.AppId)
	if !found {
		return &types.QueryTotalTVLByAppResponse{}, nil
	}
	for _, data := range appExtendedPairVaultData.ExtendedPairVaults {
		extPairVault, _ := q.GetPairsVault(ctx, data.ExtendedPairId)
		pairID, _ := q.GetPair(ctx, extPairVault.PairId)

		rate, _ := q.GetPriceForAsset(ctx, pairID.AssetIn)
		locked = data.CollateralLockedAmount.Mul(sdk.NewIntFromUint64(rate)).Add(locked)
	}
	locked = locked.Quo(sdk.NewInt(1000000))

	return &types.QueryTotalTVLByAppResponse{
		CollateralLocked: locked,
	}, nil
}

func (q Querier) QueryUserMyPositionByApp(c context.Context, req *types.QueryUserMyPositionByAppRequest) (*types.QueryUserMyPositionByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx             = sdk.UnwrapSDKContext(c)
		vaultsIds       []string
		totalLocked     = sdk.ZeroInt()
		totalDue        = sdk.ZeroInt()
		availableBorrow = sdk.ZeroInt()
		averageCr       sdk.Dec
		totalCr         = sdk.ZeroDec()
	)

	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Address is not correct")
	}

	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return &types.QueryUserMyPositionByAppResponse{}, nil
	}

	userVaultAssetData, found := q.GetUserVaultExtendedPairMapping(ctx, req.Owner)
	if !found {
		return &types.QueryUserMyPositionByAppResponse{}, nil
	}

	for _, data := range userVaultAssetData.UserVaultApp {
		if data.AppMappingId == req.AppId {
			for _, inData := range data.UserExtendedPairVault {
				vaultsIds = append(vaultsIds, inData.VaultId)
			}
		}
	}
	var count = len(vaultsIds)

	for _, data := range vaultsIds {
		vault, found := q.GetVault(ctx, data)
		if !found {
			count--
			continue
		}

		extPairVault, _ := q.GetPairsVault(ctx, vault.ExtendedPairVaultID)
		pairID, _ := q.GetPair(ctx, extPairVault.PairId)

		assetInPrice, _ := q.GetPriceForAsset(ctx, pairID.AssetIn)
		var assetOutPrice uint64
		totalLocked = vault.AmountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).Add(totalLocked)

		if extPairVault.AssetOutOraclePrice {
			assetOutPrice, _ = q.GetPriceForAsset(ctx, pairID.AssetOut)
		} else {
			assetOutPrice = extPairVault.AssetOutPrice
		}
		totalDue = vault.AmountOut.Mul(sdk.NewIntFromUint64(assetOutPrice)).Add(totalDue)

		collaterlizationRatio, err := q.CalculateCollaterlizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, vault.AmountOut)
		if err != nil {
			return nil, err
		}

		totalCr = collaterlizationRatio.Add(totalCr)
		var minCr = extPairVault.MinCr

		AmtIn := vault.AmountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).ToDec()
		AmtOut := vault.AmountOut.Mul(sdk.NewIntFromUint64(assetOutPrice)).ToDec()

		av := sdk.Int(AmtIn.Quo(minCr))
		av = av.Sub(sdk.Int(AmtOut))

		availableBorrow = av.Quo(sdk.Int(sdk.OneDec())).Add(availableBorrow)
	}
	if count == 0 {
		return &types.QueryUserMyPositionByAppResponse{}, nil
	}
	totalLocked = totalLocked.Quo(sdk.NewInt(1000000))
	totalDue = totalDue.Quo(sdk.NewInt(1000000))
	availableBorrow = availableBorrow.Quo(sdk.NewInt(1000000))
	t, _ := sdk.NewDecFromStr(strconv.Itoa(len(vaultsIds)))
	averageCr = totalCr.Quo(t)

	return &types.QueryUserMyPositionByAppResponse{
		CollateralLocked:  totalLocked,
		TotalDue:          totalDue,
		AvailableToBorrow: availableBorrow,
		AverageCrRatio:    averageCr,
	}, nil
}

func (q Querier) QueryUserExtendedPairTotalData(c context.Context, req *types.QueryUserExtendedPairTotalDataRequest) (*types.QueryUserExtendedPairTotalDataResponse, error) {
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	userVaultAssetData, found := q.GetUserVaultExtendedPairMapping(ctx, req.Owner)
	if !found {
		return &types.QueryUserExtendedPairTotalDataResponse{}, nil
	}

	return &types.QueryUserExtendedPairTotalDataResponse{
		UserTotalData: &userVaultAssetData,
	}, nil
}
