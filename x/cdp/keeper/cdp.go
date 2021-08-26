package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/cdp/types"
)

func (k *Keeper) SetID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.IDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.IDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetCDP(ctx sdk.Context, cdp types.CDP) {
	var (
		store = k.Store(ctx)
		key   = types.CDPKey(cdp.ID)
		value = k.cdc.MustMarshal(&cdp)
	)

	store.Set(key, value)
}

func (k *Keeper) GetCDP(ctx sdk.Context, id uint64) (cdp types.CDP, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.CDPKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return cdp, false
	}

	k.cdc.MustUnmarshal(value, &cdp)
	return cdp, true
}

func (k *Keeper) GetCDPs(ctx sdk.Context) (cdps []types.CDP) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.CDPKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var cdp types.CDP
		k.cdc.MustUnmarshal(iter.Value(), &cdp)
		cdps = append(cdps, cdp)
	}

	return cdps
}

func (k *Keeper) SetCDPForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.CDPForAddressByPair(address, pairID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) HasCDPForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.CDPForAddressByPair(address, pairID)
	)

	return store.Has(key)
}
