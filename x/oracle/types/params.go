package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultAdmin            = ""
	DefaultIBCPort          = "oracle"
	DefaultIBCVersion       = "comdex-ics-01"
)

var (
	KeyAdmin            = []byte("Admin")
	KeyIBCPort          = []byte("IBCPort")
	KeyIBCVersion       = []byte("IBCVersion")
	KeyOracleAskCount   = []byte("OracleAskCount")
	KeyOracleMinCount   = []byte("OracleMinCount")
	KeyOracleMultiplier = []byte("OracleMultiplier")
)

var (
	_ paramstypes.ParamSet = (*Params)(nil)
)

func NewIBCParams(port, version string) IBCParams {
	return IBCParams{
		Port:    port,
		Version: version,
	}
}

func DefaultIBCParams() IBCParams {
	return NewIBCParams(
		DefaultIBCPort,
		DefaultIBCVersion,
	)
}

func (m *IBCParams) Validate() error {
	if m.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	if m.Version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	return nil
}

func NewParams(admin string, ibc IBCParams) Params {
	return Params{
		Admin:  admin,
		IBC:    ibc,
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultAdmin,
		DefaultIBCParams(),
	)
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
		paramstypes.NewParamSetPair(
			KeyIBCPort,
			m.IBC.Port,
			func(v interface{}) error {
				value, ok := v.(string)
				if !ok {
					return fmt.Errorf("invalid parameter type %T", v)
				}

				if value == "" {
					return fmt.Errorf("ibc.port cannot be empty")
				}

				return nil
			},
		),
		paramstypes.NewParamSetPair(
			KeyIBCVersion,
			m.IBC.Version,
			func(v interface{}) error {
				value, ok := v.(string)
				if !ok {
					return fmt.Errorf("invalid parameter type %T", v)
				}

				if value == "" {
					return fmt.Errorf("ibc.version cannot be empty")
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
	if err := m.IBC.Validate(); err != nil {
		return errors.Wrapf(err, "invalid ibc params")
	}

	return nil
}
