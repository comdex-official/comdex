package types

func NewGenesisState(pairs []Pair) *GenesisState {
	return &GenesisState{
		Pairs: pairs,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(nil)
}

func ValidateGenesis(_ *GenesisState) error {
	return nil
}
