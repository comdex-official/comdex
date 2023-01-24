package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorAssetDoesNotExist = errors.Register(ModuleName, 1001, "asset does not exist")
	ErrorUnknownMsgType    = errors.Register(ModuleName, 1002, "unknown message type")
	ErrorPriceNotActive    = errors.Register(ModuleName, 1003, "Price inactive")
)
