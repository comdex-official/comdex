package types

import (
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	_ paramstypes.ParamSet = (*Params)(nil)
)

func (m *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return nil
}

func NewParams() Params {
	return Params{}
}
