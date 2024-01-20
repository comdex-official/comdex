package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrorAssetDoesNotExist                    = errorsmod.Register(ModuleName, 1201, "Asset Does not exists")
	ErrorAppMappingDoesNotExists              = errorsmod.Register(ModuleName, 1202, "App Mapping Does Not exists")
	ErrorAssetNotWhiteListedForGenesisMinting = errorsmod.Register(ModuleName, 1203, "Asset not added in appMapping data for genesis minting")
	ErrorGenesisMintingForTokenAlreadyDone    = errorsmod.Register(ModuleName, 1204, "Asset minting already done for the given app")
	ErrorBurningMakesSupplyLessThanZero       = errorsmod.Register(ModuleName, 1205, "Burning request reduces the the supply to 0 or less than 0 tokens")
	ErrorMintDataNotFound                     = errorsmod.Register(ModuleName, 1206, "minted data not found")
	ErrorInvalidAppID                         = errorsmod.Register(ModuleName, 1207, "app id can not be zero")
	ErrorInvalidAssetID                       = errorsmod.Register(ModuleName, 1208, "asset id can not be zero")
	ErrorInvalidFrom                          = errorsmod.Register(ModuleName, 1209, "invalid from")
)
