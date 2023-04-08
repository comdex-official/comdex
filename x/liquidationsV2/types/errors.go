package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/liquidationsV2 module sentinel errors
var (
	ErrVaultIDInvalid = sdkerrors.Register(ModuleName, 1501, "Vault Id invalid")
)
