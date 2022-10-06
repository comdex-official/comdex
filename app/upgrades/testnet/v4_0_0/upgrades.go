package v4_0_0

import (
	auctionkeeper "github.com/comdex-official/comdex/x/auction/keeper"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	liquiditykeeper "github.com/comdex-official/comdex/x/liquidity/keeper"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"
	rewardskeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v4_0_0
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// This change is only for testnet upgrade

		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return newVM, err
		}
		return newVM, err
	}
}

func CreateSwapFeeGauge(
	ctx sdk.Context,
	rewardsKeeper rewardskeeper.Keeper,
	liquidityKeeper liquiditykeeper.Keeper,
	appID, poolID uint64,
) {
	params, _ := liquidityKeeper.GetGenericParams(ctx, appID)
	pool, _ := liquidityKeeper.GetPool(ctx, appID, poolID)
	pair, _ := liquidityKeeper.GetPair(ctx, appID, pool.PairId)
	newGauge := rewardstypes.NewMsgCreateGauge(
		appID,
		pair.GetSwapFeeCollectorAddress(),
		ctx.BlockTime(),
		rewardstypes.LiquidityGaugeTypeID,
		liquiditytypes.DefaultSwapFeeDistributionDuration,
		sdk.NewCoin(params.SwapFeeDistrDenom, sdk.NewInt(0)),
		1,
	)
	newGauge.Kind = &rewardstypes.MsgCreateGauge_LiquidityMetaData{
		LiquidityMetaData: &rewardstypes.LiquidtyGaugeMetaData{
			PoolId:       pool.Id,
			IsMasterPool: false,
			ChildPoolIds: []uint64{},
		},
	}
	_ = rewardsKeeper.CreateNewGauge(ctx, newGauge, true)
}

// CreateUpgradeHandler creates an SDK upgrade handler for v4_1_0
func CreateUpgradeHandlerV410(
	mm *module.Manager,
	configurator module.Configurator,
	rewardskeeper rewardskeeper.Keeper,
	liquiditykeeper liquiditykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// This change is only for testnet upgrade

		CreateSwapFeeGauge(ctx, rewardskeeper, liquiditykeeper, 1, 1)
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return newVM, err
		}
		return newVM, err
	}
}

// CreateUpgradeHandler creates an SDK upgrade handler for v4_2_0
func CreateUpgradeHandlerV420(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// This change is only for testnet upgrade

		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return newVM, err
		}
		return newVM, err
	}
}

func EditAndSetPair(
	ctx sdk.Context,
	assetkeeper assetkeeper.Keeper,
) {
	pair1 := assettypes.Pair{
		Id:       1,
		AssetIn:  1,
		AssetOut: 3,
	}
	assetkeeper.SetPair(ctx, pair1)
	assetkeeper.SetPairID(ctx, 3)
}

// CreateUpgradeHandler creates an SDK upgrade handler for v4_3_0
func CreateUpgradeHandlerV430(
	mm *module.Manager,
	configurator module.Configurator,
	assetkeeper assetkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// This change is only for testnet upgrade

		EditAndSetPair(ctx, assetkeeper)
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return newVM, err
		}
		return newVM, err
	}
}
func UpdateLends(
	ctx sdk.Context,
	lendKeeper lendkeeper.Keeper,
) {
	amt11 := sdk.ZeroInt()
	amt12 := sdk.ZeroInt()
	amt31 := sdk.ZeroInt()
	amt32 := sdk.ZeroInt()
	amt21 := sdk.ZeroInt()
	amt42 := sdk.ZeroInt()
	lends, _ := lendKeeper.GetLends(ctx)
	for _, id := range lends.LendIDs {
		lend, _ := lendKeeper.GetLend(ctx, id)
		newLend := lendtypes.LendAsset{
			ID:                  lend.ID,
			AssetID:             lend.AssetID,
			PoolID:              lend.PoolID,
			Owner:               lend.Owner,
			AmountIn:            lend.AmountIn,
			LendingTime:         lend.LendingTime,
			AvailableToBorrow:   lend.AvailableToBorrow,
			AppID:               lend.AppID,
			GlobalIndex:         sdk.Dec{},
			LastInteractionTime: lend.LastInteractionTime,
			CPoolName:           lend.CPoolName,
		}
		lendKeeper.SetLend(ctx, newLend)

		if lend.AssetID == 1 && lend.PoolID == 1 {
			amt11 = amt11.Add(lend.AmountIn.Amount)
		}
		if lend.AssetID == 1 && lend.PoolID == 2 {
			amt12 = amt12.Add(lend.AmountIn.Amount)
		}
		if lend.AssetID == 3 && lend.PoolID == 1 {
			amt31 = amt31.Add(lend.AmountIn.Amount)
		}
		if lend.AssetID == 3 && lend.PoolID == 2 {
			amt32 = amt32.Add(lend.AmountIn.Amount)
		}
		if lend.AssetID == 2 && lend.PoolID == 1 {
			amt21 = amt21.Add(lend.AmountIn.Amount)
		}
		if lend.AssetID == 4 && lend.PoolID == 2 {
			amt42 = amt42.Add(lend.AmountIn.Amount)
		}
	}
	assetStats11, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 1, 1)
	assetStats11.TotalLend = amt11
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats11)

	assetStats12, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 1, 2)
	assetStats12.TotalLend = amt12
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats12)

	assetStats31, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 3, 1)
	assetStats31.TotalLend = amt31
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats31)

	assetStats32, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 3, 2)
	assetStats32.TotalLend = amt32
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats32)

	assetStats21, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 2, 1)
	assetStats21.TotalLend = amt21
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats21)

	assetStats42, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 4, 2)
	assetStats42.TotalLend = amt42
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats42)

}

