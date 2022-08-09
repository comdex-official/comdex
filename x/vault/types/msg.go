package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateRequest)(nil)
	_ sdk.Msg = (*MsgDepositRequest)(nil)
	_ sdk.Msg = (*MsgWithdrawRequest)(nil)
	_ sdk.Msg = (*MsgDrawRequest)(nil)
	_ sdk.Msg = (*MsgRepayRequest)(nil)
	_ sdk.Msg = (*MsgCloseRequest)(nil)
	_ sdk.Msg = (*MsgCreateStableMintRequest)(nil)
	_ sdk.Msg = (*MsgDepositStableMintRequest)(nil)
	_ sdk.Msg = (*MsgWithdrawStableMintRequest)(nil)
)

func NewMsgCreateRequest(
	// nolint
	from sdk.AccAddress,
	appID uint64, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) *MsgCreateRequest {
	return &MsgCreateRequest{
		From:                from.String(),
		AppId:               appID,
		ExtendedPairVaultId: extendedPairVaultID,
		AmountIn:            amountIn,
		AmountOut:           amountOut,
	}
}

func (m *MsgCreateRequest) Route() string {
	return RouterKey
}

func (m *MsgCreateRequest) Type() string {
	return TypeMsgCreateRequest
}

func (m *MsgCreateRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.AmountIn.IsNil() {
		return errors.Wrap(ErrorInvalidAmountIn, "amount_in cannot be nil")
	}
	if m.AmountIn.IsNegative() {
		return errors.Wrap(ErrorInvalidAmountIn, "amount_in cannot be negative")
	}
	if m.AmountIn.IsZero() {
		return errors.Wrap(ErrorInvalidAmountIn, "amount_in cannot be zero")
	}
	if m.AmountOut.IsNil() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be nil")
	}
	if m.AmountOut.IsNegative() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be negative")
	}
	if m.AmountOut.IsZero() {
		return errors.Wrap(ErrorInvalidAmountOut, "amount_out cannot be zero")
	}

	return nil
}

