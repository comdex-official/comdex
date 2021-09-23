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
	ErrorAssetDoesNotExist             = errors.Register(ModuleName, 201, "asset does not exist")
	ErrorVaultDoesNotExist               = errors.Register(ModuleName, 202, "vault does not exist")
	ErrorUnauthorized                  = errors.Register(ModuleName, 203, "unauthorized")
	ErrorDuplicateVault                  = errors.Register(ModuleName, 204, "duplicate vault")
	ErrorInvalidCollateralizationRatio = errors.Register(ModuleName, 205, "invalid collateralization ratio")
	ErrorPairDoesNotExist              = errors.Register(ModuleName, 206, "pair does not exist")
	ErrorPriceDoesNotExist             = errors.Register(ModuleName, 207, "price does not exist")
)

var (
	ErrorUnknownMsgType = errors.Register(ModuleName, 301, "unknown message type")
)
