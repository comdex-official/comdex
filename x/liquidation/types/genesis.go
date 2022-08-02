package types

func NewGenesisState(lockedVault []LockedVault, lockedVaultToAppMapping []LockedVaultToAppMapping, whitelistedAppIds WhitelistedAppIds, params Params) *GenesisState {
	return &GenesisState{
		LockedVault: lockedVault,
		LockedVaultToAppMapping: lockedVaultToAppMapping,
		WhitelistedAppIds: whitelistedAppIds,
		Params: params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]LockedVault{},
		[]LockedVaultToAppMapping{},
		WhitelistedAppIds{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
