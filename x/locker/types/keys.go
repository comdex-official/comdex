package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "lockerV1"
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
	LockerIDPrefix                            = []byte{0x17}
)

func LockerProductAssetMappingKey(appID, assetID uint64) []byte {
	return append(append(LockerProductAssetMappingKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func LockerProductAssetMappingByAppKey(appID uint64) []byte {
	return append(LockerProductAssetMappingKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

func LockerTotalRewardsByAssetAppWiseKey(appID, assetID uint64) []byte {
	return append(append(LockerTotalRewardsByAssetAppWiseKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func LockerTotalRewardsByAppWiseKey(appID uint64) []byte {
	return append(LockerTotalRewardsByAssetAppWiseKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

func UserAppAssetLockerMappingKey(address string, appID uint64, assetID uint64) []byte {
	return append(append(append(UserLockerAssetMappingKeyPrefix, address...), sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func UserAppLockerMappingKey(address string, appID uint64) []byte {
	return append(append(UserLockerAssetMappingKeyPrefix, address...), sdk.Uint64ToBigEndian(appID)...)
}

func UserLockerMappingKey(address string) []byte {
	return append(UserLockerAssetMappingKeyPrefix, address...)
}

func LockerKey(lockerID uint64) []byte {
	return append(LockerKeyPrefix, sdk.Uint64ToBigEndian(lockerID)...)
}

func LockerLookupTableKey(appID uint64, assetID uint64) []byte {
	return append(append(LockerLookupTableKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func LockerLookupTableByAppKey(appID uint64) []byte {
	return append(LockerLookupTableKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}
