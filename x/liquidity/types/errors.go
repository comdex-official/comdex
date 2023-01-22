package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DONTCOVER

var (
	ErrInsufficientDepositAmount = sdkerrors.Register(ModuleName, 2, "insufficient deposit amount")
	ErrPairAlreadyExists         = sdkerrors.Register(ModuleName, 3, "pair already exists")
	ErrPoolAlreadyExists         = sdkerrors.Register(ModuleName, 4, "pool already exists")
	ErrWrongPoolCoinDenom        = sdkerrors.Register(ModuleName, 5, "wrong pool coin denom")
	ErrInvalidCoinDenom          = sdkerrors.Register(ModuleName, 6, "invalid coin denom")
	ErrNoLastPrice               = sdkerrors.Register(ModuleName, 8, "cannot make a market order to a pair with no last price")
	ErrInsufficientOfferCoin     = sdkerrors.Register(ModuleName, 9, "insufficient offer coin")
	ErrPriceOutOfRange           = sdkerrors.Register(ModuleName, 10, "price out of range limit")
	ErrTooLongOrderLifespan      = sdkerrors.Register(ModuleName, 11, "order lifespan is too long")
	ErrDisabledPool              = sdkerrors.Register(ModuleName, 12, "disabled pool")
	ErrWrongPair                 = sdkerrors.Register(ModuleName, 13, "wrong denom pair")
	ErrSameBatch                 = sdkerrors.Register(ModuleName, 14, "cannot cancel an order within the same batch")
	ErrAlreadyCanceled           = sdkerrors.Register(ModuleName, 15, "the order is already canceled")
	ErrDuplicatePairID           = sdkerrors.Register(ModuleName, 16, "duplicate pair id presents in the pair id list")
	ErrTooSmallOrder             = sdkerrors.Register(ModuleName, 17, "too small order")
	ErrTooLargePool              = sdkerrors.Register(ModuleName, 18, "too large pool")
	ErrInvalidPoolID             = sdkerrors.Register(ModuleName, 19, "invalid pool id")
	ErrorFarmerNotFound          = sdkerrors.Register(ModuleName, 20, "farmer not found")
	ErrInvalidUnfarmAmount       = sdkerrors.Register(ModuleName, 21, "invalid unfarm amount")

	ErrDepletedPool                  = sdkerrors.Register(ModuleName, 23, "pool is depleted")
	ErrCalculatedPoolAmountIsZero    = sdkerrors.Register(ModuleName, 24, "calculated provided pool supply with pool tokens is zero or something went wrong while calculation")
	ErrOraclePricesNotFound          = sdkerrors.Register(ModuleName, 25, "oracle prices not found")
	ErrSupplyValueCalculationInvalid = sdkerrors.Register(ModuleName, 26, "something went wrong while calculation supply values")

	ErrInvalidPairID                   = sdkerrors.Register(ModuleName, 28, "invalid pair id")
	ErrAssetNotWhiteListed             = sdkerrors.Register(ModuleName, 29, "asset not whitelisted")
	ErrInvalidAppID                    = sdkerrors.Register(ModuleName, 30, "app id invalid")
	ErrorUnknownProposalType           = sdkerrors.Register(ModuleName, 31, "unknown proposal type")
	ErrorEmptyKeyValueForGenericParams = sdkerrors.Register(ModuleName, 32, "keys or values empty for update generic-params")
	ErrorLengthMismatch                = sdkerrors.Register(ModuleName, 33, "keys and values list length mismatch")

	ErrorNotPositiveAmont = sdkerrors.Register(ModuleName, 34, "amount should be positive")
	ErrTooManyPools       = sdkerrors.Register(ModuleName, 35, "too many pools in the pair")
	ErrPriceNotOnTicks    = sdkerrors.Register(ModuleName, 36, "price is not on ticks")
)
