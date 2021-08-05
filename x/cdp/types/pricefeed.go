package types

import "fmt"

type PricefeedType string

const (
	Spot        PricefeedType = "spot"
	liquidation PricefeedType = "liquidation"
)

func (pft PricefeedType) IsValid() error {
	switch pft {
	case Spot, liquidation:
		return nil
	}
	return fmt.Errorf("invalid pricefeed type: %s", pft)
}
