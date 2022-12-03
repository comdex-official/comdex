package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	LockedVaultDoesNotExist                         = sdkerrors.Register(ModuleName, 201, "locked vault does not exist with given id")
	BorrowDoesNotExist                              = sdkerrors.Register(ModuleName, 202, "borrow position does not exist with given id")
	BorrowPosAlreadyLiquidated                      = sdkerrors.Register(ModuleName, 203, "borrow position already liquidated")
	ErrAppIDExists                                  = sdkerrors.Register(ModuleName, 1101, "App Id exists")
	ErrAppIDDoesNotExists                           = sdkerrors.Register(ModuleName, 1102, "App Id does not exist")
	ErrorPriceDoesNotExist                          = sdkerrors.Register(ModuleName, 1103, "Price does not exist")
	SendCoinsFromModuleToAccountInLiquidationIsZero = sdkerrors.Register(ModuleName, 1104, "Coin value in module to account transfer in liquidation is zero")
	ErrAppIDInvalid                                 = sdkerrors.Register(ModuleName, 1105, "App Id invalid")
	ErrVaultIDInvalid                               = sdkerrors.Register(ModuleName, 1106, "Vault Id invalid")
	ErrorUnknownMsgType                             = sdkerrors.Register(ModuleName, 1107, "Unknown msg type")
	ErrSellOffAmtLessThanExpected                   = sdkerrors.Register(ModuleName, 1108, "Sell Off Amt Less Than Expected $ value")
)
