package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "rewardsV1"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for slashing.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key.
	MemStoreKey = "mem_rewards"

	SecondsPerYear = 31557600
	SecondsPerDay  = 86400
	DaysInYear     = "365.242"
)

var (
	RewardsKeyPrefix                    = []byte{0x05}
	AppIdsVaultKeyPrefix                = []byte{0x12}
	ExternalRewardsLockerKeyPrefix      = []byte{0x13}
	ExternalRewardsVaultKeyPrefix       = []byte{0x14}
	ExtRewardsLockerIDKey               = []byte{0x15}
	ExtRewardsVaultIDKey                = []byte{0x16}
	EpochTimeIDKey                      = []byte{0x17}
	EpochForLockerKeyPrefix             = []byte{0x20}
	ExternalRewardsLendKeyPrefix        = []byte{0x27}
	ExtRewardsLendIDKey                 = []byte{0x28}
	ExternalRewardsStableVaultKeyPrefix = []byte{0x29}
	ExtRewardsStableVaultIDKey          = []byte{0x30}

	// EpochInfoByDurationKeyPrefix defines the prefix to store EpochInfo by duration.
	EpochInfoByDurationKeyPrefix = []byte{0x21}

	// GaugeIDKey defines key to store the next Gauge ID to be used.
	GaugeIDKey = []byte{0x22}

	// GaugeKeyPrefix defines the prefix to store Gauge.
	GaugeKeyPrefix = []byte{0x23}

	// GaugeIdsByTriggerDurationKeyPrefix defines the prefix to store GaugeIds by duration.
	GaugeIdsByTriggerDurationKeyPrefix = []byte{0x24}
	LockerRewardsTrackerKeyPrefix      = []byte{0x25}
	VaultInterestTrackerKeyPrefix      = []byte{0x26}
)

// GetEpochInfoByDurationKey returns the indexing key for EpochInfo by duration.
func GetEpochInfoByDurationKey(duration time.Duration) []byte {
	return append(EpochInfoByDurationKeyPrefix, sdk.Uint64ToBigEndian(uint64(duration))...)
}

// GetGaugeKey return the indexing key for Gauge.
func AppIDKeyPrefix(appID uint64) []byte {
	return append(AppIdsVaultKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

// GetGaugeIdsByTriggerDurationKey returns indexing key for GaugeIDs by duration.
func GetGaugeIdsByTriggerDurationKey(duration time.Duration) []byte {
	return append(GaugeIdsByTriggerDurationKeyPrefix, sdk.Uint64ToBigEndian(uint64(duration))...)
}

func RewardsKey(appID, assetID uint64) []byte {
	return append(append(RewardsKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func RewardsKeyByApp(appID uint64) []byte {
	return append(RewardsKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}

func ExternalRewardsLockerMappingKey(appMappingID uint64) []byte {
	return append(ExternalRewardsLockerKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}

func EpochForLockerKey(denom uint64) []byte {
	return append(EpochForLockerKeyPrefix, sdk.Uint64ToBigEndian(denom)...)
}

func ExternalRewardsVaultMappingKey(appMappingID uint64) []byte {
	return append(ExternalRewardsVaultKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}

func ExternalRewardsLendMappingKey(appMappingID uint64) []byte {
	return append(ExternalRewardsLendKeyPrefix, sdk.Uint64ToBigEndian(appMappingID)...)
}

func LockerRewardsTrackerKey(id, appID uint64) []byte {
	return append(append(LockerRewardsTrackerKeyPrefix, sdk.Uint64ToBigEndian(id)...), sdk.Uint64ToBigEndian(appID)...)
}

func VaultInterestTrackerKey(id, appID uint64) []byte {
	return append(append(VaultInterestTrackerKeyPrefix, sdk.Uint64ToBigEndian(id)...), sdk.Uint64ToBigEndian(appID)...)
}

func GetGaugeKey(id uint64) []byte {
	return append(GaugeKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func ExternalRewardsStableVaultMappingKey(id uint64) []byte {
	return append(ExternalRewardsStableVaultKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
