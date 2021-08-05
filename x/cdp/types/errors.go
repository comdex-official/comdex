package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidField           = sdkerrors.Register(ModuleName, 101, "invalid field")
	ErrorInvalidFrom            = sdkerrors.Register(ModuleName, 102, "invalid from")
	ErrorInvalidReceiver        = sdkerrors.Register(ModuleName, 103, "invalid receiver")
	ErrorInvalidCoins           = sdkerrors.Register(ModuleName, 104, "invalid coins")
	ErrorInvalidAmount          = sdkerrors.Register(ModuleName, 105, "invalid amount")
	ErrorAccountNotFound        = sdkerrors.Register(ModuleName, 107, "account not found")
	ErrorInsufficientBalance    = sdkerrors.Register(ModuleName, 108, "insufficient balance")
	ErrorInvalidCollateralRatio = sdkerrors.Register(ModuleName, 109, "proposed collateral ratio is below liquidation ratio")
	ErrorInvalidCollateral      = sdkerrors.Register(ModuleName, 110, "collateral does not exist")
	ErrorCdpNotFound            = sdkerrors.Register(ModuleName, 111, "cdp not found")
	ErrorInvalidCDP             = sdkerrors.Register(ModuleName, 112, "cdp type does not exist")
	ErrorDenomPrefixNotFound    = sdkerrors.Register(ModuleName, 113, "denom prefix not found")
	ErrorInvalidDebt            = sdkerrors.Register(ModuleName, 114, "only one principal type per cdp")
	ErrorDebtNotSupported       = sdkerrors.Register(ModuleName, 115, "debt not supported")
	ErrorInvalidPayment         = sdkerrors.Register(ModuleName, 116, "invalid payment")
	ErrorInvalidWithdrawAmount  = sdkerrors.Register(ModuleName, 117, "withdrawal amount exceeds deposit")
)
