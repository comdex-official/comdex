package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgLend(lender sdk.AccAddress, pairID uint64, amount sdk.Coin) *MsgLend {
	return &MsgLend{
		Lender: lender.String(),
		PairId: pairID,
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

func NewMsgFundModuleAccounts(moduleName string, lender sdk.AccAddress, amount sdk.Coin) *MsgFundModuleAccounts {
	return &MsgFundModuleAccounts{
		ModuleName: moduleName,
		Lender:     lender.String(),
		Amount:     amount,
	}
}

func (msg MsgFundModuleAccounts) Route() string { return ModuleName }
func (msg MsgFundModuleAccounts) Type() string  { return EventTypeLoanAsset }

func (msg *MsgFundModuleAccounts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}

	if asset := msg.GetAmount(); !asset.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAsset, asset.String())
	}

	return nil
}

func (msg *MsgFundModuleAccounts) GetSigners() []sdk.AccAddress {
	lender, _ := sdk.AccAddressFromBech32(msg.GetLender())
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg *MsgFundModuleAccounts) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
