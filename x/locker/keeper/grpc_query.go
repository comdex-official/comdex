package keeper

import (
	"context"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (q QueryServer) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	var (
		ctx    = sdk.UnwrapSDKContext(c)
		params = q.GetParams(ctx)
	)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

func (q QueryServer) QueryLockerInfo(c context.Context, req *types.QueryLockerInfoRequest) (*types.QueryLockerInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := q.GetLocker(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "locker-info does not exist for id %d", req.Id)
	}

	return &types.QueryLockerInfoResponse{
		LockerInfo: item,
	}, nil
}

func (q QueryServer) QueryLockersByAppToAssetID(c context.Context, request *types.QueryLockersByAppToAssetIDRequest) (*types.QueryLockersByAppToAssetIDResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	lockerLookupData, found := q.GetLockerLookupTable(ctx, request.AppId, request.AssetId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exists appID %d", request.AppId)
	}

	return &types.QueryLockersByAppToAssetIDResponse{
		LockerIds: lockerLookupData.LockerIds,
	}, nil
}

func (q QueryServer) QueryLockerInfoByAppID(c context.Context, request *types.QueryLockerInfoByAppIDRequest) (*types.QueryLockerInfoByAppIDResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	app, found := q.GetApp(ctx, request.AppId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.AppId)
	}

	lockerLookupData, found := q.GetLockerLookupTableByApp(ctx, app.Id)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", app.Id)
	}

	var lockerIds []uint64
	for _, locker := range lockerLookupData {
		lockerIds = append(lockerIds, locker.LockerIds...)
	}
	return &types.QueryLockerInfoByAppIDResponse{
		LockerIds: lockerIds,
	}, nil
}

func (q QueryServer) QueryTotalDepositByAppAndAssetID(c context.Context, request *types.QueryTotalDepositByAppAndAssetIDRequest) (*types.QueryTotalDepositByAppAndAssetIDResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	lockerLookupData, found := q.GetLockerLookupTable(ctx, request.AppId, request.AssetId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", request.AppId)
	}

	return &types.QueryTotalDepositByAppAndAssetIDResponse{
		TotalDeposit: lockerLookupData.DepositedAmount.Uint64(),
	}, nil
}

func (q QueryServer) QueryLockerByAppByOwner(c context.Context,
	request *types.QueryLockerByAppByOwnerRequest) (*types.QueryLockerByAppByOwnerResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	lockerLookupData, found := q.GetUserLockerAppMapping(ctx, request.Owner, request.AppId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", request.AppId)
	}

	var lockerInfos []types.Locker
	for _, locker := range lockerLookupData {
		locker1, _ := q.GetLocker(ctx, locker.LockerId)
		lockerInfos = append(lockerInfos, locker1)
	}
	return &types.QueryLockerByAppByOwnerResponse{
		LockerInfo: lockerInfos,
	}, nil
}

func (q QueryServer) QueryOwnerLockerByAppIDbyOwner(c context.Context, request *types.QueryOwnerLockerByAppIDbyOwnerRequest) (*types.QueryOwnerLockerByAppIDbyOwnerResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	lockerLookupData, found := q.GetUserLockerAppMapping(ctx, request.Owner, request.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", request.AppId)
	}

	var lockerIds []uint64
	for _, locker := range lockerLookupData {
		lockerIds = append(lockerIds, locker.LockerId)
	}

	return &types.QueryOwnerLockerByAppIDbyOwnerResponse{
		LockerIds: lockerIds,
	}, nil
}

func (q QueryServer) QueryOwnerLockerOfAllAppsByOwner(c context.Context, request *types.QueryOwnerLockerOfAllAppsByOwnerRequest) (*types.QueryOwnerLockerOfAllAppsByOwnerResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	userLockerLookupData, _ := q.GetUserLockerMapping(ctx, request.Owner)

	var lockerIds []uint64
	for _, locker := range userLockerLookupData{
		lockerIds = append(lockerIds, locker.LockerId)
	}

	return &types.QueryOwnerLockerOfAllAppsByOwnerResponse{
		LockerIds: lockerIds,
	}, nil
}

func (q QueryServer) QueryOwnerTxDetailsLockerOfAppByOwnerByAsset(c context.Context, request *types.QueryOwnerTxDetailsLockerOfAppByOwnerByAssetRequest) (*types.QueryOwnerTxDetailsLockerOfAppByOwnerByAssetResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx        = sdk.UnwrapSDKContext(c)
	)
	userLockerLookupData, found := q.GetUserLockerAssetMapping(ctx, request.Owner, request.AppId, request.AssetId)
	if !found {
		return &types.QueryOwnerTxDetailsLockerOfAppByOwnerByAssetResponse{}, nil
	}
	
	return &types.QueryOwnerTxDetailsLockerOfAppByOwnerByAssetResponse{
		UserTxData: userLockerLookupData.UserData,
	}, nil
}

func (q QueryServer) QueryOwnerLockerByAppToAssetIDbyOwner(c context.Context, request *types.QueryOwnerLockerByAppToAssetIDbyOwnerRequest) (*types.QueryOwnerLockerByAppToAssetIDbyOwnerResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	lockerLookupData, found := q.GetUserLockerAssetMapping(ctx, request.Owner, request.AppId, request.AssetId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", request.AppId)
	}

	lockerData, _ := q.GetLocker(ctx, lockerLookupData.LockerId)
	
		
	return &types.QueryOwnerLockerByAppToAssetIDbyOwnerResponse{
		LockerInfo: lockerData,
	}, nil
}

