package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/rewards module sentinel errors .
var (
	ErrInvalidGaugeStartTime   = errorsmod.Register(ModuleName, 1101, "start time smaller than current time")
	ErrInvalidGaugeTypeID      = errorsmod.Register(ModuleName, 1102, "gauge type id invalid")
	ErrInvalidDuration         = errorsmod.Register(ModuleName, 1103, "duration should be positive")
	ErrInvalidDepositAmount    = errorsmod.Register(ModuleName, 1104, "deposit amount should be positive")
	ErrInvalidPoolID           = errorsmod.Register(ModuleName, 1105, "invalid pool id")
	ErrInvalidGaugeID          = errorsmod.Register(ModuleName, 1106, "invalid gauge id")
	ErrNoGaugeForDuration      = errorsmod.Register(ModuleName, 1107, "no gauges found for given duration")
	ErrDepositSmallThanEpoch   = errorsmod.Register(ModuleName, 1108, "deposit amount smaller than total epochs/triggers")
	ErrInvalidCalculatedAMount = errorsmod.Register(ModuleName, 1109, "available distribution coins smaller than calculated distribution amount")
	ErrSamePoolID              = errorsmod.Register(ModuleName, 1110, "same pool id cannot exists in child pool ids")
	ErrAssetIDDoesNotExist     = errorsmod.Register(ModuleName, 1111, "Asset Id does not exist in locker for App_Mapping")
	ErrNegativeTimeElapsed     = errorsmod.Register(ModuleName, 1112, "negative time elapsed since last interest time")
	ErrAppIDExists             = errorsmod.Register(ModuleName, 1113, "App id exists")
	ErrAppIDDoesNotExists      = errorsmod.Register(ModuleName, 1114, "App  Id does not exist in rewards for interest accrual")
	ErrPairNotExists           = errorsmod.Register(ModuleName, 1115, "pair does not exists")
	ErrPriceNotFound           = errorsmod.Register(ModuleName, 1116, "price not found")
	ErrInvalidAppID            = errorsmod.Register(ModuleName, 1117, "invalid app id")
	ErrInternalRewardsNotFound = errorsmod.Register(ModuleName, 1118, "Internal rewards not found")
	ErrStablemintVaultFound    = errorsmod.Register(ModuleName, 1119, "Can't give reward to stablemint vault")
	ErrDisabledPool            = errorsmod.Register(ModuleName, 1120, "diabled pool")
)
