package keeper

import (
	sdkmath "cosmossdk.io/math"
	"context"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	"github.com/comdex-official/comdex/x/locker/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k msgServer) MsgCreateLocker(c context.Context, msg *types.MsgCreateLockerRequest) (*types.MsgCreateLockerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	klwsParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	asset, found := k.asset.GetAsset(ctx, msg.AssetId)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	// Checking if user mapping exists
	userDataForLocker, _ := k.GetUserLockerAssetMapping(ctx, msg.Depositor, msg.AppId, msg.AssetId)
	if userDataForLocker.LockerId != 0 {
		return nil, types.ErrorUserLockerAlreadyExists
	}
	Collector, found := k.collector.GetCollectorLookupTable(ctx, msg.AppId, msg.AssetId)
	if !found {
		return nil, types.ErrorCollectorLookupDoesNotExists
	}

	lockerProductAssetMapping, found := k.GetLockerProductAssetMapping(ctx, appMapping.Id, msg.AssetId)
	if !found {
		return nil, types.ErrorLockerProductAssetMappingDoesNotExists
	}
	// This asset is accepted by the app
	// Create a new instance of locker

	// call Lookup table to get relevant data
	lookupTableData, exists := k.GetLockerLookupTable(ctx, lockerProductAssetMapping.AppId, msg.AssetId)
	if !exists {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	// Transferring amount from user to module
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	if msg.Amount.GT(sdkmath.ZeroInt()) {
		if err := k.bank.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Denom, msg.Amount))); err != nil {
			return nil, err
		}
	}
	blockHeight := ctx.BlockHeight()
	blockTime := ctx.BlockTime()
	if Collector.LockerSavingRate.IsZero() {
		blockHeight = 0
	}

	// Creating locker instance
	id := k.GetIDForLocker(ctx)
	var userLocker types.Locker
	userLocker.LockerId = id + 1
	userLocker.Depositor = msg.Depositor
	userLocker.AssetDepositId = asset.Id
	userLocker.CreatedAt = ctx.BlockTime()
	userLocker.IsLocked = false
	userLocker.NetBalance = msg.Amount
	userLocker.ReturnsAccumulated = sdkmath.ZeroInt()
	userLocker.AppId = appMapping.Id
	userLocker.BlockHeight = blockHeight
	userLocker.BlockTime = blockTime
	k.SetLocker(ctx, userLocker)
	k.SetIDForLocker(ctx, id+1)

	// Create a new instance
	var userMappingData types.UserAppAssetLockerMapping

	var userTxData types.UserTxData
	userMappingData.AppId = appMapping.Id
	userMappingData.AssetId = asset.Id
	userMappingData.LockerId = userLocker.LockerId
	userMappingData.Owner = msg.Depositor

	userTxData.TxType = "Create"
	userTxData.Amount = msg.Amount
	userTxData.Balance = msg.Amount
	userTxData.TxTime = ctx.BlockTime()

	userMappingData.UserData = append(userMappingData.UserData, &userTxData)

	k.SetUserLockerAssetMapping(ctx, userMappingData)

	lookupTableData.DepositedAmount = lookupTableData.DepositedAmount.Add(userLocker.NetBalance)
	lookupTableData.LockerIds = append(lookupTableData.LockerIds, userLocker.LockerId)
	k.SetLockerLookupTable(ctx, lookupTableData)

	ctx.GasMeter().ConsumeGas(types.CreateLockerGas, "CreateLockerGas")

	return &types.MsgCreateLockerResponse{}, nil
}

