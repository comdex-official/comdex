package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidAmountIn                        = errors.Register(ModuleName, 901, "invalid amount_in")
	ErrorInvalidAmountOut                       = errors.Register(ModuleName, 902, "invalid amount_out")
	ErrorInvalidFrom                            = errors.Register(ModuleName, 903, "invalid depositer adddress")
	ErrorInvalidAssetID                         = errors.Register(ModuleName, 904, "Invalid AssetID")
	ErrorInvalidLockerID                        = errors.Register(ModuleName, 905, "Invalid LockerID")
	ErrorAssetDoesNotExist                      = errors.Register(ModuleName, 906, "asset does not exist")
	ErrorUnauthorized                           = errors.Register(ModuleName, 907, "unauthorized")
	ErrorAppMappingDoesNotExist                 = errors.Register(ModuleName, 908, "App Mapping Id does not exists")
	ErrorLockerProductAssetMappingExists        = errors.Register(ModuleName, 909, "Product mapping to this asset id already exists")
	ErrorLockerProductAssetMappingDoesNotExists = errors.Register(ModuleName, 910, "Product mapping to this asset id  does not exists")
	ErrorUserLockerAlreadyExists                = errors.Register(ModuleName, 911, "User Locker already exists, try deposit command")
	ErrorLockerDoesNotExists                    = errors.Register(ModuleName, 912, "Locker does not exists")
	ErrorRequestedAmountExceedsDepositAmount    = errors.Register(ModuleName, 913, "Not enough balance in locker")
	ErrorCollectorLookupDoesNotExists           = errors.Register(ModuleName, 913, "Collector lookup does not exists")
	ErrorAppMappingIDMismatch                   = errors.Register(ModuleName, 914, "App Mapping Id mismatch")
	ErrorUnknownMsgType                         = errors.Register(ModuleName, 915, "unknown message type")
)
