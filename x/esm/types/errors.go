package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/esm module sentinel errors
var (
	ErrInvalidAsset                = sdkerrors.Register(ModuleName, 1100, "invalid asset")
	ErrorUnknownProposalType       = sdkerrors.Register(ModuleName, 1101, "unknown proposal type")
	ErrorDuplicateESMTriggerParams = sdkerrors.Register(ModuleName, 1102, "Duplicate ESM Trigger Params for AppID")
	ErrAppDataNotFound             = sdkerrors.Register(ModuleName, 1103, "App Data Not Found")
	ErrBadOfferCoinType            = sdkerrors.Register(ModuleName, 1104, "Bad Offer Coin Type")
	ErrESMTriggerParamsNotFound    = sdkerrors.Register(ModuleName, 1105, "ESM Trigger Params Not Found")
	ErrAmtExceedsTargetValue       = sdkerrors.Register(ModuleName, 1106, "Amt Exceeds Target Value")
	ErrDepositForAppReached        = sdkerrors.Register(ModuleName, 1107, "Deposit For AppID Reached")
	ErrESMAlreadyExecuted          = sdkerrors.Register(ModuleName, 1108, "ESM Already Executed")
)
