package types

func NewGenesisState(borrowAsset []BorrowAsset, userBorrowIDMapping []UserBorrowIdMapping, borrowIDByOwnerAndPoolMapping []BorrowIdByOwnerAndPoolMapping, borrowMapping BorrowMapping, lendAsset []LendAsset, pool []Pool, assetToPairMapping []AssetToPairMapping, userLendIDMapping []UserLendIdMapping, lendIDByOwnerAndPoolMapping []LendIdByOwnerAndPoolMapping, lendIDToBorrowIdMapping []LendIdToBorrowIdMapping, assetStats []AssetStats, lendMapping LendMapping, userdepositStats DepositStats, reservedepositStats DepositStats, buybackdepositStats DepositStats, borrowdepositStats DepositStats, extended_Pair []Extended_Pair, assetRatesStats []AssetRatesStats, auctionParams []AuctionParams, params Params) *GenesisState { //nolint:revive,stylecheck // seems like it may make sense to refactor this line later (pretty complex)
	return &GenesisState{
		BorrowAsset:                   borrowAsset,
		UserBorrowIdMapping:           userBorrowIDMapping,
		BorrowIdByOwnerAndPoolMapping: borrowIDByOwnerAndPoolMapping,
		BorrowMapping:                 borrowMapping,
		LendAsset:                     lendAsset,
		Pool:                          pool,
		AssetToPairMapping:            assetToPairMapping,
		UserLendIdMapping:             userLendIDMapping,
		LendIdByOwnerAndPoolMapping:   lendIDByOwnerAndPoolMapping,
		LendIdToBorrowIdMapping:       lendIDToBorrowIdMapping,
		AssetStats:                    assetStats,
		LendMapping:                   lendMapping,
		UserDepositStats:              userdepositStats,
		ReserveDepositStats:           reservedepositStats,
		BuyBackDepositStats:           buybackdepositStats,
		BorrowDepositStats:            borrowdepositStats,
		Extended_Pair:                 extended_Pair,
		AssetRatesStats:               assetRatesStats,
		AuctionParams:                 auctionParams,
		Params:                        params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]BorrowAsset{},
		[]UserBorrowIdMapping{},
		[]BorrowIdByOwnerAndPoolMapping{},
		BorrowMapping{},
		[]LendAsset{},
		[]Pool{},
		[]AssetToPairMapping{},
		[]UserLendIdMapping{},
		[]LendIdByOwnerAndPoolMapping{},
		[]LendIdToBorrowIdMapping{},
		[]AssetStats{},
		LendMapping{},
		DepositStats{},
		DepositStats{},
		DepositStats{},
		DepositStats{},
		[]Extended_Pair{},
		[]AssetRatesStats{},
		[]AuctionParams{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
