package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	CreateVaultGas         = sdk.Gas(96329)
	DepositVaultGas        = sdk.Gas(53554)
	WithdrawVaultGas       = sdk.Gas(56473)
	DrawVaultGas           = sdk.Gas(77559)
	RepayVaultGas          = sdk.Gas(77559)
	CloseVaultGas          = sdk.Gas(77559)
	CreateStableVaultGas   = sdk.Gas(96329)
	DepositStableVaultGas  = sdk.Gas(53554)
	WithdrawStableVaultGas = sdk.Gas(56473)
)

func (m *Vault) Validate() error {
	if m.ExtendedPairVaultID == 0 {
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
