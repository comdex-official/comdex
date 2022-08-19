package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "liquidation"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
	MemStoreKey    = ModuleName
)

var (
	LockedVaultIDKey                 = []byte{0x01}
	LockedVaultKeyPrefix             = []byte{0x11}
	LockedVaultKeyHistory            = []byte{0x12}
	AppIdsKeyPrefix                  = []byte{0x15}
	AppLockedVaultMappingKeyPrefix   = []byte{0x13}
	AppIDLockedVaultMappingKeyPrefix = []byte{0x14}
	LiquidationOffsetHolderKeyPrefix = []byte{0x15}
)

// LengthPrefixString returns length-prefixed bytes representation
// of a string.
func LengthPrefixString(s string) []byte {
	bz := []byte(s)
	bzLen := len(bz)
	if bzLen == 0 {
		return bz
	}
	return append([]byte{byte(bzLen)}, bz...)
}

func LockedVaultKey(id uint64) []byte {
	return append(LockedVaultKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
func LockedVaultHistoryKey(id uint64) []byte {
	return append(LockedVaultKeyHistory, sdk.Uint64ToBigEndian(id)...)
}

func AppLockedVaultMappingKey(appMappingID uint64) []byte {
	return append(AppLockedVaultMappingKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}

func AppIDLockedVaultMappingKey(appMappingID uint64) []byte {
	return append(AppIDLockedVaultMappingKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}

// GetLiquidationOffsetHolderKey returns the index key to look offset value for liquidation.
func GetLiquidationOffsetHolderKey(appID uint64, liquidatonForPrefix string) []byte {
	return append(append(LiquidationOffsetHolderKeyPrefix, sdk.Uint64ToBigEndian(appID)...), LengthPrefixString(liquidatonForPrefix)...)
}
