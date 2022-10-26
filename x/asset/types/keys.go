package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "assetv1"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
)

var (
	AssetIDKey      = []byte{0x01}
	PairIDKey       = []byte{0x03}
	AppIDKey        = []byte{0x04}
	PairsVaultIDKey = []byte{0x05}

	AssetKeyPrefix      = []byte{0x11}
	PairKeyPrefix       = []byte{0x14}
	AppKeyPrefix        = []byte{0x26}
	PairsVaultKeyPrefix = []byte{0x16}

	AssetForDenomKeyPrefix = []byte{0x21}
	AppForShortNamePrefix  = []byte{0x22}
	AppForNamePrefix       = []byte{0x23}
	GenesisForAppPrefix    = []byte{0x24}
	AssetForNameKeyPrefix  = []byte{0x25}
)

func AppKey(id uint64) []byte {
	return append(AppKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func PairsKey(id uint64) []byte {
	return append(PairsVaultKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func AssetKey(id uint64) []byte {
	return append(AssetKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func AssetForDenomKey(denom string) []byte {
	return append(AssetForDenomKeyPrefix, []byte(denom)...)
}

func AssetForNameKey(denom string) []byte {
	return append(AssetForNameKeyPrefix, []byte(denom)...)
}

func AssetForShortNameKey(shortName string) []byte {
	return append(AppForShortNamePrefix, []byte(shortName)...)
}

func AppAssetForNameKey(Name string) []byte {
	return append(AppForNamePrefix, []byte(Name)...)
}

func GenesisForApp(appID uint64) []byte {
	return append(GenesisForAppPrefix, sdk.Uint64ToBigEndian(appID)...)
}

func PairKey(id uint64) []byte {
	return append(PairKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
