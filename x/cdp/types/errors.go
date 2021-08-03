package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidField         = sdkerrors.Register(ModuleName, 101, "invalid field")
	ErrorInvalidFrom          = sdkerrors.Register(ModuleName, 102, "invalid from")
	ErrorInvalidReceiver      = sdkerrors.Register(ModuleName, 103, "invalid receiver")
	ErrorInvalidCoins         = sdkerrors.Register(ModuleName, 104, "invalid coins")
	ErrorInvalidAmount        = sdkerrors.Register(ModuleName, 105, "invalid amount")
	ErrSample                 = sdkerrors.Register(ModuleName, 106, "sample error")
	ErrAccountNotFound        = sdkerrors.Register(ModuleName, 107, "account not found")
	ErrInsufficientBalance    = sdkerrors.Register(ModuleName, 108, "insufficient balance")
	ErrInvalidCollateralRatio = sdkerrors.Register(ModuleName, 109, "proposed collateral ratio is below liquidation ratio")
	ErrInvalidCollateral      = sdkerrors.Register(ModuleName, 110, "collateral does not exist")
	ErrCdpNotFound            = sdkerrors.Register(ModuleName, 111, "cdp not found")
	ErrInvalidCDP             = sdkerrors.Register(ModuleName, 112, "cdp type does not exist")
	ErrDenomPrefixNotFound    = sdkerrors.Register(ModuleName, 113, "denom prefix not found")
	ErrInvalidDebtRequest     = sdkerrors.Register(ModuleName, 114, "only one principal type per cdp")
	ErrDebtNotSupported       = sdkerrors.Register(ModuleName, 115, "debt not supported")
	ErrInvalidPayment         = sdkerrors.Register(ModuleName, 116, "invalid payment")
	ErrInvalidWithdrawAmount  = sdkerrors.Register(ModuleName, 117, "withdrawal amount exceeds deposit")
)
