package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	storetypes "cosmossdk.io/store/types"
)

const (
	CreateVaultGas         = storetypes.Gas(36329)
	DepositVaultGas        = storetypes.Gas(23554)
	WithdrawVaultGas       = storetypes.Gas(26473)
	DrawVaultGas           = storetypes.Gas(37559)
	RepayVaultGas          = storetypes.Gas(37559)
	CloseVaultGas          = storetypes.Gas(37559)
	DepositDrawVaultGas    = storetypes.Gas(26329)
	CreateStableVaultGas   = storetypes.Gas(36329)
	DepositStableVaultGas  = storetypes.Gas(23554)
	WithdrawStableVaultGas = storetypes.Gas(26473)
)

func (m *Vault) Validate() error {
	if m.ExtendedPairVaultID == 0 {
		return fmt.Errorf("pair_id cannot be empty")
	}
	if m.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return errorsmod.Wrapf(err, "invalid owner %s", m.Owner)
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
