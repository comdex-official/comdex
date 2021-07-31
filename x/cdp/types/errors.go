package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cdp module sentinel errors
var (
	ErrorInvalidField    = sdkerrors.Register(ModuleName, 101, "invalid field")
	ErrorInvalidFrom     = sdkerrors.Register(ModuleName, 102, "invalid from")
	ErrorInvalidReceiver = sdkerrors.Register(ModuleName, 103, "invalid receiver")
	ErrorInvalidCoins      = sdkerrors.Register(ModuleName, 104, "invalid coins")
	ErrorInvalidAmount   = sdkerrors.Register(ModuleName, 105, "invalid amount")
)
