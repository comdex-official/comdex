package types

func NewGenesisState(assets []Asset, markets []Market, pairs []Pair, params Params) *GenesisState {
	return &GenesisState{
		Assets:  assets,
		Markets: markets,
		Pairs:   pairs,
		Params:  params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		nil,
		nil,
		nil,
		DefaultParams(),
	)
}

func ValidateGenesis(_ *GenesisState) error {
	return nil
}
