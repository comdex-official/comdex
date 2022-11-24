package v510

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v5
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	slashingkeeper slashingkeeper.Keeper,
	mintkeeper mintkeeper.Keeper,
	bankkeeper bankkeeper.Keeper,
	stakingkeeper stakingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Running revert of tombstoning")

		err := RevertCosTombstoning(
			ctx,
			slashingkeeper,
			mintkeeper,
			bankkeeper,
			stakingkeeper,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to revert tombstoning: %s", err))
		}

		ctx.Logger().Info("Running module migrations for v5.1.0...")
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return newVM, err
		}

		return newVM, err
	}
}
