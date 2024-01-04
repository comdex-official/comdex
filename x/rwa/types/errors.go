package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrUserNotFound         = sdkerrors.Register(ModuleName, 101, "User not found")
	ErrCounterPartyNotFound = sdkerrors.Register(ModuleName, 102, "Counterparty not found")
	ErrAccountAddressEmpty  = sdkerrors.Register(ModuleName, 103, "Account address should not be empty")
)
