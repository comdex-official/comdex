package keeper

import (
	"fmt"
	"strconv"

	"github.com/bandprotocol/bandchain-packet/obi"
	bandpacket "github.com/bandprotocol/bandchain-packet/packet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/oracle/types"
)

func (k *Keeper) SetMarket(ctx sdk.Context, market types.Market) {
	var (
		store = k.Store(ctx)
		key   = types.MarketKey(market.Symbol)
		value = k.cdc.MustMarshal(&market)
	)

	store.Set(key, value)
}

func (k *Keeper) HasMarket(ctx sdk.Context, symbol string) bool {
	var (
		store = k.Store(ctx)
		key   = types.MarketKey(symbol)
	)

	return store.Has(key)
}

func (k *Keeper) GetMarket(ctx sdk.Context, symbol string) (market types.Market, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.MarketKey(symbol)
		value = store.Get(key)
	)

	if value == nil {
		return market, found
	}

	k.cdc.MustUnmarshal(value, &market)
	return market, true
}

func (k *Keeper) GetMarkets(ctx sdk.Context) (markets []types.Market) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.MarketKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var market types.Market
		k.cdc.MustUnmarshal(iter.Value(), &market)
		markets = append(markets, market)
	}

	return markets
}

func (k *Keeper) SetPriceForMarket(ctx sdk.Context, symbol string, price uint64) {
	var (
		store = k.Store(ctx)
		key   = types.PriceForMarketKey(symbol)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: price,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetPriceForMarket(ctx sdk.Context, symbol string) (uint64, bool) {
	var (
		store = k.Store(ctx)
		key   = types.PriceForMarketKey(symbol)
		value = store.Get(key)
	)

	if value == nil {
		return 0, false
	}

	var price protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &price)

	return price.GetValue(), true
}

func (k *Keeper) SetMarketForAsset(ctx sdk.Context, id uint64, symbol string) {
	var (
		store = k.Store(ctx)
		key   = types.MarketForAssetKey(id)
		value = k.cdc.MustMarshal(
			&protobuftypes.StringValue{
				Value: symbol,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) HasMarketForAsset(ctx sdk.Context, id uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.MarketForAssetKey(id)
	)

	return store.Has(key)
}

func (k *Keeper) GetMarketForAsset(ctx sdk.Context, id uint64) (market types.Market, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.MarketForAssetKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return market, false
	}

	var symbol protobuftypes.StringValue
	k.cdc.MustUnmarshal(value, &symbol)

	return k.GetMarket(ctx, symbol.GetValue())
}

func (k *Keeper) DeleteMarketForAsset(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.MarketForAssetKey(id)
	)

	store.Delete(key)
}

func (k *Keeper) SetCalldataID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.CalldataIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetCalldataID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.CalldataIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetCalldata(ctx sdk.Context, id uint64, calldata types.Calldata) {
	var (
		store = k.Store(ctx)
		key   = types.CalldataKey(id)
		value = k.cdc.MustMarshal(&calldata)
	)

	store.Set(key, value)
}

func (k *Keeper) GetCalldata(ctx sdk.Context, id uint64) (calldata types.Calldata, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.CalldataKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return calldata, false
	}

	k.cdc.MustUnmarshal(value, &calldata)
	return calldata, true
}

func (k *Keeper) DeleteCalldata(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.CalldataKey(id)
	)

	store.Delete(key)
}

func (k *Keeper) OnRecvPacket(ctx sdk.Context, res bandpacket.OracleResponsePacketData) error {
	id, err := strconv.ParseUint(res.ClientID, 10, 64)
	if err != nil {
		return err
	}

	if res.ResolveStatus == bandpacket.RESOLVE_STATUS_SUCCESS {
		calldata, found := k.GetCalldata(ctx, id)
		if !found {
			return fmt.Errorf("calldata does not exist for id %d", id)
		}

		var result types.Result
		if err := obi.Decode(res.Result, &result); err != nil {
			return err
		}

		for i := range calldata.Symbols {
			k.SetPriceForMarket(ctx, calldata.Symbols[i], result.Rates[i])
		}
	}

	k.DeleteCalldata(ctx, id)
	return nil
}

func (k *Keeper) HasAsset(ctx sdk.Context, id uint64) bool {
	var (
		store = k.Store(ctx)
		key   = assettypes.AssetKey(id)
	)

	return store.Has(key)
}