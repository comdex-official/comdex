package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "incentives"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_incentives"
)

var (
	EpochInfoByDurationKeyPrefix       = []byte{0x00011}
	GaugeIdKey                         = []byte{0x00012}
	GaugeKeyPrefix                     = []byte{0x00013}
	GaugeIdsByTriggerDurationKeyPrefix = []byte{0x00014}
)

func GetEpochInfoByDurationKey(duration time.Duration) []byte {
	return append(EpochInfoByDurationKeyPrefix, sdk.Uint64ToBigEndian(uint64(duration))...)
}

func GetGaugeKey(id uint64) []byte {
	return append(GaugeKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func GetGaugeIdsByTriggerDurationKey(duration time.Duration) []byte {
	return append(GaugeIdsByTriggerDurationKeyPrefix, sdk.Uint64ToBigEndian(uint64(duration))...)
}
