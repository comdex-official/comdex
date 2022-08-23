package types

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidFrom = errors.Register(ModuleName, 104, "invalid from")
)

var (
	ErrorAssetDoesNotExist                        = errors.Register(ModuleName, 201, "Asset Does not exists")
	ErrorAppMappingDoesNotExists                  = errors.Register(ModuleName, 202, "App Mapping Does Not exists")
	ErrorAssetNotWhiteListedForGenesisMinting     = errors.Register(ModuleName, 203, "Asset not added in appMapping data for genesis minting")
	ErrorGenesisMintingForTokenAlreadyDone        = errors.Register(ModuleName, 204, "Asset minting already done for the given app")
	ErrorBurningMakesSupplyLessThanZero           = errors.Register(ModuleName, 205, "Burning request reduces the the supply to 0 or less than 0 tokens")
	ErrorMintDataNotFound                         = errors.Register(ModuleName, 206, "minted data not found")
	ErrorMintingGenesisSupplyLessThanOne          = errors.Register(ModuleName, 207, "Mint genesis supply should be greater than zero")
	BurnCoinValueInTokenmintIsZero                = errors.Register(ModuleName, 208, "Burn Coin value in tokenmint is zero")
	SendCoinsFromModuleToModuleInTokenmintIsZero  = errors.Register(ModuleName, 209, "Coin value in module to module transfer in tokenmint is zero")
	SendCoinsFromModuleToAccountInTokenmintIsZero = errors.Register(ModuleName, 210, "Coin value in module to account transfer in tokenmint is zero")
	SendCoinsFromAccountToModuleInTokenmintIsZero = errors.Register(ModuleName, 211, "Coin value in account to module transfer in tokenmint is zero")
	ErrorInvalidAppID                             = errors.Register(ModuleName, 212, "app id can not be zero")
	ErrorInvalidAssetID                           = errors.Register(ModuleName, 213, "asset id can not be zero")
)
