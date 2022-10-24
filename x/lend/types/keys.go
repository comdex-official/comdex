package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "lendV2"

	// ModuleAcc1 , ModuleAcc2 & ModuleAcc3  defines different module accounts to store selected pairs of asset.
	ModuleAcc1 = "cmdx"
	ModuleAcc2 = "atom"
	ModuleAcc3 = "osmo"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for slashing.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key.
	MemStoreKey = "mem_lend"

	SecondsPerYear = 31557600
)

var (
	TypeLendAssetRequest                   = ModuleName + ":lend"
	TypeWithdrawAssetRequest               = ModuleName + ":withdraw"
	TypeBorrowAssetRequest                 = ModuleName + ":borrow"
	TypeRepayAssetRequest                  = ModuleName + ":repay"
	TypeFundModuleAccountRequest           = ModuleName + ":fund-module"
	TypeDepositAssetRequest                = ModuleName + ":deposit"
	TypeCloseLendAssetRequest              = ModuleName + ":close-lend"
	TypeCloseBorrowAssetRequest            = ModuleName + ":close-borrow"
	TypeDrawAssetRequest                   = ModuleName + ":draw"
	TypeDepositBorrowAssetRequest          = ModuleName + ":deposit-borrow"
	TypeBorrowAlternateAssetRequest        = ModuleName + ":borrow-alternate"
	TypeCalculateInterestAndRewardsRequest = ModuleName + ":calculate-interest-rewards"
)

var (
	NewPoolKeyPrefix         = []byte{0x61}
	NewLendUserPrefix        = []byte{0x62}
	NewLendCounterIDPrefix   = []byte{0x63}
	NewPoolIDPrefix          = []byte{0x64}
	NewLendPairIDKey         = []byte{0x18}
	NewLendPairKeyPrefix     = []byte{0x19}
	NewBorrowCounterIDPrefix = []byte{0x67}
	NewBorrowPairKeyPrefix   = []byte{0x68}
	NewAuctionParamPrefix    = []byte{0x69}

	NewAssetToPairMappingKeyPrefix           = []byte{0x20}
	NewAssetStatsByPoolIDAndAssetIDKeyPrefix = []byte{0x73}
	NewAssetRatesParamsKeyPrefix             = []byte{0x30}
	NewLendRewardsTrackerKeyPrefix           = []byte{0x75}
	NewBorrowInterestTrackerKeyPrefix        = []byte{0x76}
	NewUserLendBorrowMappingKeyPrefix        = []byte{0x77}
	NewReserveBuybackAssetDataKeyPrefix      = []byte{0x78}
)

func NewLendUserKey(ID uint64) []byte {
	return append(NewLendUserPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func NewPoolKey(ID uint64) []byte {
	return append(NewPoolKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func NewLendPairKey(ID uint64) []byte {
	return append(NewLendPairKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func NewAuctionParamKey(ID uint64) []byte {
	return append(NewAuctionParamPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func NewAssetRatesParamsKey(ID uint64) []byte {
	return append(NewAssetRatesParamsKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func NewBorrowUserKey(ID uint64) []byte {
	return append(NewBorrowPairKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func NewAssetToPairMappingKey(assetID, poolID uint64) []byte {
	return append(append(NewAssetToPairMappingKeyPrefix, sdk.Uint64ToBigEndian(assetID)...), sdk.Uint64ToBigEndian(poolID)...)
}

func NewLendRewardsTrackerKey(ID uint64) []byte {
	return append(NewLendRewardsTrackerKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func NewBorrowInterestTrackerKey(ID uint64) []byte {
	return append(NewBorrowInterestTrackerKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func NewSetAssetStatsByPoolIDAndAssetID(assetID, pairID uint64) []byte {
	v := append(NewAssetStatsByPoolIDAndAssetIDKeyPrefix, sdk.Uint64ToBigEndian(assetID)...)
	return append(v, sdk.Uint64ToBigEndian(pairID)...)
}

func NewUserLendBorrowMappingKey(owner string, lendID uint64) []byte {
	return append(append(NewUserLendBorrowMappingKeyPrefix, owner...), sdk.Uint64ToBigEndian(lendID)...)
}

func NewUserLendBorrowKey(owner string) []byte {
	return append(NewUserLendBorrowMappingKeyPrefix, owner...)
}

func NewReserveBuybackAssetDataKey(ID uint64) []byte {
	return append(NewReserveBuybackAssetDataKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}
