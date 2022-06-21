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
	ErrorInvalidDebtFloor        = errors.Register(ModuleName, 107, "invalid Debt Floor")
	ErrorInvalidDebtCeiling      = errors.Register(ModuleName, 108, "invalid Debt Ceiling")
	ErrorInvalidGenesisSupply    = errors.Register(ModuleName, 109, "invalid Genesis Supply")
	ErrorInvalidMinGovSupply     = errors.Register(ModuleName, 110, "invalid min gov supply")
)

var (
	ErrorAssetDoesNotExist                 = errors.Register(ModuleName, 201, "asset does not exist")
	ErrorDuplicateAsset                    = errors.Register(ModuleName, 202, "duplicate asset")
	ErrorPairDoesNotExist                  = errors.Register(ModuleName, 203, "pair does not exist")
	ErrorDuplicatePair                     = errors.Register(ModuleName, 204, "duplicate pair")
	ErrorDuplicateApp                      = errors.Register(ModuleName, 206, "duplicate app")
	ErrorPairNameForID                     = errors.Register(ModuleName, 207, "already has pair name for id in this app")
	AppIdsDoesntExist                      = errors.Register(ModuleName, 209, "app ids does not exist")
	ErrorAssetAlreadyExistingApp           = errors.Register(ModuleName, 210, "asset already exist in App")
	ErrorAssetIsOffChain                   = errors.Register(ModuleName, 211, "asset has been marked off chain")
	ErrorDebtFloorIsGreaterThanDebtCeiling = errors.Register(ModuleName, 212, "Debt Floor Is Greater Than Debt Ceiling")
	ErrorGenesisTokenExistForApp           = errors.Register(ModuleName, 213, "Genesis Token already Exist For App")
	ErrorFeeShouldNotBeGTOne               = errors.Register(ModuleName, 215, "Fee Should Not Be Greater than One and less than zero")
	ErrorExtendedPairDoesNotExistForTheApp = errors.Register(ModuleName, 216, "extended pair does not exist for the app")
)

var (
	ErrorUnknownMsgType = errors.Register(ModuleName, 301, "unknown message type")
)

var (
	ErrorUnknownProposalType        = errors.Register(ModuleName, 401, "unknown proposal type")
	ErrorEmptyProposalAssets        = errors.Register(ModuleName, 402, "empty assets in proposal")
	ErrorUnknownAppType             = errors.Register(ModuleName, 403, "unknown app type")
	ErrorProposalTitleMissing       = errors.Register(ModuleName, 404, "proposal title missing")
	ErrorProposalDescriptionMissing = errors.Register(ModuleName, 405, "proposal description missing")
)
