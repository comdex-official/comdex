package types

// DONTCOVER
import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var ErrorUnknownProposalType = errors.Register(ModuleName, 401, "unknown proposal type")

var (
	ErrRequestIDNotAvailable = errors.Register(ModuleName, 1100, "Request ID not available")
	ErrInvalidVersion        = errors.Register(ModuleName, 1501, "invalid version")
	ErrUnrecognisedPacket    = errors.Register(ModuleName, 1502, "Unrecognised packet")
)
