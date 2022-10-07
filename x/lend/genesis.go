package lend

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	// k.SetParams(ctx, state.Params)

	// for _, item := range state.BorrowAsset {
	// 	k.SetBorrow(ctx, item)
	// }

	// for _, item := range state.UserBorrowIdMapping {
	// 	k.SetUserBorrows(ctx, item)
	// }

	// for _, item := range state.BorrowIdByOwnerAndPoolMapping {
	// 	k.SetBorrowIDByOwnerAndPool(ctx, item)
	// }

	// k.SetBorrows(ctx, state.BorrowMapping)

	// for _, item := range state.LendAsset {
	// 	k.SetLend(ctx, item)
	// }

	// for _, item := range state.Pool {
	// 	k.SetPool(ctx, item)
	// }

	// for _, item := range state.AssetToPairMapping {
	// 	k.SetAssetToPair(ctx, item)
	// }

	// for _, item := range state.UserLendIdMapping {
	// 	k.SetUserLends(ctx, item)
	// }

	// for _, item := range state.LendIdByOwnerAndPoolMapping {
	// 	k.SetLendIDByOwnerAndPool(ctx, item)
	// }

	// for _, item := range state.LendIdToBorrowIdMapping {
	// 	k.SetLendIDToBorrowIDMapping(ctx, item)
	// }

	// for _, item := range state.AssetStats {
	// 	k.SetAssetStatsByPoolIDAndAssetID(ctx, item)
	// }

	// k.SetLends(ctx, state.LendMapping)

	// k.SetUserDepositStats(ctx, state.UserDepositStats)

	// k.SetReserveDepositStats(ctx, state.ReserveDepositStats)

	// k.SetBuyBackDepositStats(ctx, state.BuyBackDepositStats)

	// k.SetBorrowStats(ctx, state.BorrowDepositStats)

	// for _, item := range state.Extended_Pair {
	// 	k.SetLendPair(ctx, item)
	// }

	// for _, item := range state.AssetRatesStats {
	// 	k.SetAssetRatesStats(ctx, item)
	// }

	// for _, item := range state.AuctionParams {
	// 	err := k.AddAuctionParamsData(ctx, item)
	// 	if err != nil {
	// 		return
	// 	}
	// }
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	// borrowMap, _ := k.GetBorrows(ctx)
	// lends, _ := k.GetLends(ctx)
	// userDeposit, _ := k.GetUserDepositStats(ctx)
	// reserveDeposit, _ := k.GetReserveDepositStats(ctx)
	// buyBackDeposit, _ := k.GetBuyBackDepositStats(ctx)
	// borrowDeposit, _ := k.GetBorrowStats(ctx)
	// return types.NewGenesisState(
	// 	k.GetAllBorrow(ctx),
	// 	k.GetAllUserBorrows(ctx),
	// 	k.GetAllBorrowIDByOwnerAndPool(ctx),
	// 	borrowMap,
	// 	k.GetAllLend(ctx),
	// 	k.GetPools(ctx),
	// 	k.GetAllAssetToPair(ctx),
	// 	k.GetAllUserLends(ctx),
	// 	k.GetAllLendIDByOwnerAndPool(ctx),
	// 	k.GetAllLendIDToBorrowIDMapping(ctx),
	// 	k.GetAllAssetStatsByPoolIDAndAssetID(ctx),
	// 	lends,
	// 	userDeposit,
	// 	reserveDeposit,
	// 	buyBackDeposit,
	// 	borrowDeposit,
	// 	k.GetLendPairs(ctx),
	// 	k.GetAllAssetRatesStats(ctx),
	// 	k.GetAllAddAuctionParamsData(ctx),
	// 	k.GetParams(ctx),
	// )
	return nil
}
