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

////////////

func (k Keeper) SetTwa(ctx sdk.Context, twa types.TimeWeightedAverage) {
	var (
		store = k.Store(ctx)
		key   = types.TwaKey(twa.AssetID)
		value = k.cdc.MustMarshal(&twa)
	)

	store.Set(key, value)
}

func (k Keeper) GetTwa(ctx sdk.Context, id uint64) (twa types.TimeWeightedAverage, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.TwaKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return twa, false
	}

	k.cdc.MustUnmarshal(value, &twa)
	return twa, true
}

func (k Keeper) GetAllTwa(ctx sdk.Context) (twa []types.TimeWeightedAverage) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.TwaKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var data types.TimeWeightedAverage
		k.cdc.MustUnmarshal(iter.Value(), &data)
		twa = append(twa, data)
	}

	return twa
}

func (k Keeper) UpdatePriceList(ctx sdk.Context, id, scriptID, rate uint64) {
	twa, found := k.GetTwa(ctx, id)
	if !found {
		twa.AssetID = id
		twa.ScriptID = scriptID
		twa.Twa = 0
		twa.IsPriceActive = false
		twa.PriceValue = append(twa.PriceValue, rate)
		twa.CurrentIndex = 1
		k.SetTwa(ctx, twa)
	} else {
		if twa.IsPriceActive {
			twa.PriceValue[twa.CurrentIndex] = rate
			twa.CurrentIndex = twa.CurrentIndex + 1
			twa.Twa = k.CalculateTwa(ctx, twa)
			if twa.CurrentIndex == 30 {
				twa.CurrentIndex = 0
			}
			k.SetTwa(ctx, twa)
		} else {
			twa.PriceValue = append(twa.PriceValue, rate)
			twa.CurrentIndex = twa.CurrentIndex + 1
			if twa.CurrentIndex == 30 {
				twa.IsPriceActive = true
				twa.CurrentIndex = 0
				twa.Twa = k.CalculateTwa(ctx, twa)
			}
			k.SetTwa(ctx, twa)
		}
	}
}

func (k Keeper) CalculateTwa(ctx sdk.Context, twa types.TimeWeightedAverage) uint64 {
	var sum uint64
	for _, price := range twa.PriceValue {
		sum += price
	}
	twa.Twa = sum / 30
	return twa.Twa
}

func (k Keeper) GetLatestPrice(ctx sdk.Context, id uint64) (price uint64, err error) {
	twa, found := k.GetTwa(ctx, id)
	if found && twa.IsPriceActive {
		return twa.PriceValue[twa.CurrentIndex], nil
	}
	return 0, types.ErrorPriceNotActive
}

func (k Keeper) CalcAssetPrice(ctx sdk.Context, id uint64, amt sdk.Int) (price sdk.Int, err error) {
	asset, found := k.GetAsset(ctx, id)
	if !found {
		return sdk.ZeroInt(), assetTypes.ErrorAssetDoesNotExist
	}
	twa, found := k.GetTwa(ctx, id)
	if found && twa.IsPriceActive {
		numerator := sdk.NewDecFromInt(amt).Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(twa.Twa)))
		denominator := sdk.NewDecFromInt(sdk.NewIntFromUint64(uint64(asset.Decimals)))
		result := numerator.Quo(denominator)
		return result.TruncateInt(), nil
	}
	return sdk.ZeroInt(), types.ErrorPriceNotActive
}
