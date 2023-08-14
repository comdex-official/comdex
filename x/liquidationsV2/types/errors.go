package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/liquidationsV2 module sentinel errors
var (
	ErrVaultIDInvalid          = sdkerrors.Register(ModuleName, 1501, "Vault Id invalid")
	ErrorUnknownMsgType        = sdkerrors.Register(ModuleName, 1502, "Unknown msg type")
	ErrorUnknownProposalType   = sdkerrors.Register(ModuleName, 1503, "unknown proposal type")
	ErrorInvalidAppOrAssetData = sdkerrors.Register(ModuleName, 1504, "Invalid data of app , or asset has not been added to the app , or low funds")
	ErrEnglishAuctionDisabled  = sdkerrors.Register(ModuleName, 1505, "English auction not enabled for the app")
)