func (m *MsgCreateRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgCreateRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgDepositRequest(
	// nolint
	from sdk.AccAddress,
	appID uint64, extendedPairVaultID uint64, userVaultid uint64, amount sdk.Int) *MsgDepositRequest {
	return &MsgDepositRequest{
		From:                from.String(),
		AppId:               appID,
		ExtendedPairVaultId: extendedPairVaultID,
		UserVaultId:         userVaultid,
		Amount:              amount,
	}
}

func (m *MsgDepositRequest) Route() string {
	return RouterKey
}

func (m *MsgDepositRequest) Type() string {
	return TypeMsgDepositRequest
}

func (m *MsgDepositRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.UserVaultId == 0 {
		return errors.Wrap(ErrorInvalidID, "id cannot be null")
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be zero")
	}

	return nil
}

func (m *MsgDepositRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgDepositRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgWithdrawRequest(
	// nolint
	from sdk.AccAddress,
	appID uint64, extendedPairVaultID uint64, userVaultid uint64, amount sdk.Int) *MsgWithdrawRequest {
	return &MsgWithdrawRequest{
		From:                from.String(),
		AppId:               appID,
		ExtendedPairVaultId: extendedPairVaultID,
		UserVaultId:         userVaultid,
		Amount:              amount,
	}
}

func (m *MsgWithdrawRequest) Route() string {
	return RouterKey
}

func (m *MsgWithdrawRequest) Type() string {
	return TypeMsgWithdrawRequest
}

func (m *MsgWithdrawRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.UserVaultId == 0 {
		return errors.Wrap(ErrorInvalidID, "id cannot be zero")
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be zero")
	}

	return nil
}

func (m *MsgWithdrawRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgWithdrawRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgDrawRequest(
	// nolint
	from sdk.AccAddress,
	appID uint64, extendedPairVaultID uint64, userVaultid uint64, amount sdk.Int) *MsgDrawRequest {
	return &MsgDrawRequest{
		From:                from.String(),
		AppId:               appID,
		ExtendedPairVaultId: extendedPairVaultID,
		UserVaultId:         userVaultid,
		Amount:              amount,
	}
}

func (m *MsgDrawRequest) Route() string {
	return RouterKey
}

func (m *MsgDrawRequest) Type() string {
	return TypeMsgDrawRequest
}

func (m *MsgDrawRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.UserVaultId == 0 {
		return errors.Wrap(ErrorInvalidID, "id cannot be null")
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be zero")
	}

	return nil
}

func (m *MsgDrawRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgDrawRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgRepayRequest(
	// nolint
	from sdk.AccAddress,
	appID uint64, extendedPairVaultID uint64, userVaultid uint64, amount sdk.Int) *MsgRepayRequest {
	return &MsgRepayRequest{
		From:                from.String(),
		AppId:               appID,
		ExtendedPairVaultId: extendedPairVaultID,
		UserVaultId:         userVaultid,
		Amount:              amount,
	}
}

func (m *MsgRepayRequest) Route() string {
	return RouterKey
}

func (m *MsgRepayRequest) Type() string {
	return TypeMsgRepayRequest
}

func (m *MsgRepayRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.UserVaultId == 0 {
		return errors.Wrap(ErrorInvalidID, "id cannot be null")
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be zero")
	}

	return nil
}

func (m *MsgRepayRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgRepayRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgLiquidateRequest(
	// nolint
	from sdk.AccAddress,
	appID uint64, extendedPairVaultID uint64, userVaultid uint64) *MsgCloseRequest {
	return &MsgCloseRequest{
		From:                from.String(),
		AppId:               appID,
		ExtendedPairVaultId: extendedPairVaultID,
		UserVaultId:         userVaultid,
	}
}

func (m *MsgCloseRequest) Route() string {
	return RouterKey
}

func (m *MsgCloseRequest) Type() string {
	return TypeMsgLiquidateRequest
}

func (m *MsgCloseRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.UserVaultId == 0 {
		return errors.Wrap(ErrorInvalidID, "id cannot be null")
	}

	return nil
}

func (m *MsgCloseRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgCloseRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgCreateStableMintRequest(
	// nolint
	from sdk.AccAddress,
	appID uint64, extendedPairVaultID uint64, amount sdk.Int) *MsgCreateStableMintRequest {
	return &MsgCreateStableMintRequest{
		From:                from.String(),
		AppId:               appID,
		ExtendedPairVaultId: extendedPairVaultID,
		Amount:              amount,
	}
}

func (m *MsgCreateStableMintRequest) Route() string {
	return RouterKey
}

func (m *MsgCreateStableMintRequest) Type() string {
	return TypeMsgCreateStableMintRequest
}

func (m *MsgCreateStableMintRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be zero")
	}

	return nil
}

func (m *MsgCreateStableMintRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgCreateStableMintRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgDepositStableMintRequest(
	// nolint
	from sdk.AccAddress,
	appID uint64, extendedPairVaultID uint64, amount sdk.Int, stablemintID uint64) *MsgDepositStableMintRequest {
	return &MsgDepositStableMintRequest{
		From:                from.String(),
		AppId:               appID,
		ExtendedPairVaultId: extendedPairVaultID,
		Amount:              amount,
		StableVaultId:       stablemintID,
	}
}

func (m *MsgDepositStableMintRequest) Route() string {
	return RouterKey
}

func (m *MsgDepositStableMintRequest) Type() string {
	return TypeMsgDepositStableMintRequest
}

func (m *MsgDepositStableMintRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be zero")
	}

	return nil
}

func (m *MsgDepositStableMintRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgDepositStableMintRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgWithdrawStableMintRequest(
	// nolint
	from sdk.AccAddress,
	appID uint64, extendedPairVaultID uint64, amount sdk.Int, stablemintID uint64) *MsgWithdrawStableMintRequest {
	return &MsgWithdrawStableMintRequest{
		From:                from.String(),
		AppId:               appID,
		ExtendedPairVaultId: extendedPairVaultID,
		Amount:              amount,
		StableVaultId:       stablemintID,
	}
}

func (m *MsgWithdrawStableMintRequest) Route() string {
	return RouterKey
}

func (m *MsgWithdrawStableMintRequest) Type() string {
	return TypeMsgWithdrawStableMintRequest
}

func (m *MsgWithdrawStableMintRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.Amount.IsNil() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be nil")
	}
	if m.Amount.IsNegative() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be negative")
	}
	if m.Amount.IsZero() {
		return errors.Wrap(ErrorInvalidAmount, "amount cannot be zero")
	}

	return nil
}

func (m *MsgWithdrawStableMintRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgWithdrawStableMintRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
