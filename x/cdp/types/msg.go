package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateRequest)(nil)
	_ sdk.Msg = (*MsgDepositRequest)(nil)
	_ sdk.Msg = (*MsgWithdrawRequest)(nil)
	_ sdk.Msg = (*MsgDrawRequest)(nil)
	_ sdk.Msg = (*MsgRepayRequest)(nil)
	_ sdk.Msg = (*MsgLiquidateRequest)(nil)
)

func NewMsgCreateRequest(sender sdk.AccAddress, pairId uint64, amountIn, amountOut sdk.Int) *MsgCreateRequest {
	return &MsgCreateRequest{
		Sender:    sender.String(),
		PairId:    pairId,
		AmountIn:  amountIn,
		AmountOut: amountOut,
	}
}

func (m *MsgCreateRequest) Route() string {
	return RouterKey
}

func (m *MsgCreateRequest) Type() string {
	return TypeMsgCreateRequest
}

func (m *MsgCreateRequest) ValidateBasic() error {
	if m.Sender == "" {
		return errors.Wrap(ErrorInvalidSender, "sender cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errors.Wrapf(ErrorInvalidSender, "%s", err)
	}
	if m.AmountIn.IsNil() {
		return errors.Wrap(ErrorInvalidAmountIn, "amount_in cannot be nil")
	}
	if m.AmountIn.IsNegative() {
		return errors.Wrap(ErrorInvalidAmountIn, "amount_in cannot be negative")
	}
	if m.AmountIn.IsZero() {
		return errors.Wrap(ErrorInvalidAmountIn, "amount_in cannot be zero")
	}
	if m.AmountOut.IsNil() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be nil")
	}
	if m.AmountOut.IsNegative() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be negative")
	}
	if m.AmountOut.IsZero() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be zero")
	}

	return nil
}

func (m *MsgCreateRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgCreateRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgDepositRequest(sender sdk.AccAddress, id uint64, amount sdk.Int) *MsgDepositRequest {
	return &MsgDepositRequest{
		Sender: sender.String(),
		Id:     id,
		Amount: amount,
	}
}

func (m *MsgDepositRequest) Route() string {
	return RouterKey
}

func (m *MsgDepositRequest) Type() string {
	return TypeMsgDepositRequest
}

func (m *MsgDepositRequest) ValidateBasic() error {
	if m.Sender == "" {
		return errors.Wrap(ErrorInvalidSender, "sender cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errors.Wrapf(ErrorInvalidSender, "%s", err)
	}
	if m.Id == 0 {
		return errors.Wrap(ErrorInvalidId, "id cannot be zero")
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be zero")
	}

	return nil
}

func (m *MsgDepositRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgDepositRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgWithdrawRequest(sender sdk.AccAddress, id uint64, amount sdk.Int) *MsgWithdrawRequest {
	return &MsgWithdrawRequest{
		Sender: sender.String(),
		Id:     id,
		Amount: amount,
	}
}

func (m *MsgWithdrawRequest) Route() string {
	return RouterKey
}

func (m *MsgWithdrawRequest) Type() string {
	return TypeMsgWithdrawRequest
}

func (m *MsgWithdrawRequest) ValidateBasic() error {
	if m.Sender == "" {
		return errors.Wrap(ErrorInvalidSender, "sender cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errors.Wrapf(ErrorInvalidSender, "%s", err)
	}
	if m.Id == 0 {
		return errors.Wrap(ErrorInvalidId, "id cannot be zero")
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be zero")
	}

	return nil
}

func (m *MsgWithdrawRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgWithdrawRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgDrawRequest(sender sdk.AccAddress, id uint64, amount sdk.Int) *MsgDrawRequest {
	return &MsgDrawRequest{
		Sender: sender.String(),
		Id:     id,
		Amount: amount,
	}
}

func (m *MsgDrawRequest) Route() string {
	return RouterKey
}

func (m *MsgDrawRequest) Type() string {
	return TypeMsgDrawRequest
}

func (m *MsgDrawRequest) ValidateBasic() error {
	if m.Sender == "" {
		return errors.Wrap(ErrorInvalidSender, "sender cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errors.Wrapf(ErrorInvalidSender, "%s", err)
	}
	if m.Id == 0 {
		return errors.Wrap(ErrorInvalidId, "id cannot be zero")
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be zero")
	}

	return nil
}

func (m *MsgDrawRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgDrawRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgRepayRequest(sender sdk.AccAddress, id uint64, amount sdk.Int) *MsgRepayRequest {
	return &MsgRepayRequest{
		Sender: sender.String(),
		Id:     id,
		Amount: amount,
	}
}

func (m *MsgRepayRequest) Route() string {
	return RouterKey
}

func (m *MsgRepayRequest) Type() string {
	return TypeMsgRepayRequest
}

func (m *MsgRepayRequest) ValidateBasic() error {
	if m.Sender == "" {
		return errors.Wrap(ErrorInvalidSender, "sender cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errors.Wrapf(ErrorInvalidSender, "%s", err)
	}
	if m.Id == 0 {
		return errors.Wrap(ErrorInvalidId, "id cannot be zero")
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be zero")
	}

	return nil
}

func (m *MsgRepayRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgRepayRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgLiquidateRequest(sender sdk.AccAddress, id uint64) *MsgLiquidateRequest {
	return &MsgLiquidateRequest{
		Sender: sender.String(),
		Id:     id,
	}
}

func (m *MsgLiquidateRequest) Route() string {
	return RouterKey
}

func (m *MsgLiquidateRequest) Type() string {
	return TypeMsgLiquidateRequest
}

func (m *MsgLiquidateRequest) ValidateBasic() error {
	if m.Sender == "" {
		return errors.Wrap(ErrorInvalidSender, "sender cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errors.Wrapf(ErrorInvalidSender, "%s", err)
	}
	if m.Id == 0 {
		return errors.Wrap(ErrorInvalidId, "id cannot be zero")
	}

	return nil
}

func (m *MsgLiquidateRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgLiquidateRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
