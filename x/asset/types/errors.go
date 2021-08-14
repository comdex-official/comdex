package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorUnknownProposalType = errors.Register(ModuleName, 401, "unknown proposal type")
)
