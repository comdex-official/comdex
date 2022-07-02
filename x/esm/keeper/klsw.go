package keeper

import (
	"fmt"
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
	// if switchParams.BreakerEnable && switchParams.ProtocolControl{
	// 	return types.ErrBothCantbeOpen
	// }

	fmt.Println("condition setting")
	fmt.Println(value)

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
