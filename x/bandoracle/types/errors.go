package types

// DONTCOVER
import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorAssetDoesNotExist = errors.Register(ModuleName, 201, "asset does not exist")
)

var (
	ErrorUnknownProposalType = errors.Register(ModuleName, 401, "unknown proposal type")
)

var (
	ErrSample         = errors.Register(ModuleName, 1100, "sample error")
	ErrInvalidVersion = errors.Register(ModuleName, 1501, "invalid version")
)
