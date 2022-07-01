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
)
