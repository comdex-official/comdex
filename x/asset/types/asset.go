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
	if m.ID == 0 {
		return fmt.Errorf("id cannot be zero")
	}
	if m.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(m.Name) > MaxAssetNameLength {
		return fmt.Errorf("name length cannot be greater than %d", MaxAssetNameLength)
	}
	if m.Denom == "" {
		return fmt.Errorf("denom cannot be empty")
	}
	if err := sdk.ValidateDenom(m.Denom); err != nil {
		return errors.Wrapf(err, "invalid denom %s", m.Denom)
	}
	if m.Decimals < 0 {
		return fmt.Errorf("decimals cannot be less than zero")
	}

	return nil
}
