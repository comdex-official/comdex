package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DONTCOVER

var (
	ErrorUnknownProposalType = sdkerrors.Register(ModuleName, 10000, "unknown proposal type")
)
