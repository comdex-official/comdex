package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorInvalidSurplusAuctionID      = sdkerrors.Register(ModuleName, 101, "surplus auction does not exist with given id")
	ErrorInvalidBiddingDenom          = sdkerrors.Register(ModuleName, 102, "given asset type is not accepted for bidding")
	ErrorLowBidAmount                 = sdkerrors.Register(ModuleName, 103, "bidding amount is lower than expected")
	ErrorMaxBidAmount                 = sdkerrors.Register(ModuleName, 104, "bidding amount is greater than maximum bidding amount")
	ErrorInvalidAddress               = sdkerrors.Register(ModuleName, 108, "invalid source address")
	ErrorInvalidDebtAuctionID         = sdkerrors.Register(ModuleName, 109, "debt auction does not exist with given id")
	ErrorInvalidDebtUserExpectedDenom = sdkerrors.Register(ModuleName, 110, "given asset type is not accepted for debt auction user expected token")
	ErrorDebtExpectedUserAmount       = sdkerrors.Register(ModuleName, 112, "invalid user amount")
	ErrorInvalidDebtMintedDenom       = sdkerrors.Register(ModuleName, 113, "given asset type is not accepted for debt auction user mint token")
	ErrorInvalidDutchPrice            = sdkerrors.Register(ModuleName, 119, "user max price cannot be less than collateral token price")

	ErrorPrices                   = sdkerrors.Register(ModuleName, 122, "unable to get fetches prices for asset from oracle")
	ErrorVaultNotFound            = sdkerrors.Register(ModuleName, 123, "vault not found for given id")
	ErrorInvalidLockedVault       = sdkerrors.Register(ModuleName, 125, "locked vault not found for given id")
	ErrorUnableToMakeFlagsFalse   = sdkerrors.Register(ModuleName, 127, "Unable To Make Flags False after auction closed")
	ErrorUnableToSetNetFees       = sdkerrors.Register(ModuleName, 128, "Unable To set net fees collected after auction closed")
	ErrorInvalidPair              = sdkerrors.Register(ModuleName, 130, "pair not found for extended pair id")
	ErrorAssetNotFound            = sdkerrors.Register(ModuleName, 131, "asset not found for given id")
	ErrorInvalidExtendedPairVault = sdkerrors.Register(ModuleName, 132, "extended pair vault not found for given id")
	ErrorInvalidAuctionParams     = sdkerrors.Register(ModuleName, 136, "auction params not found for given app id")
	BurnCoinValueInCloseAuctionIsZero           = sdkerrors.Register(ModuleName, 137, "Burn Coin value in close auction is zero")
	SendCoinsFromModuleToModuleInAuctionIsZero  = sdkerrors.Register(ModuleName, 138, "Coin value in module to module transfer in auction is zero")
	SendCoinsFromModuleToAccountInAuctionIsZero = sdkerrors.Register(ModuleName, 139, "Coin value in module to account transfer in auction is zero")
	SendCoinsFromAccountToModuleInAuctionIsZero = sdkerrors.Register(ModuleName, 140, "Coin value in account to module transfer in auction is zero")

)

var (
	ErrorUnknownMsgType = sdkerrors.Register(ModuleName, 301, "unknown message type")
)
