package keeper

import (
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetGenSurplusAuction(ctx sdk.Context, auction auctiontypes.SurplusAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(auction.AppId, "surplus", auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
}

func (k Keeper) SetGenDebtAuction(ctx sdk.Context, auction auctiontypes.DebtAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(auction.AppId, "debt", auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
}

func (k Keeper) SetGenDutchAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(auction.AppId, "dutch", auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)

	store.Set(key, value)
}

func (k Keeper) SetGenProtocolStatistics(ctx sdk.Context, appID, assetID uint64, amount sdk.Dec) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.ProtocolStatisticsKey(appID, assetID)
	)
	var stats auctiontypes.ProtocolStatistics
	stats.AppId = appID
	stats.AssetId = assetID
	stats.Loss = amount
	value := k.cdc.MustMarshal(&stats)
	store.Set(key, value)
}

func (k Keeper) SetGenLendDutchLendAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.LendAuctionKey(auction.AppId, "dutch", auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)

	store.Set(key, value)
}