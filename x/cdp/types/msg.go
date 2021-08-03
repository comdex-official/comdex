package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateCDPRequest)(nil)
	_ sdk.Msg = (*MsgDepositRequest)(nil)
	_ sdk.Msg = (*MsgWithdrawRequest)(nil)
	_ sdk.Msg = (*MsgDrawDebtRequest)(nil)
	_ sdk.Msg = (*MsgRepayDebtRequest)(nil)
	_ sdk.Msg = (*MsgLiquidateRequest)(nil)
)

// returns a new NewMsgCreateCDPRequest.
func NewMsgCreateCDPRequest(sender sdk.AccAddress, collateral sdk.Coin, principal sdk.Coin, collateralType string) *MsgCreateCDPRequest {
	return &MsgCreateCDPRequest{
		Sender:         sender.String(),
		Collateral:     collateral,
		Principal:      principal,
		CollateralType: collateralType,
	}
}

func (msg *MsgCreateCDPRequest) Route() string { return RouterKey }

func (msg *MsgCreateCDPRequest) Type() string { return TypeMsgCreateCDPRequest }

func (msg MsgCreateCDPRequest) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if msg.Collateral.IsZero() || !msg.Collateral.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "collateral amount %s", msg.Collateral)
	}
	if msg.Principal.IsZero() || !msg.Principal.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "principal amount %s", msg.Principal)
	}
	if strings.TrimSpace(msg.CollateralType) == "" {
		return fmt.Errorf("collateral type cannot be empty")
	}
	return nil
}

func (msg *MsgCreateCDPRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgCreateCDPRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

//returns new deposit request
func NewMsgDepositRequest(sender sdk.AccAddress, collateral sdk.Coin, collateralType string) *MsgDepositRequest {
	return &MsgDepositRequest{
		Sender:         sender.String(),
		Collateral:     collateral,
		CollateralType: collateralType,
	}
}

func (msg *MsgDepositRequest) Route() string { return RouterKey }

func (msg *MsgDepositRequest) Type() string { return TypeMsgDepositRequest }

func (msg MsgDepositRequest) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if msg.Collateral.IsZero() || !msg.Collateral.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "collateral amount %s", msg.Collateral)
	}
	if strings.TrimSpace(msg.CollateralType) == "" {
		return fmt.Errorf("collateral type cannot be empty")
	}
	return nil
}

func (msg *MsgDepositRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgDepositRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

//returns new withdraw request
func NewMsgWithdrawRequest(sender sdk.AccAddress, collateral sdk.Coin, collateralType string) *MsgWithdrawRequest {
	return &MsgWithdrawRequest{
		Sender:         sender.String(),
		Collateral:     collateral,
		CollateralType: collateralType,
	}
}

func (msg *MsgWithdrawRequest) Route() string { return RouterKey }

func (msg *MsgWithdrawRequest) Type() string { return TypeMsgWithdrawRequest }

func (msg MsgWithdrawRequest) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if msg.Collateral.IsZero() || !msg.Collateral.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "collateral amount %s", msg.Collateral)
	}
	if strings.TrimSpace(msg.CollateralType) == "" {
		return fmt.Errorf("collateral type cannot be empty")
	}
	return nil
}

func (msg *MsgWithdrawRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgWithdrawRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

//returns new draw debt request
func NewMsgDrawDebtRequest(sender sdk.AccAddress, collateralType string, principal sdk.Coin) *MsgDrawDebtRequest {
	return &MsgDrawDebtRequest{
		Sender:         sender.String(),
		CollateralType: collateralType,
		Principal:      principal,
	}
}

func (msg *MsgDrawDebtRequest) Route() string { return RouterKey }

func (msg *MsgDrawDebtRequest) Type() string { return TypeMsgDrawDebtRequest }

func (msg MsgDrawDebtRequest) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if msg.Principal.IsZero() || !msg.Principal.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "Principal amount %s", msg.Principal)
	}
	if strings.TrimSpace(msg.CollateralType) == "" {
		return fmt.Errorf("collateral type cannot be empty")
	}
	return nil
}

func (msg *MsgDrawDebtRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgDrawDebtRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

//returns new repay debt request
func NewMsgRepayDebtRequest(sender sdk.AccAddress, collateralType string, payment sdk.Coin) *MsgRepayDebtRequest {
	return &MsgRepayDebtRequest{
		Sender:         sender.String(),
		CollateralType: collateralType,
		Payment:        payment,
	}
}

func (msg *MsgRepayDebtRequest) Route() string { return RouterKey }

func (msg *MsgRepayDebtRequest) Type() string { return TypeMsgRepayDebtRequest }

func (msg MsgRepayDebtRequest) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if msg.Payment.IsZero() || !msg.Payment.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "payment amount %s", msg.Payment)
	}
	if strings.TrimSpace(msg.CollateralType) == "" {
		return fmt.Errorf("collateral type cannot be empty")
	}
	return nil
}

func (msg *MsgRepayDebtRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgRepayDebtRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

//returns new liquidate request
func NewMsgLiquidateRequest(sender sdk.AccAddress, collateralType string) *MsgLiquidateRequest {
	return &MsgLiquidateRequest{
		Sender:         sender.String(),
		CollateralType: collateralType,
	}
}

func (msg *MsgLiquidateRequest) Route() string { return RouterKey }

func (msg *MsgLiquidateRequest) Type() string { return TypeMsgLiquidateRequest }

func (msg MsgLiquidateRequest) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if strings.TrimSpace(msg.CollateralType) == "" {
		return fmt.Errorf("collateral type cannot be empty")
	}
	return nil
}

func (msg *MsgLiquidateRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgLiquidateRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
