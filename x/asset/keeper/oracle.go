package keeper

import (
	"github.com/bandprotocol/bandchain-packet/obi"
	bandpacket "github.com/bandprotocol/bandchain-packet/packet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
	"strconv"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) SetMarket(ctx sdk.Context, market types.Market) {
	var (
		store = k.Store(ctx)
		key   = types.MarketKey(market.Symbol)
		value = k.cdc.MustMarshal(&market)
	)

	store.Set(key, value)
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

func (k *Keeper) SetMarketSymbolForAsset(ctx sdk.Context, symbol string, id uint64) {
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

func (k *Keeper) GetMarketSymbolForAsset(ctx sdk.Context, id uint64) string {
	var (
		store = k.Store(ctx)
		key   = types.MarketForAssetKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return ""
	}

	var symbol protobuftypes.StringValue
	k.cdc.MustUnmarshal(value, &symbol)

	return symbol.GetValue()
}

func (k *Keeper) GetMarketForAsset(ctx sdk.Context, id uint64) (market types.Market, found bool) {
	symbol := k.GetMarketSymbolForAsset(ctx, id)
	if symbol == "" {
		return market, false
	}

	return k.GetMarket(ctx, symbol)
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

func (k *Keeper) SetCalldata(ctx sdk.Context, id uint64, data types.Calldata) {
	var (
		store = k.Store(ctx)
		key   = types.CalldataKey(id)
		value = k.cdc.MustMarshal(&data)
	)

	store.Set(key, value)
}

func (k *Keeper) GetCalldata(ctx sdk.Context, id uint64) (data types.Calldata, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.CalldataKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return data, false
	}

	k.cdc.MustUnmarshal(value, &data)
	return data, true
}

func (k *Keeper) SetPrice(ctx sdk.Context, id uint64, price uint64) {
	var (
		store = k.Store(ctx)
		key   = types.PriceKey(id)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: price,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetPrice(ctx sdk.Context, id uint64) (uint64, bool) {
	var (
		store = k.Store(ctx)
		key   = types.PriceKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return 0, false
	}

	var price protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &price)

	return price.GetValue(), true
}

func (k *Keeper) OnRecvPacket(ctx sdk.Context, res bandpacket.OracleResponsePacketData) error {
	if res.ResolveStatus == bandpacket.RESOLVE_STATUS_SUCCESS {
		id, err := strconv.ParseUint(res.ClientID, 10, 64)
		if err != nil {
			return err
		}

		calldata, found := k.GetCalldata(ctx, id)
		if !found {
			return nil
		}

		var result types.Result
		if err := obi.Decode(res.Result, &result); err != nil {
			return err
		}

		for i := range calldata.Symbols {
			id := k.GetAssetIDForMarket(ctx, calldata.Symbols[i])
			k.SetPrice(ctx, id, result.Rates[0])
		}
	}

	return nil
}
