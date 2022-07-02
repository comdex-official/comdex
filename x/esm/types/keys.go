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
	ESMTriggerParamsKeyPrefix = []byte{0x01}
	CurrentDepositStatsPrefix = []byte{0x02}
)

func ESMTriggerParamsKey(id uint64) []byte {
	return append(ESMTriggerParamsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func CurrentDepositStatsKey(id uint64) []byte {
	return append(CurrentDepositStatsPrefix, sdk.Uint64ToBigEndian(id)...)
}
