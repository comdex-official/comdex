package keeper

import (
	"context"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = (*queryServer)(nil)

type queryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServer{
		Keeper: k,
	}
}

func (q *queryServer) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	var (
		ctx    = sdk.UnwrapSDKContext(c)
		params = q.GetParams(ctx)
	)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

func (q *queryServer) QueryLockerInfo(c context.Context, req *types.QueryLockerInfoRequest) (*types.QueryLockerInfoResponse, error) {
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

func (q *queryServer) QueryLockersByProductToAssetID(c context.Context, request *types.QueryLockersByProductToAssetIDRequest) (*types.QueryLockersByProductToAssetIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	app, found := q.GetApp(ctx, request.ProductId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.ProductId)
	}

	lockerLookupData, found := q.GetLockerLookupTable(ctx, app.Id)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", app.Id)
	}

	var lockerIds []string
	for _, locker := range lockerLookupData.Lockers {
		if request.AssetId == locker.AssetId {
			lockerIds = append(lockerIds, locker.LockerIds...)
		}
	}
	return &types.QueryLockersByProductToAssetIDResponse{
		LockerIds: lockerIds,
	}, nil
}

func (q *queryServer) QueryLockerInfoByProductID(c context.Context, request *types.QueryLockerInfoByProductIDRequest) (*types.QueryLockerInfoByProductIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	app, found := q.GetApp(ctx, request.ProductId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.ProductId)
	}

	lockerLookupData, found := q.GetLockerLookupTable(ctx, app.Id)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", app.Id)
	}

	var lockerIds []string
	for _, locker := range lockerLookupData.Lockers {
		lockerIds = locker.LockerIds
	}
	return &types.QueryLockerInfoByProductIDResponse{
		LockerIds: lockerIds,
	}, nil

}

func (q *queryServer) QueryTotalDepositByProductAssetID(c context.Context, request *types.QueryTotalDepositByProductAssetIDRequest) (*types.QueryTotalDepositByProductAssetIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	app, found := q.GetApp(ctx, request.ProductId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.ProductId)
	}

	lockerLookupData, found := q.GetLockerLookupTable(ctx, app.Id)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", app.Id)
	}

	var totalDeposit uint64
	for _, locker := range lockerLookupData.Lockers {

		if request.AssetId == locker.AssetId {
			totalDeposit += locker.DepositedAmount.Uint64()
		}
	}
	return &types.QueryTotalDepositByProductAssetIDResponse{
		TotalDeposit: totalDeposit,
	}, nil
}

func (q *queryServer) QueryLockerByProductByOwner(c context.Context,
	request *types.QueryLockerByProductByOwnerRequest) (*types.QueryLockerByProductByOwnerResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	app, found := q.GetApp(ctx, request.ProductId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.ProductId)
	}

	lockerLookupData, found := q.GetLockerLookupTable(ctx, app.Id)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", app.Id)
	}

	var lockerInfos []types.Locker
	for _, locker := range lockerLookupData.Lockers {
		for _, lockerID := range locker.LockerIds {
			locker1, _ := q.GetLocker(ctx, lockerID)
			if request.Owner == locker1.Depositor {
				lockerInfos = append(lockerInfos, locker1)
			}
		}
	}
	return &types.QueryLockerByProductByOwnerResponse{
		LockerInfo: lockerInfos,
	}, nil
}

func (q *queryServer) QueryOwnerLockerByProductIDbyOwner(c context.Context, request *types.QueryOwnerLockerByProductIDbyOwnerRequest) (*types.QueryOwnerLockerByProductIDbyOwnerResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	app, found := q.GetApp(ctx, request.ProductId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.ProductId)
	}

	lockerLookupData, _ := q.GetLockerLookupTable(ctx, app.Id)

	var lockerIds []string
	for _, locker := range lockerLookupData.Lockers {
		for _, lockerID := range locker.LockerIds {
			locker1, _ := q.GetLocker(ctx, lockerID)
			if request.Owner == locker1.Depositor {
				lockerIds = append(lockerIds, locker1.LockerId)
			}
		}
	}

	return &types.QueryOwnerLockerByProductIDbyOwnerResponse{
		LockerIds: lockerIds,
	}, nil
}

