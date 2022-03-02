package types

import (
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultAdmin = ""
	DefaultLiqRatio = "1.5"
)

var (
	KeyAdmin = []byte("Admin")
	KeyLiqRatio = []byte("LiqRatio")
)

var (
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
