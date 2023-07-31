package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/oracle/types"
)

// IterateAllHistoricPrices iterates over all historic prices.
// Iterator stops when exhausting the source, or when the handler returns `true`.
func (k Keeper) IterateAllHistoricPrices(
	ctx sdk.Context,
	handler func(types.PriceStamp) bool,
) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixHistoricPrice)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var decProto sdk.DecProto
		k.cdc.MustUnmarshal(iter.Value(), &decProto)
		denom, blockNum := types.ParseDenomAndBlockFromKey(iter.Key(), types.KeyPrefixHistoricPrice)
		historicPrice := types.PriceStamp{
			ExchangeRate: &sdk.DecCoin{Denom: denom, Amount: decProto.Dec},
			BlockNum:     blockNum,
		}
		if handler(historicPrice) {
			break
		}
	}
}

// AllHistoricPrices is a helper function that collects and returns all
// median prices using the IterateAllHistoricPrices iterator
func (k Keeper) AllHistoricPrices(ctx sdk.Context) types.PriceStamps {
	prices := types.PriceStamps{}
	k.IterateAllHistoricPrices(ctx, func(median types.PriceStamp) (stop bool) {
		prices = append(prices, median)
		return false
	})
	return prices
}

// IterateAllMedianPrices iterates over all median prices.
// Iterator stops when exhausting the source, or when the handler returns `true`.
func (k Keeper) IterateAllMedianPrices(
	ctx sdk.Context,
	handler func(types.PriceStamp) bool,
) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixMedian)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var decProto sdk.DecProto
		k.cdc.MustUnmarshal(iter.Value(), &decProto)
		denom, blockNum := types.ParseDenomAndBlockFromKey(iter.Key(), types.KeyPrefixMedian)
		median := types.PriceStamp{
			ExchangeRate: &sdk.DecCoin{Denom: denom, Amount: decProto.Dec},
			BlockNum:     blockNum,
		}

		if handler(median) {
			break
		}
	}
}

// AllMedianPrices is a helper function that collects and returns all
// median prices using the IterateAllMedianPrices iterator
func (k Keeper) AllMedianPrices(ctx sdk.Context) types.PriceStamps {
	prices := types.PriceStamps{}
	k.IterateAllMedianPrices(ctx, func(median types.PriceStamp) (stop bool) {
		prices = append(prices, median)
		return false
	})
	return prices
}

// IterateAllMedianDeviationPrices iterates over all median deviation prices.
// Iterator stops when exhausting the source, or when the handler returns `true`.
func (k Keeper) IterateAllMedianDeviationPrices(
	ctx sdk.Context,
	handler func(types.PriceStamp) bool,
) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixMedianDeviation)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var decProto sdk.DecProto
		k.cdc.MustUnmarshal(iter.Value(), &decProto)
		denom, blockNum := types.ParseDenomAndBlockFromKey(iter.Key(), types.KeyPrefixMedianDeviation)
		medianDeviation := types.PriceStamp{
			ExchangeRate: &sdk.DecCoin{Denom: denom, Amount: decProto.Dec},
			BlockNum:     blockNum,
		}

		if handler(medianDeviation) {
			break
		}
	}
}

// AllMedianDeviationPrices is a helper function that collects and returns
// all median prices using the IterateAllMedianDeviationPrices iterator
func (k Keeper) AllMedianDeviationPrices(ctx sdk.Context) types.PriceStamps {
	prices := types.PriceStamps{}
	k.IterateAllMedianDeviationPrices(ctx, func(median types.PriceStamp) (stop bool) {
		prices = append(prices, median)
		return false
	})
	return prices
}
