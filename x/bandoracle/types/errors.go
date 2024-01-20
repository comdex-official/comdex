package types

// DONTCOVER
import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrorUnknownProposalType = errorsmod.Register(ModuleName, 301, "unknown proposal type")
	ErrRequestIDNotAvailable = errorsmod.Register(ModuleName, 302, "Request ID not available")
	ErrInvalidVersion        = errorsmod.Register(ModuleName, 303, "invalid version")
	ErrUnrecognisedPacket    = errorsmod.Register(ModuleName, 304, "Unrecognised packet")
)
