package types

import (
	sdkmath "cosmossdk.io/math"
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
	MemStoreKey             = "mem_newauc"
	SurplusAuctionInitiator = "surplus"
	DebtAuctionInitiator    = "debt"
	MaxPremiumDiscount      = 30
)

var (
	TypePlaceMarketBidRequest              = ModuleName + ":market-bid-request"
	TypePlaceLimitBidRequest               = ModuleName + ":limit-bid-request"
	TypeCancelLimitBidRequest              = ModuleName + ":cancel-limit-bid-request"
	TypeWithdrawLimitBidRequest            = ModuleName + ":withdraw-limit-bid-request"
	AuctionIDKey                           = []byte{0x01}
	AuctionKeyPrefix                       = []byte{0x02}
	LimitAuctionBidIDKey                   = []byte{0x03}
	AuctionParamsKey                       = []byte{0x04}
	UserBidIDKey                           = []byte{0x05}
	UserBidKeyPrefix                       = []byte{0x06}
	AuctionHistoricalKeyPrefix             = []byte{0x07}
	UserLimitBidMappingKeyPrefix           = []byte{0x08}
	UserLimitBidMappingKeyForAddressPrefix = []byte{0x09}
	AuctionLimitBidFeeKeyPrefix            = []byte{0x10}
	ExternalAuctionLimitBidFeeKeyPrefix    = []byte{0x11}
	BidHistoricalKeyPrefix                 = []byte{0x12}
	UserBidHistoricalKeyPrefix             = []byte{0x13}
	MarketBidProtocolKeyPrefix             = []byte{0x14}
)

func AuctionKey(auctionID uint64) []byte {
	return append(append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...))
}
func AuctionLimitBidFeeKey(assetID uint64) []byte {
	return append(append(AuctionLimitBidFeeKeyPrefix, sdk.Uint64ToBigEndian(assetID)...))
}
func AuctionHistoricalKey(auctionID uint64) []byte {
	return append(append(AuctionHistoricalKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...))
}

func BidHistoricalKey(bidID uint64, address string) []byte {
	return append(append(BidHistoricalKeyPrefix, address...), sdk.Uint64ToBigEndian(bidID)...)
}

func GetBidHistoricalKey(address string) []byte {
	return append(BidHistoricalKeyPrefix, address...)
}

func UserBidKey(userBidId uint64) []byte {
	return append(append(UserBidKeyPrefix, sdk.Uint64ToBigEndian(userBidId)...))
}

func UserBidHistoricalKey(bidID uint64, address string) []byte {
	return append(append(UserBidHistoricalKeyPrefix, address...), sdk.Uint64ToBigEndian(bidID)...)
}

func GetUserBidHistoricalKey(address string) []byte {
	return append(UserBidHistoricalKeyPrefix, address...)
}

func UserLimitBidKey(debtTokenID, collateralTokenID uint64, premium sdkmath.Int, address string) []byte {
	return append(append(append(append(UserLimitBidMappingKeyPrefix, sdk.Uint64ToBigEndian(debtTokenID)...), sdk.Uint64ToBigEndian(collateralTokenID)...), sdk.Uint64ToBigEndian((premium.Uint64()))...), address...)
}

func MarketBidProtocolKey(debtTokenID, collateralTokenID uint64) []byte {
	return append(append(MarketBidProtocolKeyPrefix, sdk.Uint64ToBigEndian(debtTokenID)...), sdk.Uint64ToBigEndian(collateralTokenID)...)
}

func UserLimitBidKeyForPremium(debtTokenID, collateralTokenID uint64, premium sdkmath.Int) []byte {
	return append(append(append(UserLimitBidMappingKeyPrefix, sdk.Uint64ToBigEndian(debtTokenID)...), sdk.Uint64ToBigEndian(collateralTokenID)...), sdk.Uint64ToBigEndian((premium.Uint64()))...)
}

func LimitBidKeyForAssetID(debtTokenID, collateralTokenID uint64) []byte {
	return append(append(append(UserLimitBidMappingKeyPrefix, sdk.Uint64ToBigEndian(debtTokenID)...), sdk.Uint64ToBigEndian(collateralTokenID)...))
}

func UserLimitBidKeyForAddress(address string) []byte {
	return append(UserLimitBidMappingKeyForAddressPrefix, address...)
}

func ExternalAuctionLimitBidFeeKey(assetID uint64) []byte {
	return append(append(ExternalAuctionLimitBidFeeKeyPrefix, sdk.Uint64ToBigEndian(assetID)...))
}
