package keeper

import (
	"fmt"
	"time"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"

	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	"github.com/comdex-official/comdex/x/rewards/expected"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/comdex-official/comdex/x/rewards/types"
)

type (
	Keeper struct {
		cdc             codec.BinaryCodec
		storeKey        sdk.StoreKey
		memKey          sdk.StoreKey
		paramstore      paramtypes.Subspace
		locker          expected.LockerKeeper
		collector       expected.CollectorKeeper
		vault           expected.VaultKeeper
		asset           expected.AssetKeeper
		bank            expected.BankKeeper
		liquidityKeeper expected.LiquidityKeeper
		marketKeeper    expected.MarketKeeper
		esm             expected.EsmKeeper
		lend            expected.LendKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	locker expected.LockerKeeper,
	collector expected.CollectorKeeper,
	vault expected.VaultKeeper,
	asset expected.AssetKeeper,
	bank expected.BankKeeper,
	liquidityKeeper expected.LiquidityKeeper,
	marketKeeper expected.MarketKeeper,
	esm expected.EsmKeeper,
	lend expected.LendKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		memKey:          memKey,
		paramstore:      ps,
		locker:          locker,
		collector:       collector,
		vault:           vault,
		asset:           asset,
		bank:            bank,
		liquidityKeeper: liquidityKeeper,
		marketKeeper:    marketKeeper,
		esm:             esm,
		lend:            lend,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func uint64InSlice(a uint64, list []uint64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// WhitelistAssetForInternalRewards of an app for internal rewards
func (k Keeper) WhitelistAssetForInternalRewards(ctx sdk.Context, appMappingID uint64, assetID uint64) error {
	_, found := k.locker.GetLockerProductAssetMapping(ctx, appMappingID, assetID)
	if !found {
		return types.ErrAssetIDDoesNotExist
	}
	internalReward, found := k.GetReward(ctx, appMappingID, assetID)
	if !found {
		internalReward.AppMappingId = appMappingID
		internalReward.AssetId = assetID
		k.SetReward(ctx, internalReward)
	}

	return nil
}

func (k Keeper) WhitelistAppIDVault(ctx sdk.Context, appMappingID uint64) error {
	found := uint64InSlice(appMappingID, k.GetAppIDs(ctx))
	if found {
		return types.ErrAppIDExists
	}

	k.SetAppByAppID(ctx, appMappingID)
	return nil
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) ActExternalRewardsLockers(
	ctx sdk.Context,
	appMappingID uint64,
	assetID uint64,
	totalRewards sdk.Coin,
	durationDays int64,
	depositor sdk.AccAddress,
	minLockupTimeSeconds int64,
) error {
	id := k.GetExternalRewardsLockersID(ctx)
	_, found := k.locker.GetLockerProductAssetMapping(ctx, appMappingID, assetID)
	if !found {
		return types.ErrAssetIDDoesNotExist
	}

	endTime := ctx.BlockTime().Add(time.Second * time.Duration(durationDays*types.SecondsPerDay))

	epochID := k.GetEpochTimeID(ctx)
	epoch := types.EpochTime{
		Id:           epochID + 1,
		AppMappingId: appMappingID,
		StartingTime: ctx.BlockTime().Unix() + types.SecondsPerDay,
	}

	msg := types.LockerExternalRewards{
		Id:                   id + 1,
		AppMappingId:         appMappingID,
		AssetId:              assetID,
		TotalRewards:         totalRewards,
		DurationDays:         durationDays,
		IsActive:             true,
		AvailableRewards:     totalRewards,
		Depositor:            depositor.String(),
		StartTimestamp:       ctx.BlockTime(),
		EndTimestamp:         endTime,
		MinLockupTimeSeconds: minLockupTimeSeconds,
		EpochId:              epoch.Id,
	}
	if err := k.bank.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(totalRewards)); err != nil {
		return err
	}

	k.SetEpochTimeID(ctx, msg.EpochId)
	k.SetExternalRewardsLockers(ctx, msg)
	k.SetExternalRewardsLockersID(ctx, msg.Id)
	k.SetEpochTime(ctx, epoch)
	return nil
}

func (k Keeper) ActExternalRewardsVaults(
	ctx sdk.Context,
	appMappingID uint64, extendedPairID uint64,
	durationDays, minLockupTimeSeconds int64,
	totalRewards sdk.Coin,
	depositor sdk.AccAddress,
) error {
	id := k.GetExternalRewardsVaultID(ctx)

	appExtPairVaultData, found := k.vault.GetAppMappingData(ctx, appMappingID)
	if !found {
		return types.ErrAppIDDoesNotExists
	}
	for _, v := range appExtPairVaultData {
		if extendedPairID != v.ExtendedPairId {
			return types.ErrPairNotExists
		}
	}

	endTime := ctx.BlockTime().Add(time.Second * time.Duration(durationDays*types.SecondsPerDay))

	epochID := k.GetEpochTimeID(ctx)
	epoch := types.EpochTime{
		Id:           epochID + 1,
		AppMappingId: appMappingID,
		StartingTime: ctx.BlockTime().Unix() + types.SecondsPerDay,
	}

	msg := types.VaultExternalRewards{
		Id:                   id + 1,
		AppMappingId:         appMappingID,
		ExtendedPairId:       extendedPairID,
		TotalRewards:         totalRewards,
		DurationDays:         durationDays,
		IsActive:             true,
		AvailableRewards:     totalRewards,
		Depositor:            depositor.String(),
		StartTimestamp:       ctx.BlockTime(),
		EndTimestamp:         endTime,
		MinLockupTimeSeconds: minLockupTimeSeconds,
		EpochId:              epoch.Id,
	}

	if err := k.bank.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(totalRewards)); err != nil {
		return err
	}

	k.SetEpochTimeID(ctx, msg.EpochId)
	k.SetExternalRewardVault(ctx, msg)
	k.SetExternalRewardsVaultID(ctx, msg.Id)
	k.SetEpochTime(ctx, epoch)
	return nil
}

// WasmRemoveWhitelistAssetLocker tx and query binding functions
func (k Keeper) WasmRemoveWhitelistAssetLocker(ctx sdk.Context, appMappingID uint64, assetID uint64) error {
	klwsParams, _ := k.esm.GetKillSwitchData(ctx, appMappingID)
	if klwsParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := k.esm.GetESMStatus(ctx, appMappingID)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return esmtypes.ErrESMAlreadyExecuted
	}

	_, found1 := k.GetReward(ctx, appMappingID, assetID)
	if !found1 {
		return types.ErrInternalRewardsNotFound
	}
	k.DeleteReward(ctx, appMappingID, assetID)
	return nil
}

