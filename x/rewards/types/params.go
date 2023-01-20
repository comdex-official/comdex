package types

import (
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramstypes.ParamSet = (*Params)(nil)

const (
	Int64SecondsInADay                   = int64(84600)
	UInt64One                            = uint64(1)
	Float64One                           = float64(1)
	Int64Zero                            = int64(0)
	UInt64Zero                           = uint64(0)
	DefaultAllowedBlocksForPriceInactive = 600
)

// ParamKeyTable for incentives module.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default params for the liquidity module.
func DefaultParams() Params {
	return Params{}
}

// ParamSetPairs implements ParamSet.
func (params *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{}
}

// Validate validates Params.
func (params Params) Validate() error {
	for _, field := range []struct {
		val          interface{}
		validateFunc func(i interface{}) error
	}{} {
		if err := field.validateFunc(field.val); err != nil {
			return err
		}
	}
	return nil
}
