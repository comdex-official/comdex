package types

import (
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	_ paramstypes.ParamSet = (*Params)(nil)
)

func NewParams() Params {
	return Params{}
}

func DefaultParams() Params {
	return NewParams()
}

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (m *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{}
}

func validateLiqRatio(i interface{}) error {
	return nil
}
