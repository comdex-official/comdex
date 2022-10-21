package types

func NewGenesisState(surplusAuction []SurplusAuction, debtAuction []DebtAuction, dutchAuction []DutchAuction, protocolStatistics []ProtocolStatistics, auctionParams []AuctionParams, params Params, userBiddingID uint64) *GenesisState {
	return &GenesisState{
		SurplusAuction:     surplusAuction,
		DebtAuction:        debtAuction,
		DutchAuction:       dutchAuction,
		ProtocolStatistics: protocolStatistics,
		AuctionParams:      auctionParams,
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
		DefaultParams(),
		UserBiddingID,
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
