package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

var (
	ErrInsufficientDepositAmount       = errorsmod.Register(ModuleName, 801, "insufficient deposit amount")
	ErrPairAlreadyExists               = errorsmod.Register(ModuleName, 802, "pair already exists")
	ErrPoolAlreadyExists               = errorsmod.Register(ModuleName, 803, "pool already exists")
	ErrWrongPoolCoinDenom              = errorsmod.Register(ModuleName, 804, "wrong pool coin denom")
	ErrInvalidCoinDenom                = errorsmod.Register(ModuleName, 805, "invalid coin denom")
	ErrNoLastPrice                     = errorsmod.Register(ModuleName, 806, "cannot make a market order to a pair with no last price")
	ErrInsufficientOfferCoin           = errorsmod.Register(ModuleName, 807, "insufficient offer coin")
	ErrPriceOutOfRange                 = errorsmod.Register(ModuleName, 808, "price out of range limit")
	ErrTooLongOrderLifespan            = errorsmod.Register(ModuleName, 809, "order lifespan is too long")
	ErrDisabledPool                    = errorsmod.Register(ModuleName, 810, "disabled pool")
	ErrWrongPair                       = errorsmod.Register(ModuleName, 811, "wrong denom pair")
	ErrSameBatch                       = errorsmod.Register(ModuleName, 812, "cannot cancel an order within the same batch")
	ErrAlreadyCanceled                 = errorsmod.Register(ModuleName, 813, "the order is already canceled")
	ErrDuplicatePairID                 = errorsmod.Register(ModuleName, 814, "duplicate pair id presents in the pair id list")
	ErrTooSmallOrder                   = errorsmod.Register(ModuleName, 815, "too small order")
	ErrTooLargePool                    = errorsmod.Register(ModuleName, 816, "too large pool")
	ErrInvalidPoolID                   = errorsmod.Register(ModuleName, 817, "invalid pool id")
	ErrorFarmerNotFound                = errorsmod.Register(ModuleName, 818, "farmer not found")
	ErrInvalidUnfarmAmount             = errorsmod.Register(ModuleName, 819, "invalid unfarm amount")
	ErrDepletedPool                    = errorsmod.Register(ModuleName, 820, "pool is depleted")
	ErrCalculatedPoolAmountIsZero      = errorsmod.Register(ModuleName, 821, "calculated provided pool supply with pool tokens is zero or something went wrong while calculation")
	ErrOraclePricesNotFound            = errorsmod.Register(ModuleName, 822, "oracle prices not found")
	ErrSupplyValueCalculationInvalid   = errorsmod.Register(ModuleName, 823, "something went wrong while calculation supply values")
	ErrInvalidPairID                   = errorsmod.Register(ModuleName, 824, "invalid pair id")
	ErrAssetNotWhiteListed             = errorsmod.Register(ModuleName, 825, "asset not whitelisted")
	ErrInvalidAppID                    = errorsmod.Register(ModuleName, 826, "app id invalid")
	ErrorUnknownProposalType           = errorsmod.Register(ModuleName, 827, "unknown proposal type")
	ErrorEmptyKeyValueForGenericParams = errorsmod.Register(ModuleName, 828, "keys or values empty for update generic-params")
	ErrorLengthMismatch                = errorsmod.Register(ModuleName, 829, "keys and values list length mismatch")
	ErrorNotPositiveAmont              = errorsmod.Register(ModuleName, 830, "amount should be positive")
	ErrTooManyPools                    = errorsmod.Register(ModuleName, 831, "too many pools in the pair")
	ErrPriceNotOnTicks                 = errorsmod.Register(ModuleName, 832, "price is not on ticks")
)
