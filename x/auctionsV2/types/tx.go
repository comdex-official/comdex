package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgPlaceMarketBid(bidder string, auctionId uint64, amount sdk.Coin) *MsgPlaceMarketBidRequest {
	return &MsgPlaceMarketBidRequest{
		AuctionId: auctionId,
		Bidder:    bidder,
		Amount:    amount,
	}
}

func (msg MsgPlaceMarketBidRequest) Route() string { return ModuleName }
func (msg MsgPlaceMarketBidRequest) Type() string  { return TypePlaceMarketBidRequest }

func (msg *MsgPlaceMarketBidRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return err
	}
	if msg.AuctionId == 0 {
		return fmt.Errorf("Auction id should not be 0: %d ", msg.AuctionId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}
	return nil
}

func (msg *MsgPlaceMarketBidRequest) GetSigners() []sdk.AccAddress {
	Bidder, _ := sdk.AccAddressFromBech32(msg.Bidder)
	return []sdk.AccAddress{Bidder}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgPlaceMarketBidRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgDepositLimitBid(bidder string, collateralTokenId, debtTokenId uint64, premiumDiscount sdkmath.Int, amount sdk.Coin) *MsgDepositLimitBidRequest {
	return &MsgDepositLimitBidRequest{
		CollateralTokenId: collateralTokenId,
		DebtTokenId:       debtTokenId,
		PremiumDiscount:   premiumDiscount,
		Bidder:            bidder,
		Amount:            amount,
	}
}

func (msg MsgDepositLimitBidRequest) Route() string { return ModuleName }
func (msg MsgDepositLimitBidRequest) Type() string  { return TypePlaceLimitBidRequest }

func (msg *MsgDepositLimitBidRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return err
	}
	if msg.CollateralTokenId == 0 || msg.DebtTokenId == 0 {
		return fmt.Errorf("CollateralToken id or DebtToken should not be 0: %d ", msg.CollateralTokenId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}
	return nil
}

func (msg *MsgDepositLimitBidRequest) GetSigners() []sdk.AccAddress {
	Bidder, _ := sdk.AccAddressFromBech32(msg.Bidder)
	return []sdk.AccAddress{Bidder}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgDepositLimitBidRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgCancelLimitBid(bidder string, collateralTokenId, debtTokenId uint64, premiumDiscount sdkmath.Int) *MsgCancelLimitBidRequest {
	return &MsgCancelLimitBidRequest{
		CollateralTokenId: collateralTokenId,
		DebtTokenId:       debtTokenId,
		PremiumDiscount:   premiumDiscount,
		Bidder:            bidder,
	}
}

func (msg MsgCancelLimitBidRequest) Route() string { return ModuleName }
func (msg MsgCancelLimitBidRequest) Type() string  { return TypeCancelLimitBidRequest }

func (msg *MsgCancelLimitBidRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return err
	}
	if msg.CollateralTokenId == 0 || msg.DebtTokenId == 0 {
		return fmt.Errorf("CollateralToken id or DebtToken should not be 0: %d ", msg.CollateralTokenId)
	}
	return nil
}

func (msg *MsgCancelLimitBidRequest) GetSigners() []sdk.AccAddress {
	Bidder, _ := sdk.AccAddressFromBech32(msg.Bidder)
	return []sdk.AccAddress{Bidder}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgCancelLimitBidRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func NewMsgWithdrawLimitBid(bidder string, collateralTokenId, debtTokenId uint64, premiumDiscount sdkmath.Int, amount sdk.Coin) *MsgWithdrawLimitBidRequest {
	return &MsgWithdrawLimitBidRequest{
		CollateralTokenId: collateralTokenId,
		DebtTokenId:       debtTokenId,
		PremiumDiscount:   premiumDiscount,
		Bidder:            bidder,
		Amount:            amount,
	}
}

func (msg MsgWithdrawLimitBidRequest) Route() string { return ModuleName }
func (msg MsgWithdrawLimitBidRequest) Type() string  { return TypeWithdrawLimitBidRequest }

func (msg *MsgWithdrawLimitBidRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return err
	}
	if msg.CollateralTokenId == 0 || msg.DebtTokenId == 0 {
		return fmt.Errorf("CollateralToken id or DebtToken should not be 0: %d ", msg.CollateralTokenId)
	}
	if msg.Amount.Amount.IsNegative() || msg.Amount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %s < 0", msg.Amount.Amount)
	}
	return nil
}

func (msg *MsgWithdrawLimitBidRequest) GetSigners() []sdk.AccAddress {
	Bidder, _ := sdk.AccAddressFromBech32(msg.Bidder)
	return []sdk.AccAddress{Bidder}
}

// GetSignBytes get the bytes for the message signer to sign on.
func (msg *MsgWithdrawLimitBidRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
