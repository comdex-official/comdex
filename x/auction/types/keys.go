package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName     = "auction"
	ParamsSubspace = ModuleName
	QuerierRoute   = ModuleName
	RouterKey      = ModuleName
	StoreKey       = ModuleName
	MemStoreKey    = ModuleName
)

var (
	AuctionIDKey     = []byte{0x01}
	AuctionKeyPrefix = []byte{0x11}
)

func AuctionKey(id uint64) []byte {
	return append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
