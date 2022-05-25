package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "incentives"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for slashing.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key.
	MemStoreKey = "mem_incentives"
)

var (
	// EpochInfoByDurationKeyPrefix defines the prefix to store EpochInfo by duration.
	EpochInfoByDurationKeyPrefix = []byte{0x00011}

	// GaugeIDKey defines key to store the next Gauge ID to be used.
	GaugeIDKey = []byte{0x00012}

	// GaugeKeyPrefix defines the prefix to store Gauge.
	GaugeKeyPrefix = []byte{0x00013}

	// GaugeIdsByTriggerDurationKeyPrefix defines the prefix to store GaugeIds by duration.
	GaugeIdsByTriggerDurationKeyPrefix = []byte{0x00014}
)

// GetEpochInfoByDurationKey returns the indexing key for EpochInfo by duration.
func GetEpochInfoByDurationKey(duration time.Duration) []byte {
	return append(EpochInfoByDurationKeyPrefix, sdk.Uint64ToBigEndian(uint64(duration))...)
}

// GetGaugeKey return the indexing key for Gauge.
func GetGaugeKey(id uint64) []byte {
	return append(GaugeKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

// GetGaugeIdsByTriggerDurationKey returns indexing key for GaugeIDs by duration.
func GetGaugeIdsByTriggerDurationKey(duration time.Duration) []byte {
	return append(GaugeIdsByTriggerDurationKeyPrefix, sdk.Uint64ToBigEndian(uint64(duration))...)
}
