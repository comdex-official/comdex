package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func (m *Pair) Validate() error {
	if m.Id == 0 {
		return fmt.Errorf("id cannot be zero")
	}
	if m.DenomIn == "" {
		return fmt.Errorf("denom_in cannot be empty")
	}
	if err := sdk.ValidateDenom(m.DenomIn); err != nil {
		return errors.Wrapf(err, "invalid denom_in %s", m.DenomIn)
	}
	if m.DenomOut == "" {
		return fmt.Errorf("denom_out cannot be empty")
	}
	if err := sdk.ValidateDenom(m.DenomOut); err != nil {
		return errors.Wrapf(err, "invalid denom_out %s", m.DenomOut)
	}
	if m.LiquidationRatio.IsNil() {
		return fmt.Errorf("liquidation_ratio cannot be nil")
	}
	if m.LiquidationRatio.IsNegative() {
		return fmt.Errorf("liquidation_ratio cannot be negative")
	}

	return nil
}
