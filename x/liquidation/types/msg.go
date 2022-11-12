package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgLiquidateVaultRequest)(nil)
)

func NewMsgLiquidateRequest(
	from sdk.AccAddress,
	appID, vaultID uint64,
) *MsgLiquidateVaultRequest {
	return &MsgLiquidateVaultRequest{
		From:    from.String(),
		AppId:   appID,
		VaultId: vaultID,
	}
}

func (m *MsgLiquidateVaultRequest) Route() string {
	return RouterKey
}

func (m *MsgLiquidateVaultRequest) Type() string {
	return TypeMsgLiquidateRequest
}

func (m *MsgLiquidateVaultRequest) ValidateBasic() error {

	if m.AppId == 0 {
		return errors.Wrap(ErrAppIDInvalid, "app_id cannot be zero")
	}
	if m.VaultId == 0 {
		return errors.Wrap(ErrVaultIDInvalid, "vault_id cannot be nil")
	}

	return nil
}

func (m *MsgLiquidateVaultRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgLiquidateVaultRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