// MsgDepositAsset Remove asset id from Deposit & Withdraw redundant.
func (k msgServer) MsgDepositAsset(c context.Context, msg *types.MsgDepositAssetRequest) (*types.MsgDepositAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	klwsParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	asset, found := k.asset.GetAsset(ctx, msg.AssetId)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	lockerData, found := k.GetLocker(ctx, msg.LockerId)

	if !found {
		return nil, types.ErrorLockerDoesNotExists
	}
	if lockerData.AssetDepositId != asset.Id {
		return nil, types.ErrorInvalidAssetID
	}
	if msg.Depositor != lockerData.Depositor {
		return nil, types.ErrorUnauthorized
	}
	if appMapping.Id != lockerData.AppId {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	lookupTableData, exists := k.GetLockerLookupTable(ctx, appMapping.Id, asset.Id)
	if !exists {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	err1 := k.rewards.CalculateLockerRewards(ctx, appMapping.Id, msg.AssetId, lockerData.LockerId, string(depositor), lockerData.NetBalance, lockerData.BlockHeight, lockerData.BlockTime.Unix())
	if err1 != nil {
		return nil, err1
	}
	lockerData, _ = k.GetLocker(ctx, msg.LockerId)

	if msg.Amount.GT(sdkmath.ZeroInt()) {
		if err := k.bank.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Denom, msg.Amount))); err != nil {
			return nil, err
		}
	}

	// calculating user locker rewards
	// diffHeight := ctx.BlockHeight() - lockerData.BlockHeight
	// lockRewards, err := k.CalculateLockerRewards(ctx, appMapping, lockerData.LockerId)
	lockerData.BlockHeight = ctx.BlockHeight()
	lockerData.BlockTime = ctx.BlockTime()
	lockerData.NetBalance = lockerData.NetBalance.Add(msg.Amount)
	k.SetLocker(ctx, lockerData)

	// Update  Amount in Locker Mapping
	k.UpdateAmountLockerMapping(ctx, lookupTableData.AppId, asset.Id, msg.Amount, true)

	userLockerAssetMappingData, _ := k.GetUserLockerAssetMapping(ctx, msg.Depositor, msg.AppId, msg.AssetId)
	var userHisData types.UserTxData
	userHisData.TxType = "Deposit"
	userHisData.Amount = msg.Amount
	userHisData.Balance = lockerData.NetBalance
	userHisData.TxTime = ctx.BlockTime()

	userLockerAssetMappingData.UserData = append(userLockerAssetMappingData.UserData, &userHisData)

	k.SetUserLockerAssetMapping(ctx, userLockerAssetMappingData)

	ctx.GasMeter().ConsumeGas(types.DepositLockerGas, "DepositLockerGas")

	return &types.MsgDepositAssetResponse{}, nil
}

