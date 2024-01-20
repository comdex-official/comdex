package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrorInvalidAmountIn                        = errorsmod.Register(ModuleName, 901, "invalid amount_in")
	ErrorInvalidAmountOut                       = errorsmod.Register(ModuleName, 902, "invalid amount_out")
	ErrorInvalidFrom                            = errorsmod.Register(ModuleName, 903, "invalid depositer adddress")
	ErrorInvalidAssetID                         = errorsmod.Register(ModuleName, 904, "Invalid AssetID")
	ErrorInvalidLockerID                        = errorsmod.Register(ModuleName, 905, "Invalid LockerID")
	ErrorAssetDoesNotExist                      = errorsmod.Register(ModuleName, 906, "asset does not exist")
	ErrorUnauthorized                           = errorsmod.Register(ModuleName, 907, "unauthorized")
	ErrorAppMappingDoesNotExist                 = errorsmod.Register(ModuleName, 908, "App Mapping Id does not exists")
	ErrorLockerProductAssetMappingExists        = errorsmod.Register(ModuleName, 909, "Product mapping to this asset id already exists")
	ErrorLockerProductAssetMappingDoesNotExists = errorsmod.Register(ModuleName, 910, "Product mapping to this asset id  does not exists")
	ErrorUserLockerAlreadyExists                = errorsmod.Register(ModuleName, 911, "User Locker already exists, try deposit command")
	ErrorLockerDoesNotExists                    = errorsmod.Register(ModuleName, 912, "Locker does not exists")
	ErrorRequestedAmountExceedsDepositAmount    = errorsmod.Register(ModuleName, 913, "Not enough balance in locker")
	ErrorCollectorLookupDoesNotExists           = errorsmod.Register(ModuleName, 914, "Collector lookup does not exists")
	ErrorAppMappingIDMismatch                   = errorsmod.Register(ModuleName, 915, "App Mapping Id mismatch")
	ErrorUnknownMsgType                         = errorsmod.Register(ModuleName, 916, "unknown message type")
)
