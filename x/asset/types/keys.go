package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "asset"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
)

var (
	CountKey      = []byte{0x00}
	PoolKeyPrefix = []byte{0x10}
)

func PoolKey(id uint64) []byte {
	return append(PoolKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
