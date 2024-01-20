package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrorUnknownProposalType             = errorsmod.Register(ModuleName, 401, "unknown proposal type")
	ErrorAssetDoesNotExist               = errorsmod.Register(ModuleName, 403, "asset does not exist")
	ErrorDuplicateCollectorDenomForApp   = errorsmod.Register(ModuleName, 404, "Collector Duplicate Denom For App")
	ErrorDuplicateAssetDenoms            = errorsmod.Register(ModuleName, 405, "Duplicate Asset Denoms")
	ErrorDataDoesNotExists               = errorsmod.Register(ModuleName, 406, "Data does not exists")
	ErrorRequestedAmtExceedsCollectedFee = errorsmod.Register(ModuleName, 407, "Requested Amt Exceeds CollectedFee")
	ErrorAppDoesNotExist                 = errorsmod.Register(ModuleName, 408, "app does not exist")
	ErrorAssetNotAddedForGenesisMinting  = errorsmod.Register(ModuleName, 409, "Asset Not Added For Genesis Minting")
	ErrorAuctionParamsNotSet             = errorsmod.Register(ModuleName, 410, "Auction Params Not Set")
	ErrorAmountCanNotBeNegative          = errorsmod.Register(ModuleName, 411, "amount cannot be negative")
	ErrorNetFeesCanNotBeNegative         = errorsmod.Register(ModuleName, 412, "NetFees cannot be negative")
	SendCoinFromModuleToModuleIsZero     = errorsmod.Register(ModuleName, 413, "Send coin from module to module is zero")
	ErrorSurplusDistributerCantbeTrue    = errorsmod.Register(ModuleName, 414, "Surplus and distributer can't be true at same time")
	ErrorSurplusDebtrCantbeTrueSameTime  = errorsmod.Register(ModuleName, 415, "Surplus and debt can't be true at same time")
	ErrorInsufficientBalance             = errorsmod.Register(ModuleName, 416, "collector module account does not have enough balance to refund")
	ErrorRefundCompleted                 = errorsmod.Register(ModuleName, 417, "refund already processed")
)
