package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "vault"
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
	StoreKey     = ModuleName
)

var (
	TypeMsgCreateRequest    = ModuleName + ":create"
	TypeMsgDepositRequest   = ModuleName + ":deposit"
	TypeMsgWithdrawRequest  = ModuleName + ":withdraw"
	TypeMsgDrawRequest      = ModuleName + ":draw"
	TypeMsgRepayRequest     = ModuleName + ":repay"
	TypeMsgLiquidateRequest = ModuleName + ":liquidate"
)

var (
	IDKey                          = []byte{0x00}
	VaultKeyPrefix                 = []byte{0x10}
	VaultForAddressByPairKeyPrefix = []byte{0x20}
	UserVaultsForAddressKeyPrefix  = []byte{0x30}
	CAssetMintStatisticsKeyPrefix  = []byte{0x50}
	VaultIdsForCollateralKeyPrefix = []byte{0x60}
)

func VaultKey(id uint64) []byte {
	return append(VaultKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func VaultForAddressByPair(address sdk.AccAddress, pairID uint64) []byte {
	v := append(VaultForAddressByPairKeyPrefix, address.Bytes()...)
	if len(v) != 1+20 {
		panic(fmt.Errorf("invalid key length %d; expected %d", len(v), 1+20))
	}

	return append(v, sdk.Uint64ToBigEndian(pairID)...)
}

func UserVaultsForAddressKey(address string) []byte {

	return append(UserVaultsForAddressKeyPrefix, address...)
}

func VaultIdsForCollateral(collateralDenom string) []byte {

	return append(VaultIdsForCollateralKeyPrefix, collateralDenom...)
}

func CAssetMintRecordsKey(collateralDenom string) []byte {

	return append(CAssetMintStatisticsKeyPrefix, collateralDenom...)
}
