package types

func NewGenesisState(surplusAuction []SurplusAuction, debtAuction []DebtAuction, dutchAuction []DutchAuction, protocolStatistics []ProtocolStatistics, auctionParams []AuctionParams, dutchLendAuction []DutchAuction, params Params, userBiddingID uint64) *GenesisState {
	return &GenesisState{
		SurplusAuction:     surplusAuction,
		DebtAuction:        debtAuction,
		DutchAuction:       dutchAuction,
		ProtocolStatistics: protocolStatistics,
		AuctionParams:      auctionParams,
		DutchLendAuction:   dutchLendAuction,
		Params:             params,
		UserBiddingID:      userBiddingID,
	}
}

func DefaultGenesisState() *GenesisState {
	var UserBiddingID uint64
	return NewGenesisState(
		[]SurplusAuction{},
		[]DebtAuction{},
		[]DutchAuction{},
		[]ProtocolStatistics{},
		[]AuctionParams{},
		[]DutchAuction{},
		DefaultParams(),
		UserBiddingID,
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
