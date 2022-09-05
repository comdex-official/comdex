package types

func NewGenesisState(internalRewards []InternalRewards, lockerRewardsTracker []LockerRewardsTracker, vaultInterestTracker []VaultInterestTracker, lockerExternalRewards []LockerExternalRewards, vaultExternalRewards []VaultExternalRewards, appIDs []uint64, epochInfo []EpochInfo, gauge []Gauge, gaugeDuration []GaugeByTriggerDuration, params Params) *GenesisState {
	return &GenesisState{
		InternalRewards:        internalRewards,
		LockerRewardsTracker:   lockerRewardsTracker,
		VaultInterestTracker:   vaultInterestTracker,
		LockerExternalRewards:  lockerExternalRewards,
		VaultExternalRewards:   vaultExternalRewards,
		AppIDs:                 appIDs,
		EpochInfo:              epochInfo,
		Gauge:                  gauge,
		GaugeByTriggerDuration: gaugeDuration,
		Params:                 params,
	}
}
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]InternalRewards{},
		[]LockerRewardsTracker{},
		[]VaultInterestTracker{},
		[]LockerExternalRewards{},
		[]VaultExternalRewards{},
		[]uint64{},
		[]EpochInfo{},
		[]Gauge{},
		[]GaugeByTriggerDuration{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
