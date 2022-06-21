package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	LockedVaultDoesNotExist             = sdkerrors.Register(ModuleName, 201, "locked vault does not exist with given id")
	ErrorExtendedPairVaultDoesNotExists = sdkerrors.Register(ModuleName, 202, "Extended pair vault does not exists for the given id")
	ErrorPriceDoesNotExist              = sdkerrors.Register(ModuleName, 203, "Price does not exist")
	ErrAppIdExists                      = sdkerrors.Register(ModuleName, 1101, "Asset Id does not exist in locker for App_Mapping")
	ErrAppIdDoesNotExists               = sdkerrors.Register(ModuleName, 1102, "Asset Id does not exist in locker for App_Mapping")
)
