package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/esm/types"
)

func (k *Keeper) SetKillSwitchData(ctx sdk.Context, switchParams types.KillSwitchParams) error {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.KillSwitchData(switchParams.AppId)
		value = k.cdc.MustMarshal(&switchParams)
	)

	_, found :=k.GetApp(ctx, switchParams.AppId)
	if !found{
		return types.ErrorAppDoesNotExists
	}

	store.Set(key, value)
	return nil
}

func (k *Keeper) GetKillSwitchData(ctx sdk.Context, app_id uint64) (switchParams types.KillSwitchParams,found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.KillSwitchData(app_id)
		value = store.Get(key)
	)

	if value == nil {
		return switchParams, false
	}

	k.cdc.MustUnmarshal(value, &switchParams)

	return switchParams, true
}

func (k *Keeper) Admin(ctx sdk.Context, from string) bool {
	var from_address = []string{"comdex1gvcsuex523fcwuzcpaqys99r70hajf8ffg6322", "comdex1mska4sk59e7t23r2vv3mvzljujxf9j08frl2tg", ""}
	for _, addr := range from_address{
		if addr == from{
			return true
		}
	}
	return false
}