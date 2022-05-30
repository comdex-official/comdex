package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "rewards"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_rewards"

	SecondsPerYear = 31536000
)

var (
	RewardsKeyPrefix                      = []byte{0x05}
	KeyPrefixLastInterestTime             = []byte{0x06}
	AppIdsVaultKeyPrefix                  = []byte{0x12}
	ExternalRewardsLockerKeyPrefix        = []byte{0x13}
	ExternalRewardsVaultKeyPrefix         = []byte{0x14}
	ExtRewardsLockerIDKey                 = []byte{0x15}
	ExtRewardsVaultIDKey                  = []byte{0x16}
	EpochTimeIDKey                        = []byte{0x17}
	ExternalRewardsLockerCounterKeyPrefix = []byte{0x18}
	AssetForDenomKeyPrefix                = []byte{0x19}
	EpochForLockerKeyPrefix               = []byte{0x20}
)

func RewardsKey(id uint64) []byte {
	return append(RewardsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func CreateLastInterestTimeKey() []byte {
	var key []byte
	key = append(key, KeyPrefixLastInterestTime...)
	return key
}

func ExternalRewardsLockerMappingKey(appMappingID uint64) []byte {
	return append(ExternalRewardsLockerKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}

func ExternalRewardsLockersCounter(appMappingID uint64) []byte {
	return append(ExternalRewardsLockerCounterKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}

func AssetForDenomKey(denom uint64) []byte {
	return append(AssetForDenomKeyPrefix, sdk.Uint64ToBigEndian(denom)...)
}

func EpochForLockerKey(denom uint64) []byte {
	return append(EpochForLockerKeyPrefix, sdk.Uint64ToBigEndian(denom)...)
}

func ExternalRewardsVaultMappingKey(appMappingID uint64) []byte {
	return append(ExternalRewardsVaultKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}
