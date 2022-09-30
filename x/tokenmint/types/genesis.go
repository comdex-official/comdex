package types

func NewGenesisState(tokenMint []TokenMint, params Params) *GenesisState {
	return &GenesisState{
		TokenMint: tokenMint,
		Params:    params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]TokenMint{},
		DefaultParams(),
	)
}

func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
