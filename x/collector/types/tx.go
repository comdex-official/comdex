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

func NewMsgRefund(addr string) *MsgRefund {
	return &MsgRefund{
		Addr: addr,
	}
}

func (msg MsgRefund) Route() string { return ModuleName }
func (msg MsgRefund) Type() string  { return TypeRefundRequest }

func (msg *MsgRefund) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetAddr())
	if err != nil {
		return err
	}
	return nil
}

func (msg *MsgRefund) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.GetAddr())
	return []sdk.AccAddress{addr}
}

func (msg *MsgRefund) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgUpdateDebtparams(addr string, appID, assetID, slots uint64, debtThreshold, lotSize, debtLotSize sdk.Int, isDebtAuction bool) *MsgUpdateDebtParams {
	return &MsgUpdateDebtParams{
		Addr:          addr,
		AppId:         appID,
		AssetId:       assetID,
		Slots:         slots,
		DebtThreshold: debtThreshold,
		LotSize:       lotSize,
		DebtLotSize:   debtLotSize,
		IsDebtAuction: isDebtAuction,
	}
}

func (msg MsgUpdateDebtParams) Route() string { return ModuleName }
func (msg MsgUpdateDebtParams) Type() string  { return TypeUpdateDebtParamsRequest }

func (msg *MsgUpdateDebtParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.GetAddr())
	if err != nil {
		return err
	}
	if msg.AppId == 0 {
		return fmt.Errorf("app id should not be 0: %d ", msg.AppId)
	}
	if msg.AssetId == 0 {
		return fmt.Errorf("asset id should not be 0: %d ", msg.AssetId)
	}
	if msg.Slots == 0 {
		return fmt.Errorf("slots should not be 0: %d ", msg.Slots)
	}
	if msg.DebtThreshold.IsNegative() {
		return fmt.Errorf("debt threshold should not be negative: %s ", msg.DebtThreshold)
	}
	if msg.LotSize.IsNegative() {
		return fmt.Errorf("lot size should not be negative: %s ", msg.LotSize)
	}
	if msg.DebtLotSize.IsNegative() {
		return fmt.Errorf("debt lot size should not be negative: %s ", msg.DebtLotSize)
	}
	return nil
}

func (msg *MsgUpdateDebtParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.GetAddr())
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpdateDebtParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

