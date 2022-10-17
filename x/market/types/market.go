package types

import (
	"fmt"
)

const (
	MaxMarketSymbolLength = 8
)

func (m *TimeWeightedAverage) Validate() error {
	if m.ScriptID == 0 {
		return fmt.Errorf("script_id cannot be zero")
	}

	return nil
}
