package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/rewards module sentinel errors
var (
	ErrSample              = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrAssetIdDoesNotExist = sdkerrors.Register(ModuleName, 1101, "Asset Id does not exist in locker for App_Mapping")
	ErrNegativeTimeElapsed = sdkerrors.Register(ModuleName, 1102, "negative time elapsed since last interest time")
	ErrAppIdExists         = sdkerrors.Register(ModuleName, 1103, "Asset Id does not exist in locker for App_Mapping")
	ErrAppIdDoesNotExists  = sdkerrors.Register(ModuleName, 1104, "Asset Id does not exist in locker for App_Mapping")
)
