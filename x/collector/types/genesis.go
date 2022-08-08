package types

func NewGenesisState(netFeeCollectedData []NetFeeCollectedData, appIDToAssetCollectorMapping []AppIdToAssetCollectorMapping, collectorLookup []CollectorLookup, collectorAuctionLookupTable []CollectorAuctionLookupTable, appToDenomsMapping []AppToDenomsMapping, params Params) *GenesisState {
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
		[]NetFeeCollectedData{},
		[]AppIdToAssetCollectorMapping{},
		[]CollectorLookup{},
		[]CollectorAuctionLookupTable{},
		[]AppToDenomsMapping{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
