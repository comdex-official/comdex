package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateCDPRequest)(nil)
	_ sdk.Msg = (*MsgDepositCollateralRequest)(nil)
	_ sdk.Msg = (*MsgWithdrawCollateralRequest)(nil)
	_ sdk.Msg = (*MsgDrawDebtRequest)(nil)
	_ sdk.Msg = (*MsgRepayDebtRequest)(nil)
	_ sdk.Msg = (*MsgLiquidateCDPRequest)(nil)
)

// returns a new NewMsgCreateCDPRequest.
func NewMsgCreateCDPRequest(sender sdk.AccAddress, collateral sdk.Coin, debt sdk.Coin, collateralType string) *MsgCreateCDPRequest {
	return &MsgCreateCDPRequest{
		Sender:         sender.String(),
		Collateral:     collateral,
		Debt:           debt,
		CollateralType: collateralType,
	}
}

func (msg *MsgCreateCDPRequest) Route() string { return RouterKey }

func (msg *MsgCreateCDPRequest) Type() string { return TypeMsgCreateCDPRequest }

func (msg MsgCreateCDPRequest) ValidateBasic() error {
	if msg.Sender == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if !msg.Collateral.IsPositive() || !msg.Collateral.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "collateral amount %s", msg.Collateral)
	}
	if !msg.Debt.IsPositive() || !msg.Debt.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "principal amount %s", msg.Debt)
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
func NewMsgDepositCollateralRequest(owner sdk.AccAddress, collateral sdk.Coin, collateralType string) *MsgDepositCollateralRequest {
	return &MsgDepositCollateralRequest{
		Owner:          owner.String(),
		Collateral:     collateral,
		CollateralType: collateralType,
	}
}

func (msg *MsgDepositCollateralRequest) Route() string { return RouterKey }

func (msg *MsgDepositCollateralRequest) Type() string { return TypeMsgDepositCollateralRequest }

func (msg MsgDepositCollateralRequest) ValidateBasic() error {
	if msg.Owner == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if !msg.Collateral.IsPositive() || !msg.Collateral.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "collateral amount %s", msg.Collateral)
	}
	if strings.TrimSpace(msg.CollateralType) == "" {
		return fmt.Errorf("collateral type cannot be empty")
	}
	return nil
}

func (msg *MsgDepositCollateralRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgDepositCollateralRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

//returns new withdraw request
func NewMsgWithdrawCollateralRequest(owner sdk.AccAddress, collateral sdk.Coin, collateralType string) *MsgWithdrawCollateralRequest {
	return &MsgWithdrawCollateralRequest{
		Owner:          owner.String(),
		Collateral:     collateral,
		CollateralType: collateralType,
	}
}

func (msg *MsgWithdrawCollateralRequest) Route() string { return RouterKey }

func (msg *MsgWithdrawCollateralRequest) Type() string { return TypeMsgWithdrawCollateralRequest }

func (msg MsgWithdrawCollateralRequest) ValidateBasic() error {
	if msg.Owner == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if !msg.Collateral.IsPositive() || !msg.Collateral.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "collateral amount %s", msg.Collateral)
	}
	if strings.TrimSpace(msg.CollateralType) == "" {
		return fmt.Errorf("collateral type cannot be empty")
	}
	return nil
}

func (msg *MsgWithdrawCollateralRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgWithdrawCollateralRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

//returns new draw debt request
func NewMsgDrawDebtRequest(owner sdk.AccAddress, collateralType string, debt sdk.Coin) *MsgDrawDebtRequest {
	return &MsgDrawDebtRequest{
		Owner:          owner.String(),
		CollateralType: collateralType,
		Debt:           debt,
	}
}

func (msg *MsgDrawDebtRequest) Route() string { return RouterKey }

func (msg *MsgDrawDebtRequest) Type() string { return TypeMsgDrawDebtRequest }

func (msg MsgDrawDebtRequest) ValidateBasic() error {
	if msg.Owner == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if !msg.Debt.IsPositive() || !msg.Debt.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "Debt amount %s", msg.Debt)
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
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

//returns new repay debt request
func NewMsgRepayDebtRequest(owner sdk.AccAddress, collateralType string, debt sdk.Coin) *MsgRepayDebtRequest {
	return &MsgRepayDebtRequest{
		Owner:          owner.String(),
		CollateralType: collateralType,
		Debt:           debt,
	}
}

func (msg *MsgRepayDebtRequest) Route() string { return RouterKey }

func (msg *MsgRepayDebtRequest) Type() string { return TypeMsgRepayDebtRequest }

func (msg MsgRepayDebtRequest) ValidateBasic() error {
	if msg.Owner == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if !msg.Debt.IsPositive() || !msg.Debt.IsValid() {
		return errors.Wrapf(ErrorInvalidCoins, "payment amount %s", msg.Debt)
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
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

//returns new liquidate request
func NewMsgLiquidateCDPRequest(sender sdk.AccAddress, collateralType string) *MsgLiquidateCDPRequest {
	return &MsgLiquidateCDPRequest{
		Owner:          sender.String(),
		CollateralType: collateralType,
	}
}

func (msg *MsgLiquidateCDPRequest) Route() string { return RouterKey }

func (msg *MsgLiquidateCDPRequest) Type() string { return TypeMsgLiquidateCDPRequest }

func (msg MsgLiquidateCDPRequest) ValidateBasic() error {
	if msg.Owner == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "sender address cannot be empty")
	}
	if strings.TrimSpace(msg.CollateralType) == "" {
		return fmt.Errorf("collateral type cannot be empty")
	}
	return nil
}

func (msg *MsgLiquidateCDPRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgLiquidateCDPRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
