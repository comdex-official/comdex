package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgLend(lender string, assetID uint64, amount sdk.Coin, poolID, appID uint64) *MsgLend {
	return &MsgLend{
		Lender:  lender,
		AssetId: assetID,
		Amount:  amount,
		PoolId:  poolID,
		AppId:   appID,
	}
}

func (msg MsgLend) Route() string { return ModuleName }
func (msg MsgLend) Type() string  { return TypeLendAssetRequest }

func (msg *MsgLend) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}
	if msg.AssetId <= 0 {
		return fmt.Errorf("asset id should be positive: %d > 0", msg.AssetId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}
	if msg.PoolId <= 0 {
		return fmt.Errorf("pool id should be positive: %d > 0", msg.AssetId)
	}
	if msg.AppId <= 0 {
		return fmt.Errorf("app id should be positive: %d > 0", msg.AppId)
	}

	return nil
}

func (msg *MsgLend) GetSigners() []sdk.AccAddress {
	lender, _ := sdk.AccAddressFromBech32(msg.GetLender())
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgLend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgWithdraw(lender string, lendID uint64, amount sdk.Coin) *MsgWithdraw {
	return &MsgWithdraw{
		Lender: lender,
		LendId: lendID,
		Amount: amount,
	}
}

func (msg MsgWithdraw) Route() string { return ModuleName }
func (msg MsgWithdraw) Type() string  { return TypeWithdrawAssetRequest }

func (msg *MsgWithdraw) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}

	if msg.LendId <= 0 {
		return fmt.Errorf("lend id should be positive: %d > 0", msg.LendId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}

	return nil
}

func (msg *MsgWithdraw) GetSigners() []sdk.AccAddress {
	lender, _ := sdk.AccAddressFromBech32(msg.GetLender())
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgBorrow(borrower string, lendID, pairID uint64, isStableBorrow bool, amountIn, amountOut sdk.Coin) *MsgBorrow {
	return &MsgBorrow{
		Borrower:       borrower,
		LendId:         lendID,
		PairId:         pairID,
		IsStableBorrow: isStableBorrow,
		AmountIn:       amountIn,
		AmountOut:      amountOut,
	}
}

func (msg MsgBorrow) Route() string { return ModuleName }
func (msg MsgBorrow) Type() string  { return TypeBorrowAssetRequest }

func (msg *MsgBorrow) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	if msg.LendId <= 0 {
		return fmt.Errorf("lend id should be positive: %d > 0", msg.LendId)
	}
	if msg.PairId <= 0 {
		return fmt.Errorf("pair id should be positive: %d > 0", msg.PairId)
	}
	if msg.AmountIn.Amount.IsNegative() || msg.AmountIn.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.AmountIn.Amount)
	}
	if msg.AmountOut.Amount.IsNegative() || msg.AmountOut.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.AmountOut.Amount)
	}

	return nil
}

func (msg *MsgBorrow) GetSigners() []sdk.AccAddress {
	borrower, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{borrower}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgBorrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgRepay(borrower string, borrowID uint64, amount sdk.Coin) *MsgRepay {
	return &MsgRepay{
		Borrower: borrower,
		BorrowId: borrowID,
		Amount:   amount,
	}
}

func (msg MsgRepay) Route() string { return ModuleName }
func (msg MsgRepay) Type() string  { return TypeRepayAssetRequest }

func (msg *MsgRepay) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	if msg.BorrowId <= 0 {
		return fmt.Errorf("borrower id should be positive: %d > 0", msg.BorrowId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}

	return nil
}

func (msg *MsgRepay) GetSigners() []sdk.AccAddress {
	borrower, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{borrower}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgRepay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgFundModuleAccounts(moduleName string, assetID uint64, lender string, amount sdk.Coin) *MsgFundModuleAccounts {
	return &MsgFundModuleAccounts{
		ModuleName: moduleName,
		AssetId:    assetID,
		Lender:     lender,
		Amount:     amount,
	}
}

func (msg MsgFundModuleAccounts) Route() string { return ModuleName }
func (msg MsgFundModuleAccounts) Type() string  { return TypeFundModuleAccountRequest }

func (msg *MsgFundModuleAccounts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}
	if msg.ModuleName == "" {
		return fmt.Errorf("module name can not be empty")
	}

	if msg.AssetId <= 0 {
		return fmt.Errorf("asset id should be positive: %d > 0", msg.AssetId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}

	return nil
}

func (msg *MsgFundModuleAccounts) GetSigners() []sdk.AccAddress {
	lender, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgFundModuleAccounts) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgDeposit(lender string, lendID uint64, amount sdk.Coin) *MsgDeposit {
	return &MsgDeposit{
		Lender: lender,
		LendId: lendID,
		Amount: amount,
	}
}

func (msg MsgDeposit) Route() string { return ModuleName }
func (msg MsgDeposit) Type() string  { return TypeDepositAssetRequest }

func (msg *MsgDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}

	if msg.LendId <= 0 {
		return fmt.Errorf("lend id should be positive: %d > 0", msg.LendId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}

	return nil
}

func (msg *MsgDeposit) GetSigners() []sdk.AccAddress {
	lender, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgCloseLend(lender string, lendID uint64) *MsgCloseLend {
	return &MsgCloseLend{
		Lender: lender,
		LendId: lendID,
	}
}

func (msg MsgCloseLend) Route() string { return ModuleName }
func (msg MsgCloseLend) Type() string  { return TypeCloseLendAssetRequest }

func (msg *MsgCloseLend) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}
	if msg.LendId <= 0 {
		return fmt.Errorf("lend id should be positive: %d > 0", msg.LendId)
	}

	return nil
}

func (msg *MsgCloseLend) GetSigners() []sdk.AccAddress {
	lender, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgCloseLend) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgDraw(borrower string, borrowID uint64, amount sdk.Coin) *MsgDraw {
	return &MsgDraw{
		Borrower: borrower,
		BorrowId: borrowID,
		Amount:   amount,
	}
}

func (msg MsgDraw) Route() string { return ModuleName }
func (msg MsgDraw) Type() string  { return TypeDrawAssetRequest }

func (msg *MsgDraw) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}
	if msg.BorrowId <= 0 {
		return fmt.Errorf("borrow id should be positive: %d > 0", msg.BorrowId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}
	return nil
}

