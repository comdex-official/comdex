package v5_0_0 //nolint:revive,stylecheck

import (
	"context"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateUpgradeHandlerV5Beta(
	mm *module.Manager,
	configurator module.Configurator,
	lk lendkeeper.Keeper,
	liqk liquidationkeeper.Keeper,
	vaultkeeper vaultkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		SetVaultLengthCounter(sdk.UnwrapSDKContext(ctx), vaultkeeper)
		err = FuncMigrateLiquidatedBorrow(sdk.UnwrapSDKContext(ctx), lk, liqk)
		if err != nil {
			return nil, err
		}
		return vm, err
	}
}
