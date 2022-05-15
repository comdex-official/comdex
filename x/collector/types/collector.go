package types

func NewMsgLookupTableRecords(
	AccumulatorTokenDenom, SecondaryTokenDenom string, SurplusThreshold, DebtThreshold, LockerSavingRate uint64,
) *AccumulatorLookupTable {
	return &AccumulatorLookupTable{
		AccumulatorTokenDenom: AccumulatorTokenDenom,
		SecondaryTokenDenom:   SecondaryTokenDenom,
		SurplusThreshold:      SurplusThreshold,
		DebtThreshold:         DebtThreshold,
		LockerSavingRate:      LockerSavingRate,
	}
}