func (msg *MsgDraw) GetSigners() []sdk.AccAddress {
	lender, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgDraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgDepositBorrow(borrower string, borrowID uint64, amount sdk.Coin) *MsgDepositBorrow {
	return &MsgDepositBorrow{
		Borrower: borrower,
		BorrowId: borrowID,
		Amount:   amount,
	}
}

func (msg MsgDepositBorrow) Route() string { return ModuleName }
func (msg MsgDepositBorrow) Type() string  { return TypeDepositBorrowdAssetRequest }

func (msg *MsgDepositBorrow) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	if msg.BorrowId <= 0 {
		return fmt.Errorf("borrow id should be positive: %d > 0", msg.BorrowId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}
	return nil
}

func (msg *MsgDepositBorrow) GetSigners() []sdk.AccAddress {
	lender, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgDepositBorrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgCloseBorrow(borrower string, borrowID uint64) *MsgCloseBorrow {
	return &MsgCloseBorrow{
		Borrower: borrower,
		BorrowId: borrowID,
	}
}

func (msg MsgCloseBorrow) Route() string { return ModuleName }
func (msg MsgCloseBorrow) Type() string  { return TypeCloseBorrowAssetRequest }

func (msg *MsgCloseBorrow) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}
	if msg.BorrowId <= 0 {
		return fmt.Errorf("borrow id should be positive: %d > 0", msg.BorrowId)
	}

	return nil
}

func (msg *MsgCloseBorrow) GetSigners() []sdk.AccAddress {
	lender, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgCloseBorrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgBorrowAlternate(lender string, assetID, poolID uint64, amountIn sdk.Coin, pairID uint64, stableBorrow bool, amountOut sdk.Coin, appID uint64) *MsgBorrowAlternate {
	return &MsgBorrowAlternate{
		Lender:         lender,
		AssetId:        assetID,
		PoolId:         poolID,
		AmountIn:       amountIn,
		PairId:         pairID,
		IsStableBorrow: stableBorrow,
		AmountOut:      amountOut,
		AppId:          appID,
	}
}

func (msg MsgBorrowAlternate) Route() string { return ModuleName }
func (msg MsgBorrowAlternate) Type() string  { return TypeBorrowAlternateAssetRequest }

func (msg *MsgBorrowAlternate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}

	if msg.AssetId <= 0 {
		return fmt.Errorf("asset id should be positive: %d > 0", msg.AssetId)
	}
	if msg.PoolId <= 0 {
		return fmt.Errorf("pool id should be positive: %d > 0", msg.PoolId)
	}
	if msg.PairId <= 0 {
		return fmt.Errorf("pair id should be positive: %d > 0", msg.PairId)
	}
	if msg.AppId <= 0 {
		return fmt.Errorf("pair id should be positive: %d > 0", msg.AppId)
	}
	if msg.AmountIn.Amount.IsNegative() || msg.AmountIn.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.AmountIn.Amount)
	}
	if msg.AmountOut.Amount.IsNegative() || msg.AmountOut.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.AmountOut.Amount)
	}

	return nil
}

func (msg *MsgBorrowAlternate) GetSigners() []sdk.AccAddress {
	lender, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgBorrowAlternate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgCalculateBorrowInterest(borrower string, borrowID uint64) *MsgCalculateBorrowInterest {
	return &MsgCalculateBorrowInterest{
		Borrower: borrower,
		BorrowId: borrowID,
	}
}

func (msg MsgCalculateBorrowInterest) Route() string { return ModuleName }
func (msg MsgCalculateBorrowInterest) Type() string  { return TypeCalculateBorrowInterestRequest }

func (msg *MsgCalculateBorrowInterest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	if msg.BorrowId <= 0 {
		return fmt.Errorf("borrow id should be positive: %d > 0", msg.BorrowId)
	}
	return nil
}

func (msg *MsgCalculateBorrowInterest) GetSigners() []sdk.AccAddress {
	lender, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgCalculateBorrowInterest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
