package types

func NewGenesisState(twa []TimeWeightedAverage) *GenesisState {
	return &GenesisState{
		TimeWeightedAverage: twa,
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
