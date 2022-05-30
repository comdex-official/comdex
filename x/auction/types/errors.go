package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/auction module sentinel errors
var (
	ErrorInvalidSurplusAuctionId            = sdkerrors.Register(ModuleName, 101, "surplus auction does not exist with given id")
	ErrorInvalidBiddingDenom                = sdkerrors.Register(ModuleName, 102, "given asset type is not accepted for bidding")
	ErrorLowBidAmount                       = sdkerrors.Register(ModuleName, 103, "bidding amount is lower than expected")
	ErrorMaxBidAmount                       = sdkerrors.Register(ModuleName, 104, "bidding amount is greater than maximum bidding amount")
	ErrorBidAlreadyExists                   = sdkerrors.Register(ModuleName, 105, "bid with given amount already placed, Please try with higher bid")
	ErrorInvalidAuctioningCollateral        = sdkerrors.Register(ModuleName, 106, "collateral to be auctioned <= 0")
	ErrorInvalidAmountInAddress             = sdkerrors.Register(ModuleName, 107, "there is not sufficient balance in given address for a given denom")
	ErrorInvalidAddress                     = sdkerrors.Register(ModuleName, 108, "invalid source address")
	ErrorInvalidDebtAuctionId               = sdkerrors.Register(ModuleName, 109, "debt auction does not exist with given id")
	ErrorInvalidDebtUserExpectedDenom       = sdkerrors.Register(ModuleName, 110, "given asset type is not accepted for debt auction user expected token")
	ErrorDebtMoreBidAmount                  = sdkerrors.Register(ModuleName, 111, "can not bid more minted amount")
	ErrorDebtExpectedUserAmount             = sdkerrors.Register(ModuleName, 112, "invalid user amount")
	ErrorInvalidDebtMintedDenom             = sdkerrors.Register(ModuleName, 113, "given asset type is not accepted for debt auction user mint token")
	ErrorInvalidOutFlowTokenBalance         = sdkerrors.Register(ModuleName, 114, "outflow tokens balance in outflow token address is less than what it claims")
	ErrorInvalidDutchAuctionId              = sdkerrors.Register(ModuleName, 115, "dutch auction does not exist with given id")
	ErrorInvalidDutchUserbidDenom           = sdkerrors.Register(ModuleName, 116, "given asset type is not accepted for dutch auction user bid denom")
	ErrorDutchMoreBidAmount                 = sdkerrors.Register(ModuleName, 117, "can not bid more amount")
	ErrorDutchinsufficientUserBalance       = sdkerrors.Register(ModuleName, 118, "user doesnt have balance to buy chost also")
	ErrorInvalidDutchPrice                  = sdkerrors.Register(ModuleName, 119, "user max price cannot be less than collateral token price")
	ErrorInvalidBidId                       = sdkerrors.Register(ModuleName, 120, "invalid bid id")
	ErrorInvalidBurn                        = sdkerrors.Register(ModuleName, 121, "invalid burn")
	ErrorPrices                             = sdkerrors.Register(ModuleName, 122, "unable to get fetches prices for asset from oracle")
	ErrorVaultNotFound                      = sdkerrors.Register(ModuleName, 123, "vault not found for given id")
	ErrorInvalidCollectorAuctionLookupTable = sdkerrors.Register(ModuleName, 124, "CollectorAuctionLookupTable not found for appId")
	ErrorInvalidLockedVault                 = sdkerrors.Register(ModuleName, 125, "locked vault not found for given id")
	ErrorInvalidAppIdAssetId                = sdkerrors.Register(ModuleName, 126, "invalid appId assetId")
	ErrorUnableToMakeFlagsFalse             = sdkerrors.Register(ModuleName, 127, "Unable To Make Flags False after auction closed")
	ErrorUnableToSetNetfees                 = sdkerrors.Register(ModuleName, 128, "Unable To set net fees collected after auction closed")
)

var (
	ErrorUnknownMsgType = sdkerrors.Register(ModuleName, 301, "unknown message type")
)
