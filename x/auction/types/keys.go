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

const NoAuction uint64 = 0
const StartedSurplusAuction uint64 = 1
const StartedDebtAuction uint64 = 2
const SurplusString = "surplus"
const DebtString = "debt"
const DutchString = "dutch"

var (
	AuctionKeyPrefix            = []byte{0x11}
	UserKeyPrefix               = []byte{0x12}
	AuctionIDKey                = []byte{0x13}
	UserBiddingsIDKey           = []byte{0x14}
	HistoryAuctionKeyPrefix     = []byte{0x15}
	HistoryUserKeyPrefix        = []byte{0x16}
	ProtocolStatisticsPrefixKey = []byte{0x17}
	AuctionParamsKeyPrefix      = []byte{0x18}
)

func AuctionKey(appID uint64, auctionType string, auctionID uint64) []byte {
	return append(append(append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(appID)...), auctionType...), sdk.Uint64ToBigEndian(auctionID)...)
}

func UserKey(bidder string, appID uint64, auctionType string, bidId uint64) []byte {
	return append(append(append(append(UserKeyPrefix, bidder...), sdk.Uint64ToBigEndian(appID)...), auctionType...), sdk.Uint64ToBigEndian(bidId)...)
}

func UserAuctionTypeKey(bidder string, appID uint64, auctionType string) []byte {
	return append(append(append(UserKeyPrefix, bidder...), sdk.Uint64ToBigEndian(appID)...), auctionType...)
}

func AuctionTypeKey(appID uint64, auctionType string) []byte {
	return append(append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(appID)...), auctionType...)
}

func HistoryAuctionKey(appID uint64, auctionType string, auctionID uint64) []byte {
	return append(append(append(HistoryAuctionKeyPrefix, sdk.Uint64ToBigEndian(appID)...), auctionType...), sdk.Uint64ToBigEndian(auctionID)...)
}

func HistoryUserKey(bidder string, appID uint64, auctionType string, bidId uint64) []byte {
	return append(append(append(append(HistoryUserKeyPrefix, bidder...), sdk.Uint64ToBigEndian(appID)...), auctionType...), sdk.Uint64ToBigEndian(bidId)...)
}

func HistoryUserAuctionTypeKey(bidder string, appID uint64, auctionType string) []byte {
	return append(append(append(HistoryUserKeyPrefix, bidder...), sdk.Uint64ToBigEndian(appID)...), auctionType...)
}

func HistoryAuctionTypeKey(appID uint64, auctionType string) []byte {
	return append(append(HistoryAuctionKeyPrefix, sdk.Uint64ToBigEndian(appID)...), auctionType...)
}

func ProtocolStatisticsKey(appID, assetId uint64) []byte {
	return append(append(ProtocolStatisticsPrefixKey, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetId)...)
}

func ProtocolStatisticsAppIdKey(appId uint64) []byte {
	return append(ProtocolStatisticsPrefixKey, sdk.Uint64ToBigEndian(appId)...)
}

func AuctionParamsKey(id uint64) []byte {
	return append(AuctionParamsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
