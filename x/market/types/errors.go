package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorAssetDoesNotExist          = errors.Register(ModuleName, 201, "asset does not exist")
	ErrorMarketForAssetDoesNotExist = errors.Register(ModuleName, 202, "market for asset does not exist")
	ErrorUnknownMsgType             = errors.Register(ModuleName, 203, "unknown message type")
	ErrorPriceNotActive             = errors.Register(ModuleName, 204, "Price inactive")
)
