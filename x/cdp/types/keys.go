package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "cdp"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_cdp"
)

var (
	TypeMsgCreateCDPRequest          = ModuleName + ":create_cdp"
	TypeMsgDepositCollateralRequest  = ModuleName + ":deposit_collateral"
	TypeMsgWithdrawCollateralRequest = ModuleName + ":withdraw_collateral"
	TypeMsgDrawDebtRequest           = ModuleName + ":draw_debt"
	TypeMsgRepayDebtRequest          = ModuleName + ":repay_debt"
	TypeMsgCloseCDPRequest           = ModuleName + ":close_cdp"
)

var (
	CdpIdIndexKeyPrefix = []byte{0x01}
	CdpKeyPrefix        = []byte{0x02}
	CdpIdKey            = []byte{0x03}
)

func GetCdpIDBytes(cdpID uint64) (cdpIDBz []byte) {
	return sdk.Uint64ToBigEndian(cdpID)
}

func GetCdpIDFromBytes(bz []byte) (cdpID uint64) {
	return sdk.BigEndianToUint64(bz)
}

func CdpKey(cdpID uint64) []byte {
	return append(CdpKeyPrefix, GetCdpIDBytes(cdpID)...)
}
