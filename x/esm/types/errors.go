package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/esm module sentinel errors
var (
	ErrInvalidAsset                = errorsmod.Register(ModuleName, 501, "invalid asset")
	ErrorDuplicateESMTriggerParams = errorsmod.Register(ModuleName, 502, "Duplicate ESM Trigger Params for AppID")
	ErrAppDataNotFound             = errorsmod.Register(ModuleName, 503, "App Data Not Found")
	ErrBadOfferCoinType            = errorsmod.Register(ModuleName, 504, "Bad Offer Coin Type")
	ErrESMTriggerParamsNotFound    = errorsmod.Register(ModuleName, 505, "ESM Trigger Params Not Found")
	ErrAmtExceedsTargetValue       = errorsmod.Register(ModuleName, 506, "Amount Exceeds Target Value")
	ErrDepositForAppReached        = errorsmod.Register(ModuleName, 507, "Deposit For AppID Reached")
	ErrESMAlreadyExecuted          = errorsmod.Register(ModuleName, 508, "ESM Already Executed")
	ErrCircuitBreakerEnabled       = errorsmod.Register(ModuleName, 509, "Circuit breaker is triggered")
	ErrorUnauthorized              = errorsmod.Register(ModuleName, 510, "Unauthorized")
	ErrorAppDoesNotExists          = errorsmod.Register(ModuleName, 511, "App Does Not Exists")
	ErrAppIDDoesNotExists          = errorsmod.Register(ModuleName, 512, "App Id Does Not exist")
	ErrCoolOffPeriodPassed         = errorsmod.Register(ModuleName, 513, "Cool off period has passed")
	ErrCurrentDepositNotReached    = errorsmod.Register(ModuleName, 514, "Current Deposit Not Reached for App")
	ErrMarketDataNotFound          = errorsmod.Register(ModuleName, 515, "MarketData not found for App")
	ErrCoolOffPeriodRemains        = errorsmod.Register(ModuleName, 516, "Cool off period remaining")
	ErrorInvalidAmount             = errorsmod.Register(ModuleName, 517, "invalid amount")
	ErrorInvalidFrom               = errorsmod.Register(ModuleName, 518, "invalid from")
	ErrESMParamsNotFound           = errorsmod.Register(ModuleName, 519, "ESM Params Not Found")
	ErrDepositForAppNotFound       = errorsmod.Register(ModuleName, 520, "Deposit For AppID not found")
	ErrPriceNotFound               = errorsmod.Register(ModuleName, 521, "Price not found")
)
