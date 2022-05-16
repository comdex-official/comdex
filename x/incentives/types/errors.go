package types

// DONTCOVER

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/incentives module sentinel errors.
var (
	ErrInvalidGaugeStartTime = sdkerrors.Register(ModuleName, 1, "start_time should be greater than current time")
	ErrInvalidGaugeTypeId    = sdkerrors.Register(ModuleName, 2, fmt.Sprintf("gauge_type_id invalid, available gauge_type_ids are %v", ValidGaugeTypeIds))
	ErrInvalidDuration       = sdkerrors.Register(ModuleName, 3, "duration should be positive")
	ErrInvalidDepositAmount  = sdkerrors.Register(ModuleName, 4, "deposit amount should be positive")
	ErrInvalidPoolId         = sdkerrors.Register(ModuleName, 5, "invalid pool id")
	ErrInvalidGaugeId        = sdkerrors.Register(ModuleName, 6, "invalid gauge id")
)
