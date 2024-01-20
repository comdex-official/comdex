package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = (*MsgLiquidateVaultRequest)(nil)
	_ sdk.Msg = (*MsgLiquidateBorrowRequest)(nil)
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
		return errorsmod.Wrap(ErrAppIDInvalid, "app_id cannot be zero")
	}
	if m.VaultId == 0 {
		return errorsmod.Wrap(ErrVaultIDInvalid, "vault_id cannot be nil")
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

func NewMsgLiquidateBorrowRequest(
	from sdk.AccAddress,
	borrowID uint64,
) *MsgLiquidateBorrowRequest {
	return &MsgLiquidateBorrowRequest{
		From:     from.String(),
		BorrowId: borrowID,
	}
}

func (m *MsgLiquidateBorrowRequest) Route() string {
	return RouterKey
}

func (m *MsgLiquidateBorrowRequest) Type() string {
	return TypeMsgLiquidateBorrowRequest
}

func (m *MsgLiquidateBorrowRequest) ValidateBasic() error {
	if m.BorrowId == 0 {
		return errorsmod.Wrap(ErrVaultIDInvalid, "borrow_id cannot be zero")
	}

	return nil
}

func (m *MsgLiquidateBorrowRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgLiquidateBorrowRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
