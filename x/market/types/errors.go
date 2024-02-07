package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrorAssetDoesNotExist = errorsmod.Register(ModuleName, 1001, "asset does not exist")
	ErrorUnknownMsgType    = errorsmod.Register(ModuleName, 1002, "unknown message type")
	ErrorPriceNotActive    = errorsmod.Register(ModuleName, 1003, "Price inactive")
)