// MsgWithdrawAsset Remove asset id from Deposit & Withdraw-redundant.
func (k msgServer) MsgWithdrawAsset(c context.Context, msg *types.MsgWithdrawAssetRequest) (*types.MsgWithdrawAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	asset, found := k.asset.GetAsset(ctx, msg.AssetId)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	lockerData, found := k.GetLocker(ctx, msg.LockerId)

	if !found {
		return nil, types.ErrorLockerDoesNotExists
	}
	if lockerData.AssetDepositId != asset.Id {
		return nil, types.ErrorInvalidAssetID
	}
	if msg.Depositor != lockerData.Depositor {
		return nil, types.ErrorUnauthorized
	}
	if appMapping.Id != lockerData.AppId {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	lookupTableData, exists := k.GetLockerLookupTable(ctx, appMapping.Id, asset.Id)
	if !exists {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	if lockerData.NetBalance.LT(msg.Amount) {
		return nil, types.ErrorRequestedAmountExceedsDepositAmount
	}
	err1 := k.rewards.CalculateLockerRewards(ctx, appMapping.Id, msg.AssetId, lockerData.LockerId, string(depositor), lockerData.NetBalance, lockerData.BlockHeight, lockerData.BlockTime.Unix())
	if err1 != nil {
		return nil, err1
	}
	lockerData, _ = k.GetLocker(ctx, msg.LockerId)

	lockerData.NetBalance = lockerData.NetBalance.Sub(msg.Amount)

	if msg.Amount.GT(sdkmath.ZeroInt()) {
		if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoins(sdk.NewCoin(asset.Denom, msg.Amount))); err != nil {
			return nil, err
		}
	}

	// diffHeight := ctx.BlockHeight() - lockerData.BlockHeight

	lockerData.BlockHeight = ctx.BlockHeight()
	lockerData.BlockTime = ctx.BlockTime()
	k.SetLocker(ctx, lockerData)

	// Update  Amount in Locker Mapping
	k.UpdateAmountLockerMapping(ctx, lookupTableData.AppId, asset.Id, msg.Amount, false)

	userLockerAssetMappingData, _ := k.GetUserLockerAssetMapping(ctx, msg.Depositor, msg.AppId, msg.AssetId)

	var userTxData types.UserTxData

	userTxData.TxType = "Withdraw"
	userTxData.Amount = msg.Amount
	userTxData.Balance = lockerData.NetBalance
	userTxData.TxTime = ctx.BlockTime()
	userLockerAssetMappingData.UserData = append(userLockerAssetMappingData.UserData, &userTxData)

	k.SetUserLockerAssetMapping(ctx, userLockerAssetMappingData)

	ctx.GasMeter().ConsumeGas(types.WithdrawLockerGas, "WithdrawLockerGas")

	return &types.MsgWithdrawAssetResponse{}, nil
}

func (k msgServer) MsgCloseLocker(c context.Context, msg *types.MsgCloseLockerRequest) (*types.MsgCloseLockerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	asset, found := k.asset.GetAsset(ctx, msg.AssetId)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	lockerData, found := k.GetLocker(ctx, msg.LockerId)

	if !found {
		return nil, types.ErrorLockerDoesNotExists
	}
	if lockerData.AssetDepositId != asset.Id {
		return nil, types.ErrorInvalidAssetID
	}
	if msg.Depositor != lockerData.Depositor {
		return nil, types.ErrorUnauthorized
	}
	if appMapping.Id != lockerData.AppId {
		return nil, types.ErrorAppMappingDoesNotExist
	}

	lookupTableData, exists := k.GetLockerLookupTable(ctx, appMapping.Id, asset.Id)
	if !exists {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}

	err1 := k.rewards.CalculateLockerRewards(ctx, appMapping.Id, msg.AssetId, lockerData.LockerId, string(depositor), lockerData.NetBalance, lockerData.BlockHeight, lockerData.BlockTime.Unix())
	if err1 != nil {
		return nil, err1
	}

	lockerData, _ = k.GetLocker(ctx, msg.LockerId)

	if lockerData.NetBalance.GT(sdkmath.ZeroInt()) {
		if err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, sdk.NewCoins(sdk.NewCoin(asset.Denom, lockerData.NetBalance))); err != nil {
			return nil, err
		}
	}

	userLockerAssetMappingData, _ := k.GetUserLockerAssetMapping(ctx, msg.Depositor, msg.AppId, msg.AssetId)
	userLockerAssetMappingData.LockerId = 0
	var userTxData types.UserTxData

	userTxData.TxType = "Close"
	userTxData.Amount = lockerData.NetBalance
	userTxData.Balance = sdkmath.ZeroInt()
	userTxData.TxTime = ctx.BlockTime()
	userLockerAssetMappingData.UserData = append(userLockerAssetMappingData.UserData, &userTxData)

	k.UpdateAmountLockerMapping(ctx, lookupTableData.AppId, asset.Id, lockerData.NetBalance, false)
	k.SetUserLockerAssetMapping(ctx, userLockerAssetMappingData)
	lookupTableData, _ = k.GetLockerLookupTable(ctx, appMapping.Id, asset.Id)

	lengthOfVaults := len(lookupTableData.LockerIds)
	dataIndex := sort.Search(lengthOfVaults, func(i int) bool { return lookupTableData.LockerIds[i] >= lockerData.LockerId })

	if dataIndex < lengthOfVaults && lookupTableData.LockerIds[dataIndex] == lockerData.LockerId {
		lookupTableData.LockerIds = append(lookupTableData.LockerIds[:dataIndex], lookupTableData.LockerIds[dataIndex+1:]...)
		k.SetLockerLookupTable(ctx, lookupTableData)
	}

	k.DeleteLocker(ctx, lockerData.LockerId)
	var rewards rewardstypes.LockerRewardsTracker
	rewards.AppMappingId = appMapping.Id
	rewards.LockerId = lockerData.LockerId
	k.rewards.DeleteLockerRewardTracker(ctx, rewards)

	ctx.GasMeter().ConsumeGas(types.CloseLockerGas, "CloseLockerGas")

	return &types.MsgCloseLockerResponse{}, nil
}

func (k msgServer) MsgLockerRewardCalc(c context.Context, msg *types.MsgLockerRewardCalcRequest) (*types.MsgLockerRewardCalcResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	depositor, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	lockerData, found := k.GetLocker(ctx, msg.LockerId)

	if !found {
		return nil, types.ErrorLockerDoesNotExists
	}
	if lockerData.AppId != appMapping.Id {
		return nil, types.ErrorAppMappingIDMismatch
	}
	err1 := k.rewards.CalculateLockerRewards(ctx, appMapping.Id, lockerData.AssetDepositId, lockerData.LockerId, string(depositor), lockerData.NetBalance, lockerData.BlockHeight, lockerData.BlockTime.Unix())
	if err1 != nil {
		return nil, err1
	}

	return &types.MsgLockerRewardCalcResponse{}, nil
}
