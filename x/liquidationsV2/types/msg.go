package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgLiquidateInternalKeeperRequest(
	from sdk.AccAddress,
	liqType, id uint64,
) *MsgLiquidateInternalKeeperRequest {
	return &MsgLiquidateInternalKeeperRequest{
		From:    from.String(),
		LiqType: liqType,
		Id:      id,
	}
}

func (m *MsgLiquidateInternalKeeperRequest) Route() string {
	return RouterKey
}

func (m *MsgLiquidateInternalKeeperRequest) Type() string {
	return TypeMsgLiquidateRequest
}

func (m *MsgLiquidateInternalKeeperRequest) ValidateBasic() error {
	if m.Id == 0 {
		return errors.Wrap(ErrVaultIDInvalid, "id cannot be zero")
	}

	return nil
}

func (m *MsgLiquidateInternalKeeperRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgLiquidateInternalKeeperRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
