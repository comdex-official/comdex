package types

func NewGenesisState(whitelistedContracts []WhitelistedContract, params Params) *GenesisState {
	return &GenesisState{
		WhitelistedContracts: whitelistedContracts,
		Params:               params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]WhitelistedContract{},
		DefaultParams(),
	)
}

func (m *GenesisState) ValidateGenesis() error {
	return nil
}
