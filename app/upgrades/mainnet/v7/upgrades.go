package v7

import (
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func InitializeLendReservesStates(
	ctx sdk.Context,
	lendKeeper lendkeeper.Keeper,
) {
	dataAsset1, _ := lendKeeper.GetReserveBuybackAssetData(ctx, 1)
	dataAsset2, _ := lendKeeper.GetReserveBuybackAssetData(ctx, 2)
	dataAsset3, _ := lendKeeper.GetReserveBuybackAssetData(ctx, 3)
	reserveStat1 := types.AllReserveStats{
		AssetID:                        1,
		AmountOutFromReserveToLenders:  sdk.ZeroInt(),
		AmountOutFromReserveForAuction: sdk.ZeroInt(),
		AmountInFromLiqPenalty:         sdk.ZeroInt(),
		AmountInFromRepayments:         dataAsset1.BuybackAmount.Add(dataAsset1.ReserveAmount),
		TotalAmountOutToLenders:        sdk.ZeroInt(),
	}

	reserveStat2 := types.AllReserveStats{
		AssetID:                        2,
		AmountOutFromReserveToLenders:  sdk.ZeroInt(),
		AmountOutFromReserveForAuction: sdk.ZeroInt(),
		AmountInFromLiqPenalty:         sdk.ZeroInt(),
		AmountInFromRepayments:         dataAsset2.BuybackAmount.Add(dataAsset2.ReserveAmount),
		TotalAmountOutToLenders:        sdk.ZeroInt(),
	}

	reserveStat3 := types.AllReserveStats{
		AssetID:                        3,
		AmountOutFromReserveToLenders:  sdk.ZeroInt(),
		AmountOutFromReserveForAuction: sdk.ZeroInt(),
		AmountInFromLiqPenalty:         sdk.ZeroInt(),
		AmountInFromRepayments:         dataAsset3.BuybackAmount.Add(dataAsset3.ReserveAmount),
		TotalAmountOutToLenders:        sdk.ZeroInt(),
	}
	lendKeeper.SetAllReserveStatsByAssetID(ctx, reserveStat1)
	lendKeeper.SetAllReserveStatsByAssetID(ctx, reserveStat2)
	lendKeeper.SetAllReserveStatsByAssetID(ctx, reserveStat3)
}

func CreateUpgradeHandler700(
	mm *module.Manager,
	configurator module.Configurator,
	lendKeeper lendkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		InitializeLendReservesStates(ctx, lendKeeper)
		return newVM, err
	}
}
