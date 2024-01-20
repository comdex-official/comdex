package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/auctions module sentinel errors
var (
	ErrDutchAuctionDisabled             = errorsmod.Register(ModuleName, 701, "Dutch auction not enabled for the app")
	ErrEnglishAuctionDisabled           = errorsmod.Register(ModuleName, 702, "English auction not enabled for the app")
	ErrCannotLeaveDebtLessThanDust      = errorsmod.Register(ModuleName, 703, "You need to leave debt atleast equal to dust value or greater. Try making a full bid, or a smaller bid. Your current bid is just short of the dust value , hence it fails.")
	ErrorPriceNotFound                  = errorsmod.Register(ModuleName, 704, "price not found")
	ErrBidCannotBeZero                  = errorsmod.Register(ModuleName, 705, "Bid amount can't be Zero")
	ErrorLowBidAmount                   = errorsmod.Register(ModuleName, 706, "bidding amount is lower than expected")
	ErrorMaxBidAmount                   = errorsmod.Register(ModuleName, 707, "bidding amount is greater than maximum bidding amount")
	ErrLiquidationNotFound              = errorsmod.Register(ModuleName, 708, "Liquidation data not found for the auction")
	ErrBidNotFound                      = errorsmod.Register(ModuleName, 709, "There exists no active bid for the user with given params")
	ErrAuctionParamsNotFound            = errorsmod.Register(ModuleName, 710, "There exists no auction params")
	ErrorUnknownProposalType            = errorsmod.Register(ModuleName, 711, "unknown proposal type")
	ErrorUnknownDebtToken               = errorsmod.Register(ModuleName, 712, "Bid token is not the debt token")
	ErrorDiscountGreaterThanMaxDiscount = errorsmod.Register(ModuleName, 713, "Premium discount entered is greater than max discount")
	ErrAuctionLookupTableNotFound       = errorsmod.Register(ModuleName, 714, "auctionLookupTable not found")
	ErrorUnableToSetNetFees             = errorsmod.Register(ModuleName, 715, "Unable To set net fees collected after auction closed")
	ErrorInGettingLockedVault           = errorsmod.Register(ModuleName, 716, "error in bid dutch auction - locked vault not found")
)
