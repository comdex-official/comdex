package types

import "time"

const (
	LiquidityGaugeTypeID = 1
	MinimumEpochDuration = time.Hour * 12
)

// ValidGaugeTypeIds stores all the gauge types ids
// It is used while message validations.
var ValidGaugeTypeIds = []uint64{
	LiquidityGaugeTypeID,
}
