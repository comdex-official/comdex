package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/auctionsV2 module sentinel errors
var (
	ErrDutchAuctionDisabled = sdkerrors.Register(ModuleName, 701, "Dutch auction not enabled for the app")
)
