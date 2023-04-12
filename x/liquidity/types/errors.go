package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DONTCOVER

var (
	ErrInsufficientDepositAmount       = sdkerrors.Register(ModuleName, 801, "insufficient deposit amount")
	ErrPairAlreadyExists               = sdkerrors.Register(ModuleName, 802, "pair already exists")
	ErrPoolAlreadyExists               = sdkerrors.Register(ModuleName, 803, "pool already exists")
	ErrWrongPoolCoinDenom              = sdkerrors.Register(ModuleName, 804, "wrong pool coin denom")
	ErrInvalidCoinDenom                = sdkerrors.Register(ModuleName, 805, "invalid coin denom")
	ErrNoLastPrice                     = sdkerrors.Register(ModuleName, 806, "cannot make a market order to a pair with no last price")
	ErrInsufficientOfferCoin           = sdkerrors.Register(ModuleName, 807, "insufficient offer coin")
	ErrPriceOutOfRange                 = sdkerrors.Register(ModuleName, 808, "price out of range limit")
	ErrTooLongOrderLifespan            = sdkerrors.Register(ModuleName, 809, "order lifespan is too long")
	ErrDisabledPool                    = sdkerrors.Register(ModuleName, 810, "disabled pool")
	ErrWrongPair                       = sdkerrors.Register(ModuleName, 811, "wrong denom pair")
	ErrSameBatch                       = sdkerrors.Register(ModuleName, 812, "cannot cancel an order within the same batch")
	ErrAlreadyCanceled                 = sdkerrors.Register(ModuleName, 813, "the order is already canceled")
	ErrDuplicatePairID                 = sdkerrors.Register(ModuleName, 814, "duplicate pair id presents in the pair id list")
	ErrTooSmallOrder                   = sdkerrors.Register(ModuleName, 815, "too small order")
	ErrTooLargePool                    = sdkerrors.Register(ModuleName, 816, "too large pool")
	ErrInvalidPoolID                   = sdkerrors.Register(ModuleName, 817, "invalid pool id")
	ErrorFarmerNotFound                = sdkerrors.Register(ModuleName, 818, "farmer not found")
	ErrInvalidUnfarmAmount             = sdkerrors.Register(ModuleName, 819, "invalid unfarm amount")
	ErrDepletedPool                    = sdkerrors.Register(ModuleName, 820, "pool is depleted")
	ErrCalculatedPoolAmountIsZero      = sdkerrors.Register(ModuleName, 821, "calculated provided pool supply with pool tokens is zero or something went wrong while calculation")
	ErrOraclePricesNotFound            = sdkerrors.Register(ModuleName, 822, "oracle prices not found")
	ErrSupplyValueCalculationInvalid   = sdkerrors.Register(ModuleName, 823, "something went wrong while calculation supply values")
	ErrInvalidPairID                   = sdkerrors.Register(ModuleName, 824, "invalid pair id")
	ErrAssetNotWhiteListed             = sdkerrors.Register(ModuleName, 825, "asset not whitelisted")
	ErrInvalidAppID                    = sdkerrors.Register(ModuleName, 826, "app id invalid")
	ErrorUnknownProposalType           = sdkerrors.Register(ModuleName, 827, "unknown proposal type")
	ErrorEmptyKeyValueForGenericParams = sdkerrors.Register(ModuleName, 828, "keys or values empty for update generic-params")
	ErrorLengthMismatch                = sdkerrors.Register(ModuleName, 829, "keys and values list length mismatch")
	ErrorNotPositiveAmont              = sdkerrors.Register(ModuleName, 830, "amount should be positive")
	ErrTooManyPools                    = sdkerrors.Register(ModuleName, 831, "too many pools in the pair")
	ErrPriceNotOnTicks                 = sdkerrors.Register(ModuleName, 832, "price is not on ticks")
	ErrNotActiveFarmer                 = sdkerrors.Register(ModuleName, 833, "inactive farmer")
	ErrInvalidFarmAmount               = sdkerrors.Register(ModuleName, 834, "invalid farm amount")
)
