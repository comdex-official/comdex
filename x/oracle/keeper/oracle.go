package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

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
		return market, false
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

func (k *Keeper) getRates(ctx sdk.Context, symbol string) (uint64, bool){

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

func (k *Keeper) SetRates(ctx sdk.Context, symbol string) {
	var (
		store = k.Store(ctx)
		key   = types.PriceForMarketKey(symbol)
	)
	data, _ := k.bandoraclekeeper.GetFetchPriceResult(ctx, 1)


	switch symbol {
	case "ATOM":
		value, _:= k.cdc.Marshal(&protobuftypes.UInt64Value{
			Value: data.Rates[0],
		},
		)
		store.Set(key, value)

	case "cGOLD":
		value, _:= k.cdc.Marshal(&protobuftypes.UInt64Value{
			Value: data.Rates[1],
		},
		)
		store.Set(key, value)

	case "cSILVER":
		value, _:= k.cdc.Marshal(&protobuftypes.UInt64Value{
			Value: data.Rates[2],
		},
		)
		store.Set(key, value)

	case "cOIL":
		value, _:= k.cdc.Marshal(&protobuftypes.UInt64Value{
			Value: data.Rates[3],
		},
		)
		store.Set(key, value)

	case "UST":
		value, _:= k.cdc.Marshal(&protobuftypes.UInt64Value{
			Value: data.Rates[4],
		},
		)
		store.Set(key, value)

	default:
	}

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

func (k *Keeper) GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool) {
	market, found := k.GetMarketForAsset(ctx, id)
	if !found {
		return 0, false
	}

	return k.GetPriceForMarket(ctx, market.Symbol)
}
