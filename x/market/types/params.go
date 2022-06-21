package types

import (
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
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

func (m *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return nil
}

func NewParams() Params {
	return Params{}
}
