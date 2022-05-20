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
