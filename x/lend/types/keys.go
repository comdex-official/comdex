package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "lend"

	// ModuleAcc1 , ModuleAcc2 & ModuleAcc3  defines different module accounts to store selected pairs of asset
	ModuleAcc1 = "cmdx"
	ModuleAcc2 = "atom"
	ModuleAcc3 = "osmo"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_lend"

	CTokenPrefix = "c/"
)

var (
	KeyPrefixCollateralAmount         = []byte{0x01}
	KeyPrefixReserveAmount            = []byte{0x02}
	WhitelistedAssetIDKey             = []byte{0x03}
	WhitelistedPairIDKey              = []byte{0x04}
	LendKeyPrefix                     = []byte{0x05}
	WhitelistedAssetForDenomKeyPrefix = []byte{0x06}
	WhitelistedRecordKey              = []byte{0x07}
	PairIDKey                         = []byte{0x08}
	PairKeyPrefix                     = []byte{0x09}
	LendIDKey                         = []byte{0x10}
	KeyPrefixRegisteredToken          = []byte{0x11}
	KeyPrefixCtokenSupply             = []byte{0x12}
	LendForAddressByPairKeyPrefix     = []byte{0x13}
	LendPairKeyPrefix                 = []byte{0x14}
	LendPairIDKey                     = []byte{0x15}
)

func LendKey(id uint64) []byte {
	return append(LendKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func AssetForDenomKey(denom string) []byte {
	return append(WhitelistedAssetForDenomKeyPrefix, []byte(denom)...)
}

func CreateCollateralAmountKey(lenderAddr sdk.AccAddress, uTokenDenom string) []byte {
	key := CreateCollateralAmountKeyNoDenom(lenderAddr)
	key = append(key, []byte(uTokenDenom)...)
	return append(key, 0) // append 0 for null-termination
}

func CreateCollateralAmountKeyNoDenom(lenderAddr sdk.AccAddress) []byte {
	key := CreateCollateralAmountKeyNoAddress()
	key = append(key, address.MustLengthPrefix(lenderAddr)...)
	return key
}

func CreateCollateralAmountKeyNoAddress() []byte {
	var key []byte
	key = append(key, KeyPrefixCollateralAmount...)
	return key
}

func ReserveFundsKey(tokenDenom string) []byte {
	key := CreateReserveAmountKeyNoDenom()
	key = append(key, []byte(tokenDenom)...)
	return append(key, 0) // append 0 for null-termination
}

func CreateReserveAmountKeyNoDenom() []byte {
	var key []byte
	key = append(key, KeyPrefixReserveAmount...)
	return key
}

func PairKey(id uint64) []byte {
	return append(PairKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func CreateRegisteredTokenKey(baseTokenDenom string) []byte {
	var key []byte
	key = append(key, KeyPrefixRegisteredToken...)
	key = append(key, []byte(baseTokenDenom)...)
	return append(key, 0) // append 0 for null-termination
}

func CreateCTokenSupplyKey(uTokenDenom string) []byte {
	// supplyprefix | denom | 0x00
	var key []byte
	key = append(key, KeyPrefixCtokenSupply...)
	key = append(key, []byte(uTokenDenom)...)
	return append(key, 0) // append 0 for null-termination
}

func LendForAddressByPair(address sdk.AccAddress, pairID uint64) []byte {
	v := append(LendForAddressByPairKeyPrefix, address.Bytes()...)
	if len(v) != 1+20 {
		panic(fmt.Errorf("invalid key length %d; expected %d", len(v), 1+20))
	}

	return append(v, sdk.Uint64ToBigEndian(pairID)...)
}

func LendPairKey(id uint64) []byte {
	return append(LendPairKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
