package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "market"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
)

var (
	MarketKeyPrefix = []byte{0x13}

	MarketForAssetKeyPrefix = []byte{0x22}
	PriceForMarketKeyPrefix = []byte{0x23}
)

func MarketKey(symbol string) []byte {
	return append(MarketKeyPrefix, []byte(symbol)...)
}

func MarketForAssetKey(id uint64) []byte {
	return append(MarketForAssetKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func PriceForMarketKey(symbol string) []byte {
	return append(PriceForMarketKeyPrefix, []byte(symbol)...)
}
