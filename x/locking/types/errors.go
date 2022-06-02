package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/locking module sentinel errors.
var (
	ErrInvalidLockID             = sdkerrors.Register(ModuleName, 2, "lock id invalid")
	ErrInvalidLockOwner          = sdkerrors.Register(ModuleName, 3, "msg sender is not the owner of specified lock")
	ErrInvalidUnlockingCoinDenom = sdkerrors.Register(ModuleName, 4, "provided coin denom does not match with locked coin denom")
	ErrInvalidUnlockingAmount    = sdkerrors.Register(ModuleName, 5, "locked coin amount is smaller than provided coin amount")
	ErrorDuplicateLockExists     = sdkerrors.Register(ModuleName, 6, "multiple Lock Exists for same Owner,Denom and Duration")
)
