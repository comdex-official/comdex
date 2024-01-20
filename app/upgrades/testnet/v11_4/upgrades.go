package v11_4 //nolint:revive,stylecheck

import (
	"context"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateUpgradeHandlerV114(
	mm *module.Manager,
	configurator module.Configurator,
	assetKeeper assetkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		assetKeeper.SetParams(sdk.UnwrapSDKContext(ctx), assettypes.NewParams())

		sdk.UnwrapSDKContext(ctx).Logger().Info("Applying test net upgrade - v.11.4.0")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
