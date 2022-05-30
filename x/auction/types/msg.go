package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgPlaceSurplusBidRequest)(nil)
)

func NewMsgPlaceSurplusBid(from sdk.AccAddress, auctionID uint64, amt sdk.Coin, appId, auctionMappingId uint64) *MsgPlaceSurplusBidRequest {
	return &MsgPlaceSurplusBidRequest{
		Bidder:           from.String(),
		AuctionId:        auctionID,
		Amount:           amt,
		AppId:            appId,
		AuctionMappingId: auctionMappingId,
	}
}

func (m *MsgPlaceSurplusBidRequest) ValidateBasic() error {
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

func (m *MsgPlaceSurplusBidRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgPlaceDebtBid(from sdk.AccAddress, auctionID uint64, bid, amt sdk.Coin, appId, auctionMappingId uint64) *MsgPlaceDebtBidRequest {
	return &MsgPlaceDebtBidRequest{
		Bidder:            from.String(),
		AuctionId:         auctionID,
		Bid:               bid,
		ExpectedUserToken: amt,
		AppId:             appId,
		AuctionMappingId:  auctionMappingId,
	}
}

func (m *MsgPlaceDebtBidRequest) ValidateBasic() error {
	if m.AuctionId == 0 {
		return errors.New("auction id cannot be zero")
	}
	_, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "--from address cannot be empty or invalid")
	}
	return nil
}

func (m *MsgPlaceDebtBidRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgPlaceDutchBid(from sdk.AccAddress, auctionID uint64, amt sdk.Coin, max sdk.Dec, appId, auctionMappingId uint64) *MsgPlaceDutchBidRequest {
	return &MsgPlaceDutchBidRequest{
		Bidder:           from.String(),
		AuctionId:        auctionID,
		Amount:           amt,
		Max:              max,
		AppId:            appId,
		AuctionMappingId: auctionMappingId,
	}
}

func (m *MsgPlaceDutchBidRequest) ValidateBasic() error {
	if m.AuctionId == 0 {
		return errors.New("auction id cannot be zero")
	}
	_, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "--from address cannot be empty or invalid")
	}
	return nil
}

func (m *MsgPlaceDutchBidRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
