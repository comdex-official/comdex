package keeper

import (
	"github.com/comdex-official/comdex/x/esm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k *Keeper) SetTriggerEsm(ctx sdk.Context, isActive types.EsmActive) error {

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.TriggerEsmKey(isActive.AppId)
		value = k.cdc.MustMarshal(&isActive)
	)

	store.Set(key, value)
	return nil
}

func (k *Keeper) GetTriggerEsm(ctx sdk.Context, app_id uint64) (isActive types.EsmActive, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.TriggerEsmKey(app_id)
		value = store.Get(key)
	)

	if value == nil {
		return isActive, false
	}

	k.cdc.MustUnmarshal(value, &isActive)
	return isActive, true
}