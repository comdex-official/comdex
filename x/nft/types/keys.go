package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "nft"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_nft"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	DenomKey      = "Denom-value-"
	DenomCountKey = "Denom-count-"
)

var (
	PrefixDenom = []byte{0x01}
	PrefixDenomSymbol = []byte{0x02}
	PrefixCreator = []byte{0x03}
	PrefixNFT = []byte{0x04}
	PrefixOwners = []byte{0x05}

	delimiter = []byte("/")
)

func KeyDenomID(id string) []byte {
	key := append(PrefixDenom, delimiter...)
	return append(key, []byte(id)...)
}

func KeyDenomSymbol(symbol string) []byte {
	key := append(PrefixDenomSymbol, delimiter...)
	return append(key, []byte(symbol)...)
}

func KeyDenomCreator( address sdk.AccAddress, denomId string) []byte {
	key := append(PrefixCreator, delimiter...)
	if address != nil {
		key = append(key, []byte(address)...)
		key = append(key, delimiter...)
	}
	if address != nil && len(denomId) > 0 {
		key = append(key, []byte(denomId)...)
		key = append(key, delimiter...)
	}
	return key
}

func KeyNFT(denomId, nftID string) []byte {
	key := append(PrefixNFT, delimiter...)
	if len(denomId) > 0 {
		key = append(key, []byte(denomId)...)
		key = append(key, delimiter...)
	}
	if len(nftID) > 0 && len(denomId) > 0 {
		key = append(key, []byte(nftID)...)
	}
	return key
}

func KeyOwner( address sdk.AccAddress, denomId, nftID string) []byte {
	key := append(PrefixOwners, delimiter...)
	if address != nil {
		key = append(key, []byte(address.String())...)
		key = append(key, delimiter...)
	}
	if address != nil && len(denomId) > 0 {
		key = append(key, []byte(denomId)...)
		key = append(key, delimiter...)
	}
	if address != nil && len(denomId) > 0 && len(nftID) > 0 {
		key = append(key, []byte(nftID)...)
	}
	return key
}