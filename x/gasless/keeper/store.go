package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/comdex-official/comdex/x/gasless/types"
)

func (k Keeper) GetTxGPIDS(ctx sdk.Context, txPathOrContractAddress string) (txGPIDS types.TxGPIDS, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTxGPIDSKey(txPathOrContractAddress))
	if bz == nil {
		return
	}
	txGPIDS = types.MustUnmarshalTxGPIDS(k.cdc, bz)
	return txGPIDS, true
}

func (k Keeper) IterateAllTxGPIDS(ctx sdk.Context, cb func(txGPIDS types.TxGPIDS) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllTxGPIDSKey())
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		txGPIDS := types.MustUnmarshalTxGPIDS(k.cdc, iter.Value())
		stop, err := cb(txGPIDS)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func (k Keeper) GetAllTxGPIDS(ctx sdk.Context) (txGPIDSs []types.TxGPIDS) {
	txGPIDSs = []types.TxGPIDS{}
	_ = k.IterateAllTxGPIDS(ctx, func(txGPIDS types.TxGPIDS) (stop bool, err error) {
		txGPIDSs = append(txGPIDSs, txGPIDS)
		return false, nil
	})
	return txGPIDSs
}

func (k Keeper) SetTxGPIDS(ctx sdk.Context, txGPIDS types.TxGPIDS) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalTxGPIDS(k.cdc, txGPIDS)
	store.Set(types.GetTxGPIDSKey(txGPIDS.TxPathOrContractAddress), bz)
}

func (k Keeper) GetLastGasProviderID(ctx sdk.Context) (id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastGasProviderIDKey())
	if bz == nil {
		id = 0 // initialize the GasProviderID
	} else {
		var val gogotypes.UInt64Value
		k.cdc.MustUnmarshal(bz, &val)
		id = val.GetValue()
	}
	return
}

func (k Keeper) SetLastGasProviderID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})
	store.Set(types.GetLastGasProviderIDKey(), bz)
}

func (k Keeper) GetNextGasProviderIDWithUpdate(ctx sdk.Context) uint64 {
	id := k.GetLastGasProviderID(ctx) + 1
	k.SetLastGasProviderID(ctx, id)
	return id
}

func (k Keeper) GetGasProvider(ctx sdk.Context, id uint64) (gasProvider types.GasProvider, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGasProviderKey(id))
	if bz == nil {
		return
	}
	gasProvider = types.MustUnmarshalGasProvider(k.cdc, bz)
	return gasProvider, true
}

func (k Keeper) IterateAllGasProviders(ctx sdk.Context, cb func(gasProvider types.GasProvider) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllGasProvidersKey())
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		gasProvider := types.MustUnmarshalGasProvider(k.cdc, iter.Value())
		stop, err := cb(gasProvider)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func (k Keeper) GetAllGasProviders(ctx sdk.Context) (gasProviders []types.GasProvider) {
	gasProviders = []types.GasProvider{}
	_ = k.IterateAllGasProviders(ctx, func(gasProvider types.GasProvider) (stop bool, err error) {
		gasProviders = append(gasProviders, gasProvider)
		return false, nil
	})
	return gasProviders
}

func (k Keeper) SetGasProvider(ctx sdk.Context, gasProvider types.GasProvider) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalGasProvider(k.cdc, gasProvider)
	store.Set(types.GetGasProviderKey(gasProvider.Id), bz)
}

func (k Keeper) GetGasConsumer(ctx sdk.Context, consumer sdk.AccAddress) (gasConsumer types.GasConsumer, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGasConsumerKey(consumer))
	if bz == nil {
		return
	}
	gasConsumer = types.MustUnmarshalGasConsumer(k.cdc, bz)
	return gasConsumer, true
}

func (k Keeper) IterateAllGasConsumers(ctx sdk.Context, cb func(gasConsumer types.GasConsumer) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllGasConsumersKey())
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		gasConsumer := types.MustUnmarshalGasConsumer(k.cdc, iter.Value())
		stop, err := cb(gasConsumer)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func (k Keeper) GetAllGasConsumers(ctx sdk.Context) (gasConsumers []types.GasConsumer) {
	gasConsumers = []types.GasConsumer{}
	_ = k.IterateAllGasConsumers(ctx, func(gasConsumer types.GasConsumer) (stop bool, err error) {
		gasConsumers = append(gasConsumers, gasConsumer)
		return false, nil
	})
	return gasConsumers
}

func (k Keeper) SetGasConsumer(ctx sdk.Context, gasConsumer types.GasConsumer) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalGasConsumer(k.cdc, gasConsumer)
	store.Set(types.GetGasConsumerKey(sdk.MustAccAddressFromBech32(gasConsumer.Consumer)), bz)
}

func (k Keeper) GetOrCreateGasConsumer(ctx sdk.Context, consumer sdk.AccAddress, gasProvider types.GasProvider) types.GasConsumer {
	gasConsumer, found := k.GetGasConsumer(ctx, consumer)
	if !found {
		gasConsumer = types.NewGasConsumer(consumer)
	}
	if gasConsumer.Consumption == nil {
		gasConsumer.Consumption = make(map[uint64]*types.ConsumptionDetail)
	}
	if _, ok := gasConsumer.Consumption[gasProvider.Id]; !ok {
		gasConsumer.Consumption[gasProvider.Id] = types.NewConsumptionDetail(
			gasProvider.MaxTxsCountPerConsumer,
			sdk.NewCoin(gasProvider.FeeDenom, gasProvider.MaxFeeUsagePerConsumer),
		)
		k.SetGasConsumer(ctx, gasConsumer)
	}
	return gasConsumer
}
