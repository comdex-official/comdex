package types

func NewGenesisState(surplusAuction []SurplusAuction, debtAuction []DebtAuction, dutchAuction []DutchAuction, protocolStatistics []ProtocolStatistics, auctionParams []AuctionParams, params Params) *GenesisState {
	return &GenesisState{
		SurplusAuction:     surplusAuction,
		DebtAuction:        debtAuction,
		DutchAuction:       dutchAuction,
		ProtocolStatistics: protocolStatistics,
		AuctionParams:      auctionParams,
		Params:             params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]SurplusAuction{},
		[]DebtAuction{},
		[]DutchAuction{},
		[]ProtocolStatistics{},
		[]AuctionParams{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
