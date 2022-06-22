package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewMsgWhitelistAppID(appMappingId uint64, from sdk.AccAddress) *WhitelistAppId {
	return &WhitelistAppId{
		AppMappingId: appMappingId,
		From:         from.String(),
	}
}

func (m *WhitelistAppId) Route() string {
	return RouterKey
}

func (m *WhitelistAppId) Type() string {
	return ModuleName
}

func (m *WhitelistAppId) ValidateBasic() error {
	return nil
}

func (m *WhitelistAppId) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *WhitelistAppId) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgRemoveWhitelistAsset(appMappingID uint64, from sdk.AccAddress) *RemoveWhitelistAppId {
	return &RemoveWhitelistAppId{
		AppMappingId: appMappingID,
		From:         from.String(),
	}
}

func (m *RemoveWhitelistAppId) Route() string {
	return RouterKey
}

func (m *RemoveWhitelistAppId) Type() string {
	return ModuleName
}

func (m *RemoveWhitelistAppId) ValidateBasic() error {
	return nil
}

func (m *RemoveWhitelistAppId) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *RemoveWhitelistAppId) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
