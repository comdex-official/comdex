package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func (m *Collateral) Validate() error {
	if m.Id == 0 {
		return errors.New("id cannot be zero")
	}
	if m.DenomIn == "" {
		return errors.New("denom_in cannot be empty")
	}
	if err := sdk.ValidateDenom(m.DenomIn); err != nil {
		return errors.Wrapf(err, "invalid denom_in %s", m.DenomIn)
	}
	if m.DenomOut == "" {
		return errors.New("denom_out cannot be empty")
	}
	if err := sdk.ValidateDenom(m.DenomOut); err != nil {
		return errors.Wrapf(err, "invalid denom_out %s", m.DenomOut)
	}
	if m.LiquidationRatio.IsNil() {
		return errors.New("liquidation_ratio cannot be nil")
	}
	if m.LiquidationRatio.IsNegative() {
		return errors.New("liquidation_ratio cannot be negative")
	}

	return nil
}

type (
	Collaterals []Collateral
)
