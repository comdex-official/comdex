package types

// DONTCOVER
import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorUnknownProposalType = errors.Register(ModuleName, 301, "unknown proposal type")
	ErrRequestIDNotAvailable = errors.Register(ModuleName, 302, "Request ID not available")
	ErrInvalidVersion        = errors.Register(ModuleName, 303, "invalid version")
	ErrUnrecognisedPacket    = errors.Register(ModuleName, 304, "Unrecognised packet")
)
