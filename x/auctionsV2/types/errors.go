package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/auctions module sentinel errors
var (
	ErrDutchAuctionDisabled        = sdkerrors.Register(ModuleName, 701, "Dutch auction not enabled for the app")
	ErrEnglishAuctionDisabled      = sdkerrors.Register(ModuleName, 702, "English auction not enabled for the app")
	ErrCannotLeaveDebtLessThanDust = sdkerrors.Register(ModuleName, 703, "You need to leave debt atleast equal to dust value or greater. Try making a full bid, or a smaller bid. Your current bid is just short of the dust value , hence it fails.")
	ErrorPriceNotFound             = sdkerrors.Register(ModuleName, 704, "price not found")
	ErrBidCannotBeZero             = sdkerrors.Register(ModuleName, 705, "Bid amount can't be Zero")
	ErrorLowBidAmount              = sdkerrors.Register(ModuleName, 706, "bidding amount is lower than expected")
	ErrorMaxBidAmount              = sdkerrors.Register(ModuleName, 707, "bidding amount is greater than maximum bidding amount")
	ErrLiquidationNotFound         = sdkerrors.Register(ModuleName, 708, "Liquidation data not found for the auction")
	ErrBidNotFound                 = sdkerrors.Register(ModuleName, 709, "There exists no active bid for the user with given params")
	ErrAuctionParamsNotFound       = sdkerrors.Register(ModuleName, 710, "There exists no auction params")
	ErrorUnknownProposalType       = sdkerrors.Register(ModuleName, 711, "unknown proposal type")
)
