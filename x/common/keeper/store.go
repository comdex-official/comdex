package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/comdex-official/comdex/x/common/types"
)

func (k Keeper) SetContract(ctx sdk.Context, msg types.MsgRegisterContract) error {
	var (
		store = k.Store(ctx)
		key   = types.ContractKey(msg.GameId)
		value = k.cdc.MustMarshal(&msg)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) GetContract(ctx sdk.Context, gameID uint64) (contract types.MsgRegisterContract, found bool) {
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

func (k Keeper) GetAllContract(ctx sdk.Context) (contracts []types.MsgRegisterContract) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.SetContractKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var contract types.MsgRegisterContract
		k.cdc.MustUnmarshal(iter.Value(), &contract)
		contracts = append(contracts, contract)
	}
	return contracts
}