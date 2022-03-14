package types

import (
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultLiqRatio = "1.5"
	KeyLiqRatio = []byte("LiqRatio")
	_ paramstypes.ParamSet = (*Params)(nil)
)

func NewParams( ratio string) Params {
	return Params{
		LiquidationRatio: ratio,
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultLiqRatio,
	)
}

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (m *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyLiqRatio, &m.LiquidationRatio, validateLiqRatio),
	}
}

func validateLiqRatio(i interface{}) error {
	return nil
}