func (k Keeper) WasmRemoveWhitelistAssetLockerQuery(ctx sdk.Context, appMappingID uint64, assetID uint64) (bool, string) {
	_, found := k.GetReward(ctx, appMappingID, assetID)
	if !found {
		return false, "app Id not found"
	}
	return true, ""
}

func (k Keeper) WasmRemoveWhitelistAppIDVaultInterest(ctx sdk.Context, appMappingID uint64) error {
	klwsParams, _ := k.esm.GetKillSwitchData(ctx, appMappingID)
	if klwsParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := k.esm.GetESMStatus(ctx, appMappingID)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return esmtypes.ErrESMAlreadyExecuted
	}

	k.DeleteAppIDByApp(ctx, appMappingID)

	return nil
}

func (k Keeper) WasmRemoveWhitelistAppIDVaultInterestQuery(ctx sdk.Context, appMappingID uint64) (bool, string) {
	found := uint64InSlice(appMappingID, k.GetAppIDs(ctx))
	if !found {
		return false, types.ErrAppIDDoesNotExists.Error()
	}
	return true, ""
}

func (k Keeper) AddLendExternalRewards(ctx sdk.Context, msg types.ActivateExternalRewardsLend) error {
	id := k.GetExternalRewardsLendID(ctx)
	endTime := ctx.BlockTime().Add(time.Second * time.Duration(msg.DurationDays*types.SecondsPerDay))
	epochID := k.GetEpochTimeID(ctx)
	_, found := k.liquidityKeeper.GetPool(ctx, msg.CSwapAppId, uint64(msg.MasterPoolId))
	if !found {
		return liquiditytypes.ErrInvalidPoolID
	}
	for _, assetID := range msg.AssetId {
		_, found = k.asset.GetAsset(ctx, assetID)
		if !found {
			return assettypes.ErrorAssetDoesNotExist
		}
	}

	_, found = k.lend.GetPool(ctx, msg.CPoolId)
	if !found {
		return lendtypes.ErrPoolNotFound
	}

	_, found = k.asset.GetAssetForDenom(ctx, msg.TotalRewards.Denom)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	if msg.DurationDays == 0 {
		return types.ErrInvalidDuration
	}

	RewardsAssetPoolData := types.RewardsAssetPoolData{
		CPoolId:            msg.CPoolId,
		AssetId:            msg.AssetId,
		CSwapAppId:         msg.CSwapAppId,
		CSwapMinLockAmount: msg.CSwapMinLockAmount,
	}

	epoch := types.EpochTime{
		Id:           epochID + 1,
		AppMappingId: msg.AppMappingId,
		StartingTime: ctx.BlockTime().Unix() + 84600,
	}
	newMsg := types.LendExternalRewards{
		Id:                   id + 1,
		AppMappingId:         msg.AppMappingId,
		RewardsAssetPoolData: &RewardsAssetPoolData,
		TotalRewards:         msg.TotalRewards,
		MasterPoolId:         msg.MasterPoolId,
		DurationDays:         msg.DurationDays,
		IsActive:             true,
		AvailableRewards:     msg.TotalRewards,
		Depositor:            msg.Depositor,
		StartTimestamp:       ctx.BlockTime(),
		EndTimestamp:         endTime,
		MinLockupTimeSeconds: msg.MinLockupTimeSeconds,
		EpochId:              epoch.Id,
	}
	depositor, _ := sdk.AccAddressFromBech32(msg.Depositor)

	if err := k.bank.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(msg.TotalRewards)); err != nil {
		return err
	}

	k.SetEpochTimeID(ctx, newMsg.EpochId)
	k.SetExternalRewardLend(ctx, newMsg)
	k.SetExternalRewardsLendID(ctx, newMsg.Id)
	k.SetEpochTime(ctx, epoch)
	return nil
}

