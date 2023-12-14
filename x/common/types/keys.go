package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "common"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_common"
)

var (
	SetContractKeyPrefix = []byte{0x11}
	GameIDKey            = []byte{0x12}
	ParamsKey            = []byte{0x13}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func ContractKey(gameId uint64) []byte {
	return append(SetContractKeyPrefix, sdk.Uint64ToBigEndian(gameId)...)
}
