package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/auctions module sentinel errors
var (
	ErrDutchAuctionDisabled             = sdkerrors.Register(ModuleName, 701, "Dutch auction not enabled for the app")
	ErrEnglishAuctionDisabled           = sdkerrors.Register(ModuleName, 702, "English auction not enabled for the app")
	ErrCannotLeaveDebtLessThanDust      = sdkerrors.Register(ModuleName, 703, "You need to leave debt atleast equal to dust value or greater. Try making a full bid, or a smaller bid. Your current bid is just short of the dust value , hence it fails.")
	ErrorPriceNotFound                  = sdkerrors.Register(ModuleName, 704, "price not found")
	ErrBidCannotBeZero                  = sdkerrors.Register(ModuleName, 705, "Bid amount can't be Zero")
	ErrorLowBidAmount                   = sdkerrors.Register(ModuleName, 706, "bidding amount is lower than expected")
	ErrorMaxBidAmount                   = sdkerrors.Register(ModuleName, 707, "bidding amount is greater than maximum bidding amount")
	ErrLiquidationNotFound              = sdkerrors.Register(ModuleName, 708, "Liquidation data not found for the auction")
	ErrBidNotFound                      = sdkerrors.Register(ModuleName, 709, "There exists no active bid for the user with given params")
	ErrAuctionParamsNotFound            = sdkerrors.Register(ModuleName, 710, "There exists no auction params")
	ErrorUnknownProposalType            = sdkerrors.Register(ModuleName, 711, "unknown proposal type")
	ErrorUnknownDebtToken               = sdkerrors.Register(ModuleName, 712, "Bid token is not the debt token")
	ErrorDiscountGreaterThanMaxDiscount = sdkerrors.Register(ModuleName, 713, "Premium discount entered is greater than max discount")
	ErrAuctionLookupTableNotFound       = sdkerrors.Register(ModuleName, 714, "auctionLookupTable not found")
	ErrorUnableToSetNetFees             = sdkerrors.Register(ModuleName, 715, "Unable To set net fees collected after auction closed")
	ErrorInGettingLockedVault           = sdkerrors.Register(ModuleName, 716, "error in bid dutch auction - locked vault not found")
	ErrorInsufficientReserveBalance     = sdkerrors.Register(ModuleName, 717, "Insufficient Reserve Balance for this transaction")
)
