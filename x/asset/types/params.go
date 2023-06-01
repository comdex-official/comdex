package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultAssetRegistrationFee = sdk.NewCoin("ucmdx", sdk.NewInt(100_000_000))
)

var (
	KeyAssetRegistrationFee = []byte("AssetRegistrationFee")
)

var _ paramstypes.ParamSet = (*Params)(nil)

func NewParams() Params {
	return Params{
		AssetRegisrationFee: DefaultAssetRegistrationFee,
	}
}

func DefaultParams() Params {
	return NewParams()
}

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (m *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyAssetRegistrationFee, &m.AssetRegisrationFee, validateAssetRegistrationFee),
	}
}

func (m *Params) Validate() error {
	for _, field := range []struct {
		val          interface{}
		validateFunc func(i interface{}) error
	}{
		{m.AssetRegisrationFee, validateAssetRegistrationFee},
	} {
		if err := field.validateFunc(field.val); err != nil {
			return err
		}
	}
	return nil
}

func validateAssetRegistrationFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := v.Validate(); err != nil {
		return fmt.Errorf("invalid pair creation fee: %w", err)
	}

	return nil
}
