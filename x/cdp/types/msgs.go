package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

const (
	TypeMsgCreateCDP = "create_cdp"
)

var _ sdk.Msg = &MsgCreateCDPRequest{}

func (msg MsgCreateCDPRequest) Route() string { return RouterKey }
func (msg MsgCreateCDPRequest) Type() string  { return TypeMsgCreateCDP }
func (msg MsgCreateCDPRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}
	if msg.Collateral.IsZero() || !msg.Collateral.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "collateral amount %s", msg.Collateral)
	}
	if msg.Principal.IsZero() || !msg.Principal.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "principal amount %s", msg.Principal)
	}
	if strings.TrimSpace(msg.CollateralType) == "" {
		return fmt.Errorf("collateral type cannot be empty")
	}
	return nil
}
func (msg MsgCreateCDPRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgCreateCDPRequest) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgDepositRequest{}

func (msg MsgDepositRequest) Route() string { return RouterKey }
func (msg MsgDepositRequest) Type() string  { return TypeMsgCreateCDP }
func (msg MsgDepositRequest) ValidateBasic() error {

	return nil
}
func (msg MsgDepositRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgDepositRequest) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgWithdrawRequest{}

func (msg MsgWithdrawRequest) Route() string { return RouterKey }
func (msg MsgWithdrawRequest) Type() string  { return TypeMsgCreateCDP }
func (msg MsgWithdrawRequest) ValidateBasic() error {

	return nil
}
func (msg MsgWithdrawRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgWithdrawRequest) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgDrawDebtRequest{}

func (msg MsgDrawDebtRequest) Route() string { return RouterKey }
func (msg MsgDrawDebtRequest) Type() string  { return TypeMsgCreateCDP }
func (msg MsgDrawDebtRequest) ValidateBasic() error {

	return nil
}
func (msg MsgDrawDebtRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgDrawDebtRequest) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgRepayDebtRequest{}

func (msg MsgRepayDebtRequest) Route() string { return RouterKey }
func (msg MsgRepayDebtRequest) Type() string  { return TypeMsgCreateCDP }
func (msg MsgRepayDebtRequest) ValidateBasic() error {

	return nil
}
func (msg MsgRepayDebtRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgRepayDebtRequest) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgLiquidateRequest{}

func (msg MsgLiquidateRequest) Route() string { return RouterKey }
func (msg MsgLiquidateRequest) Type() string  { return TypeMsgCreateCDP }
func (msg MsgLiquidateRequest) ValidateBasic() error {

	return nil
}
func (msg MsgLiquidateRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgLiquidateRequest) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
