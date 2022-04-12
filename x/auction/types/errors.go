package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/auction module sentinel errors
var (
	ErrorInvalidAuctionId            = sdkerrors.Register(ModuleName, 101, "auction does not exist with given id")
	ErrorInvalidBiddingDenom         = sdkerrors.Register(ModuleName, 102, "given asset type is not accepted for bidding")
	ErrorLowBidAmount                = sdkerrors.Register(ModuleName, 103, "bidding amount is lower than expected")
	ErrorMaxBidAmount                = sdkerrors.Register(ModuleName, 104, "bidding amount is greater than maximum bidding amount")
	ErrorBidAlreadyExists            = sdkerrors.Register(ModuleName, 105, "bid with given amount already placed, Please try with higher bid")
	ErrorInvalidAuctioningCollateral = sdkerrors.Register(ModuleName, 106, "collateral to be auctioned <= 0")
)

var (
	ErrorUnknownMsgType = sdkerrors.Register(ModuleName, 301, "unknown message type")
)
