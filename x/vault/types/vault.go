package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func (m *Vault) Validate() error {
	if m.ID == 0 {
		return fmt.Errorf("id cannot be empty")
	}
	if m.PairID == 0 {
		return fmt.Errorf("pair_id cannot be empty")
	}
	if m.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return errors.Wrapf(err, "invalid owner %s", m.Owner)
	}
	if m.AmountIn.IsNil() {
		return fmt.Errorf("amount_in cannot be nil")
	}
	if m.AmountIn.IsNegative() {
		return fmt.Errorf("amount_in cannot be negative")
	}
	if m.AmountOut.IsNil() {
		return fmt.Errorf("amount_out cannot be nil")
	}
	if m.AmountOut.IsNegative() {
		return fmt.Errorf("amount_out cannot be negative")
	}

	return nil
}
