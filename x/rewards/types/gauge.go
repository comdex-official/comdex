package types

// should be unique.
const (
	LiquidityGaugeTypeID = 1
)

// ValidGaugeTypeIds stores all the gauge types ids
// It is used while message validations.
var ValidGaugeTypeIds = []uint64{
	LiquidityGaugeTypeID,
}
