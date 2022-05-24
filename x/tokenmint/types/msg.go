package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgMintNewTokensRequest)(nil)
	

)

func NewMsgMintNewTokensRequest(from sdk.AccAddress, app_mapping_id uint64,asset_id uint64) *MsgMintNewTokensRequest {
	return &MsgMintNewTokensRequest{
		From:                from.String(),
		AppMappingId:        app_mapping_id,
		AssetId: asset_id,
	}
}

func (m *MsgMintNewTokensRequest) Route() string {
	return RouterKey
}

func (m *MsgMintNewTokensRequest) Type() string {
	return TypeMsgMintNewTokensRequest
}

func (m *MsgMintNewTokensRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}

	return nil
}

func (m *MsgMintNewTokensRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgMintNewTokensRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
