package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidField         = errors.Register(ModuleName, 101, "invalid field")
	ErrorInvalidFrom          = errors.Register(ModuleName, 102, "invalid from")
	ErrorInvalidSourcePort    = errors.Register(ModuleName, 103, "invalid source_port")
	ErrorInvalidSourceChannel = errors.Register(ModuleName, 104, "invalid source_channel")
	ErrorInvalidSymbols       = errors.Register(ModuleName, 105, "invalid symbols")
	ErrorInvalidScriptID      = errors.Register(ModuleName, 106, "invalid script_id")
)

var (
	ErrorUnknownProposalType = errors.Register(ModuleName, 401, "unknown proposal type")
)

var (
	ErrorInvalidVersion   = errors.Register(ModuleName, 501, "invalid version")
	ErrorMaxAssetChannels = errors.Register(ModuleName, 502, "max asset channels")
)
