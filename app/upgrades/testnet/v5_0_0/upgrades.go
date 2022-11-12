package v5_0_0 //nolint:revive,stylecheck

import (
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandlerV5Beta(
	mm *module.Manager,
	configurator module.Configurator,
	lk lendkeeper.Keeper,
	liqk liquidationkeeper.Keeper,
	vaultkeeper vaultkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		SetVaultLengthCounter(ctx, vaultkeeper)
		err = FuncMigrateLiquidatedBorrow(ctx, lk, liqk)
		if err != nil {
			return nil, err
		}
		return vm, err
	}
}

func CreateUpgradeHandlerV51Beta(
	mm *module.Manager,
	configurator module.Configurator,
	lk lendkeeper.Keeper,
	vaultkeeper vaultkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		SetVaultLengthCounter(ctx, vaultkeeper)
		FixAssetStatsData(ctx, lk)
		return vm, err
	}
}

func FixAssetStatsData(
	ctx sdk.Context,
	lk lendkeeper.Keeper,
) {
	// first deleting asset-stats of previous state and adding data again
	// same with borrow
	// lastly delete previous asset-pair mappings and add new asset-pair mapping
	lbMapping11 := lendtypes.PoolAssetLBMapping{
		PoolID:                   1,
		AssetID:                  1,
		TotalBorrowed:            sdk.ZeroInt(),
		TotalStableBorrowed:      sdk.ZeroInt(),
		TotalLend:                sdk.ZeroInt(),
		TotalInterestAccumulated: sdk.ZeroInt(),
	}

	lbMapping12 := lendtypes.PoolAssetLBMapping{
		PoolID:                   1,
		AssetID:                  2,
		TotalBorrowed:            sdk.ZeroInt(),
		TotalStableBorrowed:      sdk.ZeroInt(),
		TotalLend:                sdk.ZeroInt(),
		TotalInterestAccumulated: sdk.ZeroInt(),
	}
	lbMapping13 := lendtypes.PoolAssetLBMapping{
		PoolID:                   1,
		AssetID:                  3,
		TotalBorrowed:            sdk.ZeroInt(),
		TotalStableBorrowed:      sdk.ZeroInt(),
		TotalLend:                sdk.ZeroInt(),
		TotalInterestAccumulated: sdk.ZeroInt(),
	}
	lbMapping21 := lendtypes.PoolAssetLBMapping{
		PoolID:                   2,
		AssetID:                  1,
		TotalBorrowed:            sdk.ZeroInt(),
		TotalStableBorrowed:      sdk.ZeroInt(),
		TotalLend:                sdk.ZeroInt(),
		TotalInterestAccumulated: sdk.ZeroInt(),
	}
	lbMapping24 := lendtypes.PoolAssetLBMapping{
		PoolID:                   2,
		AssetID:                  4,
		TotalBorrowed:            sdk.ZeroInt(),
		TotalStableBorrowed:      sdk.ZeroInt(),
		TotalLend:                sdk.ZeroInt(),
		TotalInterestAccumulated: sdk.ZeroInt(),
	}
	lbMapping23 := lendtypes.PoolAssetLBMapping{
		PoolID:                   2,
		AssetID:                  3,
		TotalBorrowed:            sdk.ZeroInt(),
		TotalStableBorrowed:      sdk.ZeroInt(),
		TotalLend:                sdk.ZeroInt(),
		TotalInterestAccumulated: sdk.ZeroInt(),
	}

	lk.DeleteAssetStatsByPoolIDAndAssetID(ctx, lbMapping11)
	lk.DeleteAssetStatsByPoolIDAndAssetID(ctx, lbMapping12)
	lk.DeleteAssetStatsByPoolIDAndAssetID(ctx, lbMapping13)
	lk.DeleteAssetStatsByPoolIDAndAssetID(ctx, lbMapping21)
	lk.DeleteAssetStatsByPoolIDAndAssetID(ctx, lbMapping24)
	lk.DeleteAssetStatsByPoolIDAndAssetID(ctx, lbMapping23)

	lk.SetAssetStatsByPoolIDAndAssetID(ctx, lbMapping11)
	lk.SetAssetStatsByPoolIDAndAssetID(ctx, lbMapping12)
	lk.SetAssetStatsByPoolIDAndAssetID(ctx, lbMapping13)
	lk.SetAssetStatsByPoolIDAndAssetID(ctx, lbMapping21)
	lk.SetAssetStatsByPoolIDAndAssetID(ctx, lbMapping24)
	lk.SetAssetStatsByPoolIDAndAssetID(ctx, lbMapping23)

	lends := lk.GetAllLend(ctx)
	for _, v := range lends {
		lbmap, _ := lk.GetAssetStatsByPoolIDAndAssetID(ctx, v.PoolID, v.AssetID)
		lbmap.TotalLend = lbmap.TotalLend.Add(v.AmountIn.Amount)
		lbmap.LendIds = append(lbmap.LendIds, v.ID)
		lk.SetAssetStatsByPoolIDAndAssetID(ctx, lbmap)
	}

	borrows := lk.GetAllBorrow(ctx)
	for _, v := range borrows {
		pair, _ := lk.GetLendPair(ctx, v.PairID)
		lbmap, _ := lk.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
		lbmap.TotalBorrowed = lbmap.TotalBorrowed.Add(v.AmountOut.Amount)
		lbmap.BorrowIds = append(lbmap.BorrowIds, v.ID)
		lk.SetAssetStatsByPoolIDAndAssetID(ctx, lbmap)
	}

	assetToPair1 := lendtypes.AssetToPairMapping{
		PoolID:  1,
		AssetID: 1,
	}
	assetToPair2 := lendtypes.AssetToPairMapping{
		PoolID:  2,
		AssetID: 1,
	}
	assetToPair3 := lendtypes.AssetToPairMapping{
		PoolID:  3,
		AssetID: 1,
	}
	assetToPair4 := lendtypes.AssetToPairMapping{
		PoolID:  1,
		AssetID: 2,
	}
	assetToPair5 := lendtypes.AssetToPairMapping{
		PoolID:  4,
		AssetID: 2,
	}
	assetToPair6 := lendtypes.AssetToPairMapping{
		PoolID:  3,
		AssetID: 2,
	}
	lk.DeleteAssetToPair(ctx, assetToPair1)
	lk.DeleteAssetToPair(ctx, assetToPair2)
	lk.DeleteAssetToPair(ctx, assetToPair3)
	lk.DeleteAssetToPair(ctx, assetToPair4)
	lk.DeleteAssetToPair(ctx, assetToPair5)
	lk.DeleteAssetToPair(ctx, assetToPair6)

	newAssetPairMapping1 := lendtypes.AssetToPairMapping{
		PoolID:  1,
		AssetID: 1,
		PairID:  []uint64{3, 4, 15},
	}
	newAssetPairMapping2 := lendtypes.AssetToPairMapping{
		PoolID:  1,
		AssetID: 2,
		PairID:  []uint64{1, 2, 13},
	}
	newAssetPairMapping3 := lendtypes.AssetToPairMapping{
		PoolID:  1,
		AssetID: 3,
		PairID:  []uint64{5, 6, 14},
	}
	newAssetPairMapping4 := lendtypes.AssetToPairMapping{
		PoolID:  2,
		AssetID: 1,
		PairID:  []uint64{9, 10, 18},
	}
	newAssetPairMapping5 := lendtypes.AssetToPairMapping{
		PoolID:  2,
		AssetID: 3,
		PairID:  []uint64{11, 12, 17},
	}
	newAssetPairMapping6 := lendtypes.AssetToPairMapping{
		PoolID:  2,
		AssetID: 4,
		PairID:  []uint64{7, 8, 16},
	}

	lk.SetAssetToPair(ctx, newAssetPairMapping1)
	lk.SetAssetToPair(ctx, newAssetPairMapping2)
	lk.SetAssetToPair(ctx, newAssetPairMapping3)
	lk.SetAssetToPair(ctx, newAssetPairMapping4)
	lk.SetAssetToPair(ctx, newAssetPairMapping5)
	lk.SetAssetToPair(ctx, newAssetPairMapping6)
}
