package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultAdmin            = ""
	DefaultIBCPort          = "asset"
	DefaultIBCVersion       = "comdex-ics-01"
	DefaultOracleAskCount   = 1
	DefaultOracleMinCount   = 1
	DefaultOracleMultiplier = 9
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

func NewOracleParams(askCount, minCount, multiplier uint64) OracleParams {
	return OracleParams{
		AskCount:   askCount,
		MinCount:   minCount,
		Multiplier: multiplier,
	}
}

func DefaultOracleParams() OracleParams {
	return NewOracleParams(
		DefaultOracleAskCount,
		DefaultOracleMinCount,
		DefaultOracleMultiplier,
	)
}

func (m *OracleParams) Validate() error {
	if m.AskCount == 0 {
		return fmt.Errorf("ask_count cannot be zero")
	}
	if m.MinCount == 0 {
		return fmt.Errorf("min_count cannot be zero")
	}

	return nil
}

func NewParams(admin string, ibc IBCParams, oracle OracleParams) Params {
	return Params{
		Admin:  admin,
		IBC:    ibc,
		Oracle: oracle,
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultAdmin,
		DefaultIBCParams(),
		DefaultOracleParams(),
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
		paramstypes.NewParamSetPair(
			KeyOracleAskCount,
			m.Oracle.AskCount,
			func(v interface{}) error {
				value, ok := v.(uint64)
				if !ok {
					return fmt.Errorf("invalid parameter type %T", v)
				}

				if value == 0 {
					return fmt.Errorf("oracle.ask_count cannot be zero")
				}

				return nil
			},
		),
		paramstypes.NewParamSetPair(
			KeyOracleMinCount,
			m.Oracle.MinCount,
			func(v interface{}) error {
				value, ok := v.(uint64)
				if !ok {
					return fmt.Errorf("invalid parameter type %T", v)
				}

				if value == 0 {
					return fmt.Errorf("oracle.min_count cannot be zero")
				}

				return nil
			},
		),
		paramstypes.NewParamSetPair(
			KeyOracleMultiplier,
			m.Oracle.Multiplier,
			func(v interface{}) error {
				_, ok := v.(uint64)
				if !ok {
					return fmt.Errorf("invalid parameter type %T", v)
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
	if err := m.Oracle.Validate(); err != nil {
		return errors.Wrapf(err, "invalid oracle params")
	}

	return nil
}
