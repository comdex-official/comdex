package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "liquidationV1"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
	MemStoreKey    = ModuleName
)

var (
	TypeMsgLiquidateRequest          = ModuleName + ":liquidate"
	LockedVaultIDKey                 = []byte{0x01}
	LockedVaultKeyPrefix             = []byte{0x11}
	LockedVaultKeyHistory            = []byte{0x12}
	AppIdsKeyPrefix                  = []byte{0x15}
	LiquidationOffsetHolderKeyPrefix = []byte{0x16}
	LockedVaultDataKeyHistory        = []byte{0x17}
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

func LockedVaultKey(appID, lockedVaultID uint64) []byte {
	return append(append(LockedVaultKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(lockedVaultID)...)
}

func LockedVaultKeyByApp(appID uint64) []byte {
	return append(LockedVaultKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

func LockedVaultHistoryKey(appID, lockedVaultID uint64) []byte {
	return append(append(LockedVaultDataKeyHistory, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(lockedVaultID)...)
}

// GetLiquidationOffsetHolderKey returns the index key to look offset value for liquidation.
func GetLiquidationOffsetHolderKey(appID uint64, liquidatonForPrefix string) []byte {
	return append(append(LiquidationOffsetHolderKeyPrefix, sdk.Uint64ToBigEndian(appID)...), LengthPrefixString(liquidatonForPrefix)...)
}

// whitelisting kv

func WhitelistAppKeyByApp(appID uint64) []byte {
	return append(AppIdsKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}
