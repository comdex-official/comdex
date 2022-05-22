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

	ActiveAuctionStatus = "active"
	ClosedAuctionStatus = "inactive"

	PlacedBiddingStatus   = "placed"
	RejectedBiddingStatus = "rejected"
	SuccessBiddingStatus  = "success"
)

const AuctionStartNoBids uint64 = 0
const AuctionGoingOn uint64 = 1
const AuctionEnded uint64 = 2

var (
	CollateralAuctionIdKey     = []byte{0x01}
	CollateralAuctionKeyPrefix = []byte{0x11}
	BiddingsIdKey              = []byte{0x02}
	BiddingsKeyPrefix          = []byte{0x22}
	UserBiddingsIdKey          = []byte{0x03}
	UserBiddingsKeyPrefix      = []byte{0x33}
	DebtAuctionIdKey           = []byte{0x41}
	DebtAuctionKeyPrefix       = []byte{0x42}
	DebtBiddingsIdKey          = []byte{0x43}
	DebtBiddingsKeyPrefix      = []byte{0x44}
	DebtUserBiddingsIdKey      = []byte{0x45}
	DebtUserBiddingsKeyPrefix  = []byte{0x46}
	DutchAuctionIdKey          = []byte{0x51}
	DutchAuctionKeyPrefix      = []byte{0x52}
	DutchBiddingsIdKey         = []byte{0x53}
	DutchBiddingsKeyPrefix     = []byte{0x54}
	DutchUserBiddingsIdKey     = []byte{0x55}
	DutchUserBiddingsKeyPrefix = []byte{0x56}
)

func CollateralAuctionKey(id uint64) []byte {
	return append(CollateralAuctionKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func BiddingsKey(id uint64) []byte {
	return append(BiddingsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func UserBiddingsKey(bidder string) []byte {
	return append(UserBiddingsKeyPrefix, bidder...)
}

func DebtAuctionKey(id uint64) []byte {
	return append(DebtAuctionKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func DebtBiddingsKey(id uint64) []byte {
	return append(DebtBiddingsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func DebtUserBiddingsKey(bidder string) []byte {
	return append(DebtUserBiddingsKeyPrefix, bidder...)
}

func DutchAuctionKey(id uint64) []byte {
	return append(DutchAuctionKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func DutchBiddingsKey(id uint64) []byte {
	return append(DutchBiddingsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func DutchUserBiddingsKey(bidder string) []byte {
	return append(DutchUserBiddingsKeyPrefix, bidder...)
}
