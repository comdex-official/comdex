package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/incentives module sentinel errors.
var (
	ErrInvalidGaugeStartTime = sdkerrors.Register(ModuleName, 2, "start time smaller than current time")
	ErrInvalidGaugeTypeId    = sdkerrors.Register(ModuleName, 3, "gauge type id invalid")
	ErrInvalidDuration       = sdkerrors.Register(ModuleName, 4, "duration should be positive")
	ErrInvalidDepositAmount  = sdkerrors.Register(ModuleName, 5, "deposit amount should be positive")
	ErrInvalidPoolId         = sdkerrors.Register(ModuleName, 6, "invalid pool id")
	ErrInvalidGaugeId        = sdkerrors.Register(ModuleName, 7, "invalid gauge id")
)
