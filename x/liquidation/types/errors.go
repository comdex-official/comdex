package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	LockedVaultDoesNotExist    = errorsmod.Register(ModuleName, 701, "locked vault does not exist with given id")
	BorrowDoesNotExist         = errorsmod.Register(ModuleName, 702, "borrow position does not exist with given id")
	BorrowPosAlreadyLiquidated = errorsmod.Register(ModuleName, 703, "borrow position already liquidated")
	ErrAppIDExists             = errorsmod.Register(ModuleName, 704, "App Id exists")
	ErrAppIDDoesNotExists      = errorsmod.Register(ModuleName, 705, "App Id does not exist")
	ErrAppIDInvalid            = errorsmod.Register(ModuleName, 706, "App Id invalid")
	ErrVaultIDInvalid          = errorsmod.Register(ModuleName, 707, "Vault Id invalid")
	ErrorUnknownMsgType        = errorsmod.Register(ModuleName, 708, "Unknown msg type")
)
