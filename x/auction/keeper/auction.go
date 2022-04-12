package keeper

import (
	"time"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) GetCollateralAuctionID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetCollateralAuctionID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetCollateralAuction(ctx sdk.Context, auction auctiontypes.CollateralAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionKey(auction.Id)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteCollateralAuction(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) GetCollateralAuction(ctx sdk.Context, id uint64) (auction auctiontypes.CollateralAuction, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetCollateralAuctions(ctx sdk.Context) (auctions []auctiontypes.CollateralAuction) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.CollateralAuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.CollateralAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetBidding(ctx sdk.Context, bidding auctiontypes.Biddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsKey(bidding.Id)
		value = k.cdc.MustMarshal(&bidding)
	)
	store.Set(key, value)
}

func (k *Keeper) GetBidding(ctx sdk.Context, id uint64) (bidding auctiontypes.Biddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return bidding, false
	}

	k.cdc.MustUnmarshal(value, &bidding)
	return bidding, true
}

func (k *Keeper) GetBiddings(ctx sdk.Context) (biddings []auctiontypes.Biddings) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.BiddingsKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var bidding auctiontypes.Biddings
		k.cdc.MustUnmarshal(iter.Value(), &bidding)
		biddings = append(biddings, bidding)
	}

	return biddings
}

func (k *Keeper) GetUserBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetUserBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetUserBidding(ctx sdk.Context, userBiddings auctiontypes.UserBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsKey(userBiddings.Bidder)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
}

func (k *Keeper) GetUserBiddings(ctx sdk.Context, bidder string) (userBiddings auctiontypes.UserBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsKey(bidder)
		value = store.Get(key)
	)

	if value == nil {
		return userBiddings, false
	}

	k.cdc.MustUnmarshal(value, &userBiddings)
	return userBiddings, true
}

func (k Keeper) CreateNewAuctions(ctx sdk.Context) {
	locked_vaults := k.GetLockedVaults(ctx)
	for _, locked_vault := range locked_vaults {
		pair, found := k.GetPair(ctx, locked_vault.PairId)
		if !found {
			continue
		}
		assetIn, found := k.GetAsset(ctx, pair.AssetIn)
		if !found {
			continue
		}

		assetOut, found := k.GetAsset(ctx, pair.AssetOut)
		if !found {
			continue
		}
		collateralizationRatio, err := k.CalculateCollaterlizationRatio(ctx, locked_vault.AmountIn, assetIn, locked_vault.AmountOut, assetOut)
		if err != nil {
			continue
		}
		if sdk.Dec.LT(collateralizationRatio, pair.LiquidationRatio) && !locked_vault.IsAuctionInProgress {
			k.StartCollateralAuction(ctx, locked_vault, pair, assetIn, assetOut)
		}
	}
}

func (k Keeper) CloseAuctions(ctx sdk.Context) {
	collateral_auctions := k.GetCollateralAuctions(ctx)
	for _, collateral_auction := range collateral_auctions {
		if ctx.BlockTime().After(collateral_auction.EndTime) {
			k.CloseCollateralAuction(ctx, collateral_auction)
		}
	}
}

func (k Keeper) StartCollateralAuction(
	ctx sdk.Context,
	locked_vault liquidationtypes.LockedVault,
	pair assettypes.Pair,
	assetIn assettypes.Asset,
	assetOut assettypes.Asset,
) error {

	assetInPrice, found := k.GetPriceForAsset(ctx, pair.AssetIn)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}

	assetOutPrice, found := k.GetPriceForAsset(ctx, pair.AssetOut)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}

	if locked_vault.CollateralToBeAuctioned.LTE(sdk.NewDec(0)) {
		return auctiontypes.ErrorInvalidAuctioningCollateral
	}

	auctionParams := k.GetParams(ctx)

	liquidatedQuantity := sdk.NewDec(locked_vault.CollateralToBeAuctioned.Quo(sdk.NewDec(int64(assetInPrice))).RoundInt64())

	penaltyQuantity := liquidatedQuantity.Mul(sdk.MustNewDecFromStr(auctionParams.LiquidationPenaltyPercent).Mul(sdk.NewDec(100))).Quo(sdk.NewDec(100))
	discountedQuantity := liquidatedQuantity.Mul(sdk.MustNewDecFromStr(auctionParams.AuctionDiscountPercent).Mul(sdk.NewDec(100))).Quo(sdk.NewDec(100))
	auctioningQuantity := liquidatedQuantity.Sub(penaltyQuantity.Add(discountedQuantity))
	minBid := auctioningQuantity.Mul(sdk.NewDec(int64(assetInPrice))).Quo(sdk.NewDec(int64(assetOutPrice))).Ceil().RoundInt()
	maxBid := auctioningQuantity.Add(discountedQuantity).Mul(sdk.NewDec(int64(assetInPrice))).Quo(sdk.NewDec(int64(assetOutPrice))).Ceil().RoundInt()

	auction := auctiontypes.CollateralAuction{
		LockedVaultId:       locked_vault.LockedVaultId,
		AuctionedCollateral: sdk.NewCoin(assetIn.Denom, sdk.NewInt(auctioningQuantity.RoundInt64())),
		DiscountQuantity:    sdk.NewCoin(assetIn.Denom, sdk.NewInt(discountedQuantity.RoundInt64())),
		ActiveBiddingId:     0,
		Bidder:              nil,
		Bid:                 sdk.NewCoin(assetOut.Denom, sdk.NewInt(0)),
		MinBid:              sdk.NewCoin(assetOut.Denom, minBid),
		MaxBid:              sdk.NewCoin(assetOut.Denom, maxBid),
		EndTime:             ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		Pair:                pair,
		BiddingIds:          []uint64{},
	}
	auction.Id = k.GetCollateralAuctionID(ctx) + 1
	k.SetCollateralAuctionID(ctx, auction.Id)
	k.SetCollateralAuction(ctx, auction)
	k.SetFlagIsAuctionInProgress(ctx, locked_vault.LockedVaultId, true)
	return nil
}

