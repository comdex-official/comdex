package v4_0_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
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
