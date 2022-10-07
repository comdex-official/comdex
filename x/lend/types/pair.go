package types

import (
	"fmt"
)

func (m *Extended_Pair) Validate() error {
	if m.AssetIn == 0 {
		return fmt.Errorf("asset_in cannot be zero")
	}
	if m.AssetOut == 0 {
		return fmt.Errorf("asset_out cannot be zero")
	}

	return nil
}

func (m *Pool) Validate() error {
	if len(m.CPoolName) >= 16 {
		return ErrInvalidLengthCPoolName
	}
	return nil
}

func (m *AssetToPairMapping) Validate() error {
	return nil
}

func (m *AssetRatesParams) Validate() error {
	return nil
}
