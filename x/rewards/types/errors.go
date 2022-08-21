package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/rewards module sentinel errors .
var (
	ErrInvalidGaugeStartTime   = sdkerrors.Register(ModuleName, 2, "start time smaller than current time")
	ErrInvalidGaugeTypeID      = sdkerrors.Register(ModuleName, 3, "gauge type id invalid")
	ErrInvalidDuration         = sdkerrors.Register(ModuleName, 4, "duration should be positive")
	ErrInvalidDepositAmount    = sdkerrors.Register(ModuleName, 5, "deposit amount should be positive")
	ErrInvalidPoolID           = sdkerrors.Register(ModuleName, 6, "invalid pool id")
	ErrInvalidGaugeID          = sdkerrors.Register(ModuleName, 7, "invalid gauge id")
	ErrNoGaugeForDuration      = sdkerrors.Register(ModuleName, 8, "no gauges found for given duration")
	ErrDepositSmallThanEpoch   = sdkerrors.Register(ModuleName, 9, "deposit amount smaller than total epochs/triggers")
	ErrInvalidCalculatedAMount = sdkerrors.Register(ModuleName, 10, "available distribution coins smaller than calculated distribution amount")

	ErrAssetIDDoesNotExist = sdkerrors.Register(ModuleName, 1101, "Asset Id does not exist in locker for App_Mapping")
	ErrNegativeTimeElapsed = sdkerrors.Register(ModuleName, 1102, "negative time elapsed since last interest time")
	ErrAppIDExists         = sdkerrors.Register(ModuleName, 1103, "Asset Id does not exist in locker for App_Mapping")
	ErrAppIDDoesNotExists  = sdkerrors.Register(ModuleName, 1104, "App  Id does not exist in rewards for interest accrual")
	ErrPairNotExists       = sdkerrors.Register(ModuleName, 1105, "pair does not exists")
	ErrPriceNotFound       = sdkerrors.Register(ModuleName, 1106, "price not found")
	ErrInvalidAppID        = sdkerrors.Register(ModuleName, 1107, "invalid app id")
	ErrVaultNotFound       = sdkerrors.Register(ModuleName, 1108, "vault not found")

	BurnCoinValueInRewardsIsZero                = sdkerrors.Register(ModuleName, 1109, "Burn Coin value in rewards is zero")
	MintCoinValueInRewardsIsZero                = sdkerrors.Register(ModuleName, 1110, "Mint Coin value in rewards is zero")
	SendCoinsFromModuleToModuleInRewardsIsZero  = sdkerrors.Register(ModuleName, 1111, "Coin value in module to module transfer in rewards is zero")
	SendCoinsFromModuleToAccountInRewardsIsZero = sdkerrors.Register(ModuleName, 1112, "Coin value in module to account transfer in rewards is zero")
	SendCoinsFromAccountToModuleInRewardsIsZero = sdkerrors.Register(ModuleName, 1113, "Coin value in account to module transfer in rewards is zero")
	ErrInternalRewardsNotFound                  = sdkerrors.Register(ModuleName, 1114, "Internal rewards not found")
)
