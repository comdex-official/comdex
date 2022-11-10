package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidAmountIn  = errors.Register(ModuleName, 102, "invalid amount_in")
	ErrorInvalidAmountOut = errors.Register(ModuleName, 103, "invalid amount_out")
	ErrorInvalidFrom      = errors.Register(ModuleName, 104, "invalid depositer adddress")

	ErrorInvalidAssetID  = errors.Register(ModuleName, 107, "Invalid AssetID")
	ErrorInvalidLockerID = errors.Register(ModuleName, 108, "Invalid LockerID")
)

var (
	ErrorAssetDoesNotExist                      = errors.Register(ModuleName, 201, "asset does not exist")
	ErrorUnauthorized                           = errors.Register(ModuleName, 203, "unauthorized")
	ErrorAppMappingDoesNotExist                 = errors.Register(ModuleName, 208, "App Mapping Id does not exists")
	ErrorLockerProductAssetMappingExists        = errors.Register(ModuleName, 209, "Product mapping to this asset id already exists")
	ErrorLockerProductAssetMappingDoesNotExists = errors.Register(ModuleName, 210, "Product mapping to this asset id  does not exists")
	ErrorUserLockerAlreadyExists                = errors.Register(ModuleName, 211, "User Locker already exists, try deposit command")
	ErrorLockerDoesNotExists                    = errors.Register(ModuleName, 212, "Locker does not exists")
	ErrorRequestedAmountExceedsDepositAmount    = errors.Register(ModuleName, 213, "Not enough balance in locker")
	SendCoinsFromModuleToModuleInLockerIsZero   = errors.Register(ModuleName, 214, "Coin value in module to module transfer in locker is zero")
	SendCoinsFromModuleToAccountInLockerIsZero  = errors.Register(ModuleName, 215, "Coin value in module to account transfer in locker is zero")
	SendCoinsFromAccountToModuleInLockerIsZero  = errors.Register(ModuleName, 216, "Coin value in account to module transfer in locker is zero")
	ErrorCollectorLookupDoesNotExists           = errors.Register(ModuleName, 217, "Collector lookup does not exists")
	ErrorAppMappingIDMismatch                   = errors.Register(ModuleName, 218, "App Mapping Id mismatch")
)

var ErrorUnknownMsgType = errors.Register(ModuleName, 301, "unknown message type")
