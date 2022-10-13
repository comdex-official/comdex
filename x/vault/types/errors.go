package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidAmount    = errors.Register(ModuleName, 101, "invalid amount")
	ErrorInvalidAmountIn  = errors.Register(ModuleName, 102, "invalid amount_in")
	ErrorInvalidAmountOut = errors.Register(ModuleName, 103, "invalid amount_out")
	ErrorInvalidFrom      = errors.Register(ModuleName, 104, "invalid from")
	ErrorInvalidID        = errors.Register(ModuleName, 105, "invalid id")
)

var (

	// ErrorUnauthorized                  = errors.Register(ModuleName, 203, "unauthorized")
	// ErrorDuplicateVault                = errors.Register(ModuleName, 204, "duplicate vault")

	ErrorExtendedPairVaultDoesNotExists       = errors.Register(ModuleName, 201, "Extended pair vault does not exists for the given id")
	ErrorAppMappingDoesNotExist               = errors.Register(ModuleName, 202, "App Mapping Id does not exists")
	ErrorAppMappingIDMismatch                 = errors.Register(ModuleName, 203, "App Mapping Id mismatch, use the correct App Mapping ID in request")
	ErrorVaultCreationInactive                = errors.Register(ModuleName, 204, "Vault Creation Inactive")
	ErrorUserVaultAlreadyExists               = errors.Register(ModuleName, 205, "User vault already exists for the given extended pair vault id ")
	ErrorAmountOutLessThanDebtFloor           = errors.Register(ModuleName, 206, "Amount Out is less than Debt Floor")
	ErrorAmountOutGreaterThanDebtCeiling      = errors.Register(ModuleName, 207, "Amount Out is greater than Debt Ceiling")
	ErrorPairDoesNotExist                     = errors.Register(ModuleName, 208, "Pair does not exists")
	ErrorAssetDoesNotExist                    = errors.Register(ModuleName, 209, "Asset does not exists")
	ErrorPriceDoesNotExist                    = errors.Register(ModuleName, 210, "Price does not exist")
	ErrorInvalidCollateralizationRatio        = errors.Register(ModuleName, 211, "Invalid collateralization ratio")
	ErrorVaultDoesNotExist                    = errors.Register(ModuleName, 212, "Vault does not exist")
	ErrVaultAccessUnauthorised                = errors.Register(ModuleName, 213, "Unauthorized user for the tx")
	ErrorInvalidAppMappingData                = errors.Register(ModuleName, 214, "Invalid App Mapping data sent as compared to data exists in vault")
	ErrorInvalidExtendedPairMappingData       = errors.Register(ModuleName, 215, "Invalid Extended Pair Vault Mapping data sent as compared to data exists in vault")
	ErrorVaultInactive                        = errors.Register(ModuleName, 216, "Vault tx Inactive")
	ErrorStableMintVaultAlreadyCreated        = errors.Register(ModuleName, 217, "Stable Mint with this ExtendedPair Already exists, try deposit for stable mint")
	BurnCoinValueInVaultIsZero                = errors.Register(ModuleName, 218, "Burn Coin value in vault is zero")
	MintCoinValueInVaultIsZero                = errors.Register(ModuleName, 219, "Mint Coin value in vault is zero")
	SendCoinsFromModuleToModuleInVaultIsZero  = errors.Register(ModuleName, 220, "Coin value in module to module transfer in vault is zero")
	SendCoinsFromModuleToAccountInVaultIsZero = errors.Register(ModuleName, 221, "Coin value in module to account transfer in vault is zero")
	SendCoinsFromAccountToModuleInVaultIsZero = errors.Register(ModuleName, 222, "Coin value in account to module transfer in vault is zero")
	ErrorAppExtendedPairDataDoesNotExists     = errors.Register(ModuleName, 223, "App ExtendedPair Data Does Not Exists")
)

var ErrorUnknownMsgType = errors.Register(ModuleName, 301, "unknown message type")

var ErrorCannotCreateStableMintVault = errors.Register(ModuleName, 401, "Cannot Create Stable Mint Vault, StableMint tx command")
