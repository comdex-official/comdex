package types

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)


var (
	ErrorInvalidAmount        = errors.Register(ModuleName, 101, "invalid amount")
	ErrorInvalidAmountIn      = errors.Register(ModuleName, 102, "invalid amount_in")
	ErrorInvalidAmountOut     = errors.Register(ModuleName, 103, "invalid amount_out")
	ErrorInvalidFrom          = errors.Register(ModuleName, 104, "invalid from")
	ErrorInvalidID            = errors.Register(ModuleName, 105, "invalid id")
	ErrorAppIstoExtendedAppId = errors.Register(ModuleName, 106, "app id does not match with extended pair app id")
)

var (

	ErrorAssetDoesNotExist   = errors.Register(ModuleName, 201, "Asset Does not exists")
	ErrorAppMappingDoesNotExists= errors.Register(ModuleName, 202, "App Mapping Does Not exists")
	ErrorAssetNotWhiteListedForGenesisMinting=errors.Register(ModuleName, 203, "Asset not added in appMapping data for genesis minting")
	ErrorGensisMintingForTokenalreadyDone=errors.Register(ModuleName, 204, "Asset minting already done for the given app")
	ErrorBuringMakesSupplyLessThanZero=errors.Register(ModuleName, 205, "Burning request recudes the the supply to 0 or less than 0 tokens")
	ErrorMintDataNotFound   = errors.Register(ModuleName, 206, "minted data not found")
)

var (
	ErrorEmergencyShutdownIsActive  = errors.Register(ModuleName, 301, "Error Emergency Shutdown Is Active for this App")
)