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
	reserve, found := k.GetReserveBuybackAssetData(ctx, assetID)
	if !found {
		reserve.AssetID = assetID
		reserve.BuybackAmount = sdk.ZeroInt()
		reserve.ReserveAmount = sdk.ZeroInt()
	}

	if inc {
		reserve.BuybackAmount = reserve.BuybackAmount.Add(newAmount)
		reserve.ReserveAmount = reserve.ReserveAmount.Add(newAmount)
		if err := k.bank.SendCoinsFromModuleToModule(ctx, moduleName, types.ModuleName, sdk.NewCoins(payment)); err != nil {
			return err
		}
	} else {
		reserve.BuybackAmount = reserve.BuybackAmount.Sub(newAmount)
		reserve.ReserveAmount = reserve.ReserveAmount.Sub(newAmount)
		if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, moduleName, sdk.NewCoins(payment)); err != nil {
			return err
		}
	}
	k.SetReserveBuybackAssetData(ctx, reserve)
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
