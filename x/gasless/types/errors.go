package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// DONTCOVER

var (
	ErrorUnknownProposalType      = sdkerrors.Register(ModuleName, 10000, "unknown proposal type")
	ErrorInvalidrequest           = sdkerrors.Register(ModuleName, 10001, "invalid request")
	ErrorMaxLimitReachedByCreator = sdkerrors.Register(ModuleName, 10002, "creator reached maximum limit to create gas tanks")
)
