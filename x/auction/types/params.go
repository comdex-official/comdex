package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyLiquidationPaneltyPercent = []byte("LiquidationPaneltyPercent")
	KeyAuctionDiscountPercent    = []byte("AuctionDiscountPercent")
	KeyAuctionDurationHours      = []byte("AuctionDurationHours")
)

var (
	DefaultLiquidationPaneltyPercent = uint64(15)
	DefaultAuctionDiscountPercent    = uint64(5)
	DefaultAuctionDurationHours      = uint64(6)
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		LiquidationPaneltyPercent: DefaultLiquidationPaneltyPercent,
		AuctionDiscountPercent:    DefaultAuctionDiscountPercent,
		AuctionDurationHours:      DefaultAuctionDurationHours,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLiquidationPaneltyPercent, &p.LiquidationPaneltyPercent, validateLiquidationPanelty),
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
		{p.LiquidationPaneltyPercent, validateLiquidationPanelty},
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

func validateLiquidationPanelty(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("liquidation panelty cannot be less than 1 percent")
	}
	return nil
}

func validateAuctionDiscount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
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
		return fmt.Errorf("auction duraction cannot be less than 1 hour")
	}
	return nil
}
