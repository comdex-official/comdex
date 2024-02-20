package v14

import (
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandlerV14(
	mm *module.Manager,
	configurator module.Configurator,
	lendKeeper lendkeeper.Keeper,

) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Applying main net upgrade - v.14.0.0")
		ctx.Logger().With("upgrade", UpgradeName)

		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		UpdateLendParams(ctx, lendKeeper)
		return vm, err
	}
}

func UpdateLendParams(
	ctx sdk.Context,
	lendKeeper lendkeeper.Keeper,
) {
	assetRatesParamsStAtom, _ := lendKeeper.GetAssetRatesParams(ctx, 14)
	assetRatesParamsStAtom.CAssetID = 23
	lendKeeper.SetAssetRatesParams(ctx, assetRatesParamsStAtom)
}
