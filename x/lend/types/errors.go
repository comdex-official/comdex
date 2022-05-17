package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/lend module sentinel errors
var (
	ErrInvalidAsset            = sdkerrors.Register(ModuleName, 1100, "invalid asset")
	ErrInsufficientBalance     = sdkerrors.Register(ModuleName, 1101, "insufficient balance")
	ErrBorrowLimitLow          = sdkerrors.Register(ModuleName, 1102, "borrow limit too low")
	ErrLendingPoolInsufficient = sdkerrors.Register(ModuleName, 1103, "lending pool insufficient")
	ErrInvalidRepayment        = sdkerrors.Register(ModuleName, 1104, "invalid repayment")
	ErrInvalidAddress          = sdkerrors.Register(ModuleName, 1105, "invalid address")
	ErrNegativeTotalBorrowed   = sdkerrors.Register(ModuleName, 1106, "total borrowed was negative")
	ErrInvalidUtilization      = sdkerrors.Register(ModuleName, 1107, "invalid token utilization")
	ErrLiquidationIneligible   = sdkerrors.Register(ModuleName, 1108, "borrower not eligible for liquidation")
	ErrBadValue                = sdkerrors.Register(ModuleName, 1109, "bad USD value")
	ErrLiquidatorBalanceZero   = sdkerrors.Register(ModuleName, 1110, "liquidator base asset balance is zero")
	ErrNegativeTimeElapsed     = sdkerrors.Register(ModuleName, 1111, "negative time elapsed since last interest time")
	ErrInvalidOraclePrice      = sdkerrors.Register(ModuleName, 1112, "invalid oracle price")
	ErrNegativeAPY             = sdkerrors.Register(ModuleName, 1113, "negative APY")
	ErrInvalidExchangeRate     = sdkerrors.Register(ModuleName, 1114, "exchange rate less than one")
	ErrInconsistentTotalBorrow = sdkerrors.Register(ModuleName, 1115, "total adjusted borrow inconsistency")
	ErrInvalidInteresrScalar   = sdkerrors.Register(ModuleName, 1116, "interest scalar less than one")
	ErrEmptyAddress            = sdkerrors.Register(ModuleName, 1117, "empty address")
	ErrLiquidationRewardRatio  = sdkerrors.Register(ModuleName, 1118, "requested liquidation reward not met")
	ErrorUnknownProposalType   = sdkerrors.Register(ModuleName, 1119, "unknown proposal type")
	ErrorEmptyProposalAssets   = sdkerrors.Register(ModuleName, 1120, "empty assets in proposal")
	ErrorAssetDoesNotExist     = sdkerrors.Register(ModuleName, 1121, "asset does not exist")
	ErrorDuplicateAsset        = sdkerrors.Register(ModuleName, 1122, "duplicate asset")
	ErrorPairDoesNotExist      = sdkerrors.Register(ModuleName, 1123, "pair does not exist")
	ErrorUnauthorized          = sdkerrors.Register(ModuleName, 1124, "unauthorized")
	ErrorDuplicateLend         = sdkerrors.Register(ModuleName, 1125, "duplicate Lend")
	ErrorLendDoesNotExist      = sdkerrors.Register(ModuleName, 1126, "Lend does not exists")
	ErrorDuplicateLendPair     = sdkerrors.Register(ModuleName, 1127, "duplicate lend pair")
)
