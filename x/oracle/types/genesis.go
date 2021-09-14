package types

func NewGenesisState(markets []Market) *GenesisState {
	return &GenesisState{
		Markets: markets,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		nil,
	)
}

func ValidateGenesis(_ *GenesisState) error {
	return nil
}

func (m GenesisState) Validate() error {

	return nil
}
