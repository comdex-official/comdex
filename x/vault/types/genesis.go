package types

func NewGenesisState(vaults []Vault) *GenesisState {
	return &GenesisState{
		Vaults: vaults,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(nil)
}

func (m *GenesisState) Validate() error {
	return nil
}
