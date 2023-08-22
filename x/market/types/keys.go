package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "marketV1"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
)

var (
	TwaKeyPrefix = []byte{0x24}
	RandomPrefix = []byte{0x0A}
)

func TwaKey(id uint64) []byte {
	return append(TwaKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

// GetRandomKey returns the key for the random seed for each block
func GetRandomKey(height int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(height))
	return append(RandomPrefix, b...)
}