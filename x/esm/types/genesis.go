package types

func NewGenesisState(eSMTriggerParams []ESMTriggerParams, currentDepositStats []CurrentDepositStats, eSMStatus []ESMStatus, killSwitchParams []KillSwitchParams, usersDepositMapping []UsersDepositMapping, eSMMarketPrice []ESMMarketPrice, dataAfterCoolOff []DataAfterCoolOff, assetToAmountValue []AssetToAmountValue, appToAmountValue []AppToAmountValue, params Params) *GenesisState {
	return &GenesisState{
		ESMTriggerParams: eSMTriggerParams,
		CurrentDepositStats: currentDepositStats,
		ESMStatus: eSMStatus,
		KillSwitchParams: killSwitchParams,
		UsersDepositMapping: usersDepositMapping,
		ESMMarketPrice: eSMMarketPrice,
		DataAfterCoolOff: dataAfterCoolOff,
		AssetToAmountValue: assetToAmountValue,
		AppToAmountValue: appToAmountValue,
		Params: params,
	}
}

func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		[]ESMTriggerParams{},
		[]CurrentDepositStats{},
		[]ESMStatus{},
		[]KillSwitchParams{},
		[]UsersDepositMapping{},
		[]ESMMarketPrice{},
		[]DataAfterCoolOff{},
		[]AssetToAmountValue{},
		[]AppToAmountValue{},
		DefaultParams(),
	)
}

func (m *GenesisState) Validate() error {
	return nil
}
