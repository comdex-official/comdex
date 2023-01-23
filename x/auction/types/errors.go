package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidSurplusAuctionID      = sdkerrors.Register(ModuleName, 201, "surplus auction does not exist with given id")
	ErrorInvalidBiddingDenom          = sdkerrors.Register(ModuleName, 202, "given asset type is not accepted for bidding")
	ErrorLowBidAmount                 = sdkerrors.Register(ModuleName, 203, "bidding amount is lower than expected")
	ErrorMaxBidAmount                 = sdkerrors.Register(ModuleName, 204, "bidding amount is greater than maximum bidding amount")
	ErrorInvalidAddress               = sdkerrors.Register(ModuleName, 205, "invalid source address")
	ErrorInvalidDebtAuctionID         = sdkerrors.Register(ModuleName, 206, "debt auction does not exist with given id")
	ErrorInvalidDebtUserExpectedDenom = sdkerrors.Register(ModuleName, 207, "given asset type is not accepted for debt auction user expected token")
	ErrorDebtExpectedUserAmount       = sdkerrors.Register(ModuleName, 208, "invalid user amount")
	ErrorInvalidDebtMintedDenom       = sdkerrors.Register(ModuleName, 209, "given asset type is not accepted for debt auction user mint token")
	ErrorInvalidDutchPrice            = sdkerrors.Register(ModuleName, 210, "user max price cannot be less than collateral token price")
	ErrorPrices                       = sdkerrors.Register(ModuleName, 211, "unable to get fetches prices for asset from oracle")
	ErrorVaultNotFound                = sdkerrors.Register(ModuleName, 212, "vault not found for given id")
	ErrorInvalidLockedVault           = sdkerrors.Register(ModuleName, 213, "locked vault not found for given id")
	ErrorUnableToMakeFlagsFalse       = sdkerrors.Register(ModuleName, 214, "Unable To Make Flags False after auction closed")
	ErrorUnableToSetNetFees           = sdkerrors.Register(ModuleName, 215, "Unable To set net fees collected after auction closed")
	ErrorInvalidPair                  = sdkerrors.Register(ModuleName, 216, "pair not found for extended pair id")
	ErrorUnknownMsgType               = sdkerrors.Register(ModuleName, 217, "unknown message type")
	ErrorInvalidExtendedPairVault     = sdkerrors.Register(ModuleName, 218, "extended pair vault not found for given id")
	ErrorInvalidAuctionParams         = sdkerrors.Register(ModuleName, 219, "auction params not found for given app id")
	ErrorInStartDutchAuction          = sdkerrors.Register(ModuleName, 220, "error in start dutch auction for locked vault id")
	ErrorAssetRates                   = sdkerrors.Register(ModuleName, 221, "error in asset rates")
)
