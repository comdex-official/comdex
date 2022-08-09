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
	LockerProductAssetMappingKeyPrefix        = []byte{0x10}
	LockerLookupTableKeyPrefix                = []byte{0x12}
	UserLockerAssetMappingKeyPrefix           = []byte{0x14}
	LockerKeyPrefix                           = []byte{0x15}
	LockerTotalRewardsByAssetAppWiseKeyPrefix = []byte{0x16}
)

func LockerProductAssetMappingKey(id uint64) []byte {
	return append(LockerProductAssetMappingKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func LockerTotalRewardsByAssetAppWiseKey(appID, assetID uint64) []byte {
	return append(append(LockerTotalRewardsByAssetAppWiseKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}
func LockerLookupTableKey(id uint64) []byte {
	return append(LockerLookupTableKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func UserLockerAssetMappingKey(address string) []byte {
	return append(UserLockerAssetMappingKeyPrefix, address...)
}

func LockerKey(lockerID uint64) []byte {
	return append(LockerKeyPrefix, sdk.Uint64ToBigEndian(lockerID)...)
}
