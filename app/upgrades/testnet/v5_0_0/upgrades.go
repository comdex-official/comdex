package v5_0_0 //nolint:revive,stylecheck

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandlerV5Beta(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.JSONCodec,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// This change is only for testnet upgrade
		//fromVM[lendtypes.ModuleName] = lend.AppModule{}.ConsensusVersion()
		//fromVM[assettypes.ModuleName] = asset.AppModule{}.ConsensusVersion()
		//mm.Modules[lendtypes.ModuleName].InitGenesis(ctx, cdc, myCustomGenesisState)
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
