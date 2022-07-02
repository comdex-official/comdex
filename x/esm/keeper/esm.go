package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/esm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) AddESMTriggerParamsRecords(ctx sdk.Context, record types.ESMTriggerParams) error {
	_, found := k.GetESMTriggerParams(ctx, record.AppId)
	if found {
		return types.ErrorDuplicateESMTriggerParams
	}

	var (
		esmTriggerParams = types.ESMTriggerParams{
			AppId:         record.AppId,
			TargetValue:   record.TargetValue,
			CoolOffPeriod: record.CoolOffPeriod,
		}
	)
	fmt.Println("in prop esmTriggerParams", esmTriggerParams)
	k.SetESMTriggerParams(ctx, esmTriggerParams)

	return nil
}

func (k *Keeper) SetESMTriggerParams(ctx sdk.Context, esmTriggerParams types.ESMTriggerParams) {
	var (
		store = k.Store(ctx)
		key   = types.ESMTriggerParamsKey(esmTriggerParams.AppId)
		value = k.cdc.MustMarshal(&esmTriggerParams)
	)
	fmt.Println("key... set", key)
	fmt.Println("value... set", value)
	fmt.Println("esmTriggerParams... set", esmTriggerParams)

	store.Set(key, value)
}

func (k *Keeper) GetESMTriggerParams(ctx sdk.Context, id uint64) (esmTriggerParams types.ESMTriggerParams, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.ESMTriggerParamsKey(id)
		value = store.Get(key)
	)
	fmt.Println("key... get", key)
	fmt.Println("value... get ", value)
	fmt.Println("id... get ", id)

	if value == nil {
		return esmTriggerParams, false
	}
	fmt.Println("value...")

	k.cdc.MustUnmarshal(value, &esmTriggerParams)
	return esmTriggerParams, true
}
