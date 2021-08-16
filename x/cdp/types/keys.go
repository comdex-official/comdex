package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "cdp"
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
	CountKey                          = []byte{0x00}
	CDPKeyPrefix                      = []byte{0x10}
	CDPForAddressByAssetPairKeyPrefix = []byte{0x20}
)

func CDPKey(id uint64) []byte {
	return append(CDPKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func CDPForAddressByAssetPairKey(address sdk.AccAddress, pairId uint64) []byte {
	v := append(CDPForAddressByAssetPairKeyPrefix, address.Bytes()...)
	if len(v) != 1+sdk.AddrLen {
		panic(fmt.Errorf("invalid key length %d; expected %d", len(v), 1+sdk.AddrLen))
	}

	return append(v, sdk.Uint64ToBigEndian(pairId)...)
}
