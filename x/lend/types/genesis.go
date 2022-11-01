package types

func NewGenesisState(borrowAsset []BorrowAsset, borrowInterestTracker []BorrowInterestTracker, lendAsset []LendAsset, pool []Pool, assetToPairMapping []AssetToPairMapping, poolAssetLBMapping []PoolAssetLBMapping, lendRewardsTracker []LendRewardsTracker, userAssetLendBorrowMapping []UserAssetLendBorrowMapping, reserveBuybackAssetData []ReserveBuybackAssetData, extendedPair []Extended_Pair, auctionParams []AuctionParams, assetRatesParams []AssetRatesParams, params Params) *GenesisState {
	return &GenesisState{
		BorrowAsset:                borrowAsset,
		BorrowInterestTracker:      borrowInterestTracker,
		LendAsset:                  lendAsset,
		Pool:                       pool,
		AssetToPairMapping:         assetToPairMapping,
		PoolAssetLBMapping:         poolAssetLBMapping,
		LendRewardsTracker:         lendRewardsTracker,
		UserAssetLendBorrowMapping: userAssetLendBorrowMapping,
		ReserveBuybackAssetData:    reserveBuybackAssetData,
		Extended_Pair:              extendedPair,
		AuctionParams:              auctionParams,
		AssetRatesParams:           assetRatesParams,
		Params:                     params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]BorrowAsset{},
		[]BorrowInterestTracker{},
		[]LendAsset{},
		[]Pool{},
		[]AssetToPairMapping{},
		[]PoolAssetLBMapping{},
		[]LendRewardsTracker{},
		[]UserAssetLendBorrowMapping{},
		[]ReserveBuybackAssetData{},
		[]Extended_Pair{},
		[]AuctionParams{},
		[]AssetRatesParams{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
