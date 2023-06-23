package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MaxAssetNameLength = 16
)

func (m *Asset) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(m.Name) > MaxAssetNameLength {
		return fmt.Errorf("name length cannot be greater than %d", MaxAssetNameLength)
	}
	if err := sdk.ValidateDenom(m.Denom); err != nil {
		return errors.Wrapf(err, "invalid denom %s", m.Denom)
	}
	if m.Decimals.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("decimals cannot be less than or equal to zero")
	}

	return nil
}

func (m *AssetPair) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(m.Name) > MaxAssetNameLength {
		return fmt.Errorf("name length cannot be greater than %d", MaxAssetNameLength)
	}
	if err := sdk.ValidateDenom(m.Denom); err != nil {
		return errors.Wrapf(err, "invalid denom %s", m.Denom)
	}
	if m.Decimals.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("decimals cannot be less than or equal to zero")
	}

	return nil
}
