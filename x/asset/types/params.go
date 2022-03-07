package types

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"strings"
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
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("Liquidation Ratio can not be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}
