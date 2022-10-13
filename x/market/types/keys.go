package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "market"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
)

var (
	TwaKeyPrefix            = []byte{0x24}
)

func TwaKey(id uint64) []byte {
	return append(TwaKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
