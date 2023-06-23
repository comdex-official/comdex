package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	LockedVaultDoesNotExist    = sdkerrors.Register(ModuleName, 701, "locked vault does not exist with given id")
	BorrowDoesNotExist         = sdkerrors.Register(ModuleName, 702, "borrow position does not exist with given id")
	BorrowPosAlreadyLiquidated = sdkerrors.Register(ModuleName, 703, "borrow position already liquidated")
	ErrAppIDExists             = sdkerrors.Register(ModuleName, 704, "App Id exists")
	ErrAppIDDoesNotExists      = sdkerrors.Register(ModuleName, 705, "App Id does not exist")
	ErrAppIDInvalid            = sdkerrors.Register(ModuleName, 706, "App Id invalid")
	ErrVaultIDInvalid          = sdkerrors.Register(ModuleName, 707, "Vault Id invalid")
	ErrorUnknownMsgType        = sdkerrors.Register(ModuleName, 708, "Unknown msg type")
)
