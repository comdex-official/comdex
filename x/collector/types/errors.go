package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/collector module sentinel errors
var (
	ErrorUnknownProposalType = sdkerrors.Register(ModuleName, 401, "unknown proposal type")
)
