package types

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidFrom = errors.Register(ModuleName, 104, "invalid from")
)

var (
	ErrorAssetDoesNotExist                    = errors.Register(ModuleName, 201, "Asset Does not exists")
	ErrorAppMappingDoesNotExists              = errors.Register(ModuleName, 202, "App Mapping Does Not exists")
	ErrorAssetNotWhiteListedForGenesisMinting = errors.Register(ModuleName, 203, "Asset not added in appMapping data for genesis minting")
	ErrorGenesisMintingForTokenAlreadyDone    = errors.Register(ModuleName, 204, "Asset minting already done for the given app")
	ErrorBurningMakesSupplyLessThanZero       = errors.Register(ModuleName, 205, "Burning request reduces the the supply to 0 or less than 0 tokens")
	ErrorMintDataNotFound                     = errors.Register(ModuleName, 206, "minted data not found")
	ErrorMintingGenesisSupplyLessThanOne      = errors.Register(ModuleName, 207, "Mint genesis supply should be greater than zero")
)
