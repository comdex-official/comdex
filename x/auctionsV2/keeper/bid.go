package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) PlaceDutchAuctionBid(ctx sdk.Context, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin, auctionData types.Auctions) error {

	return nil
}

func (k Keeper) PlaceEnglishAuctionBid(ctx sdk.Context, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin, auctionData types.Auctions) error {
	return nil
}

func (k Keeper) SetLimitAuctionBidID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LimitAuctionBidIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetLimitAuctionBidID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LimitAuctionBidIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetUserLimitBidData(ctx sdk.Context, userLimitBidData types.LimitOrderBid, debtTokenID, collateralTokenID uint64, premium string) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKey(debtTokenID, collateralTokenID, premium, userLimitBidData.BidderAddress)
		value = k.cdc.MustMarshal(&userLimitBidData)
	)

	store.Set(key, value)
}

func (k Keeper) GetUserLimitBidData(ctx sdk.Context, debtTokenID, collateralTokenID uint64, premium, address string) (mappingData types.LimitOrderBid, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKey(debtTokenID, collateralTokenID, premium, address)
		value = store.Get(key)
	)

	if value == nil {
		return mappingData, false
	}

	k.cdc.MustUnmarshal(value, &mappingData)
	return mappingData, true
}

func (k Keeper) DeleteUserLimitBidData(ctx sdk.Context, debtTokenID, collateralTokenID uint64, premium, address string) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKey(debtTokenID, collateralTokenID, premium, address)
	)

	store.Delete(key)
}

func (k Keeper) DepositLimitAuctionBid(ctx sdk.Context, bidder string, CollateralTokenId, DebtTokenId uint64, PremiumDiscount string, amount sdk.Coin) error {
	id := k.GetLimitAuctionBidID(ctx)
	bidderAddr, err := sdk.AccAddressFromBech32(bidder)
	if err != nil {
		return nil
	}
	premiumDiscount, err := sdk.NewDecFromStr(PremiumDiscount)
	if err != nil {
		return err
	}
	collateralAsset, found := k.asset.GetAsset(ctx, CollateralTokenId)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	collateralAssetToken := sdk.NewCoin(collateralAsset.Denom, sdk.NewInt(0))
	userLimitBid, found := k.GetUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)
	if !found {
		userLimitBid = types.LimitOrderBid{
			LimitOrderBiddingId: id + 1,
			BidderAddress:       bidder,
			CollateralToken:     collateralAssetToken, // zero
			DebtToken:           amount,               // user's balance
			BiddingId:           nil,
			PremiumDiscount:     premiumDiscount,
		}
	} else {
		userLimitBid.DebtToken = userLimitBid.DebtToken.Add(amount)
	}

	// send tokens from user to the auction module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, bidderAddr, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	// Set ID and LimitBid Data
	k.SetLimitAuctionBidID(ctx, userLimitBid.LimitOrderBiddingId)
	k.SetUserLimitBidData(ctx, userLimitBid, DebtTokenId, CollateralTokenId, PremiumDiscount)
	return nil
}

func (k Keeper) CancelLimitAuctionBid(ctx sdk.Context, bidder string, DebtTokenId, CollateralTokenId uint64, PremiumDiscount string) error {
	userLimitBid, found := k.GetUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)
	if !found {
		// return err not found
	}

	bidderAddr, err := sdk.AccAddressFromBech32(bidder)
	if err != nil {
		return err
	}
	// return all the tokens back to the user
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, bidderAddr, sdk.NewCoins(userLimitBid.DebtToken))
	if err != nil {
		return err
	}

	// delete userLimitBid from KV store
	k.DeleteUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)

	return nil
}

func (k Keeper) WithdrawLimitAuctionBid(ctx sdk.Context, bidder string, CollateralTokenId, DebtTokenId uint64, PremiumDiscount string, amount sdk.Coin) error {
	userLimitBid, found := k.GetUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)
	if !found {
		// return err not found
	}

	bidderAddr, err := sdk.AccAddressFromBech32(bidder)
	if err != nil {
		return err
	}

	if amount.Amount.Equal(userLimitBid.DebtToken.Amount) {
		err := k.CancelLimitAuctionBid(ctx, bidder, DebtTokenId, CollateralTokenId, PremiumDiscount)
		if err != nil {
			return err
		}
		return nil
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, bidderAddr, sdk.NewCoins(amount))
	if err != nil {
		return err
	}
	userLimitBid.DebtToken.Amount = userLimitBid.DebtToken.Amount.Sub(amount.Amount)
	k.SetUserLimitBidData(ctx, userLimitBid, DebtTokenId, CollateralTokenId, PremiumDiscount)
	return nil
}
