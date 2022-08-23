package types

func NewGenesisState(lockedVault []LockedVault, whitelistedApps []uint64, params Params) *GenesisState {
	return &GenesisState{
		LockedVault:     lockedVault,
		WhitelistedApps: whitelistedApps,
		Params:          params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]LockedVault{},
		[]uint64{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
