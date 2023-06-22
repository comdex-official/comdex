package v12

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v4/keeper"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v4/types"
)

func CreateUpgradeHandlerV12(
	mm *module.Manager,
	configurator module.Configurator,
	icqkeeper *icqkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Applying main net upgrade - v.12.0.0")

		icqparams := icqtypes.DefaultParams()
		icqparams.AllowQueries = append(icqparams.AllowQueries, "/cosmwasm.wasm.v1.Query/SmartContractState")
		icqkeeper.SetParams(ctx, icqparams)

		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		return vm, err
	}
}
