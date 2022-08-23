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

func NewMsgCreateLockerRequest(from string, amount sdk.Int, assetID uint64, appMappingID uint64) *MsgCreateLockerRequest {
	return &MsgCreateLockerRequest{
		Depositor: from,
		AppId:     appMappingID,
		AssetId:   assetID,
		Amount:    amount,
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

func NewMsgDepositAssetRequest(from string, lockerID uint64, amount sdk.Int, assetID uint64, appMappingID uint64) *MsgDepositAssetRequest {
	return &MsgDepositAssetRequest{
		Depositor: from,
		LockerId:  lockerID,
		Amount:    amount,
		AssetId:   assetID,
		AppId:     appMappingID,
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

	if m.LockerId <= 0 {
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

func NewMsgWithdrawAssetRequest(from string, lockerID uint64, amount sdk.Int, assetID uint64, appMappingID uint64) *MsgWithdrawAssetRequest {
	return &MsgWithdrawAssetRequest{
		Depositor: from,
		LockerId:  lockerID,
		Amount:    amount,
		AssetId:   assetID,
		AppId:     appMappingID,
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

	if m.LockerId <= 0 {
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

func NewMsgAddWhiteListedAssetRequest(from string, appMappingID uint64, assetID uint64) *MsgAddWhiteListedAssetRequest {
	return &MsgAddWhiteListedAssetRequest{
		From:    from,
		AppId:   appMappingID,
		AssetId: assetID,
	}
}

func (m *MsgAddWhiteListedAssetRequest) Route() string {
	return RouterKey
}

func (m *MsgAddWhiteListedAssetRequest) Type() string {
	return TypeMsgAddWhiteListedAssetRequest
}

func (m *MsgAddWhiteListedAssetRequest) ValidateBasic() error {

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

func NewMsgCloseLockerRequest(from string, appID uint64, assetID uint64, lockerID uint64) *MsgCloseLockerRequest {
	return &MsgCloseLockerRequest{
		Depositor: from,
		AppId:     appID,
		AssetId:   assetID,
		LockerId:  lockerID,
	}
}

func (m *MsgCloseLockerRequest) Route() string {
	return RouterKey
}

func (m *MsgCloseLockerRequest) Type() string {
	return TypeMsgWithdrawAssetRequest
}

func (m *MsgCloseLockerRequest) ValidateBasic() error {
	if m.Depositor == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Depositor); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}

	if m.LockerId <= 0 {
		return errors.Wrap(ErrorInvalidLockerID, "lockerID  cannot be negative")
	}

	return nil
}

func (m *MsgCloseLockerRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgCloseLockerRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Depositor)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}


func NewMsgLockerRewardCalcRequest(from string, appID uint64, lockerID uint64) *MsgLockerRewardCalcRequest {
	return &MsgLockerRewardCalcRequest{
		From: from,
		AppId:     appID,
		LockerId:  lockerID,
	}
}

func (m *MsgLockerRewardCalcRequest) Route() string {
	return RouterKey
}

func (m *MsgLockerRewardCalcRequest) Type() string {
	return TypeMsgLockerRewardCalcRequest
}

func (m *MsgLockerRewardCalcRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}

	if m.LockerId <= 0 {
		return errors.Wrap(ErrorInvalidLockerID, "lockerID  cannot be negative")
	}

	return nil
}

func (m *MsgLockerRewardCalcRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgLockerRewardCalcRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

