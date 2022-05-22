package keeper

import (
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) SetReward(ctx sdk.Context, rewards types.InternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.RewardsKey(rewards.App_mapping_ID)
		value = k.cdc.MustMarshal(&rewards)
	)

	store.Set(key, value)
}

func (k *Keeper) GetReward(ctx sdk.Context, id uint64) (asset types.InternalRewards, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.RewardsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return asset, false
	}

	k.cdc.MustUnmarshal(value, &asset)
	return asset, true
}

func (k *Keeper) GetRewards(ctx sdk.Context) (lends []types.InternalRewards) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.RewardsKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var rewards types.InternalRewards
		k.cdc.MustUnmarshal(iter.Value(), &rewards)
		lends = append(lends, rewards)
	}

	return lends
}

func (k *Keeper) SetAppId(ctx sdk.Context, AppIds types.WhitelistedAppIdsVault) {
	var (
		store = k.Store(ctx)
		key   = types.AppIdsVaultKeyPrefix
		value = k.cdc.MustMarshal(&AppIds)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAppIds(ctx sdk.Context) (appIds types.WhitelistedAppIdsVault) {
	var (
		store = k.Store(ctx)
		key   = types.AppIdsVaultKeyPrefix
		value = store.Get(key)
	)

	if value == nil {
		return appIds
	}

	k.cdc.MustUnmarshal(value, &appIds)
	return appIds
}
