package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidDecimals         = errors.Register(ModuleName, 101, "invalid decimals")
	ErrorInvalidDenom            = errors.Register(ModuleName, 102, "invalid denom")
	ErrorInvalidFrom             = errors.Register(ModuleName, 103, "invalid from")
	ErrorInvalidID               = errors.Register(ModuleName, 104, "invalid id")
	ErrorInvalidLiquidationRatio = errors.Register(ModuleName, 105, "invalid liquidation ratio")
	ErrorInvalidName             = errors.Register(ModuleName, 106, "invalid name")
	ErrorInvalidScriptID         = errors.Register(ModuleName, 107, "invalid script id")
	ErrorInvalidSourceChannel    = errors.Register(ModuleName, 108, "invalid source channel")
	ErrorInvalidSourcePort       = errors.Register(ModuleName, 109, "invalid source port")
	ErrorInvalidSymbol           = errors.Register(ModuleName, 110, "invalid symbol")
	ErrorInvalidSymbols          = errors.Register(ModuleName, 111, "invalid symbols")
)

var (
	ErrorAssetDoesNotExist          = errors.Register(ModuleName, 201, "asset does not exist")
	ErrorDuplicateAsset             = errors.Register(ModuleName, 202, "duplicate asset")
	ErrorDuplicateMarket            = errors.Register(ModuleName, 203, "duplicate market")
	ErrorDuplicateMarketForAsset    = errors.Register(ModuleName, 204, "duplicate market for asset")
	ErrorMarketDoesNotExist         = errors.Register(ModuleName, 205, "market does not exist")
	ErrorMarketForAssetDoesNotExist = errors.Register(ModuleName, 206, "market for asset does not exist")
	ErrorPairDoesNotExist           = errors.Register(ModuleName, 207, "pair does not exist")
	ErrorScriptIDMismatch           = errors.Register(ModuleName, 208, "script id mismatch")
	ErrorUnauthorized               = errors.Register(ModuleName, 209, "unauthorized")
)

var (
	ErrorUnknownMsgType = errors.Register(ModuleName, 301, "unknown message type")
)

var (
	ErrorUnknownProposalType = errors.Register(ModuleName, 401, "unknown proposal type")
)

var (
	ErrorInvalidVersion   = errors.Register(ModuleName, 501, "invalid version")
	ErrorMaxAssetChannels = errors.Register(ModuleName, 502, "max asset channels")
)
var (
	ErrSample               = errors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout = errors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = errors.Register(ModuleName, 1501, "invalid version")
)