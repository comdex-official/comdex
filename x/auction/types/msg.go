package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgPlaceSurplusBidRequest)(nil)
	_ sdk.Msg = (*MsgPlaceDebtBidRequest)(nil)
	_ sdk.Msg = (*MsgPlaceDutchBidRequest)(nil)
	_ sdk.Msg = (*MsgPlaceDutchLendBidRequest)(nil)
)

const (
	TypeMsgPlaceSurplusBidRequest   = "place_surplus_bid"
	TypeMsgPlaceDebtBidRequest      = "place_debt_bid"
	TypeMsgPlaceDutchBidRequest     = "place_dutch_bid"
	TypeMsgPlaceDutchLendBidRequest = "place_dutch_lend_bid"
)

func NewMsgPlaceSurplusBid(from string, auctionID uint64, amt sdk.Coin, appID, auctionMappingID uint64) *MsgPlaceSurplusBidRequest {
	return &MsgPlaceSurplusBidRequest{
		Bidder:           from,
		AuctionId:        auctionID,
		Amount:           amt,
		AppId:            appID,
		AuctionMappingId: auctionMappingID,
	}
}

func (m MsgPlaceSurplusBidRequest) Route() string { return RouterKey }
func (m MsgPlaceSurplusBidRequest) Type() string  { return TypeMsgPlaceSurplusBidRequest }

func (m MsgPlaceSurplusBidRequest) ValidateBasic() error {
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

func (m MsgPlaceSurplusBidRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgPlaceSurplusBidRequest) GetSigners() []sdk.AccAddress {
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

func (m MsgPlaceDebtBidRequest) Route() string { return RouterKey }
func (m MsgPlaceDebtBidRequest) Type() string  { return TypeMsgPlaceDebtBidRequest }

func (m MsgPlaceDebtBidRequest) ValidateBasic() error {
	if m.AuctionId == 0 {
		return errors.New("auction id cannot be zero")
	}
	_, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "--from address cannot be empty or invalid")
	}
	return nil
}

func (m MsgPlaceDebtBidRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgPlaceDebtBidRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgPlaceDutchBid(from string, auctionID uint64, amt sdk.Coin, appID, auctionMappingID uint64) *MsgPlaceDutchBidRequest {
	return &MsgPlaceDutchBidRequest{
		Bidder:           from,
		AuctionId:        auctionID,
		Amount:           amt,
		AppId:            appID,
		AuctionMappingId: auctionMappingID,
	}
}

func (m MsgPlaceDutchBidRequest) Route() string { return RouterKey }
func (m MsgPlaceDutchBidRequest) Type() string  { return TypeMsgPlaceDutchBidRequest }

func (m MsgPlaceDutchBidRequest) ValidateBasic() error {
	if m.AuctionId == 0 {
		return errors.New("auction id cannot be zero")
	}
	_, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "--from address cannot be empty or invalid")
	}
	return nil
}

func (m MsgPlaceDutchBidRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgPlaceDutchBidRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgPlaceDutchLendBid(from string, auctionID uint64, amt sdk.Coin, appID, auctionMappingID uint64) *MsgPlaceDutchLendBidRequest {
	return &MsgPlaceDutchLendBidRequest{
		Bidder:           from,
		AuctionId:        auctionID,
		Amount:           amt,
		AppId:            appID,
		AuctionMappingId: auctionMappingID,
	}
}

func (m MsgPlaceDutchLendBidRequest) Route() string { return RouterKey }
func (m MsgPlaceDutchLendBidRequest) Type() string  { return TypeMsgPlaceDutchLendBidRequest }

func (m MsgPlaceDutchLendBidRequest) ValidateBasic() error {
	if m.AuctionId == 0 {
		return errors.New("auction id cannot be zero")
	}
	_, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "--from address cannot be empty or invalid")
	}
	return nil
}

func (m MsgPlaceDutchLendBidRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgPlaceDutchLendBidRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Bidder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
