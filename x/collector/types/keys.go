package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "collector"

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
	AddCollectorLookupKey              = []byte{0x01}
	AppIDToAssetCollectorMappingPrefix = []byte{0x03}
	AppIDToAuctionMappingPrefix        = []byte{0x05}
	CollectorForDenomKeyPrefix         = []byte{0x07}
	NetFeeCollectedDataPrefix          = []byte{0x08}
)

func CollectorLookupTableMappingKey(appID uint64) []byte {
	return append(AddCollectorLookupKey, sdk.Uint64ToBigEndian(appID)...)
}

func AppidToAssetCollectorMappingKey(appID uint64) []byte {
	return append(AppIDToAssetCollectorMappingPrefix, sdk.Uint64ToBigEndian(appID)...)
}

func AppIDToAuctionMappingKey(appID uint64) []byte {
	return append(AppIDToAuctionMappingPrefix, sdk.Uint64ToBigEndian(appID)...)
}

func NetFeeCollectedDataKey(appID uint64) []byte {
	return append(NetFeeCollectedDataPrefix, sdk.Uint64ToBigEndian(appID)...)
}

func CollectorForDenomKey(appID uint64) []byte {
	return append(CollectorForDenomKeyPrefix, sdk.Uint64ToBigEndian(appID)...)
}
