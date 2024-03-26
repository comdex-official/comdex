package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorUnknownProposalType             = sdkerrors.Register(ModuleName, 401, "unknown proposal type")
	ErrorAssetDoesNotExist               = sdkerrors.Register(ModuleName, 403, "asset does not exist")
	ErrorDuplicateCollectorDenomForApp   = sdkerrors.Register(ModuleName, 404, "Collector Duplicate Denom For App")
	ErrorDuplicateAssetDenoms            = sdkerrors.Register(ModuleName, 405, "Duplicate Asset Denoms")
	ErrorDataDoesNotExists               = sdkerrors.Register(ModuleName, 406, "Data does not exists")
	ErrorRequestedAmtExceedsCollectedFee = sdkerrors.Register(ModuleName, 407, "Requested Amt Exceeds CollectedFee")
	ErrorAppDoesNotExist                 = sdkerrors.Register(ModuleName, 408, "app does not exist")
	ErrorAssetNotAddedForGenesisMinting  = sdkerrors.Register(ModuleName, 409, "Asset Not Added For Genesis Minting")
	ErrorAuctionParamsNotSet             = sdkerrors.Register(ModuleName, 410, "Auction Params Not Set")
	ErrorAmountCanNotBeNegative          = sdkerrors.Register(ModuleName, 411, "amount cannot be negative")
	ErrorNetFeesCanNotBeNegative         = sdkerrors.Register(ModuleName, 412, "NetFees cannot be negative")
	SendCoinFromModuleToModuleIsZero     = sdkerrors.Register(ModuleName, 413, "Send coin from module to module is zero")
	ErrorSurplusDistributerCantbeTrue    = sdkerrors.Register(ModuleName, 414, "Surplus and distributer can't be true at same time")
	ErrorSurplusDebtrCantbeTrueSameTime  = sdkerrors.Register(ModuleName, 415, "Surplus and debt can't be true at same time")
	ErrorInsufficientBalance             = sdkerrors.Register(ModuleName, 416, "collector module account does not have enough balance to refund")
	ErrorRefundCompleted                 = sdkerrors.Register(ModuleName, 417, "refund already processed")
	ErrorAmountCanNotBeZero              = sdkerrors.Register(ModuleName, 418, "amount cannot be zero")
)
