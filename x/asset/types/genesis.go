package types

func NewGenesisState(pairs []Pair, params Params) *GenesisState {
	return &GenesisState{
		Pairs:  pairs,
		Params: params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(nil, Params{})
}

func ValidateGenesis(_ *GenesisState) error {
	return nil
}
