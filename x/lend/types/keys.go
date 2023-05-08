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
	ModuleAcc4 = "axlusdc"
	ModuleAcc5 = "statom"
	ModuleAcc6 = "evmos"
	ModuleAcc7 = "gusdc"

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
	TypeFundReserveAccountRequest          = ModuleName + ":fund-reserve"
)

var (
	PoolKeyPrefix         = []byte{0x13}
	LendUserPrefix        = []byte{0x15}
	LendCounterIDPrefix   = []byte{0x16}
	PoolIDPrefix          = []byte{0x17}
	LendPairIDKey         = []byte{0x18}
	LendPairKeyPrefix     = []byte{0x19}
	BorrowCounterIDPrefix = []byte{0x25}
	BorrowPairKeyPrefix   = []byte{0x26}
	AuctionParamPrefix    = []byte{0x41}

	AssetToPairMappingKeyPrefix           = []byte{0x20}
	LendForAddressByAssetKeyPrefix        = []byte{0x22}
	BorrowForAddressByPairKeyPrefix       = []byte{0x24}
	AssetStatsByPoolIDAndAssetIDKeyPrefix = []byte{0x29}
	AssetRatesParamsKeyPrefix             = []byte{0x30}
	LendRewardsTrackerKeyPrefix           = []byte{0x43}
	BorrowInterestTrackerKeyPrefix        = []byte{0x44}
	UserLendBorrowMappingKeyPrefix        = []byte{0x45}
	ReserveBuybackAssetDataKeyPrefix      = []byte{0x46}
	NewStableBorrowIDsKeyPrefix           = []byte{0x47}
	KeyFundModBal                         = []byte{0x48}
	KeyFundReserveBal                     = []byte{0x49}
	AllReserveStatsPrefix                 = []byte{0x50}
	AssetAndPoolWiseModBalKeyPrefix       = []byte{0x51}
	DepreciatedPoolPrefix                 = []byte{0x52}
)

func LendUserKey(ID uint64) []byte {
	return append(LendUserPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func PoolKey(ID uint64) []byte {
	return append(PoolKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func LendPairKey(ID uint64) []byte {
	return append(LendPairKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func AuctionParamKey(ID uint64) []byte {
	return append(AuctionParamPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func AssetRatesParamsKey(ID uint64) []byte {
	return append(AssetRatesParamsKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func BorrowUserKey(ID uint64) []byte {
	return append(BorrowPairKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func AssetToPairMappingKey(assetID, poolID uint64) []byte {
	return append(append(AssetToPairMappingKeyPrefix, sdk.Uint64ToBigEndian(assetID)...), sdk.Uint64ToBigEndian(poolID)...)
}

func LendRewardsTrackerKey(ID uint64) []byte {
	return append(LendRewardsTrackerKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func BorrowInterestTrackerKey(ID uint64) []byte {
	return append(BorrowInterestTrackerKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func SetAssetStatsByPoolIDAndAssetID(assetID, pairID uint64) []byte {
	v := append(AssetStatsByPoolIDAndAssetIDKeyPrefix, sdk.Uint64ToBigEndian(assetID)...)
	return append(v, sdk.Uint64ToBigEndian(pairID)...)
}

func UserLendBorrowMappingKey(owner string, lendID uint64) []byte {
	return append(append(UserLendBorrowMappingKeyPrefix, owner...), sdk.Uint64ToBigEndian(lendID)...)
}

func UserLendBorrowKey(owner string) []byte {
	return append(UserLendBorrowMappingKeyPrefix, owner...)
}

func ReserveBuybackAssetDataKey(ID uint64) []byte {
	return append(ReserveBuybackAssetDataKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func AllReserveStatsKey(ID uint64) []byte {
	return append(AllReserveStatsPrefix, sdk.Uint64ToBigEndian(ID)...)
}

func FundModBalanceKey(assetID, poolID uint64) []byte {
	return append(append(AssetAndPoolWiseModBalKeyPrefix, sdk.Uint64ToBigEndian(assetID)...), sdk.Uint64ToBigEndian(poolID)...)
}

func DepreciatedPoolKey(ID uint64) []byte {
	return append(DepreciatedPoolPrefix, sdk.Uint64ToBigEndian(ID)...)
}