func UpdateBorrows(
	ctx sdk.Context,
	lendKeeper lendkeeper.Keeper,
) {
	amt11 := sdk.ZeroInt()
	amt12 := sdk.ZeroInt()
	amt31 := sdk.ZeroInt()
	amt32 := sdk.ZeroInt()
	amt21 := sdk.ZeroInt()
	amt42 := sdk.ZeroInt()
	borrows, _ := lendKeeper.GetBorrows(ctx)
	for _, id := range borrows.BorrowIDs {
		borrow, _ := lendKeeper.GetBorrow(ctx, id)
		newBorrow := lendtypes.BorrowAsset{
			ID:                  borrow.ID,
			LendingID:           borrow.LendingID,
			IsStableBorrow:      borrow.IsStableBorrow,
			PairID:              borrow.PairID,
			AmountIn:            borrow.AmountIn,
			AmountOut:           borrow.AmountOut,
			BridgedAssetAmount:  borrow.BridgedAssetAmount,
			BorrowingTime:       borrow.BorrowingTime,
			StableBorrowRate:    borrow.StableBorrowRate,
			UpdatedAmountOut:    borrow.UpdatedAmountOut,
			GlobalIndex:         borrow.GlobalIndex,
			ReserveGlobalIndex:  borrow.ReserveGlobalIndex,
			LastInteractionTime: borrow.LastInteractionTime,
			CPoolName:           borrow.CPoolName,
		}
		lendKeeper.SetBorrow(ctx, newBorrow)
		pair, _ := lendKeeper.GetLendPair(ctx, borrow.PairID)

		if pair.AssetOut == 1 && pair.AssetOutPoolID == 1 {
			amt11 = amt11.Add(borrow.AmountOut.Amount)
		}
		if pair.AssetOut == 1 && pair.AssetOutPoolID == 2 {
			amt12 = amt12.Add(borrow.AmountOut.Amount)
		}
		if pair.AssetOut == 3 && pair.AssetOutPoolID == 1 {
			amt31 = amt31.Add(borrow.AmountOut.Amount)
		}
		if pair.AssetOut == 3 && pair.AssetOutPoolID == 2 {
			amt32 = amt32.Add(borrow.AmountOut.Amount)
		}
		if pair.AssetOut == 2 && pair.AssetOutPoolID == 1 {
			amt21 = amt21.Add(borrow.AmountOut.Amount)
		}
		if pair.AssetOut == 4 && pair.AssetOutPoolID == 2 {
			amt42 = amt42.Add(borrow.AmountOut.Amount)
		}
	}

	assetStats11, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 1, 1)
	assetStats11.TotalLend = amt11
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats11)

	assetStats12, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 1, 2)
	assetStats12.TotalLend = amt12
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats12)

	assetStats31, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 3, 1)
	assetStats31.TotalLend = amt31
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats31)

	assetStats32, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 3, 2)
	assetStats32.TotalLend = amt32
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats32)

	assetStats21, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 2, 1)
	assetStats21.TotalLend = amt21
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats21)

	assetStats42, _ := lendKeeper.GetAssetStatsByPoolIDAndAssetID(ctx, 4, 2)
	assetStats42.TotalLend = amt42
	lendKeeper.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats42)

}

func UpdateDutchLendAuctions(
	ctx sdk.Context,
	liquidationkeeper liquidationkeeper.Keeper,
	auctionkeeper auctionkeeper.Keeper,
) {
	lockedVaults := liquidationkeeper.GetLockedVaults(ctx)
	for _, v := range lockedVaults {
		if v.Kind != nil {
			err := auctionkeeper.LendDutchActivator(ctx, v)
			if err != nil {
				return
			}
		}
	}
}

// CreateUpgradeHandler creates an SDK upgrade handler for v4_4_0
func CreateUpgradeHandlerV440(
	mm *module.Manager,
	configurator module.Configurator,
	lendkeeper lendkeeper.Keeper,
	liquidationkeeper liquidationkeeper.Keeper,
	auctionkeeper auctionkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// This change is only for testnet upgrade

		UpdateLends(ctx, lendkeeper)
		UpdateBorrows(ctx, lendkeeper)
		UpdateDutchLendAuctions(ctx, liquidationkeeper, auctionkeeper)
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)

		if err != nil {
			return newVM, err
		}
		return newVM, err
	}
}
