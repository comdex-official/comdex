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
	IDKey                        = []byte{0x00}
	CDPKeyPrefix                 = []byte{0x10}
	CDPForAddressByPairKeyPrefix = []byte{0x20}
)

func CDPKey(id uint64) []byte {
	return append(CDPKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}

func CDPForAddressByPair(address sdk.AccAddress, pairID uint64) []byte {
	v := append(CDPForAddressByPairKeyPrefix, address.Bytes()...)
	if len(v) != 1+20 {
		panic(fmt.Errorf("invalid key length %d; expected %d", len(v), 1+20))
	}

	return append(v, sdk.Uint64ToBigEndian(pairID)...)
}
