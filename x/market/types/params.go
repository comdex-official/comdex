package types

import (
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	_ paramstypes.ParamSet = (*Params)(nil)
)

const (
	Int64Twenty    = int64(20)
	Int64TwentyOne = int64(21)
	Int64Ten       = int64(10)
	Int64One       = int64(1)
	Int64Zero      = int64(0)
)

func (m *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return nil
}

func NewParams() Params {
	return Params{}
}
