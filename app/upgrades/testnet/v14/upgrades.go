package v14

import (
	"context"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	commonkeeper "github.com/comdex-official/comdex/x/common/keeper"
	commontypes "github.com/comdex-official/comdex/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateUpgradeHandlerV14(
	mm *module.Manager,
	configurator module.Configurator,
	commonkeeper commonkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		sdk.UnwrapSDKContext(ctx).Logger().Info("Applying test net upgrade - v14.0.0")
		logger := sdk.UnwrapSDKContext(ctx).Logger().With("upgrade", UpgradeName)

		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return vm, err
		}
		logger.Info("set common module params")
		commonkeeper.SetParams(sdk.UnwrapSDKContext(ctx), commontypes.DefaultParams())

		return vm, err
	}
}
