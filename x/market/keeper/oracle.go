package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	bandoraclemoduletypes "github.com/comdex-official/comdex/x/bandoracle/types"

	"github.com/comdex-official/comdex/x/market/types"
)

func (k Keeper) SetMarket(ctx sdk.Context, market types.Market) {
	var (
		store = k.Store(ctx)
		key   = types.MarketKey(market.Symbol)
		value = k.cdc.MustMarshal(&market)
	)
	store.Set(key, value)
}

func (k Keeper) HasMarket(ctx sdk.Context, symbol string) bool {
	var (
		store = k.Store(ctx)
		key   = types.MarketKey(symbol)
	)
	return store.Has(key)
}

func (k Keeper) GetMarket(ctx sdk.Context, symbol string) (market types.Market, found bool) {
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

func (k Keeper) GetMarkets(ctx sdk.Context) (markets []types.Market) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.MarketKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var market types.Market
		k.cdc.MustUnmarshal(iter.Value(), &market)
		markets = append(markets, market)
	}

	return markets
}

func (k Keeper) GetPriceForMarket(ctx sdk.Context, symbol string) (uint64, bool) {
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

func (k Keeper) GetRates(ctx sdk.Context, symbol string) (uint64, bool) {
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

func (k Keeper) SetRates(ctx sdk.Context, _ string) {
	id := k.bandoraclekeeper.GetLastFetchPriceID(ctx)
	data, _ := k.bandoraclekeeper.GetFetchPriceResult(ctx, bandoraclemoduletypes.OracleRequestID(id))
	if data.Rates != nil {
		var sym []string
		allAssets := k.GetAssets(ctx)
		var assets []assetTypes.Asset
		for _, a := range allAssets {
			if a.IsOraclePriceRequired {
				assets = append(assets, a)
			}
		}
		for i, asset := range assets {
			if asset.IsOraclePriceRequired {
				sym = append(sym, asset.Name)
				store := k.Store(ctx)
				key := types.PriceForMarketKey(sym[i])
				value, _ := k.cdc.Marshal(&protobuftypes.UInt64Value{
					Value: data.Rates[i],
				})
				store.Set(key, value)
			}
		}
	}
}

func (k Keeper) SetMarketForAsset(ctx sdk.Context, id uint64, symbol string) {
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

func (k Keeper) HasMarketForAsset(ctx sdk.Context, id uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.MarketForAssetKey(id)
	)

	return store.Has(key)
}

func (k Keeper) GetMarketForAsset(ctx sdk.Context, id uint64) (market types.Market, found bool) {
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

func (k Keeper) GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool) {
	if id != 3 {
		market, found := k.GetMarketForAsset(ctx, id)
		if !found {
			return 0, false
		}

		rates, found := k.GetPriceForMarket(ctx, market.Symbol)
		if !found || rates == 0 {
			return 0, false
		}
		return rates, found
	}
	return 1000000, true
}
