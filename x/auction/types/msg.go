package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = (*MsgPlaceSurplusBidRequest)(nil)

func NewMsgPlaceSurplusBid(from string, auctionID uint64, amt sdk.Coin, appID, auctionMappingID uint64) *MsgPlaceSurplusBidRequest {
	return &MsgPlaceSurplusBidRequest{
		Bidder:           from,
		AuctionId:        auctionID,
		Amount:           amt,
		AppId:            appID,
		AuctionMappingId: auctionMappingID,
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

func NewMsgPlaceDebtBid(from string, auctionID uint64, bid, amt sdk.Coin, appID, auctionMappingID uint64) *MsgPlaceDebtBidRequest {
	return &MsgPlaceDebtBidRequest{
		Bidder:            from,
		AuctionId:         auctionID,
		Bid:               bid,
		ExpectedUserToken: amt,
		AppId:             appID,
		AuctionMappingId:  auctionMappingID,
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

func NewMsgPlaceDutchBid(from string, auctionID uint64, amt sdk.Coin, max sdk.Dec, appID, auctionMappingID uint64) *MsgPlaceDutchBidRequest {
	return &MsgPlaceDutchBidRequest{
		Bidder:           from,
		AuctionId:        auctionID,
		Amount:           amt,
		Max:              max,
		AppId:            appID,
		AuctionMappingId: auctionMappingID,
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
