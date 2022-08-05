package types

func NewGenesisState(locker []Locker, lockerProductAssetMapping []LockerProductAssetMapping, lockerTotalRewardsByAssetAppWise []LockerTotalRewardsByAssetAppWise, lockerLookupTable []LockerLookupTable, userLockerAssetMapping []UserLockerAssetMapping, params Params) *GenesisState {
	return &GenesisState{
		Lockers:                          locker,
		LockerProductAssetMapping:        lockerProductAssetMapping,
		LockerTotalRewardsByAssetAppWise: lockerTotalRewardsByAssetAppWise,
		LockerLookupTable:                lockerLookupTable,
		UserLockerAssetMapping:           userLockerAssetMapping,
		Params:                           params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]Locker{},
		[]LockerProductAssetMapping{},
		[]LockerTotalRewardsByAssetAppWise{},
		[]LockerLookupTable{},
		[]UserLockerAssetMapping{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
