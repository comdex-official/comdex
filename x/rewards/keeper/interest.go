package keeper

import (
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) SetLastInterestTime(ctx sdk.Context, interestTime int64) error {
	store := ctx.KVStore(k.storeKey)
	timeKey := types.CreateLastInterestTimeKey()

	bz, err := k.cdc.Marshal(&gogotypes.Int64Value{Value: interestTime})
	if err != nil {
		return err
	}

	store.Set(timeKey, bz)
	return nil
}

// GetLastInterestTime gets last time at which interest was accrued
func (k Keeper) GetLastInterestTime(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	timeKey := types.CreateLastInterestTimeKey()
	bz := store.Get(timeKey)

	val := gogotypes.Int64Value{}

	if err := k.cdc.Unmarshal(bz, &val); err != nil {
		panic(err)
	}

	return val.Value
}