func (k Keeper) CloseCollateralAuction(
	ctx sdk.Context,
	collateral_auction auctiontypes.CollateralAuction,
) error {

	if collateral_auction.Bidder != nil && collateral_auction.Bid.Amount.GTE(collateral_auction.MinBid.Amount) {

		assetIn, found := k.GetAsset(ctx, collateral_auction.Pair.AssetIn)
		if !found {
			return assettypes.ErrorAssetDoesNotExist
		}
		assetOut, found := k.GetAsset(ctx, collateral_auction.Pair.AssetOut)
		if !found {
			return assettypes.ErrorAssetDoesNotExist
		}

		highestBidReceived := collateral_auction.Bid

		collateralQuantity := collateral_auction.AuctionedCollateral.Amount.Add(collateral_auction.DiscountQuantity.Amount)

		err := k.bank.SendCoinsFromModuleToAccount(ctx, vaulttypes.ModuleName, collateral_auction.Bidder, sdk.NewCoins(sdk.NewCoin(assetIn.Denom, collateralQuantity)))
		if err != nil {
			return err
		}
		bidding, _ := k.GetBidding(ctx, collateral_auction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.SuccessBiddingStatus
		k.SetBidding(ctx, bidding)
		k.BurnCAssets(ctx, liquidationtypes.ModuleName, assetIn.Denom, assetOut.Denom, highestBidReceived.Amount)
		k.UpdateAssetQuantitiesInLockedVault(ctx, collateral_auction, collateralQuantity, assetIn, highestBidReceived.Amount, assetOut)

		for _, biddingId := range collateral_auction.BiddingIds {
			bidding, found := k.GetBidding(ctx, biddingId)
			if !found {
				continue
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			k.SetBidding(ctx, bidding)
		}
	}
	k.SetFlagIsAuctionComplete(ctx, collateral_auction.LockedVaultId, true)
	k.SetFlagIsAuctionInProgress(ctx, collateral_auction.LockedVaultId, false)
	k.DeleteCollateralAuction(ctx, collateral_auction.Id)
	return nil
}

func (k Keeper) CreateNewBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) (biddingId uint64, err error) {
	auction, found := k.GetCollateralAuction(ctx, auctionId)
	if !found {
		return 0, auctiontypes.ErrorInvalidAuctionId
	}
	bidding := auctiontypes.Biddings{
		Id:                  k.GetBiddingID(ctx) + 1,
		AuctionId:           auctionId,
		AuctionStatus:       auctiontypes.ActiveAuctionStatus,
		AuctionedCollateral: auction.AuctionedCollateral,
		Bidder:              bidder.String(),
		Bid:                 bid,
		BiddingTimestamp:    ctx.BlockTime(),
		BiddingStatus:       auctiontypes.PlacedBiddingStatus,
	}
	k.SetBiddingID(ctx, bidding.Id)
	k.SetBidding(ctx, bidding)

	userBiddings, found := k.GetUserBiddings(ctx, bidder.String())
	if !found {
		userBiddings = auctiontypes.UserBiddings{
			Id:         k.GetUserBiddingID(ctx) + 1,
			Bidder:     bidder.String(),
			BiddingIds: []uint64{},
		}
		k.SetUserBiddingID(ctx, userBiddings.Id)
	}
	userBiddings.BiddingIds = append(userBiddings.BiddingIds, bidding.Id)
	k.SetUserBidding(ctx, userBiddings)
	return bidding.Id, nil
}

func (k Keeper) PlaceBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) error {
	auction, found := k.GetCollateralAuction(ctx, auctionId)
	if !found {
		return auctiontypes.ErrorInvalidAuctionId
	}
	if bid.Denom != auction.MinBid.Denom {
		return auctiontypes.ErrorInvalidBiddingDenom
	}
	if bid.Amount.LT(auction.MinBid.Amount) {
		return auctiontypes.ErrorLowBidAmount
	}
	if bid.Amount.GT(auction.MaxBid.Amount) {
		return auctiontypes.ErrorMaxBidAmount
	}
	if bid.Amount.LT(auction.Bid.Amount.Add(sdk.NewInt(1))) {
		return auctiontypes.ErrorBidAlreadyExists
	}
	err := k.SendCoinsFromAccountToModule(ctx, bidder, liquidationtypes.ModuleName, sdk.NewCoins(bid))
	if err != nil {
		return err
	}
	biddingId, err := k.CreateNewBid(ctx, auctionId, bidder, bid)
	if err != nil {
		return err
	}
	// auction.Bidder as previous bidder
	err = k.bank.SendCoinsFromModuleToAccount(ctx, liquidationtypes.ModuleName, auction.Bidder, sdk.NewCoins(auction.Bid))
	if err != nil {
		return err
	}
	if auction.ActiveBiddingId != 0 {
		bidding, _ := k.GetBidding(ctx, auction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.RejectedBiddingStatus
		k.SetBidding(ctx, bidding)
	}
	auction.ActiveBiddingId = biddingId
	auction.BiddingIds = append(auction.BiddingIds, biddingId)
	auction.Bidder = bidder
	auction.Bid = bid
	k.SetCollateralAuction(ctx, auction)
	return nil
}
