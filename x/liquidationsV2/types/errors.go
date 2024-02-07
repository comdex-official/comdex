package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/liquidationsV2 module sentinel errors
var (
	ErrVaultIDInvalid          = errorsmod.Register(ModuleName, 1501, "Vault Id invalid")
	ErrorUnknownMsgType        = errorsmod.Register(ModuleName, 1502, "Unknown msg type")
	ErrorUnknownProposalType   = errorsmod.Register(ModuleName, 1503, "unknown proposal type")
	ErrorInvalidAppOrAssetData = errorsmod.Register(ModuleName, 1504, "Invalid data of app , or asset has not been added to the app , or low funds")
	ErrEnglishAuctionDisabled  = errorsmod.Register(ModuleName, 1505, "English auction not enabled for the app")
)
