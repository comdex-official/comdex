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
)

var (
	ErrorAssetDoesNotExist        = errors.Register(ModuleName, 201, "asset does not exist")
	ErrorDuplicateAsset           = errors.Register(ModuleName, 202, "duplicate asset")
	ErrorPairDoesNotExist         = errors.Register(ModuleName, 203, "pair does not exist")
	ErrorDuplicatePair            = errors.Register(ModuleName, 204, "duplicate pair")
	ErrorUnauthorized             = errors.Register(ModuleName, 205, "unauthorized")
	ErrorDuplicateApp             = errors.Register(ModuleName, 206, "duplicate app")
	ErrorPairNameForID            = errors.Register(ModuleName, 207, "already has pair name for id in this app")
	ErrorExtendedPairDoesNotExist = errors.Register(ModuleName, 208, "extended pair does not exist")
	AppIdsDoesntExist             = errors.Register(ModuleName, 209, "app ids does not exist")
	ErrorAssetAlreadyExistinApp   = errors.Register(ModuleName, 210, "asset already exist in App")
)

var (
	ErrorUnknownMsgType = errors.Register(ModuleName, 301, "unknown message type")
)

var (
	ErrorUnknownProposalType = errors.Register(ModuleName, 401, "unknown proposal type")
	ErrorEmptyProposalAssets = errors.Register(ModuleName, 402, "empty assets in proposal")
	ErrorUnknownAppType      = errors.Register(ModuleName, 403, "unknown app type")
)
