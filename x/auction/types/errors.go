package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrorInvalidSurplusAuctionID      = errorsmod.Register(ModuleName, 201, "surplus auction does not exist with given id")
	ErrorInvalidBiddingDenom          = errorsmod.Register(ModuleName, 202, "given asset type is not accepted for bidding")
	ErrorLowBidAmount                 = errorsmod.Register(ModuleName, 203, "bidding amount is lower than expected")
	ErrorMaxBidAmount                 = errorsmod.Register(ModuleName, 204, "bidding amount is greater than maximum bidding amount")
	ErrorInvalidAddress               = errorsmod.Register(ModuleName, 205, "invalid source address")
	ErrorInvalidDebtAuctionID         = errorsmod.Register(ModuleName, 206, "debt auction does not exist with given id")
	ErrorInvalidDebtUserExpectedDenom = errorsmod.Register(ModuleName, 207, "given asset type is not accepted for debt auction user expected token")
	ErrorDebtExpectedUserAmount       = errorsmod.Register(ModuleName, 208, "invalid user amount")
	ErrorInvalidDebtMintedDenom       = errorsmod.Register(ModuleName, 209, "given asset type is not accepted for debt auction user mint token")
	ErrorInvalidDutchPrice            = errorsmod.Register(ModuleName, 210, "user max price cannot be less than collateral token price")
	ErrorPrices                       = errorsmod.Register(ModuleName, 211, "unable to get fetches prices for asset from oracle")
	ErrorVaultNotFound                = errorsmod.Register(ModuleName, 212, "vault not found for given id")
	ErrorInvalidLockedVault           = errorsmod.Register(ModuleName, 213, "locked vault not found for given id")
	ErrorUnableToMakeFlagsFalse       = errorsmod.Register(ModuleName, 214, "Unable To Make Flags False after auction closed")
	ErrorUnableToSetNetFees           = errorsmod.Register(ModuleName, 215, "Unable To set net fees collected after auction closed")
	ErrorInvalidPair                  = errorsmod.Register(ModuleName, 216, "pair not found for extended pair id")
	ErrorUnknownMsgType               = errorsmod.Register(ModuleName, 217, "unknown message type")
	ErrorInvalidExtendedPairVault     = errorsmod.Register(ModuleName, 218, "extended pair vault not found for given id")
	ErrorInvalidAuctionParams         = errorsmod.Register(ModuleName, 219, "auction params not found for given app id")
	ErrorInStartDutchAuction          = errorsmod.Register(ModuleName, 220, "error in start dutch auction for locked vault id")
	ErrorAssetRates                   = errorsmod.Register(ModuleName, 221, "error in asset rates")
)
