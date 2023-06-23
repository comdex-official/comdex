package types

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorAssetDoesNotExist                    = errors.Register(ModuleName, 1201, "Asset Does not exists")
	ErrorAppMappingDoesNotExists              = errors.Register(ModuleName, 1202, "App Mapping Does Not exists")
	ErrorAssetNotWhiteListedForGenesisMinting = errors.Register(ModuleName, 1203, "Asset not added in appMapping data for genesis minting")
	ErrorGenesisMintingForTokenAlreadyDone    = errors.Register(ModuleName, 1204, "Asset minting already done for the given app")
	ErrorBurningMakesSupplyLessThanZero       = errors.Register(ModuleName, 1205, "Burning request reduces the the supply to 0 or less than 0 tokens")
	ErrorMintDataNotFound                     = errors.Register(ModuleName, 1206, "minted data not found")
	ErrorInvalidAppID                         = errors.Register(ModuleName, 1207, "app id can not be zero")
	ErrorInvalidAssetID                       = errors.Register(ModuleName, 1208, "asset id can not be zero")
	ErrorInvalidFrom                          = errors.Register(ModuleName, 1209, "invalid from")
)
