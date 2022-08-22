package types

func NewGenesisState(internal_rewards []InternalRewards, locker_rewards_tracker []LockerRewardsTracker, vault_interest_tracker []VaultInterestTracker, locker_external_rewards []LockerExternalRewards, vault_external_rewards []VaultExternalRewards, params Params) *GenesisState {
	return &GenesisState{
		InternalRewards: internal_rewards,
		LockerRewardsTracker: locker_rewards_tracker,
		VaultInterestTracker: vault_interest_tracker,
		LockerExternalRewards: locker_external_rewards,
		VaultExternalRewards: vault_external_rewards,
		Params:                           params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]InternalRewards{},
		[]LockerRewardsTracker{},
		[]VaultInterestTracker{},
		[]LockerExternalRewards{},
		[]VaultExternalRewards{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
