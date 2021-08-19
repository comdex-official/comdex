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
	PortKey       = []byte{0x00}
	AssetIDKey    = []byte{0x01}
	CalldataIDKey = []byte{0x02}
	PairIDKey     = []byte{0x03}

	AssetKeyPrefix    = []byte{0x10}
	CalldataKeyPrefix = []byte{0x11}
	MarketKeyPrefix   = []byte{0x12}
	PairKeyPrefix     = []byte{0x13}
	PriceKeyPrefix    = []byte{0x14}

	AssetForMarketKeyPrefix = []byte{0x20}
	MarketForAssetKeyPrefix = []byte{0x21}
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

func PriceKey(id uint64) []byte {
	return append(PriceKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func PairKey(id uint64) []byte {
	return append(PairKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func AssetForMarketKey(symbol string) []byte {
	return append(AssetForMarketKeyPrefix, []byte(symbol)...)
}

func MarketForAssetKey(id uint64) []byte {
	return append(MarketForAssetKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
