package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "auction"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	MemStoreKey  = ModuleName
)

var (
	CollateralAuctionIdKey     = []byte{0x01}
	CollateralAuctionKeyPrefix = []byte{0x11}
)

func CollateralAuctionKey(id uint64) []byte {
	return append(CollateralAuctionKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
