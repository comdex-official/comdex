package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgWhitelistAsset(appMappingId uint64, from sdk.AccAddress, assetId []uint64) *WhitelistAsset {
	return &WhitelistAsset{
		AppMappingId: appMappingId,
		From:         from.String(),
		AssetId:      assetId,
	}
}

func (m *WhitelistAsset) Route() string {
	return RouterKey
}

func (m *WhitelistAsset) Type() string {
	return ModuleName
}

func (m *WhitelistAsset) ValidateBasic() error {

	return nil
}

func (m *WhitelistAsset) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *WhitelistAsset) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgRemoveWhitelistAsset(appMappingId uint64, from sdk.AccAddress, assetId uint64) *RemoveWhitelistAsset {
	return &RemoveWhitelistAsset{
		AppMappingId: appMappingId,
		From:         from.String(),
		AssetId:      assetId,
	}
}

func (m *RemoveWhitelistAsset) Route() string {
	return RouterKey
}

func (m *RemoveWhitelistAsset) Type() string {
	return ModuleName
}

func (m *RemoveWhitelistAsset) ValidateBasic() error {

	return nil
}

func (m *RemoveWhitelistAsset) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *RemoveWhitelistAsset) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgWhitelistAppIdVault(appMappingId uint64, from sdk.AccAddress) *WhitelistAppIdVault {
	return &WhitelistAppIdVault{
		AppMappingId: appMappingId,
		From:         from.String(),
	}
}

func (m *WhitelistAppIdVault) Route() string {
	return RouterKey
}

func (m *WhitelistAppIdVault) Type() string {
	return ModuleName
}

func (m *WhitelistAppIdVault) ValidateBasic() error {

	return nil
}

func (m *WhitelistAppIdVault) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *WhitelistAppIdVault) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgRemoveWhitelistAppIdVault(appMappingId uint64, from sdk.AccAddress) *RemoveWhitelistAppIdVault {
	return &RemoveWhitelistAppIdVault{
		AppMappingId: appMappingId,
		From:         from.String(),
	}
}

func (m *RemoveWhitelistAppIdVault) Route() string {
	return RouterKey
}

func (m *RemoveWhitelistAppIdVault) Type() string {
	return ModuleName
}

func (m *RemoveWhitelistAppIdVault) ValidateBasic() error {

	return nil
}

func (m *RemoveWhitelistAppIdVault) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *RemoveWhitelistAppIdVault) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgActivateExternalRewardsLockers(appMappingId uint64, AssetId uint64, TotalRewards sdk.Coin, DurationDays, MinLockupTimeSeconds int64, from sdk.AccAddress) *ActivateExternalRewardsLockers {
	return &ActivateExternalRewardsLockers{
		AppMappingId:         appMappingId,
		AssetId:              AssetId,
		TotalRewards:         TotalRewards,
		DurationDays:         DurationDays,
		MinLockupTimeSeconds: MinLockupTimeSeconds,
		Depositor:            from.String(),
	}
}

func (m *ActivateExternalRewardsLockers) Route() string {
	return RouterKey
}

func (m *ActivateExternalRewardsLockers) Type() string {
	return ModuleName
}

func (m *ActivateExternalRewardsLockers) ValidateBasic() error {

	return nil
}

func (m *ActivateExternalRewardsLockers) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *ActivateExternalRewardsLockers) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Depositor)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgActivateExternalVaultLockers(appMappingId uint64, extendedPairId uint64, TotalRewards sdk.Coin, DurationDays, MinLockupTimeSeconds int64, from sdk.AccAddress) *ActivateExternalRewardsVault {
	return &ActivateExternalRewardsVault{
		AppMappingId:         appMappingId,
		Extended_Pair_Id:     extendedPairId,
		TotalRewards:         TotalRewards,
		DurationDays:         DurationDays,
		MinLockupTimeSeconds: MinLockupTimeSeconds,
		Depositor:            from.String(),
	}
}

func (m *ActivateExternalRewardsVault) Route() string {
	return RouterKey
}

func (m *ActivateExternalRewardsVault) Type() string {
	return ModuleName
}

func (m *ActivateExternalRewardsVault) ValidateBasic() error {

	return nil
}

func (m *ActivateExternalRewardsVault) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *ActivateExternalRewardsVault) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Depositor)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
