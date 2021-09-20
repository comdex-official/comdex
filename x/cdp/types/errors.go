package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidAmount    = errors.Register(ModuleName, 101, "invalid amount")
	ErrorInvalidAmountIn  = errors.Register(ModuleName, 102, "invalid amount_in")
	ErrorInvalidAmountOut = errors.Register(ModuleName, 103, "invalid amount_out")
	ErrorInvalidFrom      = errors.Register(ModuleName, 104, "invalid from")
	ErrorInvalidID        = errors.Register(ModuleName, 105, "invalid id")
)

var (
	ErrorAssetDoesNotExist  = errors.Register(ModuleName, 201, "asset does not exist")
	ErrorCDPDoesNotExist    = errors.Register(ModuleName, 102, "cdp does not exist")
	ErrorUnauthorized       = errors.Register(ModuleName, 102, "unauthorized")
	ErrorDuplicateCDP       = errors.Register(ModuleName, 202, "duplicate cdp")
	ErrorInvalidAmountRatio = errors.Register(ModuleName, 203, "invalid amount ratio")
	ErrorPairDoesNotExist   = errors.Register(ModuleName, 204, "pair does not exist")
	ErrorPriceDoesNotExist  = errors.Register(ModuleName, 203, "price does not exist")
)

var (
	ErrorUnknownMsgType = errors.Register(ModuleName, 301, "unknown message type")
)
