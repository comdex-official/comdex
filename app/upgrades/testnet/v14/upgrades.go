package v14

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	commonkeeper "github.com/comdex-official/comdex/x/common/keeper"
	commontypes "github.com/comdex-official/comdex/x/common/types"
)

func CreateUpgradeHandlerV14(
	mm *module.Manager,
	configurator module.Configurator,
	commonkeeper commonkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		ctx.Logger().Info("Applying test net upgrade - v14.0.0")
		logger := ctx.Logger().With("upgrade", UpgradeName)

		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return vm, err
		}
		logger.Info("set common module params")
		commonkeeper.SetParams(ctx, commontypes.DefaultParams())
		
		return vm, err
	}
}
