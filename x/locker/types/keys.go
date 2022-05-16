package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "locker"
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
	StoreKey     = ModuleName
)

var (
	TypeMsgCreateLockerRequest        = ModuleName + ":createLocker"
	TypeMsgDepositAssetRequest        = ModuleName + ":depositAsset"
	TypeMsgWithdrawAssetRequest       = ModuleName + ":withdrawAsset"
	TypeMsgAddWhiteListedAssetRequest = ModuleName + ":whitelistAsset"
)

var (
	IDKey                              = []byte{0x00}
	LockerProductAssetMappingKeyPrefix = []byte{0x10}
	LockerLookupTableKeyPrefix         = []byte{0x12}
	UserLockerAssetMappingKeyPrefix    = []byte{0x14}
	LockerKeyPrefix                    = []byte{0x15}
)

func LockerProductAssetMappingKey(id uint64) []byte {
	return append(LockerProductAssetMappingKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func LockerLookupTableKey(id uint64) []byte {
	return append(LockerLookupTableKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func UserLockerAssetMappingKey(address string) []byte {
	return append(UserLockerAssetMappingKeyPrefix, address...)
}

func LockerKey(lockerId string) []byte {
	return append(LockerKeyPrefix, lockerId...)
}
