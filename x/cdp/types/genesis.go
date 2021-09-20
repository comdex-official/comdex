package types

func NewGenesisState(cdps []CDP) *GenesisState {
	return &GenesisState{
		CDPs: cdps,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(nil)
}

func (m *GenesisState) Validate() error {
	return nil
}
