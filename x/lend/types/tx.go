package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgLend(lender sdk.AccAddress, assetId uint64, amount sdk.Coin, poolId uint64) *MsgLend {
	return &MsgLend{
		Lender:  lender.String(),
		AssetId: assetId,
		Amount:  amount,
		PoolId:  poolId,
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

func NewMsgWithdraw(lender sdk.AccAddress, lendId uint64, amount sdk.Coin) *MsgWithdraw {
	return &MsgWithdraw{
		Lender: lender.String(),
		LendId: lendId,
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

func NewMsgBorrow(borrower sdk.AccAddress, lendId, pairId uint64, isStableBorrow bool, amountIn, amountOut sdk.Coin) *MsgBorrow {
	return &MsgBorrow{
		Borrower:       borrower.String(),
		LendId:         lendId,
		PairId:         pairId,
		IsStableBorrow: isStableBorrow,
		AmountIn:       amountIn,
		AmountOut:      amountOut,
	}
}

func (msg MsgBorrow) Route() string { return ModuleName }
func (msg MsgBorrow) Type() string  { return EventTypeBorrowAsset }

func (msg *MsgBorrow) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	if asset := msg.GetAmountIn(); !asset.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAsset, asset.String())
	}
	if asset := msg.GetAmountOut(); !asset.IsValid() {
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

func NewMsgRepay(borrower sdk.AccAddress, borrowId uint64, amount sdk.Coin) *MsgRepay {
	return &MsgRepay{
		Borrower: borrower.String(),
		BorrowId: borrowId,
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

func NewMsgFundModuleAccounts(moduleName string, AssetId uint64, lender sdk.AccAddress, amount sdk.Coin) *MsgFundModuleAccounts {
	return &MsgFundModuleAccounts{
		ModuleName: moduleName,
		AssetId:    AssetId,
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

func NewMsgDeposit(lender sdk.AccAddress, lendId uint64, amount sdk.Coin) *MsgDeposit {
	return &MsgDeposit{
		Lender: lender.String(),
		LendId: lendId,
		Amount: amount,
	}
}

func (msg MsgDeposit) Route() string { return ModuleName }
func (msg MsgDeposit) Type() string  { return EventTypeWithdrawLoanedAsset }

func (msg *MsgDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}

	if asset := msg.GetAmount(); !asset.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAsset, asset.String())
	}

	return nil
}

func (msg *MsgDeposit) GetSigners() []sdk.AccAddress {
	lender, _ := sdk.AccAddressFromBech32(msg.GetLender())
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg *MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgCloseLend(lender sdk.AccAddress, lendId uint64) *MsgCloseLend {
	return &MsgCloseLend{
		Lender: lender.String(),
		LendId: lendId,
	}
}

func (msg MsgCloseLend) Route() string { return ModuleName }
func (msg MsgCloseLend) Type() string  { return EventTypeWithdrawLoanedAsset }

func (msg *MsgCloseLend) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}
	return nil
}

func (msg *MsgCloseLend) GetSigners() []sdk.AccAddress {
	lender, _ := sdk.AccAddressFromBech32(msg.GetLender())
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg *MsgCloseLend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgDraw(borrower sdk.AccAddress, borrowId uint64, amount sdk.Coin) *MsgDraw {
	return &MsgDraw{
		Borrower: borrower.String(),
		BorrowId: borrowId,
		Amount:   amount,
	}
}

func (msg MsgDraw) Route() string { return ModuleName }
func (msg MsgDraw) Type() string  { return EventTypeWithdrawLoanedAsset }

func (msg *MsgDraw) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	if asset := msg.GetAmount(); !asset.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAsset, asset.String())
	}

	return nil
}

func (msg *MsgDraw) GetSigners() []sdk.AccAddress {
	lender, _ := sdk.AccAddressFromBech32(msg.GetBorrower())
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg *MsgDraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgDepositBorrow(borrower sdk.AccAddress, borrowId uint64, amount sdk.Coin) *MsgDepositBorrow {
	return &MsgDepositBorrow{
		Borrower: borrower.String(),
		BorrowId: borrowId,
		Amount:   amount,
	}
}

func (msg MsgDepositBorrow) Route() string { return ModuleName }
func (msg MsgDepositBorrow) Type() string  { return EventTypeWithdrawLoanedAsset }

func (msg *MsgDepositBorrow) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	if asset := msg.GetAmount(); !asset.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAsset, asset.String())
	}

	return nil
}

func (msg *MsgDepositBorrow) GetSigners() []sdk.AccAddress {
	lender, _ := sdk.AccAddressFromBech32(msg.GetBorrower())
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg *MsgDepositBorrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgCloseBorrow(borrower sdk.AccAddress, borrowId uint64) *MsgCloseBorrow {
	return &MsgCloseBorrow{
		Borrower: borrower.String(),
		BorrowId: borrowId,
	}
}

func (msg MsgCloseBorrow) Route() string { return ModuleName }
func (msg MsgCloseBorrow) Type() string  { return EventTypeWithdrawLoanedAsset }

func (msg *MsgCloseBorrow) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	return nil
}

func (msg *MsgCloseBorrow) GetSigners() []sdk.AccAddress {
	lender, _ := sdk.AccAddressFromBech32(msg.GetBorrower())
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on
func (msg *MsgCloseBorrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
