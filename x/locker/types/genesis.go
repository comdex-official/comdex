package types

func NewGenesisState(locker []Locker, lockerProductAssetMapping []LockerProductAssetMapping, lockerTotalRewardsByAssetAppWise []LockerTotalRewardsByAssetAppWise, lockerLookupTable []LockerLookupTableData, userLockerAssetMapping []UserAppAssetLockerMapping, params Params) *GenesisState {
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
		[]LockerLookupTableData{},
		[]UserAppAssetLockerMapping{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
