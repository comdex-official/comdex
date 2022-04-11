package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyLiquidationPenaltyPercent = []byte("LiquidationPenaltyPercent")
	KeyAuctionDiscountPercent    = []byte("AuctionDiscountPercent")
	KeyAuctionDurationHours      = []byte("AuctionDurationHours")
)

var (
	DefaultLiquidationPenaltyPercent = "0.15"
	DefaultAuctionDiscountPercent    = "0.05"
	DefaultAuctionDurationHours      = uint64(6)
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(liquidationPenaltyPercent string, auctionDiscountPercent string, auctionDurationHours uint64) Params {
	return Params{
		LiquidationPenaltyPercent: liquidationPenaltyPercent,
		AuctionDiscountPercent:    auctionDiscountPercent,
		AuctionDurationHours:      auctionDurationHours,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultLiquidationPenaltyPercent,
		DefaultAuctionDiscountPercent,
		DefaultAuctionDurationHours,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLiquidationPenaltyPercent, &p.LiquidationPenaltyPercent, validateLiquidationPenalty),
		paramtypes.NewParamSetPair(KeyAuctionDiscountPercent, &p.AuctionDiscountPercent, validateAuctionDiscount),
		paramtypes.NewParamSetPair(KeyAuctionDurationHours, &p.AuctionDurationHours, validateAuctionDuration),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.LiquidationPenaltyPercent, validateLiquidationPenalty},
		{p.AuctionDiscountPercent, validateAuctionDiscount},
		{p.AuctionDurationHours, validateAuctionDuration},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}

// String implements the Stringer interface.
// func (p Params) String() string {
// 	out, _ := yaml.Marshal(p)
// 	return string(out)
// }

func validateLiquidationPenalty(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	q, _ := sdk.NewDecFromStr(v)
	u, _ := sdk.NewDecFromStr("0.01")
	if q.LT(u) {
		return fmt.Errorf("liquidation penalty cannot be less than 1 percent")
	}
	return nil
}

func validateAuctionDiscount(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	q, _ := sdk.NewDecFromStr(v)
	u, _ := sdk.NewDecFromStr("0.01")
	if q.LT(u) {
		return fmt.Errorf("auction discount cannot be less than 1 percent")
	}
	return nil
}

func validateAuctionDuration(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("auction duration cannot be less than 1 hour")
	}
	return nil
}
