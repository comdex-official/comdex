package types

import (
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
)

var (
	KeyPrefixRegisteredToken     = []byte{0x01}
	KeyPrefixAdjustedBorrow      = []byte{0x02}
	KeyPrefixCollateralSetting   = []byte{0x03}
	KeyPrefixCollateralAmount    = []byte{0x04}
	KeyPrefixReserveAmount       = []byte{0x05}
	KeyPrefixLastInterestTime    = []byte{0x06}
	KeyPrefixBadDebt             = []byte{0x07}
	KeyPrefixInterestScalar      = []byte{0x08}
	KeyPrefixAdjustedTotalBorrow = []byte{0x09}
	KeyPrefixUtokenSupply        = []byte{0x0A}
)

func CreateAdjustedBorrowKeyNoDenom(borrowerAddr sdk.AccAddress) []byte {
	// borrowprefix | lengthprefixed(borrowerAddr)
	var key []byte
	key = append(key, KeyPrefixAdjustedBorrow...)
	key = append(key, address.MustLengthPrefix(borrowerAddr)...)
	return key
}

// CreateCollateralAmountKey returns a KVStore key for getting and setting the amount of
// collateral stored for a lender in a given denom.
func CreateCollateralAmountKey(lenderAddr sdk.AccAddress, uTokenDenom string) []byte {
	// collateralPrefix | lengthprefixed(lenderAddr) | denom | 0x00
	key := CreateCollateralAmountKeyNoDenom(lenderAddr)
	key = append(key, []byte(uTokenDenom)...)
	return append(key, 0) // append 0 for null-termination
}

// CreateCollateralAmountKeyNoDenom returns the common prefix used by all collateral associated
// with a given lender address.
func CreateCollateralAmountKeyNoDenom(lenderAddr sdk.AccAddress) []byte {
	// collateralPrefix | lengthprefixed(lenderAddr)
	key := CreateCollateralAmountKeyNoAddress()
	key = append(key, address.MustLengthPrefix(lenderAddr)...)
	return key
}

// CreateCollateralAmountKeyNoAddress returns a safe copy of collateralPrefix
func CreateCollateralAmountKeyNoAddress() []byte {
	// collateralPrefix
	var key []byte
	key = append(key, KeyPrefixCollateralAmount...)
	return key
}

// ReserveFundsKey returns a KVStore key for getting and setting the amount reserved of a a given token.
func ReserveFundsKey(tokenDenom string) []byte {
	// reserveamountprefix | denom | 0x00
	key := CreateReserveAmountKeyNoDenom()
	key = append(key, []byte(tokenDenom)...)
	return append(key, 0) // append 0 for null-termination
}

// CreateReserveAmountKeyNoDenom returns a safe copy of reserveAmountPrefix
func CreateReserveAmountKeyNoDenom() []byte {
	// reserveAmountPrefix
	var key []byte
	key = append(key, KeyPrefixReserveAmount...)
	return key
}

func CreateBadDebtKeyNoAddress() []byte {
	// badDebtPrefix
	var key []byte
	key = append(key, KeyPrefixBadDebt...)
	return key
}
