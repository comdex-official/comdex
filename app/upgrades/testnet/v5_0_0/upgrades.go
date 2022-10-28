package v5_0_0 //nolint:revive,stylecheck

import (
	"fmt"
	"github.com/comdex-official/comdex/x/asset"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/lend"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
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
		fromVM[lendtypes.ModuleName] = lend.AppModule{}.ConsensusVersion() - 1
		fromVM[assettypes.ModuleName] = asset.AppModule{}.ConsensusVersion() - 1
		fmt.Println("fromVM[lendtypes.ModuleName]", fromVM[lendtypes.ModuleName])
		fmt.Println("fromVM[assettypes.ModuleName]", fromVM[assettypes.ModuleName])
		//mm.Modules[lendtypes.ModuleName].InitGenesis(ctx, cdc, myCustomGenesisState)
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
