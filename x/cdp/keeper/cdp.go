package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/cdp/types"
)

func (k *Keeper) SetCount(ctx sdk.Context, count uint64) {
	var (
		store = k.Store(ctx)
		key   = types.CountKey
		value = k.cdc.MustMarshalBinaryBare(
			&protobuftypes.UInt64Value{
				Value: count,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetCount(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.CountKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var count protobuftypes.UInt64Value
	k.cdc.MustUnmarshalBinaryBare(value, &count)

	return count.GetValue()
}

func (k *Keeper) SetCDP(ctx sdk.Context, cdp types.CDP) {
	var (
		store = k.Store(ctx)
		key   = types.CDPKey(cdp.Id)
		value = k.cdc.MustMarshalBinaryBare(&cdp)
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

	k.cdc.MustUnmarshalBinaryBare(value, &cdp)
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
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &cdp)
		cdps = append(cdps, cdp)
	}

	return cdps
}

func (k *Keeper) SetCDPForAddressByAssetPair(ctx sdk.Context, address sdk.AccAddress, pairId, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.CDPForAddressByAssetPairKey(address, pairId)
		value = k.cdc.MustMarshalBinaryBare(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) HasCDPForAddressByAssetPair(ctx sdk.Context, address sdk.AccAddress, pairId uint64) (yes bool) {
	var (
		store = k.Store(ctx)
		key   = types.CDPForAddressByAssetPairKey(address, pairId)
	)

	return store.Has(key)
}
