package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrSample = sdkerrors.Register(ModuleName, 1100, "sample error")

	// ErrAccountNotFound error for when no account is found for an input address
	ErrAccountNotFound = sdkerrors.Register(ModuleName, 1, "account not found")
	// ErrInsufficientBalance error for when an account does not have enough funds
	ErrInsufficientBalance = sdkerrors.Register(ModuleName, 2, "insufficient balance")
	// ErrInvalidCollateralRatio error for attempted draws that are below liquidation ratio
	ErrInvalidCollateralRatio = sdkerrors.Register(ModuleName, 3, "proposed collateral ratio is below liquidation ratio")

	ErrInvalidCollateral = sdkerrors.Register(ModuleName, 4, " collateral does not exist")

	ErrCdpNotFound = sdkerrors.Register(ModuleName, 5, "cdp not found")

	ErrInvalidCDP = sdkerrors.Register(ModuleName, 6, " cdp type does not exist")

	ErrDenomPrefixNotFound = sdkerrors.Register(ModuleName, 7, "denom prefix not found")

	ErrInvalidDebtRequest = sdkerrors.Register(ModuleName, 8, "only one principal type per cdp")

	ErrDebtNotSupported = sdkerrors.Register(ModuleName, 9, "debt not supported")

	ErrInvalidPayment = sdkerrors.Register(ModuleName, 10, "invalid payment")

	ErrInvalidWithdrawAmount = sdkerrors.Register(ModuleName, 13, "withdrawal amount exceeds deposit")
)
