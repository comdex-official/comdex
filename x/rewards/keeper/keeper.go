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

func (k Keeper) WhitelistAsset(ctx sdk.Context, appMappingId uint64, assetId []uint64) error {
	lockerAssets, _ := k.locker.GetLockerProductAssetMapping(ctx, appMappingId)
	for i := range assetId {
		found := uint64InSlice(assetId[i], lockerAssets.AssetIds)
		if !found {
			return types.ErrAssetIdDoesNotExist
		}
	}

	internalRewards := types.InternalRewards{
		App_mapping_ID: appMappingId,
		Asset_ID:       assetId,
	}

	k.SetReward(ctx, internalRewards)
	return nil
}

func (k Keeper) RemoveWhitelistAsset(ctx sdk.Context, appMappingId uint64, assetId uint64) error {

	rewards, found := k.GetReward(ctx, appMappingId)
	if found != true {
		return nil
	}
	var newAssetIds []uint64
	for i := range rewards.Asset_ID {
		if assetId != rewards.Asset_ID[i] {
			newAssetId := rewards.Asset_ID[i]
			newAssetIds = append(newAssetIds, newAssetId)
		}

	}
	newRewards := types.InternalRewards{
		App_mapping_ID: appMappingId,
		Asset_ID:       newAssetIds,
	}
	k.SetReward(ctx, newRewards)
	return nil
}

