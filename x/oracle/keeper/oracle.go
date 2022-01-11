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

func (k *Keeper) SetPriceForMarket(ctx sdk.Context, symbol string, price uint64) {
	var (
		store = k.Store(ctx)
		key1   = types.PriceForMarketKey("atom")
		key2   = types.PriceForMarketKey("btc")
		key3   = types.PriceForMarketKey("eth")
		key4   = types.PriceForMarketKey("xrp")

	)

	var prices [4]uint64
	data,_ := k.bandoraclekeeper.GetFetchPriceResult( ctx, 1)
	for i, j := range data.Rates {
	prices[i] = j
	}
	/*prices[0] = 38
	prices[1] = 40000
	prices[2] = 3000
	prices[3] = 200*/

	value1 := k.cdc.MustMarshal(
		&protobuftypes.UInt64Value{
			Value: prices[0],
		},
	)
	store.Set(key1, value1)

	value2 := k.cdc.MustMarshal(
		&protobuftypes.UInt64Value{
			Value: prices[1],
		},
	)
	store.Set(key2, value2)

	value3 := k.cdc.MustMarshal(
		&protobuftypes.UInt64Value{
			Value: prices[2],
		},
	)
	store.Set(key3, value3)

	value4 := k.cdc.MustMarshal(
		&protobuftypes.UInt64Value{
			Value: prices[3],
		},
	)
	store.Set(key4, value4)
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

func (k *Keeper) getrates(ctx sdk.Context, symbol string) (uint64, bool){

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

func (k *Keeper) setRates(ctx sdk.Context) {
	var (
		store = k.Store(ctx)
		key1   = types.PriceForMarketKey("atom")
		key2   = types.PriceForMarketKey("btc")
		key3   = types.PriceForMarketKey("eth")
		key4   = types.PriceForMarketKey("xrp")

	)

	var prices []uint64
	data, _ := k.bandoraclekeeper.GetFetchPriceResult(ctx, 1)
	//for i, j := range data.Rates {
	//prices[i] = j
	//}
	prices[0] = data.Rates[1]
	/*prices[0] = 38
	prices[1] = 40000
	prices[2] = 3000
	prices[3] = 200*/

	value1 := k.cdc.MustMarshal(
		&protobuftypes.UInt64Value{
			Value: prices[0],
		},
	)
	store.Set(key1, value1)

	value2 := k.cdc.MustMarshal(
		&protobuftypes.UInt64Value{
			Value: prices[1],
		},
	)
	store.Set(key2, value2)

	value3 := k.cdc.MustMarshal(
		&protobuftypes.UInt64Value{
			Value: prices[2],
		},
	)
	store.Set(key3, value3)

	value4 := k.cdc.MustMarshal(
		&protobuftypes.UInt64Value{
			Value: prices[3],
		},
	)
	store.Set(key4, value4)
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

func (k *Keeper) SetCalldata(ctx sdk.Context, id uint64, calldata types.FetchPriceCallData) {
	var (
		store = k.Store(ctx)
		key   = types.CalldataKey(id)
		value = k.cdc.MustMarshal(&calldata)
	)

	store.Set(key, value)
}

func (k *Keeper) GetCalldata(ctx sdk.Context, id uint64) (calldata types.FetchPriceCallData, found bool) {
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

/*func (k *Keeper) OnRecvPacket(ctx sdk.Context, res bandpacket.OracleResponsePacketData) error {
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
*/
/*func (k *Keeper) HasAsset(ctx sdk.Context, id uint64) bool {
	var (
		store = k.Store(ctx)
		key   = assettypes.AssetKey(id)
	)

	return store.Has(key)
}*/

func (k *Keeper) GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool) {
	market, found := k.GetMarketForAsset(ctx, id)
	if !found {
		return 0, false
	}

	return k.GetPriceForMarket(ctx, market.Symbol)
}
