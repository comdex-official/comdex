package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultAdmin = "comdex1pkkayn066msg6kn33wnl5srhdt3tnu2v9jjqu0"
)

var (
	KeyAdmin = []byte("Admin")
)

var (
	_ paramstypes.ParamSet = (*Params)(nil)
)

func NewParams(admin string) Params {
	return Params{
		Admin: admin,
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultAdmin,
	)
}

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (m *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(
			KeyAdmin,
			m.Admin,
			func(v interface{}) error {
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
			},
		),
	}
}

func (m *Params) Validate() error {
	if m.Admin == "" {
		return fmt.Errorf("admin cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Admin); err != nil {
		return errors.Wrapf(err, "invalid admin %s", m.Admin)
	}

	return nil
}
