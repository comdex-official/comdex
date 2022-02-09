package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "liquidation"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
	MemStoreKey    = ModuleName
)

var (
	LockedVaultIdKey     = []byte{0x01}
	LockedVaultKeyPrefix = []byte{0x11}
)

func LockedVaultKey(id uint64) []byte {
	return append(LockedVaultKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
