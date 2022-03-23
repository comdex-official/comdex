package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/auction module sentinel errors
var (
	ErrorInvalidAuctionId    = sdkerrors.Register(ModuleName, 101, "Auction does not exist with given id")
	ErrorInvalidBiddingDenom = sdkerrors.Register(ModuleName, 102, "Given asset type is not accepted for bidding")
	ErrorLowBidAmount        = sdkerrors.Register(ModuleName, 103, "Bidding amount is lower than expected")
	ErrorMaxBidAmount        = sdkerrors.Register(ModuleName, 104, "Bidding amount is greater than maximum bidding amount")
	ErrorBidAlreadyExists    = sdkerrors.Register(ModuleName, 105, "Bid with given amount already placed, Please try with higher bid")
)

var (
	ErrorUnknownMsgType = sdkerrors.Register(ModuleName, 301, "unknown message type")
)
