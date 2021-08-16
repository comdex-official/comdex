package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidField     = errors.Register(ModuleName, 101, "invalid field")
	ErrorInvalidSender    = errors.Register(ModuleName, 102, "invalid sender")
	ErrorInvalidAmountIn  = errors.Register(ModuleName, 103, "invalid amount_in")
	ErrorInvalidAmountOut = errors.Register(ModuleName, 104, "invalid amount_out")
	ErrorInvalidId        = errors.Register(ModuleName, 105, "invalid id")
	ErrorInvalidAmount    = errors.Register(ModuleName, 106, "invalid amount")
)

var (
	ErrorCDPAlreadyExists       = errors.Register(ModuleName, 201, "cdp already exists")
	ErrorAssetPairDoesNotExist  = errors.Register(ModuleName, 202, "asset pair does not exist")
	ErrorInsufficientCollateral = errors.Register(ModuleName, 203, "insufficient collateral")
)

var (
	ErrorUnknownMsgType = errors.Register(ModuleName, 301, "unknown message type")
)
