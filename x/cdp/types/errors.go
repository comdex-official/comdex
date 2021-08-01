package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cdp module sentinel errors
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

	ErrInvalidWithdrawAmount = sdkerrors.Register(ModuleName, 13, "withdrawal amount exceeds deposit")
)
