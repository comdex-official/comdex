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
)

var (
	RewardsKeyPrefix = []byte{0x05}
)

func RewardsKey(id uint64) []byte {
	return append(RewardsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}
