package types

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state

func NewGenesisState(lockedVault []LockedVault, whitelistedApps []LiquidationWhiteListing, appReserveFunds []AppReserveFunds, params Params) *GenesisState {
	return &GenesisState{
		LockedVault:             lockedVault,
		LiquidationWhiteListing: whitelistedApps,
		AppReserveFunds:         appReserveFunds,
		Params:                  params,
	}
}

func DefaultGenesis() *GenesisState {
	return NewGenesisState(
		[]LockedVault{},
		[]LiquidationWhiteListing{},
		[]AppReserveFunds{},
		DefaultParams(),
	)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
