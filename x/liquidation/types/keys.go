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
	LockedVaultIdKey                 = []byte{0x01}
	LockedVaultKeyPrefix             = []byte{0x11}
	LockedVaultKeyHistory            = []byte{0x12}
	AppIdsKeyPrefix                  = []byte{0x12}
	AppLockedVaultMappingKeyPrefix   = []byte{0x13}
	AppIDLockedVaultMappingKeyPrefix = []byte{0x14}
)

func LockedVaultKey(id uint64) []byte {
	return append(LockedVaultKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
func LockedVaultHistoryKey(id uint64) []byte {
	return append(LockedVaultKeyHistory, sdk.Uint64ToBigEndian(id)...)
}

func AppIdsKey(id uint64) []byte {
	return append(AppIdsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func AppLockedVaultMappingKey(appMappingID uint64) []byte {
	return append(AppLockedVaultMappingKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}

func AppIDLockedVaultMappingKey(appMappingID uint64) []byte {
	return append(AppIDLockedVaultMappingKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}
