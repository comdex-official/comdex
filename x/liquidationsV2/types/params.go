package types

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultLiquidationBatchSize Liquidation params default values
var (
	DefaultLiquidationBatchSize = uint64(200)
)

var KeyLiquidationBatchSize = []byte("LiquidationBatchSize")

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(liquidationBatchSize uint64) Params {
	return Params{
		LiquidationBatchSize: liquidationBatchSize,
	}
}

func DefaultParams() Params {
	return NewParams(DefaultLiquidationBatchSize)
}

func (p Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLiquidationBatchSize, &p.LiquidationBatchSize, validateLiquidationBatchSize),
	}
}

// Validate validates Params.
func (p Params) Validate() error {
	for _, field := range []struct {
		val          interface{}
		validateFunc func(i interface{}) error
	}{
		{p.LiquidationBatchSize, validateLiquidationBatchSize},
	} {
		if err := field.validateFunc(field.val); err != nil {
			return err
		}
	}
	return nil
}

func validateLiquidationBatchSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("batch size must be positive: %d", v)
	}

	return nil
}
