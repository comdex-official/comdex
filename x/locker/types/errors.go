package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidAmount       = errors.Register(ModuleName, 101, "invalid amount")
	ErrorInvalidAmountIn     = errors.Register(ModuleName, 102, "invalid amount_in")
	ErrorInvalidAmountOut    = errors.Register(ModuleName, 103, "invalid amount_out")
	ErrorInvalidFrom         = errors.Register(ModuleName, 104, "invalid depositer adddress")
	ErrorInvalidID           = errors.Register(ModuleName, 105, "invalid id")
	ErrorInvalidAppMappingId = errors.Register(ModuleName, 106, "Invalid App Mapping ID")
	ErrorInvalidAssetID      = errors.Register(ModuleName, 107, "Invalid AssetID")
	ErrorInvalidLockerId     = errors.Register(ModuleName, 108, "Invalid LockerID")
)

var (
	ErrorAssetDoesNotExist                      = errors.Register(ModuleName, 201, "asset does not exist")
	ErrorVaultDoesNotExist                      = errors.Register(ModuleName, 202, "vault does not exist")
	ErrorUnauthorized                           = errors.Register(ModuleName, 203, "unauthorized")
	ErrorDuplicateVault                         = errors.Register(ModuleName, 204, "duplicate vault")
	ErrorInvalidCollateralizationRatio          = errors.Register(ModuleName, 205, "invalid collateralization ratio")
	ErrorPairDoesNotExist                       = errors.Register(ModuleName, 206, "pair does not exist")
	ErrorPriceDoesNotExist                      = errors.Register(ModuleName, 207, "price does not exist")
	ErrorAppMappingDoesNotExist                 = errors.Register(ModuleName, 208, "App Mapping Id does not exists")
	ErrorLockerProductAssetMappingExists        = errors.Register(ModuleName, 209, "Product mapping to this asset id already exists")
	ErrorLockerProductAssetMappingDoesNotExists = errors.Register(ModuleName, 210, "Product mapping to this asset id  does not exists")
	ErrorUserLockerAlreadyExists                = errors.Register(ModuleName, 211, "User Locker already exists, try deposit command")
	ErrorLockerDoesNotExists                    = errors.Register(ModuleName, 212, "Locker does not exists")
	ErrorRequestedAmountExceedsDepositAmount    = errors.Register(ModuleName, 213, "Not enough balance in locker")
)

var (
	ErrorUnknownMsgType = errors.Register(ModuleName, 301, "unknown message type")
)

var (
	ErrorCannotCreateStableSwapVault = errors.Register(ModuleName, 401, "Cannot Create Stable Swap Vault")
	ErrorIdnotFound                  = errors.Register(ModuleName, 402, "not found")
)
