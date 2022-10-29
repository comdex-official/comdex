package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "auctionV2"
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

const (
	AuctionStartNoBids uint64 = 0
	AuctionGoingOn     uint64 = 1
	AuctionEnded       uint64 = 2
)

const (
	NoAuction             uint64 = 0
	StartedSurplusAuction uint64 = 1
	StartedDebtAuction    uint64 = 2
	SurplusString                = "surplus"
	DebtString                   = "debt"
	DutchString                  = "dutch"
)

var (
	AuctionKeyPrefix            = []byte{0x11}
	UserKeyPrefix               = []byte{0x12}
	AuctionIDKey                = []byte{0x13}
	UserBiddingsIDKey           = []byte{0x14}
	HistoryAuctionKeyPrefix     = []byte{0x15}
	HistoryUserKeyPrefix        = []byte{0x16}
	ProtocolStatisticsPrefixKey = []byte{0x17}
	AuctionParamsKeyPrefix      = []byte{0x18}
	LendAuctionIDKey            = []byte{0x19}
	LendAuctionKeyPrefix        = []byte{0x20}
	LendUserKeyPrefix           = []byte{0x21}
	LendHistoryAuctionKeyPrefix = []byte{0x22}
	LendHistoryUserKeyPrefix    = []byte{0x23}
)

func AuctionKey(appID uint64, auctionType string, auctionID uint64) []byte {
	return append(append(append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(appID)...), auctionType...), sdk.Uint64ToBigEndian(auctionID)...)
}

func UserKey(bidder string, appID uint64, auctionType string, bidID uint64) []byte {
	return append(append(append(append(UserKeyPrefix, bidder...), sdk.Uint64ToBigEndian(appID)...), auctionType...), sdk.Uint64ToBigEndian(bidID)...)
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

func HistoryUserKey(bidder string, appID uint64, auctionType string, bidID uint64) []byte {
	return append(append(append(append(HistoryUserKeyPrefix, bidder...), sdk.Uint64ToBigEndian(appID)...), auctionType...), sdk.Uint64ToBigEndian(bidID)...)
}

func HistoryUserAuctionTypeKey(bidder string, appID uint64, auctionType string) []byte {
	return append(append(append(HistoryUserKeyPrefix, bidder...), sdk.Uint64ToBigEndian(appID)...), auctionType...)
}

func HistoryAuctionTypeKey(appID uint64, auctionType string) []byte {
	return append(append(HistoryAuctionKeyPrefix, sdk.Uint64ToBigEndian(appID)...), auctionType...)
}

func ProtocolStatisticsKey(appID, assetID uint64) []byte {
	return append(append(ProtocolStatisticsPrefixKey, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func ProtocolStatisticsAppIDKey(appID uint64) []byte {
	return append(ProtocolStatisticsPrefixKey, sdk.Uint64ToBigEndian(appID)...)
}

func AuctionParamsKey(id uint64) []byte {
	return append(AuctionParamsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func LendAuctionKey(appID uint64, auctionType string, auctionID uint64) []byte {
	return append(append(append(LendAuctionKeyPrefix, sdk.Uint64ToBigEndian(appID)...), auctionType...), sdk.Uint64ToBigEndian(auctionID)...)
}

func LendUserKey(bidder string, appID uint64, auctionType string, bidID uint64) []byte {
	return append(append(append(append(LendUserKeyPrefix, bidder...), sdk.Uint64ToBigEndian(appID)...), auctionType...), sdk.Uint64ToBigEndian(bidID)...)
}

func LendUserAuctionTypeKey(bidder string, appID uint64, auctionType string) []byte {
	return append(append(append(LendUserKeyPrefix, bidder...), sdk.Uint64ToBigEndian(appID)...), auctionType...)
}

func LendAuctionTypeKey(appID uint64, auctionType string) []byte {
	return append(append(LendAuctionKeyPrefix, sdk.Uint64ToBigEndian(appID)...), auctionType...)
}

func HistoryLendAuctionKey(appID uint64, auctionType string, auctionID uint64) []byte {
	return append(append(append(LendHistoryAuctionKeyPrefix, sdk.Uint64ToBigEndian(appID)...), auctionType...), sdk.Uint64ToBigEndian(auctionID)...)
}

func HistoryLendUserKey(bidder string, appID uint64, auctionType string, bidID uint64) []byte {
	return append(append(append(append(LendHistoryUserKeyPrefix, bidder...), sdk.Uint64ToBigEndian(appID)...), auctionType...), sdk.Uint64ToBigEndian(bidID)...)
}

func HistoryLendUserAuctionTypeKey(bidder string, appID uint64, auctionType string) []byte {
	return append(append(append(LendHistoryUserKeyPrefix, bidder...), sdk.Uint64ToBigEndian(appID)...), auctionType...)
}

func HistoryLendAuctionTypeKey(appID uint64, auctionType string) []byte {
	return append(append(LendHistoryAuctionKeyPrefix, sdk.Uint64ToBigEndian(appID)...), auctionType...)
}