func (k Keeper) WhitelistAppIdVault(ctx sdk.Context, appMappingId uint64) error {
	found := uint64InSlice(appMappingId, k.GetAppIds(ctx).WhitelistedAppMappingIdsVaults)
	if found {
		return types.ErrAppIdExists
	}
	WhitelistedAppIds := append(k.GetAppIds(ctx).WhitelistedAppMappingIdsVaults, appMappingId)
	UpdatedWhitelistedAppIds := types.WhitelistedAppIdsVault{
		WhitelistedAppMappingIdsVaults: WhitelistedAppIds,
	}
	k.SetAppId(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k Keeper) RemoveWhitelistAppIdVault(ctx sdk.Context, appMappingId uint64) error {
	WhitelistedAppIds := k.GetAppIds(ctx).WhitelistedAppMappingIdsVaults
	found := uint64InSlice(appMappingId, k.GetAppIds(ctx).WhitelistedAppMappingIdsVaults)
	if !found {
		return types.ErrAppIdDoesNotExists
	}
	var newAppIds []uint64
	for i := range WhitelistedAppIds {
		if appMappingId != WhitelistedAppIds[i] {
			newAppId := WhitelistedAppIds[i]
			newAppIds = append(newAppIds, newAppId)
		}
	}
	UpdatedWhitelistedAppIds := types.WhitelistedAppIdsVault{
		WhitelistedAppMappingIdsVaults: newAppIds,
	}

	k.SetAppId(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) ActExternalRewardsLockers(ctx sdk.Context, AppMappingId uint64, AssetId uint64, TotalRewards sdk.Coin, DurationDays int64, Depositor sdk.AccAddress, MinLockupTimeSeconds int64) error {
	Id := k.GetExternalRewardsLockersId(ctx)
	lockerAssets, found := k.GetLockerProductAssetMapping(ctx, AppMappingId)
	if !found {
		return types.ErrAssetIdDoesNotExist
	}

	found = uint64InSlice(AssetId, lockerAssets.AssetIds)
	if !found {
		return types.ErrAssetIdDoesNotExist
	}
	extRewards := k.GetExternalRewardsLockers(ctx)
	for _, v := range extRewards {
		if v.AppMappingId == AppMappingId && v.AssetId == AssetId {
			return types.ErrAssetIdDoesNotExist
		}
	}

	endTime := ctx.BlockTime().Add(time.Second * time.Duration(DurationDays*86400))

	epochId := k.GetEpochTimeId(ctx)
	epoch := types.EpochTime{
		Id:           epochId,
		AppMappingId: AppMappingId,
		StartingTime: ctx.BlockTime().Unix() + 84600,
	}

	msg := types.LockerExternalRewards{
		Id:                   Id + 1,
		AppMappingId:         AppMappingId,
		AssetId:              AssetId,
		TotalRewards:         TotalRewards,
		DurationDays:         DurationDays,
		IsActive:             true,
		AvailableRewards:     TotalRewards,
		Depositor:            Depositor.String(),
		StartTimestamp:       ctx.BlockTime(),
		EndTimestamp:         endTime,
		MinLockupTimeSeconds: MinLockupTimeSeconds,
		EpochId:              epochId,
	}
	if err := k.bank.SendCoinsFromAccountToModule(ctx, Depositor, types.ModuleName, sdk.NewCoins(sdk.Coin{Amount: TotalRewards.Amount, Denom: TotalRewards.Denom})); err != nil {
		return err
	}

	k.SetEpochTimeId(ctx, epochId+1)
	k.SetExternalRewardsLockers(ctx, msg)
	k.SetExternalRewardsLockersId(ctx, msg.Id)
	k.SetEpochTime(ctx, epoch)
	return nil
}

func (k Keeper) ActExternalRewardsVaults(ctx sdk.Context, AppMappingId uint64, Extended_Pair_Id uint64, TotalRewards sdk.Coin, DurationDays int64, Depositor sdk.AccAddress, MinLockupTimeSeconds int64) error {
	Id := k.GetExternalRewardsVaultId(ctx)

	appExtPairVaultData, found := k.GetAppExtendedPairVaultMapping(ctx, AppMappingId)
	if !found {
		return types.ErrAssetIdDoesNotExist
	}
	extPairVault := appExtPairVaultData.ExtendedPairVaults
	for _, v := range extPairVault {
		if Extended_Pair_Id != v.ExtendedPairId {
			return types.ErrAssetIdDoesNotExist
		}
	}

	endTime := ctx.BlockTime().Add(time.Second * time.Duration(DurationDays*86400))

	epochId := k.GetEpochTimeId(ctx)
	epoch := types.EpochTime{
		Id:           epochId,
		AppMappingId: AppMappingId,
		StartingTime: ctx.BlockTime().Unix() + 84600,
	}

	msg := types.VaultExternalRewards{
		Id:                   Id + 1,
		AppMappingId:         AppMappingId,
		Extended_Pair_Id:     Extended_Pair_Id,
		TotalRewards:         TotalRewards,
		DurationDays:         DurationDays,
		IsActive:             true,
		AvailableRewards:     TotalRewards,
		Depositor:            Depositor.String(),
		StartTimestamp:       ctx.BlockTime(),
		EndTimestamp:         endTime,
		MinLockupTimeSeconds: MinLockupTimeSeconds,
		EpochId:              epochId,
	}
	if err := k.bank.SendCoinsFromAccountToModule(ctx, Depositor, types.ModuleName, sdk.NewCoins(sdk.Coin{Amount: TotalRewards.Amount, Denom: TotalRewards.Denom})); err != nil {
		return err
	}

	k.SetEpochTimeId(ctx, epochId+1)
	k.SetExternalRewardVault(ctx, msg)
	k.SetExternalRewardsVaultId(ctx, msg.Id)
	k.SetEpochTime(ctx, epoch)
	return nil
}

//Wasm tx and query binding functions

func (k Keeper) WasmRemoveWhitelistAssetLocker(ctx sdk.Context, appMappingId uint64, assetId uint64) error {

	rewards, _ := k.GetReward(ctx, appMappingId)

	var newAssetIds []uint64
	for i := range rewards.Asset_ID {
		if assetId != rewards.Asset_ID[i] {
			newAssetId := rewards.Asset_ID[i]
			newAssetIds = append(newAssetIds, newAssetId)
		}

	}
	newRewards := types.InternalRewards{
		App_mapping_ID: appMappingId,
		Asset_ID:       newAssetIds,
	}
	k.SetReward(ctx, newRewards)
	return nil
}

func (k Keeper) WasmRemoveWhitelistAssetLockerQuery(ctx sdk.Context, appMappingId uint64, assetId uint64) (bool, string) {
	rewards, found := k.GetReward(ctx, appMappingId)
	if found != true {
		return false, "app Id not found"
	}
	for _, j := range rewards.Asset_ID {
		if j != assetId {
			return false, types.ErrAssetIdDoesNotExist.Error()
		}
	}
	return true, ""
}

func (k Keeper) WasmRemoveWhitelistAppIdVaultInterest(ctx sdk.Context, appMappingId uint64) error {
	WhitelistedAppIds := k.GetAppIds(ctx).WhitelistedAppMappingIdsVaults

	var newAppIds []uint64
	for i := range WhitelistedAppIds {
		if appMappingId != WhitelistedAppIds[i] {
			newAppId := WhitelistedAppIds[i]
			newAppIds = append(newAppIds, newAppId)
		}
	}
	UpdatedWhitelistedAppIds := types.WhitelistedAppIdsVault{
		WhitelistedAppMappingIdsVaults: newAppIds,
	}

	k.SetAppId(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k Keeper) WasmRemoveWhitelistAppIdVaultInterestQuery(ctx sdk.Context, appMappingId uint64) (bool, string) {
	found := uint64InSlice(appMappingId, k.GetAppIds(ctx).WhitelistedAppMappingIdsVaults)
	if !found {
		return false, types.ErrAppIdDoesNotExists.Error()
	}
	return true, ""
}
