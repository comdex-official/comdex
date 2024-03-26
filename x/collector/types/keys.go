package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "collectorV1"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for slashing.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key.
	MemStoreKey = "mem_collector"
)

var (
	TypeDepositRequest                 = ModuleName + ":deposit"
	TypeRefundRequest                  = ModuleName + ":refund"
	TypeUpdateDebtParamsRequest        = ModuleName + ":update"
	AddCollectorLookupKey              = []byte{0x01}
	AppIDToAssetCollectorMappingPrefix = []byte{0x03}
	AppIDToAuctionMappingPrefix        = []byte{0x05}
	CollectorForDenomKeyPrefix         = []byte{0x07}
	NetFeeCollectedDataPrefix          = []byte{0x08}
	RefundCounterStatusPrefix          = []byte{0x09}
	SlotsKeyPrefix                     = []byte{0x10}
)

func CollectorLookupTableMappingKey(appID, assetID uint64) []byte {
	return append(append(AddCollectorLookupKey, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func CollectorLookupTableMappingByAppKey(appID uint64) []byte {
	return append(AddCollectorLookupKey, sdk.Uint64ToBigEndian(appID)...)
}

func AppidToAssetCollectorMappingKey(appID, assetID uint64) []byte {
	return append(append(AppIDToAssetCollectorMappingPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func AppIDToAuctionMappingKey(appID, assetID uint64) []byte {
	return append(append(AppIDToAuctionMappingPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func NetFeeCollectedDataKey(appID, assetID uint64) []byte {
	return append(append(NetFeeCollectedDataPrefix, sdk.Uint64ToBigEndian(appID)...), sdk.Uint64ToBigEndian(assetID)...)
}

func AppNetFeeCollectedDataKey(appID uint64) []byte {
	return append(NetFeeCollectedDataPrefix, sdk.Uint64ToBigEndian(appID)...)
}

func CollectorForDenomKey(appID uint64) []byte {
	return append(CollectorForDenomKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}
