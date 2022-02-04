package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgPlaceBidRequest)(nil)
)

func NewMsgPlaceBid(from sdk.AccAddress, auctionID uint64, amt sdk.Coin) *MsgPlaceBidRequest {
	return &MsgPlaceBidRequest{
		Bidder:    from.String(),
		AuctionId: auctionID,
		Amount:    amt,
	}
}

func (m *MsgPlaceBidRequest) ValidateBasic() error {
	if m.AuctionId == 0 {
		return errors.New("auction id cannot be zero")
	}
	_, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "--from address cannot be empty or invalid")
	}
	if !m.Amount.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "bid amount %s", m.Amount)
	}
	return nil
}

func (m *MsgPlaceBidRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
