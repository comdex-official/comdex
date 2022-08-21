package types

func NewGenesisState(lockedVault []LockedVault, params Params) *GenesisState {
	return &GenesisState{
		LockedVault: lockedVault,
		Params:      params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]LockedVault{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
