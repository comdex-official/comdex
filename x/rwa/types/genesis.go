package types

func NewGenesisState() *GenesisState {
	return &GenesisState{}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState()
}
func (m *GenesisState) Validate() error {
	return nil
}
