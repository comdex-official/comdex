package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultIBCPort          = "asset"
	DefaultIBCVersion       = "comdex-ics-01"
	DefaultOracleAskCount   = 1
	DefaultOracleMinCount   = 1
	DefaultOracleMultiplier = 6
)

var (
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
	return nil
}

func NewParams(ibc IBCParams, oracle OracleParams) Params {
	return Params{
		IBC:    ibc,
		Oracle: oracle,
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultIBCParams(),
		DefaultOracleParams(),
	)
}

func (m *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(
			KeyIBCPort,
			m.IBC.Port,
			func(value interface{}) error {
				return nil
			},
		),
		paramstypes.NewParamSetPair(
			KeyIBCVersion,
			m.IBC.Version,
			func(value interface{}) error {
				return nil
			},
		),
		paramstypes.NewParamSetPair(
			KeyOracleAskCount,
			m.Oracle.AskCount,
			func(value interface{}) error {
				return nil
			},
		),
		paramstypes.NewParamSetPair(
			KeyOracleMinCount,
			m.Oracle.MinCount,
			func(value interface{}) error {
				return nil
			},
		),
		paramstypes.NewParamSetPair(
			KeyOracleMultiplier,
			m.Oracle.Multiplier,
			func(value interface{}) error {
				return nil
			},
		),
	}
}

func (m *Params) Validate() error {
	if err := m.IBC.Validate(); err != nil {
		return errors.Wrapf(err, "invalid ibc params")
	}
	if err := m.Oracle.Validate(); err != nil {
		return errors.Wrapf(err, "invalid oracle params")
	}

	return nil
}
