package types

func NewGenesisState(eSMTriggerParams []ESMTriggerParams, currentDepositStats []CurrentDepositStats, eSMStatus []ESMStatus, killSwitchParams []KillSwitchParams, usersDepositMapping []UsersDepositMapping, dataAfterCoolOff []DataAfterCoolOff, params Params) *GenesisState {
	return &GenesisState{
		ESMTriggerParams:    eSMTriggerParams,
		CurrentDepositStats: currentDepositStats,
		ESMStatus:           eSMStatus,
		KillSwitchParams:    killSwitchParams,
		UsersDepositMapping: usersDepositMapping,
		DataAfterCoolOff:    dataAfterCoolOff,
		Params:              params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]ESMTriggerParams{},
		[]CurrentDepositStats{},
		[]ESMStatus{},
		[]KillSwitchParams{},
		[]UsersDepositMapping{},
		[]DataAfterCoolOff{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
