package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/esm module sentinel errors
var (
	ErrSample                      = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrorUnknownProposalType       = sdkerrors.Register(ModuleName, 1101, "unknown proposal type")
	ErrorDuplicateESMTriggerParams = sdkerrors.Register(ModuleName, 1102, "Duplicate ESM Trigger Params for AppID")
)
