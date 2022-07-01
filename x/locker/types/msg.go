package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateLockerRequest)(nil)
	_ sdk.Msg = (*MsgDepositAssetRequest)(nil)
	_ sdk.Msg = (*MsgWithdrawAssetRequest)(nil)
	_ sdk.Msg = (*MsgAddWhiteListedAssetRequest)(nil)
)

func NewMsgCreateLockerRequest(from sdk.AccAddress, amount sdk.Int, assetID uint64, appMappingID uint64) *MsgCreateLockerRequest {
	return &MsgCreateLockerRequest{
		Depositor:    from.String(),
		AppId: appMappingID,
		AssetId:      assetID,
		Amount:       amount,
	}
}

func (m *MsgCreateLockerRequest) Route() string {
	return RouterKey
}

func (m *MsgCreateLockerRequest) Type() string {
	return TypeMsgCreateLockerRequest
}

func (m *MsgCreateLockerRequest) ValidateBasic() error {
	if m.Depositor == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Depositor); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}

	if m.AppId < 0 {
		return errors.Wrap(ErrorInvalidAppMappingID, "app_mapping_id  cannot be negative")
	}
	if m.AssetId < 0 {
		return errors.Wrap(ErrorInvalidAssetID, "asset_id cannot be negative")
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be zero")
	}

	return nil
}

func (m *MsgCreateLockerRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgCreateLockerRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Depositor)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgDepositAssetRequest(from sdk.AccAddress, lockerID string, amount sdk.Int, assetID uint64, appMappingID uint64) *MsgDepositAssetRequest {
	return &MsgDepositAssetRequest{
		Depositor:    from.String(),
		LockerId:     lockerID,
		Amount:       amount,
		AssetId:      assetID,
		AppId: appMappingID,
	}
}

func (m *MsgDepositAssetRequest) Route() string {
	return RouterKey
}

func (m *MsgDepositAssetRequest) Type() string {
	return TypeMsgDepositAssetRequest
}

func (m *MsgDepositAssetRequest) ValidateBasic() error {
	if m.Depositor == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Depositor); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}

	if m.AppId < 0 {
		return errors.Wrap(ErrorInvalidAppMappingID, "app_mapping_id  cannot be negative")
	}
	if m.AssetId < 0 {
		return errors.Wrap(ErrorInvalidAssetID, "asset_id cannot be negative")
	}
	if len(m.LockerId) <= 0 {
		return errors.Wrap(ErrorInvalidLockerID, "lockerID  cannot be negative")
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be zero")
	}

	return nil
}

func (m *MsgDepositAssetRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgDepositAssetRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Depositor)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgWithdrawAssetRequest(from sdk.AccAddress, lockerID string, amount sdk.Int, assetID uint64, appMappingID uint64) *MsgWithdrawAssetRequest {
	return &MsgWithdrawAssetRequest{
		Depositor:    from.String(),
		LockerId:     lockerID,
		Amount:       amount,
		AssetId:      assetID,
		AppId: appMappingID,
	}
}

func (m *MsgWithdrawAssetRequest) Route() string {
	return RouterKey
}

func (m *MsgWithdrawAssetRequest) Type() string {
	return TypeMsgWithdrawAssetRequest
}

func (m *MsgWithdrawAssetRequest) ValidateBasic() error {
	if m.Depositor == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Depositor); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}

	if m.AppId < 0 {
		return errors.Wrap(ErrorInvalidAppMappingID, "app_mapping_id  cannot be negative")
	}
	if m.AssetId < 0 {
		return errors.Wrap(ErrorInvalidAssetID, "asset_id cannot be negative")
	}
	if len(m.LockerId) <= 0 {
		return errors.Wrap(ErrorInvalidLockerID, "lockerID  cannot be negative")
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be zero")
	}

	return nil
}

func (m *MsgWithdrawAssetRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgWithdrawAssetRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Depositor)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgAddWhiteListedAssetRequest(from sdk.AccAddress, appMappingID uint64, assetID uint64) *MsgAddWhiteListedAssetRequest {
	return &MsgAddWhiteListedAssetRequest{
		From:         from.String(),
		AppId: appMappingID,
		AssetId:      assetID,
	}
}

func (m *MsgAddWhiteListedAssetRequest) Route() string {
	return RouterKey
}

func (m *MsgAddWhiteListedAssetRequest) Type() string {
	return TypeMsgAddWhiteListedAssetRequest
}

func (m *MsgAddWhiteListedAssetRequest) ValidateBasic() error {
	if m.AppId < 0 {
		return errors.Wrap(ErrorInvalidAppMappingID, "app_mapping_id  cannot be negative")
	}
	if m.AssetId < 0 {
		return errors.Wrap(ErrorInvalidAssetID, "asset_id cannot be negative")
	}

	return nil
}

func (m *MsgAddWhiteListedAssetRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgAddWhiteListedAssetRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