func (q *queryServer) QueryOwnerLockerOfAllProductByOwner(c context.Context, request *types.QueryOwnerLockerOfAllProductByOwnerRequest) (*types.QueryOwnerLockerOfAllProductByOwnerResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	userLockerLookupData, _ := q.GetUserLockerAssetMapping(ctx, request.Owner)

	var lockerIds []string
	for _, locker := range userLockerLookupData.LockerAppMapping {
		for _, data := range locker.UserAssetLocker {
			lockerIds = append(lockerIds, data.LockerId)
		}
	}

	return &types.QueryOwnerLockerOfAllProductByOwnerResponse{
		LockerIds: lockerIds,
	}, nil
}

func (q *queryServer) QueryOwnerTxDetailsLockerOfProductByOwnerByAsset(c context.Context, request *types.QueryOwnerTxDetailsLockerOfProductByOwnerByAssetRequest) (*types.QueryOwnerTxDetailsLockerOfProductByOwnerByAssetResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx        = sdk.UnwrapSDKContext(c)
		userTxData []*types.UserTxData
	)
	userlockerLookupData, _ := q.GetUserLockerAssetMapping(ctx, request.Owner)
	if userlockerLookupData.Owner == request.Owner {
		for _, locker := range userlockerLookupData.LockerAppMapping {
			if locker.AppMappingId == request.ProductId {
				for _, data := range locker.UserAssetLocker {
					if data.AssetId == request.AssetId{
						userTxData = append(userTxData, data.UserData...)
					}
				}
			}
		}
	}
	return &types.QueryOwnerTxDetailsLockerOfProductByOwnerByAssetResponse{
		UserTxData: userTxData,
	}, nil
}

func (q *queryServer) QueryOwnerLockerByProductToAssetIDbyOwner(c context.Context, request *types.QueryOwnerLockerByProductToAssetIDbyOwnerRequest) (*types.QueryOwnerLockerByProductToAssetIDbyOwnerResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	app, found := q.GetApp(ctx, request.ProductId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.ProductId)
	}

	lockerLookupData, found := q.GetLockerLookupTable(ctx, app.Id)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", app.Id)
	}

	var lockerInfos []types.Locker
	for _, locker := range lockerLookupData.Lockers {

		if request.AssetId == locker.AssetId {
			for _, lockerID := range locker.LockerIds {
				locker1, _ := q.GetLocker(ctx, lockerID)
				if request.Owner == locker1.Depositor {
					lockerInfos = append(lockerInfos, locker1)
				}
			}
		}
	}
	return &types.QueryOwnerLockerByProductToAssetIDbyOwnerResponse{
		LockerInfo: lockerInfos,
	}, nil
}

func (q *queryServer) QueryLockerCountByProductID(c context.Context, request *types.QueryLockerCountByProductIDRequest) (*types.QueryLockerCountByProductIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	app, found := q.GetApp(ctx, request.ProductId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.ProductId)
	}

	lockerLookupData, found := q.GetLockerLookupTable(ctx, app.Id)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", app.Id)
	}

	var lockerCount uint64
	for _, locker := range lockerLookupData.Lockers {
		for _, lockerID := range locker.LockerIds {
			_, lockerFound := q.GetLocker(ctx, lockerID)
			if lockerFound {
				lockerCount += 1
			}
		}
	}
	return &types.QueryLockerCountByProductIDResponse{
		TotalCount: lockerCount,
	}, nil

}

func (q *queryServer) QueryLockerCountByProductToAssetID(c context.Context, request *types.QueryLockerCountByProductToAssetIDRequest) (*types.QueryLockerCountByProductToAssetIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	app, found := q.GetApp(ctx, request.ProductId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.ProductId)
	}

	lockerLookupData, found := q.GetLockerLookupTable(ctx, app.Id)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", app.Id)
	}

	var lockerCount uint64
	for _, locker := range lockerLookupData.Lockers {
		if request.AssetId == locker.AssetId {
			for _, lockerID := range locker.LockerIds {
				_, lockerFound := q.GetLocker(ctx, lockerID)
				if lockerFound {
					lockerCount += 1
				}
			}
		}
	}
	return &types.QueryLockerCountByProductToAssetIDResponse{
		TotalCount: lockerCount,
	}, nil
}

