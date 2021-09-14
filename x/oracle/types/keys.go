package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "oracle"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
)

var (
	CalldataIDKey = []byte{0x02}

	CalldataKeyPrefix = []byte{0x12}
	MarketKeyPrefix   = []byte{0x13}

	MarketForAssetKeyPrefix = []byte{0x22}
	PriceForMarketKeyPrefix = []byte{0x23}
)

func CalldataKey(id uint64) []byte {
	return append(CalldataKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func MarketKey(symbol string) []byte {
	return append(MarketKeyPrefix, []byte(symbol)...)
}

func MarketForAssetKey(id uint64) []byte {
	return append(MarketForAssetKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}


func PriceForMarketKey(symbol string) []byte {
	return append(PriceForMarketKeyPrefix, []byte(symbol)...)
}
