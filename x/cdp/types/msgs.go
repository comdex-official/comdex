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

var _ sdk.Msg = &MsgCreateCDP{}

func (msg MsgCreateCDP) Route() string { return RouterKey }
func (msg MsgCreateCDP) Type() string  { return TypeMsgCreateCDP }
func (msg MsgCreateCDP) ValidateBasic() error {
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
func (msg MsgCreateCDP) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
func (msg MsgCreateCDP) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