func (q *queryServer) QueryWhiteListedAssetIDsByProductID(c context.Context, request *types.QueryWhiteListedAssetIDsByProductIDRequest) (*types.QueryWhiteListedAssetIDsByProductIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	app, found := q.GetApp(ctx, request.ProductId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "app does not exist for appID %d", request.ProductId)
	}

	lockerLookupData, found := q.GetLockerLookupTable(ctx, app.Id)

	if !found {
		return nil, status.Errorf(codes.NotFound, "no asset exists appID %d", app.Id)
	}

	var assetIds []uint64
	for _, locker := range lockerLookupData.Lockers {
		assetIds = append(assetIds, locker.AssetId)
	}

	return &types.QueryWhiteListedAssetIDsByProductIDResponse{
		AssetIds: assetIds,
	}, nil
}

func (q *queryServer) QueryWhiteListedAssetByAllProduct(c context.Context, request *types.QueryWhiteListedAssetByAllProductRequest) (*types.QueryWhiteListedAssetByAllProductResponse, error) {

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

	var productToAll []types.ProductToAllAsset

	for _, app := range apps {
		var product types.ProductToAllAsset
		var assets []assettypes.Asset
		appData, _ := q.GetLockerProductAssetMapping(ctx, app.Id)
		for _, assetId := range appData.AssetIds {
			asset, assetFound := q.asset.GetAsset(ctx, assetId)
			if assetFound {
				assets = append(assets, asset)
			}
			product = types.ProductToAllAsset{
				ProductId: appData.AppMappingId,
				Assets:    assets,
			}
		}

		productToAll = append(productToAll, product)
	}
	return &types.QueryWhiteListedAssetByAllProductResponse{
		ProductToAllAsset: productToAll,
	}, nil
}

func (q *queryServer) QueryLockerLookupTableByApp(c context.Context, req *types.QueryLockerLookupTableByAppRequest) (*types.QueryLockerLookupTableByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := q.GetLockerLookupTable(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "locker-info does not exist for id %d", req.AppId)
	}

	return &types.QueryLockerLookupTableByAppResponse{
		TokenToLockerMapping: item.Lockers,
	}, nil
}

func (q *queryServer) QueryLockerLookupTableByAppAndAssetId(c context.Context, req *types.QueryLockerLookupTableByAppAndAssetIdRequest) (*types.QueryLockerLookupTableByAppAndAssetIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := q.GetLockerLookupTable(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "locker-info does not exist for id %d", req.AppId)
	}
	var locker types.TokenToLockerMapping

	for _, data := range item.Lockers {
		if data.AssetId == req.AssetId {
			locker = *data
		}
	}

	return &types.QueryLockerLookupTableByAppAndAssetIdResponse{
		TokenToLockerMapping: &locker,
	}, nil
}

func (q *queryServer) QueryLockerTotalDepositedByApp(c context.Context, req *types.QueryLockerTotalDepositedByAppRequest) (*types.QueryLockerTotalDepositedByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := q.GetLockerLookupTable(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "locker-info does not exist for id %d", req.AppId)
	}
	var lockedDepositedAmt []types.LockedDepositedAmountDataMap

	for _, data := range item.Lockers {
		var lockeddata types.LockedDepositedAmountDataMap
		lockeddata.AssetId = data.AssetId
		lockeddata.DepositedAmount = data.DepositedAmount
		lockedDepositedAmt = append(lockedDepositedAmt, lockeddata)

	}

	return &types.QueryLockerTotalDepositedByAppResponse{
		LockedDepositedAmountDataMap: lockedDepositedAmt,
	}, nil
}

func (q *queryServer) QueryState(c context.Context, req *types.QueryStateRequest) (*types.QueryStateResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	qs, _ := QueryState(req.Address, req.Denom, req.Height, req.Target)

	return &types.QueryStateResponse{
		Amount: *qs,
	}, nil
}

func (q *queryServer) QueryLockerTotalRewardsByAssetAppWise(c context.Context, request *types.QueryLockerTotalRewardsByAssetAppWiseRequest) (*types.QueryLockerTotalRewardsByAssetAppWiseResponse, error) {

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

	rewards_data, found := q.GetLockerTotalRewardsByAssetAppWise(ctx,request.AppId,request.AssetId)
	if !found{
		return &types.QueryLockerTotalRewardsByAssetAppWiseResponse{},nil
	}
	return &types.QueryLockerTotalRewardsByAssetAppWiseResponse{
		TotalRewards: rewards_data,
	}, nil
}