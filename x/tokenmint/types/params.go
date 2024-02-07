package types

import (
	storetypes "cosmossdk.io/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	TokenmintGas = storetypes.Gas(36329)
)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
