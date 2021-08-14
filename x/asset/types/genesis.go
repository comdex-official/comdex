package types

func NewGenesisState(pools Pools) *GenesisState {
	return &GenesisState{
		Pools: pools,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(nil)
}

func ValidateGenesis(_ *GenesisState) error {
	return nil
}