func (q QueryServer) QueryLockerCountByAppID(c context.Context, request *types.QueryLockerCountByAppIDRequest) (*types.QueryLockerCountByAppIDResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	app, found := q.GetApp(ctx, request.AppId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.AppId)
	}

	lockerLookupData, found := q.GetLockerLookupTableByApp(ctx, app.Id)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", app.Id)
	}

	var lockerCount int
	for _, locker := range lockerLookupData {
		lockerCount = lockerCount + len(locker.LockerIds)
	}
	return &types.QueryLockerCountByAppIDResponse{
		TotalCount: uint64(lockerCount),
	}, nil
}

func (q QueryServer) QueryLockerCountByAppToAssetID(c context.Context, request *types.QueryLockerCountByAppToAssetIDRequest) (*types.QueryLockerCountByAppToAssetIDResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	lockerLookupData, found := q.GetLockerLookupTable(ctx, request.AppId, request.AssetId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", request.AppId)
	}

	lockerCount := len(lockerLookupData.LockerIds)
	
	return &types.QueryLockerCountByAppToAssetIDResponse{
		TotalCount: uint64(lockerCount),
	}, nil
}

func (q QueryServer) QueryWhiteListedAssetIDsByAppID(c context.Context, request *types.QueryWhiteListedAssetIDsByAppIDRequest) (*types.QueryWhiteListedAssetIDsByAppIDResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
		assetIds   []uint64  
	)

	lockerLookupData, found := q.GetLockerProductAssetMappingByApp(ctx, request.AppId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", request.AppId)
	}

	for _, data := range lockerLookupData {
		assetIds = append(assetIds, data.AssetId)
	}

	return &types.QueryWhiteListedAssetIDsByAppIDResponse{
		AssetIds: assetIds,
	}, nil
}

func (q QueryServer) QueryWhiteListedAssetByAllApps(c context.Context, request *types.QueryWhiteListedAssetByAllAppsRequest) (*types.QueryWhiteListedAssetByAllAppsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	apps, found := q.asset.GetApps(ctx)

	if !found {
		return nil, status.Errorf(codes.NotFound, "ano apps exist")
	}

	var productToAll []types.AppToAllAsset

	for _, app := range apps {
		var product types.AppToAllAsset
		var assets []assettypes.Asset
		appData, _ := q.GetLockerProductAssetMappingByApp(ctx, app.Id)
		for _, assetID := range appData {
			asset, assetFound := q.asset.GetAsset(ctx, assetID.AssetId)
			if assetFound {
				assets = append(assets, asset)
			}
			product = types.AppToAllAsset{
				AppId:  app.Id,
				Assets: assets,
			}
		}

		productToAll = append(productToAll, product)
	}
	return &types.QueryWhiteListedAssetByAllAppsResponse{
		ProductToAllAsset: productToAll,
	}, nil
}

func (q QueryServer) QueryLockerLookupTableByApp(c context.Context, req *types.QueryLockerLookupTableByAppRequest) (*types.QueryLockerLookupTableByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := q.GetLockerLookupTableByApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "locker-info does not exist for id %d", req.AppId)
	}

	return &types.QueryLockerLookupTableByAppResponse{
		TokenToLockerMapping: item,
	}, nil
}

func (q QueryServer) QueryLockerLookupTableByAppAndAssetID(c context.Context, req *types.QueryLockerLookupTableByAppAndAssetIDRequest) (*types.QueryLockerLookupTableByAppAndAssetIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := q.GetLockerLookupTable(ctx, req.AppId, req.AssetId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "locker-info does not exist for id %d", req.AppId)
	}

	return &types.QueryLockerLookupTableByAppAndAssetIDResponse{
		TokenToLockerMapping: &item,
	}, nil
}

func (q QueryServer) QueryLockerTotalDepositedByApp(c context.Context, req *types.QueryLockerTotalDepositedByAppRequest) (*types.QueryLockerTotalDepositedByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := q.GetLockerLookupTableByApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "locker-info does not exist for id %d", req.AppId)
	}
	var lockedDepositedAmt []types.LockedDepositedAmountDataMap

	for _, data := range item {
		var lockeddata types.LockedDepositedAmountDataMap
		lockeddata.AssetId = data.AssetId
		lockeddata.DepositedAmount = data.DepositedAmount
		lockedDepositedAmt = append(lockedDepositedAmt, lockeddata)
	}

	return &types.QueryLockerTotalDepositedByAppResponse{
		LockedDepositedAmountDataMap: lockedDepositedAmt,
	}, nil
}

func (q QueryServer) QueryState(c context.Context, req *types.QueryStateRequest) (*types.QueryStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	qs, _ := QueryState(req.Address, req.Denom, req.Height, req.Target)

	return &types.QueryStateResponse{
		Amount: *qs,
	}, nil
}

func (q QueryServer) QueryLockerTotalRewardsByAssetAppWise(c context.Context, request *types.QueryLockerTotalRewardsByAssetAppWiseRequest) (*types.QueryLockerTotalRewardsByAssetAppWiseResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	_, found := q.GetApp(ctx, request.AppId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.AppId)
	}

	rewardsData, found := q.GetLockerTotalRewardsByAssetAppWise(ctx, request.AppId, request.AssetId)
	if !found {
		return &types.QueryLockerTotalRewardsByAssetAppWiseResponse{}, nil
	}
	return &types.QueryLockerTotalRewardsByAssetAppWiseResponse{
		TotalRewards: rewardsData,
	}, nil
}
