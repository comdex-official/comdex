package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "esm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_esm"
)

var (
	ESMTriggerParamsKeyPrefix   = []byte{0x01}
	CurrentDepositStatsPrefix   = []byte{0x02}
	ESMStatusPrefix             = []byte{0x03}
	KillSwitchDataKey           = []byte{0x04}
	UserDepositByAppPrefix      = []byte{0x05}
	ESMPricePrefix              = []byte{0x06}
	ESMDataAfterCoolOffPrefix   = []byte{0x07}
	AssetToAmountValueKeyPrefix = []byte{0x08}
	AppToAmountValueKeyPrefix   = []byte{0x09}
)

func ESMTriggerParamsKey(id uint64) []byte {
	return append(ESMTriggerParamsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func ESMDataAfterCoolOff(id uint64) []byte {
	return append(ESMDataAfterCoolOffPrefix, sdk.Uint64ToBigEndian(id)...)
}

func CurrentDepositStatsKey(id uint64) []byte {
	return append(CurrentDepositStatsPrefix, sdk.Uint64ToBigEndian(id)...)
}

func ESMStatusKey(id uint64) []byte {
	return append(ESMStatusPrefix, sdk.Uint64ToBigEndian(id)...)
}

func KillSwitchData(appID uint64) []byte {
	return append(KillSwitchDataKey, sdk.Uint64ToBigEndian(appID)...)
}

func UserDepositByAppKey(owner string, id uint64) []byte {
	return append(append(UserDepositByAppPrefix, sdk.Uint64ToBigEndian(id)...), owner...)
}

func ESMSPriceKey(id uint64) []byte {
	return append(ESMPricePrefix, sdk.Uint64ToBigEndian(id)...)
}

func AssetToAmountValueKey(appID, assetID uint64) []byte {
	return append(append(AssetToAmountValueKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func AppToAmountValueKey(id uint64) []byte {
	return append(AppToAmountValueKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
