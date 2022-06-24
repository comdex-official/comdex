package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidFrom     = errors.Register(ModuleName, 103, "invalid from")
	ErrorInvalidID       = errors.Register(ModuleName, 104, "invalid id")
	ErrorInvalidScriptID = errors.Register(ModuleName, 107, "invalid script id")
	ErrorInvalidSymbol   = errors.Register(ModuleName, 110, "invalid symbol")
)

var (
	ErrorAssetDoesNotExist          = errors.Register(ModuleName, 201, "asset does not exist")
	ErrorDuplicateMarket            = errors.Register(ModuleName, 203, "duplicate market")
	ErrorMarketDoesNotExist         = errors.Register(ModuleName, 205, "market does not exist")
	ErrorMarketForAssetDoesNotExist = errors.Register(ModuleName, 206, "market for asset does not exist")
)

var (
	ErrorUnknownMsgType = errors.Register(ModuleName, 301, "unknown message type")
)

var (
	ErrorUnknownProposalType = errors.Register(ModuleName, 401, "unknown proposal type")
)
