package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorUnknownProposalType = sdkerrors.Register(ModuleName, 1100, "Proposal type is invalid")
)
