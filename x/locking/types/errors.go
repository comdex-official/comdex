package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/locking module sentinel errors.
var (
	ErrInvalidLockId             = sdkerrors.Register(ModuleName, 1, "lock id invalid")
	ErrInvalidLockOwner          = sdkerrors.Register(ModuleName, 2, "msg sender is not the owner of specified lock")
	ErrInvalidUnlockingCoinDenom = sdkerrors.Register(ModuleName, 3, "provided coin denom does not match with locked coin denom")
	ErrInvalidUnlockingAmount    = sdkerrors.Register(ModuleName, 4, "locked coin amount is smaller than provided coin amount")
	ErrorDuplicateLockExists     = sdkerrors.Register(ModuleName, 5, "multiple Lock Exists for same Owner,Denom and Duration")
)
