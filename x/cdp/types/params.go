package types

import (
	paramTypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyCollateralParams = []byte("CollateralParams")
	KeyDebtParam        = []byte("DebtParam")
)

func (p *Params) ParamSetPairs() paramTypes.ParamSetPairs {
	return paramTypes.ParamSetPairs{
		paramTypes.NewParamSetPair(KeyCollateralParams, &p.CollateralParams, validateCollateralParams),
	}
}

func validateCollateralParams(i interface{}) error {
	//TODO
	return nil
}
