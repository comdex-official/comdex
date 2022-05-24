package types

import (
	"fmt"
)

func (m *Pair) Validate() error {
	// if m.Id == 0 {
	// 	return fmt.Errorf("id cannot be zero")
	// }
	if m.AssetIn == 0 {
		return fmt.Errorf("asset_in cannot be zero")
	}
	if m.AssetOut == 0 {
		return fmt.Errorf("asset_out cannot be zero")
	}

	return nil
}
