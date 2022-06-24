package keeper

import (
	"fmt"
	"time"

	"github.com/comdex-official/comdex/x/rewards/expected"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

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

func (k Keeper) WhitelistAsset(ctx sdk.Context, appMappingID uint64, assetIDs []uint64) error {
	lockerAssets, _ := k.locker.GetLockerProductAssetMapping(ctx, appMappingID)
	for i := range assetIDs {
		found := uint64InSlice(assetIDs[i], lockerAssets.AssetIds)
		if !found {
			return types.ErrAssetIDDoesNotExist
		}
	}

	internalRewards := types.InternalRewards{
		App_mapping_ID: appMappingID,
		Asset_ID:       assetIDs,
	}

	k.SetReward(ctx, internalRewards)
	return nil
}

func (k Keeper) RemoveWhitelistAsset(ctx sdk.Context, appMappingID uint64, assetID uint64) error {
	rewards, found := k.GetReward(ctx, appMappingID)
	if !found {
		return nil
	}
	var newAssetIDs []uint64
	for i := range rewards.Asset_ID {
		if assetID != rewards.Asset_ID[i] {
			newAssetID := rewards.Asset_ID[i]
			newAssetIDs = append(newAssetIDs, newAssetID)
		}
	}
	newRewards := types.InternalRewards{
		App_mapping_ID: appMappingID,
		Asset_ID:       newAssetIDs,
	}
	k.SetReward(ctx, newRewards)
	return nil
}

func (k Keeper) WhitelistAppIDVault(ctx sdk.Context, appMappingID uint64) error {
	found := uint64InSlice(appMappingID, k.GetAppIDs(ctx).WhitelistedAppMappingIdsVaults)
	if found {
		return types.ErrAppIDExists
	}
	WhitelistedAppIds := append(k.GetAppIDs(ctx).WhitelistedAppMappingIdsVaults, appMappingID)
	UpdatedWhitelistedAppIds := types.WhitelistedAppIdsVault{
		WhitelistedAppMappingIdsVaults: WhitelistedAppIds,
	}
	k.SetAppID(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k Keeper) RemoveWhitelistAppIDVault(ctx sdk.Context, appMappingID uint64) error {
	WhitelistedAppIds := k.GetAppIDs(ctx).WhitelistedAppMappingIdsVaults
	found := uint64InSlice(appMappingID, k.GetAppIDs(ctx).WhitelistedAppMappingIdsVaults)
	if !found {
		return types.ErrAppIDDoesNotExists
	}
	var newAppIds []uint64
	for i := range WhitelistedAppIds {
		if appMappingID != WhitelistedAppIds[i] {
			newAppID := WhitelistedAppIds[i]
			newAppIds = append(newAppIds, newAppID)
		}
	}
	UpdatedWhitelistedAppIds := types.WhitelistedAppIdsVault{
		WhitelistedAppMappingIdsVaults: newAppIds,
	}

	k.SetAppID(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) ActExternalRewardsLockers(
	ctx sdk.Context,
	appMappingID uint64,
	assetID uint64,
	totalRewards sdk.Coin,
	durationDays int64,
	// nolint
	depositor sdk.AccAddress,
	minLockupTimeSeconds int64,
) error {
	id := k.GetExternalRewardsLockersID(ctx)
	lockerAssets, found := k.GetLockerProductAssetMapping(ctx, appMappingID)
	if !found {
		return types.ErrAssetIDDoesNotExist
	}

	found = uint64InSlice(assetID, lockerAssets.AssetIds)
	if !found {
		return types.ErrAssetIDDoesNotExist
	}
	extRewards := k.GetExternalRewardsLockers(ctx)
	for _, v := range extRewards {
		if v.AppMappingId == appMappingID && v.AssetId == assetID {
			return types.ErrAssetIDDoesNotExist
		}
	}

	endTime := ctx.BlockTime().Add(time.Second * time.Duration(durationDays*86400))

	epochID := k.GetEpochTimeID(ctx)
	epoch := types.EpochTime{
		Id:           epochID + 1,
		AppMappingId: appMappingID,
		StartingTime: ctx.BlockTime().Unix() + 84600,
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
	// nolint
	depositor sdk.AccAddress,
) error {
	id := k.GetExternalRewardsVaultID(ctx)

	appExtPairVaultData, found := k.GetAppExtendedPairVaultMapping(ctx, appMappingID)
	if !found {
		return types.ErrAssetIDDoesNotExist
	}
	extPairVault := appExtPairVaultData.ExtendedPairVaults
	for _, v := range extPairVault {
		if extendedPairID != v.ExtendedPairId {
			return types.ErrAssetIDDoesNotExist
		}
	}

	endTime := ctx.BlockTime().Add(time.Second * time.Duration(durationDays*86400))

	epochID := k.GetEpochTimeID(ctx)
	epoch := types.EpochTime{
		Id:           epochID + 1,
		AppMappingId: appMappingID,
		StartingTime: ctx.BlockTime().Unix() + 84600,
	}

	msg := types.VaultExternalRewards{
		Id:                   id + 1,
		AppMappingId:         appMappingID,
		Extended_Pair_Id:     extendedPairID,
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

//Wasm tx and query binding functions

func (k Keeper) WasmRemoveWhitelistAssetLocker(ctx sdk.Context, appMappingID uint64, assetID uint64) error {
	rewards, _ := k.GetReward(ctx, appMappingID)

	var newAssetIDs []uint64
	for i := range rewards.Asset_ID {
		if assetID != rewards.Asset_ID[i] {
			newAssetID := rewards.Asset_ID[i]
			newAssetIDs = append(newAssetIDs, newAssetID)
		}
	}
	newRewards := types.InternalRewards{
		App_mapping_ID: appMappingID,
		Asset_ID:       newAssetIDs,
	}
	k.SetReward(ctx, newRewards)
	return nil
}

func (k Keeper) WasmRemoveWhitelistAssetLockerQuery(ctx sdk.Context, appMappingID uint64, assetID uint64) (bool, string) {
	rewards, found := k.GetReward(ctx, appMappingID)
	if !found {
		return false, "app Id not found"
	}
	for _, j := range rewards.Asset_ID {
		if j != assetID {
			return false, types.ErrAssetIDDoesNotExist.Error()
		}
	}
	return true, ""
}

func (k Keeper) WasmRemoveWhitelistAppIDVaultInterest(ctx sdk.Context, appMappingID uint64) error {
	WhitelistedAppIds := k.GetAppIDs(ctx).WhitelistedAppMappingIdsVaults

	var newAppIDs []uint64
	for i := range WhitelistedAppIds {
		if appMappingID != WhitelistedAppIds[i] {
			newAppID := WhitelistedAppIds[i]
			newAppIDs = append(newAppIDs, newAppID)
		}
	}
	UpdatedWhitelistedAppIds := types.WhitelistedAppIdsVault{
		WhitelistedAppMappingIdsVaults: newAppIDs,
	}

	k.SetAppID(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k Keeper) WasmRemoveWhitelistAppIDVaultInterestQuery(ctx sdk.Context, appMappingID uint64) (bool, string) {
	found := uint64InSlice(appMappingID, k.GetAppIDs(ctx).WhitelistedAppMappingIdsVaults)
	if !found {
		return false, types.ErrAppIDDoesNotExists.Error()
	}
	return true, ""
}
