package types

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)


var (
	ErrorInvalidAmount        = errors.Register(ModuleName, 101, "invalid amount")
	ErrorInvalidAmountIn      = errors.Register(ModuleName, 102, "invalid amount_in")
	ErrorInvalidAmountOut     = errors.Register(ModuleName, 103, "invalid amount_out")
	ErrorInvalidFrom          = errors.Register(ModuleName, 104, "invalid from")
	ErrorInvalidID            = errors.Register(ModuleName, 105, "invalid id")
	ErrorAppIstoExtendedAppId = errors.Register(ModuleName, 106, "app id does not match with extended pair app id")
)