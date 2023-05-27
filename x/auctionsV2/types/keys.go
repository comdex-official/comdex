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
	TypePlaceMarketBidRequest   = ModuleName + ":market-bid-request"
	TypePlaceLimitBidRequest    = ModuleName + ":limit-bid-request"
	TypeCancelLimitBidRequest   = ModuleName + ":cancel-limit-bid-request"
	TypeWithdrawLimitBidRequest = ModuleName + ":withdraw-limit-bid-request"
	AuctionIDKey                = []byte{0x01}
	AuctionKeyPrefix            = []byte{0x02}
	LimitAuctionBidIDKey        = []byte{0x03}
	AuctionParamsKey            = []byte{0x04}
	UserBidIDKey                = []byte{0x05}
	UserBidKeyPrefix            = []byte{0x06}
	AuctionHistoricalKeyPrefix  = []byte{0x07}
)

func AuctionKey(auctionID uint64) []byte {
	return append(append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...))
}
func AuctionHistoricalKey(auctionID uint64) []byte {
	return append(append(AuctionHistoricalKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...))
}

func UserBidKey(userBidId uint64) []byte {
	return append(append(UserBidKeyPrefix, sdk.Uint64ToBigEndian(userBidId)...))
}
