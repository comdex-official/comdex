package types

import (
	"fmt"
	"github.com/comdex-official/comdex/x/asset/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
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
	if msg.AssetId == 0 {
		return fmt.Errorf("asset id should not be 0: %d ", msg.AssetId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}
	if msg.PoolId == 0 {
		return fmt.Errorf("pool id should not be 0: %d ", msg.AssetId)
	}
	if msg.AppId == 0 {
		return fmt.Errorf("app id should not be 0: %d ", msg.AppId)
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

	if msg.LendId == 0 {
		return fmt.Errorf("lend id should not be 0: %d ", msg.LendId)
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

	if msg.LendId == 0 {
		return fmt.Errorf("lend id should not be 0: %d ", msg.LendId)
	}
	if msg.PairId == 0 {
		return fmt.Errorf("pair id should not be 0: %d ", msg.PairId)
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

	if msg.BorrowId == 0 {
		return fmt.Errorf("borrower id should not be 0: %d ", msg.BorrowId)
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

func NewMsgFundModuleAccounts(poolID, assetID uint64, lender string, amount sdk.Coin) *MsgFundModuleAccounts {
	return &MsgFundModuleAccounts{
		PoolId:  poolID,
		AssetId: assetID,
		Lender:  lender,
		Amount:  amount,
	}
}

func (msg MsgFundModuleAccounts) Route() string { return ModuleName }
func (msg MsgFundModuleAccounts) Type() string  { return TypeFundModuleAccountRequest }

func (msg *MsgFundModuleAccounts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}
	if msg.PoolId == 0 {
		return fmt.Errorf("pool id should not be 0: %d ", msg.PoolId)
	}
	if msg.AssetId == 0 {
		return fmt.Errorf("asset id should not be 0: %d ", msg.AssetId)
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

	if msg.LendId == 0 {
		return fmt.Errorf("lend id should not be 0: %d ", msg.LendId)
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
	if msg.LendId == 0 {
		return fmt.Errorf("lend id should not be 0: %d ", msg.LendId)
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
	if msg.BorrowId == 0 {
		return fmt.Errorf("borrow id should not be 0: %d ", msg.BorrowId)
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
func (msg MsgDepositBorrow) Type() string  { return TypeDepositBorrowAssetRequest }

func (msg *MsgDepositBorrow) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	if msg.BorrowId == 0 {
		return fmt.Errorf("borrow id should not be 0: %d ", msg.BorrowId)
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
	if msg.BorrowId == 0 {
		return fmt.Errorf("borrow id should not be 0: %d ", msg.BorrowId)
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

	if msg.AssetId == 0 {
		return fmt.Errorf("asset id should not be 0: %d ", msg.AssetId)
	}
	if msg.PoolId == 0 {
		return fmt.Errorf("pool id should not be 0: %d ", msg.PoolId)
	}
	if msg.PairId == 0 {
		return fmt.Errorf("pair id should not be 0: %d ", msg.PairId)
	}
	if msg.AppId == 0 {
		return fmt.Errorf("pair id should not be 0: %d ", msg.AppId)
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

func NewMsgCalculateInterestAndRewards(borrower string) *MsgCalculateInterestAndRewards {
	return &MsgCalculateInterestAndRewards{
		Borrower: borrower,
	}
}

func (msg MsgCalculateInterestAndRewards) Route() string { return ModuleName }
func (msg MsgCalculateInterestAndRewards) Type() string {
	return TypeCalculateInterestAndRewardsRequest
}

func (msg *MsgCalculateInterestAndRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		return err
	}

	return nil
}

func (msg *MsgCalculateInterestAndRewards) GetSigners() []sdk.AccAddress {
	lender, err := sdk.AccAddressFromBech32(msg.GetBorrower())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgCalculateInterestAndRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgFundReserveAccounts(assetID uint64, lender string, amount sdk.Coin) *MsgFundReserveAccounts {
	return &MsgFundReserveAccounts{
		AssetId: assetID,
		Lender:  lender,
		Amount:  amount,
	}
}

func (msg MsgFundReserveAccounts) Route() string { return ModuleName }
func (msg MsgFundReserveAccounts) Type() string  { return TypeFundReserveAccountRequest }

func (msg *MsgFundReserveAccounts) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		return err
	}
	if msg.AssetId == 0 {
		return fmt.Errorf("asset id should not be 0: %d ", msg.AssetId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}

	return nil
}

func (msg *MsgFundReserveAccounts) GetSigners() []sdk.AccAddress {
	lender, err := sdk.AccAddressFromBech32(msg.GetLender())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{lender}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgFundReserveAccounts) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgMsgLimitSupplyCapRequest(
	from sdk.AccAddress,
) *MsgLimitSupplyCapRequest {
	return &MsgLimitSupplyCapRequest{
		From: from.String(),
	}
}

func (m *MsgLimitSupplyCapRequest) Route() string {
	return RouterKey
}

func (m *MsgLimitSupplyCapRequest) Type() string {
	return TypeMsgLimitSupplyCapRequest
}

func (m *MsgLimitSupplyCapRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(types.ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(types.ErrorInvalidFrom, "%s", err)
	}

	return nil
}

func (m *MsgLimitSupplyCapRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgLimitSupplyCapRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
