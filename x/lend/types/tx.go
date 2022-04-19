package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgLend(lender sdk.AccAddress, amount sdk.Coin) *MsgLend {
	return &MsgLend{
		Lender: lender.String(),
		Amount: amount,
	}
}

func (msg MsgLend) Route() string { return ModuleName }
func (msg MsgLend) Type() string  { return EventTypeLoanAsset }

func (msg *MsgLend) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}

	if asset := msg.GetAmount(); !asset.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAsset, asset.String())
	}

	return nil
}

func (msg *MsgLend) GetSigners() []sdk.AccAddress {
	lender, _ := sdk.AccAddressFromBech32(msg.GetLender())
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg *MsgLend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgWithdraw(lender sdk.AccAddress, amount sdk.Coin) *MsgWithdraw {
	return &MsgWithdraw{
		Lender: lender.String(),
		Amount: amount,
	}
}

func (msg MsgWithdraw) Route() string { return ModuleName }
func (msg MsgWithdraw) Type() string  { return EventTypeWithdrawLoanedAsset }

func (msg *MsgWithdraw) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}

	if asset := msg.GetAmount(); !asset.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAsset, asset.String())
	}

	return nil
}

func (msg *MsgWithdraw) GetSigners() []sdk.AccAddress {
	lender, _ := sdk.AccAddressFromBech32(msg.GetLender())
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg *MsgWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgBorrow(borrower sdk.AccAddress, amount sdk.Coin) *MsgBorrow {
	return &MsgBorrow{
		Borrower: borrower.String(),
		Amount:   amount,
	}
}

func (msg MsgBorrow) Route() string { return ModuleName }
func (msg MsgBorrow) Type() string  { return EventTypeBorrowAsset }

func (msg *MsgBorrow) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	if asset := msg.GetAmount(); !asset.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAsset, asset.String())
	}

	return nil
}

func (msg *MsgBorrow) GetSigners() []sdk.AccAddress {
	borrower, _ := sdk.AccAddressFromBech32(msg.GetBorrower())
	return []sdk.AccAddress{borrower}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg *MsgBorrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgRepay(borrower sdk.AccAddress, amount sdk.Coin) *MsgRepay {
	return &MsgRepay{
		Borrower: borrower.String(),
		Amount:   amount,
	}
}

func (msg MsgRepay) Route() string { return ModuleName }
func (msg MsgRepay) Type() string  { return EventTypeRepayBorrowedAsset }

func (msg *MsgRepay) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	if asset := msg.GetAmount(); !asset.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAsset, asset.String())
	}

	return nil
}

func (msg *MsgRepay) GetSigners() []sdk.AccAddress {
	borrower, _ := sdk.AccAddressFromBech32(msg.GetBorrower())
	return []sdk.AccAddress{borrower}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg *MsgRepay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
