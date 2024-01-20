package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		return errorsmod.Wrap(ErrVaultIDInvalid, "id cannot be zero")
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

func NewMsgAppReserveFundsRequest(from string, appId, assetId uint64, TokenQuantity sdk.Coin) *MsgAppReserveFundsRequest {
	return &MsgAppReserveFundsRequest{
		AppId:         appId,
		AssetId:       assetId,
		TokenQuantity: TokenQuantity,
		From:          from,
	}
}

func (m *MsgAppReserveFundsRequest) Route() string {
	return RouterKey
}

func (m *MsgAppReserveFundsRequest) Type() string {
	return TypeAppReserveFundsRequest
}

func (m *MsgAppReserveFundsRequest) ValidateBasic() error {
	if m.AppId == 0 || m.AssetId == 0 || m.TokenQuantity.Amount == sdkmath.NewInt(0) {
		return errorsmod.Wrap(ErrVaultIDInvalid, "id cannot be zero")
	}

	return nil
}

func (m *MsgAppReserveFundsRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgAppReserveFundsRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgLiquidateExternalKeeperRequest(
	from sdk.AccAddress,
	appId uint64,
	owner string,
	collateralToken, debtToken sdk.Coin,
	collateralAssetId, debtAssetId uint64,
	isDebtCmst bool,
) *MsgLiquidateExternalKeeperRequest {
	return &MsgLiquidateExternalKeeperRequest{
		From:              from.String(),
		AppId:             appId,
		Owner:             owner,
		CollateralToken:   collateralToken,
		DebtToken:         debtToken,
		CollateralAssetId: collateralAssetId,
		DebtAssetId:       debtAssetId,
		IsDebtCmst:        isDebtCmst,
	}
}

func (m *MsgLiquidateExternalKeeperRequest) Route() string {
	return RouterKey
}

func (m *MsgLiquidateExternalKeeperRequest) Type() string {
	return TypeMsgLiquidateExternalRequest
}

func (m *MsgLiquidateExternalKeeperRequest) ValidateBasic() error {
	if m.AppId == 0 {
		return errorsmod.Wrap(ErrVaultIDInvalid, "app_id cannot be zero")
	}

	return nil
}

func (m *MsgLiquidateExternalKeeperRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgLiquidateExternalKeeperRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
