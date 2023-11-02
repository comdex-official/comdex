package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgDeposit(addr string, amount sdk.Coin, appID uint64) *MsgDeposit {
	return &MsgDeposit{
		Addr:   addr,
		Amount: amount,
		AppId:  appID,
	}
}

func (msg MsgDeposit) Route() string { return ModuleName }
func (msg MsgDeposit) Type() string  { return TypeDepositRequest }

func (msg *MsgDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetAddr())
	if err != nil {
		return err
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}
	if msg.AppId == 0 {
		return fmt.Errorf("app id should not be 0: %d ", msg.AppId)
	}

	return nil
}

func (msg *MsgDeposit) GetSigners() []sdk.AccAddress {
	lender, _ := sdk.AccAddressFromBech32(msg.GetAddr())
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
