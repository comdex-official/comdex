package v4_0_0 //nolint:revive,stylecheck

import (
	"context"
	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

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
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
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
		sdk.NewCoin(params.SwapFeeDistrDenom, sdkmath.NewInt(0)),
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
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
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
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// This change is only for testnet upgrade

		EditAndSetPair(ctx, assetkeeper)
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return newVM, err
		}
		return newVM, err
	}
}
