package lend

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	var (
		borrowID       uint64
		lendID         uint64
		poolID         uint64
		extendedPairID uint64
	)

	k.SetParams(ctx, state.Params)

	for _, item := range state.BorrowAsset {
		k.SetBorrow(ctx, item)
		borrowID = item.ID
	}

	for _, item := range state.BorrowInterestTracker {
		k.SetBorrowInterestTracker(ctx, item)
	}

	for _, item := range state.LendAsset {
		k.SetLend(ctx, item)
		lendID = item.ID
	}

	for _, item := range state.Pool {
		k.SetPool(ctx, item)
		poolID = item.PoolID
	}

	for _, item := range state.AssetToPairMapping {
		k.SetAssetToPair(ctx, item)
	}

	for _, item := range state.PoolAssetLBMapping {
		k.SetAssetStatsByPoolIDAndAssetID(ctx, item)
	}

	for _, item := range state.LendRewardsTracker {
		k.SetLendRewardTracker(ctx, item)
	}

	for _, item := range state.UserAssetLendBorrowMapping {
		k.SetUserLendBorrowMapping(ctx, item)
	}

	for _, item := range state.ReserveBuybackAssetData {
		k.SetReserveBuybackAssetData(ctx, item)
	}

	for _, item := range state.Extended_Pair {
		k.SetLendPair(ctx, item)
		extendedPairID = item.Id
	}

	for _, item := range state.AuctionParams {
		err := k.AddAuctionParamsData(ctx, item)
		if err != nil {
			continue
		}
	}
	for _, item := range state.AssetRatesParams {
		k.SetAssetRatesParams(ctx, item)
	}
	k.SetUserBorrowIDCounter(ctx, borrowID)
	k.SetUserLendIDCounter(ctx, lendID)
	k.SetPoolID(ctx, poolID)
	k.SetLendPairID(ctx, extendedPairID)
	k.SetStableBorrowIds(ctx, state.StableBorrowMapping)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetAllBorrow(ctx),
		k.GetAllBorrowInterestTracker(ctx),
		k.GetAllLend(ctx),
		k.GetPools(ctx),
		k.GetAllAssetToPair(ctx),
		k.GetAllAssetStatsByPoolIDAndAssetID(ctx),
		k.GetAllLendRewardTracker(ctx),
		k.GetAllUserTotalMappingData(ctx),
		k.GetAllReserveBuybackAssetData(ctx),
		k.GetLendPairs(ctx),
		k.GetAllAddAuctionParamsData(ctx),
		k.GetAllAssetRatesParams(ctx),
		k.GetStableBorrowIds(ctx),
		k.GetParams(ctx),
	)
}
