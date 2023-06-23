package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/esm/types"
)

func (k Keeper) SetKillSwitchData(ctx sdk.Context, switchParams types.KillSwitchParams) error {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.KillSwitchData(switchParams.AppId)
		value = k.cdc.MustMarshal(&switchParams)
	)

	_, found := k.asset.GetApp(ctx, switchParams.AppId)
	if !found {
		return types.ErrorAppDoesNotExists
	}

	store.Set(key, value)
	return nil
}

func (k Keeper) GetKillSwitchData(ctx sdk.Context, appID uint64) (switchParams types.KillSwitchParams, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.KillSwitchData(appID)
		value = store.Get(key)
	)

	if value == nil {
		return switchParams, false
	}

	k.cdc.MustUnmarshal(value, &switchParams)

	return switchParams, true
}

func (k Keeper) GetAllKillSwitchData(ctx sdk.Context) (killSwitchParams []types.KillSwitchParams) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.KillSwitchDataKey)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var esm types.KillSwitchParams
		k.cdc.MustUnmarshal(iter.Value(), &esm)
		killSwitchParams = append(killSwitchParams, esm)
	}
	return killSwitchParams
}

func (k Keeper) Admin(ctx sdk.Context, from string) bool {
	fromAddress := k.AdminParam(ctx)
	for _, addr := range fromAddress {
		if addr == from {
			return true
		}
	}
	return false
}
