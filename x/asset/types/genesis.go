package types

func NewGenesisState(assets []Asset, pairs []Pair, appData []AppData, extendedPairVault []ExtendedPairVault, params Params) *GenesisState {
	return &GenesisState{
		Assets:            assets,
		Pairs:             pairs,
		AppData:           appData,
		ExtendedPairVault: extendedPairVault,
		Params:            params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]Asset{},
		[]Pair{},
		[]AppData{},
		[]ExtendedPairVault{},
		DefaultParams(),
	)
}

func ValidateGenesis(_ *GenesisState) error {
	return nil
}
