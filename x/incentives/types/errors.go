package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/incentives module sentinel errors.
var (
	ErrInvalidLockId = sdkerrors.Register(ModuleName, 1, "lock id invalid")
)
