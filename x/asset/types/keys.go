package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "asset"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
)

var (
	AssetIDKey    = []byte{0x01}
	CalldataIDKey = []byte{0x02}
	PairIDKey     = []byte{0x03}

	AssetKeyPrefix    = []byte{0x11}
	CalldataKeyPrefix = []byte{0x12}
	MarketKeyPrefix   = []byte{0x13}
	PairKeyPrefix     = []byte{0x14}

	AssetForDenomKeyPrefix  = []byte{0x21}
	MarketForAssetKeyPrefix = []byte{0x22}
	PriceForMarketKeyPrefix = []byte{0x23}
)

func AssetKey(id uint64) []byte {
	return append(AssetKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func CalldataKey(id uint64) []byte {
	return append(CalldataKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func MarketKey(symbol string) []byte {
	return append(MarketKeyPrefix, []byte(symbol)...)
}

func AssetForDenomKey(denom string) []byte {
	return append(AssetForDenomKeyPrefix, []byte(denom)...)
}

func MarketForAssetKey(id uint64) []byte {
	return append(MarketForAssetKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func PairKey(id uint64) []byte {
	return append(PairKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func PriceForMarketKey(symbol string) []byte {
	return append(PriceForMarketKeyPrefix, []byte(symbol)...)
}
