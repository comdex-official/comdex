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
	AssetIDKey = []byte{0x01}
	PairIDKey  = []byte{0x03}

	AssetKeyPrefix = []byte{0x11}
	PairKeyPrefix  = []byte{0x14}

	AssetForDenomKeyPrefix = []byte{0x21}

	WhitelistAssetIDKey     = []byte{0x31}
	WhitelistAssetKeyPrefix = []byte{0x32}
	WhitelistPairIDKey      = []byte{0x33}
	WhitelistPairKeyPrefix  = []byte{0x34}
)

func AssetKey(id uint64) []byte {
	return append(AssetKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func WhitelistAssetKey(id uint64) []byte {
	return append(WhitelistAssetKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func AssetForDenomKey(denom string) []byte {
	return append(AssetForDenomKeyPrefix, []byte(denom)...)
}

func PairKey(id uint64) []byte {
	return append(PairKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func WhitelistPairKey(id uint64) []byte {
	return append(WhitelistPairKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
