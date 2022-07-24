package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	LockedVaultDoesNotExist = sdkerrors.Register(ModuleName, 201, "locked vault does not exist with given id")
	ErrAppIDExists          = sdkerrors.Register(ModuleName, 1101, "App Id exists")
	ErrAppIDDoesNotExists   = sdkerrors.Register(ModuleName, 1102, "App Id does not exist")
	ErrorPriceDoesNotExist  = sdkerrors.Register(ModuleName, 1103, "Price does not exist")
	SendCoinsFromModuleToAccountInLiquidationIsZero  = sdkerrors.Register(ModuleName, 1104, "Coin value in module to account transfer in liquidation is zero")

)
