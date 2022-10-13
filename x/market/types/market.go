package types

import (
	"fmt"
)

const (
	MaxMarketSymbolLength = 8
)

func (m *TimeWeightedAverage) Validate() error {
	if m.AssetID < 0 {
		return fmt.Errorf("id cannot be less than zero")
	}
	if m.ScriptID == 0 {
		return fmt.Errorf("script_id cannot be zero")
	}

	return nil
}
