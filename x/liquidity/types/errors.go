package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DONTCOVER

// x/liquidity module sentinel errors.
var (
	ErrInsufficientDepositAmount       = sdkerrors.Register(ModuleName, 2, "insufficient deposit amount")
	ErrPairAlreadyExists               = sdkerrors.Register(ModuleName, 3, "pair already exists")
	ErrPoolAlreadyExists               = sdkerrors.Register(ModuleName, 4, "pool already exists")
	ErrWrongPoolCoinDenom              = sdkerrors.Register(ModuleName, 5, "wrong pool coin denom")
	ErrInvalidCoinDenom                = sdkerrors.Register(ModuleName, 6, "invalid coin denom")
	ErrNoLastPrice                     = sdkerrors.Register(ModuleName, 8, "cannot make a market order to a pair with no last price")
	ErrInsufficientOfferCoin           = sdkerrors.Register(ModuleName, 9, "insufficient offer coin")
	ErrPriceOutOfRange                 = sdkerrors.Register(ModuleName, 10, "price out of range limit")
	ErrTooLongOrderLifespan            = sdkerrors.Register(ModuleName, 11, "order lifespan is too long")
	ErrDisabledPool                    = sdkerrors.Register(ModuleName, 12, "disabled pool")
	ErrWrongPair                       = sdkerrors.Register(ModuleName, 13, "wrong denom pair")
	ErrSameBatch                       = sdkerrors.Register(ModuleName, 14, "cannot cancel an order within the same batch")
	ErrAlreadyCanceled                 = sdkerrors.Register(ModuleName, 15, "the order is already canceled")
	ErrDuplicatePairID                 = sdkerrors.Register(ModuleName, 16, "duplicate pair id presents in the pair id list")
	ErrTooSmallOrder                   = sdkerrors.Register(ModuleName, 17, "too small order")
	ErrTooLargePool                    = sdkerrors.Register(ModuleName, 18, "too large pool")
	ErrInvalidPoolID                   = sdkerrors.Register(ModuleName, 19, "invalid pool id")
	ErrNoSoftLockPresent               = sdkerrors.Register(ModuleName, 20, "no soft lock created for given pool")
	ErrInvalidUnlockAmount             = sdkerrors.Register(ModuleName, 21, "invalid soft unlock amount")
	ErrLPDataNotExistsForPool          = sdkerrors.Register(ModuleName, 22, "liquidity providers data does not exists for givn pool id")
	ErrDepletedPool                    = sdkerrors.Register(ModuleName, 23, "pool is depleted")
	ErrCalculatedPoolAmountIsZero      = sdkerrors.Register(ModuleName, 24, "calculated provided pool supply with pool tokens is zero or something went wrong while calculation")
	ErrOraclePricesNotFound            = sdkerrors.Register(ModuleName, 25, "oracle prices not found")
	ErrSupplyValueCalculationInvalid   = sdkerrors.Register(ModuleName, 26, "something went wrong while calculation supply values")
	ErrInsufficientAvailableBalance    = sdkerrors.Register(ModuleName, 27, "insufficient available balance")
	ErrInvalidPairId                   = sdkerrors.Register(ModuleName, 28, "invalid pair id")
	ErrAssetNotWhiteListed             = sdkerrors.Register(ModuleName, 29, "asset not whitelisted")
	ErrInvalidAppID                    = sdkerrors.Register(ModuleName, 30, "app id invalid")
	ErrorUnknownProposalType           = sdkerrors.Register(ModuleName, 31, "unknown proposal type")
	ErrorEmptyKeyValueForGenericParams = sdkerrors.Register(ModuleName, 32, "keys or values empty for update generic-params")
	ErrorLengthMismatch                = sdkerrors.Register(ModuleName, 33, "keys and values list length mismatch")
)
