package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "locking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_locking"
)

var (
	LockIdKey       = []byte{0x01}
	UnlockingIdKey  = []byte{0x02}
	LockKeyPrefix   = []byte{0x03}
	UnlockKeyPrefix = []byte{0x04}

	LockByOwnerKeyPrefix   = []byte{0x05}
	UnlockByOwnerKeyPrefix = []byte{0x06}
)

func GetLockKey(id uint64) []byte {
	return append(LockKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func GetLockByOwnerKey(owner string) []byte {
	return append(LockByOwnerKeyPrefix, []byte(owner)...)
}

func GetUnlockKey(id uint64) []byte {
	return append(UnlockKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func GetUnlockByOwnerKey(owner string) []byte {
	return append(UnlockByOwnerKeyPrefix, []byte(owner)...)
}
