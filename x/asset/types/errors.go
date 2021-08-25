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
	ErrorDuplicateMarket            = errors.Register(ModuleName, 201, "duplicate market")
	ErrorMarketDoesNotExist         = errors.Register(ModuleName, 202, "market does not exist")
	ErrorDuplicateAsset             = errors.Register(ModuleName, 203, "duplicate asset")
	ErrorAssetDoesNotExist          = errors.Register(ModuleName, 204, "asset does not exist")
	ErrorDuplicateMarketForAsset    = errors.Register(ModuleName, 205, "duplicate market for asset")
	ErrorMarketForAssetDoesNotExist = errors.Register(ModuleName, 206, "market for asset does not exist")
	ErrorPairDoesNotExist           = errors.Register(ModuleName, 207, "pair does not exist")
	ErrorScriptIDMismatch           = errors.Register(ModuleName, 208, "script_id mismatch")
)

var (
	ErrorUnknownProposalType = errors.Register(ModuleName, 401, "unknown proposal type")
)

var (
	ErrorInvalidVersion   = errors.Register(ModuleName, 501, "invalid version")
	ErrorMaxAssetChannels = errors.Register(ModuleName, 502, "max asset channels")
)
