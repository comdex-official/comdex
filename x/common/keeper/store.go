package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/comdex-official/comdex/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/cosmos/gogoproto/types"
)

func (k Keeper) SetGameID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.GameIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetGameID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.GameIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetContract(ctx sdk.Context, msg types.WhitelistedContract) error {
	var (
		store = k.Store(ctx)
		key   = types.ContractKey(msg.GameId)
		value = k.cdc.MustMarshal(&msg)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) GetContract(ctx sdk.Context, gameID uint64) (contract types.WhitelistedContract, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.ContractKey(gameID)
		value = store.Get(key)
	)

	if value == nil {
		return contract, false
	}

	k.cdc.MustUnmarshal(value, &contract)
	return contract, true
}

func (k Keeper) DeleteContract(ctx sdk.Context, gameID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.ContractKey(gameID)
	)

	store.Delete(key)
}

func (k Keeper) GetAllContract(ctx sdk.Context) (contracts []types.WhitelistedContract) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.SetContractKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var contract types.WhitelistedContract
		k.cdc.MustUnmarshal(iter.Value(), &contract)
		contracts = append(contracts, contract)
	}
	return contracts
}
