package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "auctionsV2"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_newauc"
)

var (
	AuctionIDKey     = []byte{0x01}
	AuctionKeyPrefix = []byte{0x02}
)

func AuctionKey(appID uint64, auctionID uint64) []byte {
	return append(append(append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(auctionID)...))
}