func (k Keeper) ActExternalRewardsStableVaults(
	ctx sdk.Context,
	appID uint64, cswapAppID, commodoAppID uint64,
	durationDays, acceptedBlockHeight int64,
	totalRewards sdk.Coin,
	depositor sdk.AccAddress,
) error {
	id := k.GetExternalRewardsStableVault(ctx)

	endTime := ctx.BlockTime().Add(time.Second * time.Duration(durationDays*types.SecondsPerDay))

	epochID := k.GetEpochTimeID(ctx)
	epoch := types.EpochTime{
		Id:           epochID + 1,
		AppMappingId: appID,
		StartingTime: ctx.BlockTime().Unix() + types.SecondsPerDay,
	}

	msg := types.StableVaultExternalRewards{
		Id:                  id + 1,
		AppId:               appID,
		CswapAppId:          cswapAppID,
		CommodoAppId:        commodoAppID,
		TotalRewards:        totalRewards,
		DurationDays:        durationDays,
		IsActive:            true,
		AvailableRewards:    totalRewards,
		Depositor:           depositor.String(),
		StartTimestamp:      ctx.BlockTime(),
		EndTimestamp:        endTime,
		AcceptedBlockHeight: acceptedBlockHeight,
		EpochId:             epoch.Id,
	}

	if err := k.bank.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(totalRewards)); err != nil {
		return err
	}

	k.SetEpochTimeID(ctx, msg.EpochId)
	k.SetExternalRewardStableVault(ctx, msg)
	k.SetExternalRewardsStableVault(ctx, msg.Id)
	k.SetEpochTime(ctx, epoch)
	return nil
}
