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

	CTokenPrefix   = "uc"
	SecondsPerYear = 31557600
)

var (
	WhitelistedAssetKeyPrefix = []byte{0x02}
	KeyPrefixCollateralAmount = []byte{0x04}
	KeyPrefixReserveAmount    = []byte{0x05}

	KeyPrefixCtokenSupply = []byte{0x12}
	PoolKeyPrefix         = []byte{0x13}
	PairKeyPrefix         = []byte{0x14}
	LendUserPrefix        = []byte{0x15}
	LendHistoryIdPrefix   = []byte{0x16}
	PoolIdPrefix          = []byte{0x17}
	LendPairIDKey         = []byte{0x18}
	LendPairKeyPrefix     = []byte{0x19}
	BorrowHistoryIdPrefix = []byte{0x25}
	BorrowPairKeyPrefix   = []byte{0x26}
	LendsKey              = []byte{0x32}
	BorrowsKey            = []byte{0x33}

	AssetToPairMappingKeyPrefix           = []byte{0x20}
	WhitelistedAssetForDenomKeyPrefix     = []byte{0x21}
	LendForAddressByAssetKeyPrefix        = []byte{0x22}
	UserLendsForAddressKeyPrefix          = []byte{0x23}
	BorrowForAddressByPairKeyPrefix       = []byte{0x24}
	UserBorrowsForAddressKeyPrefix        = []byte{0x27}
	LendIdToBorrowIdMappingKeyPrefix      = []byte{0x28}
	AssetStatsByPoolIdAndAssetIdKeyPrefix = []byte{0x29}
	AssetRatesStatsKeyPrefix              = []byte{0x30}
	KeyPrefixLastTime                     = []byte{0x31}
)

func AssetKey(id uint64) []byte {
	return append(WhitelistedAssetKeyPrefix, sdk.Uint64ToBigEndian(id)...)
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

func LendUserKey(id uint64) []byte {
	return append(LendUserPrefix, sdk.Uint64ToBigEndian(id)...)
}

func PoolKey(id uint64) []byte {
	return append(PoolKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func CreateCTokenSupplyKey(uTokenDenom string) []byte {
	// supplyprefix | denom | 0x00
	var key []byte
	key = append(key, KeyPrefixCtokenSupply...)
	key = append(key, []byte(uTokenDenom)...)
	return append(key, 0) // append 0 for null-termination
}

func LendPairKey(id uint64) []byte {
	return append(LendPairKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func AssetRatesStatsKey(id uint64) []byte {
	return append(AssetRatesStatsKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func BorrowUserKey(id uint64) []byte {
	return append(BorrowPairKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func AssetToPairMappingKey(assetId, poolId uint64) []byte {
	return append(append(AssetToPairMappingKeyPrefix, sdk.Uint64ToBigEndian(assetId)...), sdk.Uint64ToBigEndian(poolId)...)
}

func LendForAddressByAsset(address sdk.AccAddress, assetID, poolID uint64) []byte {
	v := append(LendForAddressByAssetKeyPrefix, address.Bytes()...)
	if len(v) != 1+20 {
		panic(fmt.Errorf("invalid key length %d; expected %d", len(v), 1+20))
	}
	return append(append(v, sdk.Uint64ToBigEndian(assetID)...), sdk.Uint64ToBigEndian(poolID)...)
}

func BorrowForAddressByPair(address sdk.AccAddress, pairID uint64) []byte {
	v := append(BorrowForAddressByPairKeyPrefix, address.Bytes()...)
	if len(v) != 1+20 {
		panic(fmt.Errorf("invalid key length %d; expected %d", len(v), 1+20))
	}

	return append(v, sdk.Uint64ToBigEndian(pairID)...)
}

func UserLendsForAddressKey(address string) []byte {
	return append(UserLendsForAddressKeyPrefix, address...)
}

func UserBorrowsForAddressKey(address string) []byte {
	return append(UserBorrowsForAddressKeyPrefix, address...)
}

func LendIdToBorrowIdMappingKey(id uint64) []byte {
	return append(LendIdToBorrowIdMappingKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func SetAssetStatsByPoolIdAndAssetId(assetID, pairID uint64) []byte {
	v := append(AssetStatsByPoolIdAndAssetIdKeyPrefix, sdk.Uint64ToBigEndian(assetID)...)
	return append(v, sdk.Uint64ToBigEndian(pairID)...)
}

func CreateLastInterestTimeKey() []byte {
	var key []byte
	key = append(key, KeyPrefixLastTime...)
	return key
}
