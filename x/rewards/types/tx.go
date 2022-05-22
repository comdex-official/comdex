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
