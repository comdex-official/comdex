package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/rewards module sentinel errors .
var (
	ErrInvalidGaugeStartTime   = sdkerrors.Register(ModuleName, 1101, "start time smaller than current time")
	ErrInvalidGaugeTypeID      = sdkerrors.Register(ModuleName, 1102, "gauge type id invalid")
	ErrInvalidDuration         = sdkerrors.Register(ModuleName, 1103, "duration should be positive")
	ErrInvalidDepositAmount    = sdkerrors.Register(ModuleName, 1104, "deposit amount should be positive")
	ErrInvalidPoolID           = sdkerrors.Register(ModuleName, 1105, "invalid pool id")
	ErrInvalidGaugeID          = sdkerrors.Register(ModuleName, 1106, "invalid gauge id")
	ErrNoGaugeForDuration      = sdkerrors.Register(ModuleName, 1107, "no gauges found for given duration")
	ErrDepositSmallThanEpoch   = sdkerrors.Register(ModuleName, 1108, "deposit amount smaller than total epochs/triggers")
	ErrInvalidCalculatedAMount = sdkerrors.Register(ModuleName, 1109, "available distribution coins smaller than calculated distribution amount")
	ErrSamePoolID              = sdkerrors.Register(ModuleName, 1110, "same pool id cannot exists in child pool ids")
	ErrAssetIDDoesNotExist     = sdkerrors.Register(ModuleName, 1111, "Asset Id does not exist in locker for App_Mapping")
	ErrNegativeTimeElapsed     = sdkerrors.Register(ModuleName, 1112, "negative time elapsed since last interest time")
	ErrAppIDExists             = sdkerrors.Register(ModuleName, 1113, "App id exists")
	ErrAppIDDoesNotExists      = sdkerrors.Register(ModuleName, 1114, "App  Id does not exist in rewards for interest accrual")
	ErrPairNotExists           = sdkerrors.Register(ModuleName, 1115, "pair does not exists")
	ErrPriceNotFound           = sdkerrors.Register(ModuleName, 1116, "price not found")
	ErrInvalidAppID            = sdkerrors.Register(ModuleName, 1117, "invalid app id")
	ErrInternalRewardsNotFound = sdkerrors.Register(ModuleName, 1118, "Internal rewards not found")
	ErrStablemintVaultFound    = sdkerrors.Register(ModuleName, 1119, "Can't give reward to stablemint vault")
)
