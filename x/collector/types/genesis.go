package types

func NewGenesisState(netFeeCollectedData []AppAssetIdToFeeCollectedData, appIDToAssetCollectorMapping []AppToAssetIdCollectorMapping, collectorLookup []CollectorLookupTableData, collectorAuctionLookupTable []AppAssetIdToAuctionLookupTable, appToDenomsMapping []AppToDenomsMapping, params Params) *GenesisState {
	return &GenesisState{
		NetFeeCollectedData:          netFeeCollectedData,
		AppIdToAssetCollectorMapping: appIDToAssetCollectorMapping,
		CollectorLookup:              collectorLookup,
		CollectorAuctionLookupTable:  collectorAuctionLookupTable,
		AppToDenomsMapping:           appToDenomsMapping,
		Params:                       params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]AppAssetIdToFeeCollectedData{},
		[]AppToAssetIdCollectorMapping{},
		[]CollectorLookupTableData{},
		[]AppAssetIdToAuctionLookupTable{},
		[]AppToDenomsMapping{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
