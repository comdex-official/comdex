package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "collector"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_collector"
)

var (
	AddCollectorLookupKey = [] byte{0x01}
	AppidToAssetCollectorMappingPrefix = []byte{0x03}
	AppIdToAuctionMappingForAssetPrefix = []byte{0x04}
	AppIdToAuctionMappingPrefix = []byte{0x05}
	CollectorAuctionLookupPrefix = []byte{0x06}
	CollectorForDenomKeyPrefix = []byte{0x07}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}


func CollectorLookupTableMappingKey(app_id uint64) []byte {
	return append(AddCollectorLookupKey, sdk.Uint64ToBigEndian(app_id)...)
}

func AppidToAssetCollectorMappingKey(app_id uint64) []byte {
	return append(AppidToAssetCollectorMappingPrefix, sdk.Uint64ToBigEndian(app_id)...)
}
func AppIdToAuctionMappingForAssetKey(app_id uint64) []byte {
	return append(AppIdToAuctionMappingForAssetPrefix, sdk.Uint64ToBigEndian(app_id)...)
}
func AppIdToAuctionMappingKey(app_id uint64) []byte {
	return append(AppIdToAuctionMappingPrefix, sdk.Uint64ToBigEndian(app_id)...)
}
func CollectorAuctionLookupKey(app_id uint64) []byte {
	return append(CollectorAuctionLookupPrefix, sdk.Uint64ToBigEndian(app_id)...)
}

func CollectorForDenomKey(app_id uint64) []byte {
	return append(CollectorForDenomKeyPrefix, sdk.Uint64ToBigEndian(app_id)...)
}