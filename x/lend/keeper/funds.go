package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) GetReserveFunds(_ sdk.Context, pool types.Pool) sdk.Int {
	return sdk.NewInt(int64(pool.ReserveFunds))
}

func (k Keeper) UpdateReserveBalances(ctx sdk.Context, assetID uint64, moduleName string, payment sdk.Coin, inc bool) error {
	newAmount := payment.Amount.Quo(sdk.NewIntFromUint64(types.Uint64Two))
	buyBackStats, _ := k.GetBuyBackDepositStats(ctx)
	reserveStats, _ := k.GetReserveDepositStats(ctx)
	//var reserveBalanceStats []types.BalanceStats

	//var balanceStats []types.BalanceStats
	if inc {
		for _, v := range buyBackStats.BalanceStats {
			if v.AssetID == assetID {
				v.Amount = v.Amount.Add(newAmount)
				k.SetBuyBackDepositStats(ctx, buyBackStats)

			}
			//balanceStats = append(balanceStats, v)
			//newDepositStats := types.DepositStats{BalanceStats: balanceStats}

		}
		for _, v := range reserveStats.BalanceStats {
			if v.AssetID == assetID {
				v.Amount = v.Amount.Add(newAmount)
				k.SetReserveDepositStats(ctx, reserveStats)

			}
			//reserveBalanceStats = append(reserveBalanceStats, v)
			//newUserDepositStats := types.DepositStats{BalanceStats: reserveBalanceStats}
		}
		if err := k.bank.SendCoinsFromModuleToModule(ctx, moduleName, types.ModuleName, sdk.NewCoins(payment)); err != nil {
			return err
		}
	} else {
		for _, v := range buyBackStats.BalanceStats {
			if v.AssetID == assetID {
				v.Amount = v.Amount.Sub(newAmount)
				k.SetBuyBackDepositStats(ctx, buyBackStats)

			}
			//balanceStats = append(balanceStats, v)
			//newDepositStats := types.DepositStats{BalanceStats: balanceStats}
		}
		for _, v := range reserveStats.BalanceStats {
			if v.AssetID == assetID {
				v.Amount = v.Amount.Sub(newAmount)
				k.SetReserveDepositStats(ctx, reserveStats)

			}
			//reserveBalanceStats = append(reserveBalanceStats, v)
			//newUserDepositStats := types.DepositStats{BalanceStats: reserveBalanceStats}
		}
		if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, moduleName, sdk.NewCoins(payment)); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) UpdateLendStats(ctx sdk.Context, AssetID, PoolID uint64, amount sdk.Int, inc bool) {
	assetStats, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, PoolID, AssetID)
	if inc {
		assetStats.TotalLend = assetStats.TotalLend.Add(amount)
	} else {
		assetStats.TotalLend = assetStats.TotalLend.Sub(amount)
	}
	k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
}

func (k Keeper) UpdateBorrowStats(ctx sdk.Context, pair types.Extended_Pair, isStableBorrow bool, amount sdk.Int, inc bool) {
	assetStats, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
	if inc {
		if isStableBorrow {
			assetStats.TotalStableBorrowed = assetStats.TotalStableBorrowed.Add(amount)
		} else {
			assetStats.TotalBorrowed = assetStats.TotalBorrowed.Add(amount)
		}
	} else {
		if isStableBorrow {
			assetStats.TotalStableBorrowed = assetStats.TotalStableBorrowed.Sub(amount)
		} else {
			assetStats.TotalBorrowed = assetStats.TotalBorrowed.Sub(amount)
		}
	}
	k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
}
