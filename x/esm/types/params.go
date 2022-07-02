package types

import (
	
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	DefaultAdmin = ""
)

var (
	KeyAdmin = []byte("Admin")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams( admin string) Params {
	return Params{
		Admin: admin,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAdmin,
	)
}

// ParamSetPairs get the params.ParamSet
func (m *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(
			KeyAdmin,
			m.Admin,
			validateAdmin,
		),
	}
}

func validateAdmin(v interface{}) error {
	value, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type %T", v)
	}

	if value == "" {
		return fmt.Errorf("admin cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(value); err != nil {
		return errors.Wrapf(err, "invalid admin %s", value)
	}

	return nil
}

// Validate validates the set of params
func (m *Params) Validate() error {
	if m.Admin == "" {
		return fmt.Errorf("admin cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Admin); err != nil {
		return errors.Wrapf(err, "invalid admin %s", m.Admin)
	}

	return nil
}
