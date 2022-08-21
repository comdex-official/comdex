package v4_0_0

import (
	"github.com/comdex-official/comdex/x/collector"
	collectorKeeper "github.com/comdex-official/comdex/x/collector/keeper"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	"github.com/comdex-official/comdex/x/liquidation"
	liquidationKeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	locker "github.com/comdex-official/comdex/x/locker"
	lockerKeeper "github.com/comdex-official/comdex/x/locker/keeper"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	vault "github.com/comdex-official/comdex/x/vault"

	vaultKeeper "github.com/comdex-official/comdex/x/vault/keeper"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v4_0_0
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	vaultKeeper vaultKeeper.Keeper,
	lockerKeeper lockerKeeper.Keeper,
	collectorKeeper collectorKeeper.Keeper,
	liquidationKeeper liquidationKeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// This change is only for testnet upgrade
		vault.InitGenesis(ctx, vaultKeeper, vaulttypes.DefaultGenesisState())
		locker.InitGenesis(ctx, lockerKeeper, lockertypes.DefaultGenesisState())
		collector.InitGenesis(ctx, collectorKeeper, collectortypes.DefaultGenesisState())
		liquidation.InitGenesis(ctx, liquidationKeeper, liquidationtypes.DefaultGenesisState())
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)

		if err != nil {
			return newVM, err
		}
		return newVM, err
	}
}
