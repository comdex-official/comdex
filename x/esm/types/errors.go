package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/esm module sentinel errors
var (
	ErrInvalidAsset                = sdkerrors.Register(ModuleName, 1100, "invalid asset")
	ErrorDuplicateESMTriggerParams = sdkerrors.Register(ModuleName, 1102, "Duplicate ESM Trigger Params for AppID")
	ErrAppDataNotFound             = sdkerrors.Register(ModuleName, 1103, "App Data Not Found")
	ErrBadOfferCoinType            = sdkerrors.Register(ModuleName, 1104, "Bad Offer Coin Type")
	ErrESMTriggerParamsNotFound    = sdkerrors.Register(ModuleName, 1105, "ESM Trigger Params Not Found")
	ErrAmtExceedsTargetValue       = sdkerrors.Register(ModuleName, 1106, "Amount Exceeds Target Value")
	ErrDepositForAppReached        = sdkerrors.Register(ModuleName, 1107, "Deposit For AppID Reached")
	ErrESMAlreadyExecuted          = sdkerrors.Register(ModuleName, 1108, "ESM Already Executed")
	ErrCircuitBreakerEnabled       = sdkerrors.Register(ModuleName, 1109, "Circuit breaker is triggered")
	ErrorUnauthorized              = sdkerrors.Register(ModuleName, 1110, "Unauthorized")
	ErrorAppDoesNotExists          = sdkerrors.Register(ModuleName, 1111, "App Does Not Exists")
	ErrAppIDDoesNotExists          = sdkerrors.Register(ModuleName, 1112, "App Id Does Not exist")
	ErrCoolOffPeriodPassed         = sdkerrors.Register(ModuleName, 1113, "Cool off period has passed")
	ErrCurrentDepositNotReached    = sdkerrors.Register(ModuleName, 1114, "Current Deposit Not Reached for App")
	ErrMarketDataNotFound          = sdkerrors.Register(ModuleName, 1115, "MarketData not found for App")
	ErrCoolOffPeriodRemains        = sdkerrors.Register(ModuleName, 1116, "Cool off period remaining")
	ErrorInvalidAmount             = sdkerrors.Register(ModuleName, 1117, "invalid amount")
	ErrorInvalidFrom               = sdkerrors.Register(ModuleName, 1118, "invalid from")
	ErrESMParamsNotFound           = sdkerrors.Register(ModuleName, 1119, "ESM Params Not Found")
	ErrDepositForAppNotFound       = sdkerrors.Register(ModuleName, 1120, "Deposit For AppID not found")
	ErrPriceNotFound               = sdkerrors.Register(ModuleName, 1121, "Price not found")
)
