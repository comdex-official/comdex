package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/auctionsV2 module sentinel errors
var (
	ErrDutchAuctionDisabled = sdkerrors.Register(ModuleName, 701, "Dutch auction not enabled for the app")
	ErrEnglishAuctionDisabled = sdkerrors.Register(ModuleName, 702, "English auction not enabled for the app")
	ErrCannotLeaveDebtLessThanDust= sdkerrors.Register(ModuleName, 703, "You need to leave debt atleast equal to dust value or greater. Try making a full bid, or a smaller bid. Your current bid is just short of the dust value , hence it fails.")
	ErrorPriceNotFound= sdkerrors.Register(ModuleName, 704, "price not found")
)
