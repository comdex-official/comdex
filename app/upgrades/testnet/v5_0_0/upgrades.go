package v5_0_0 //nolint:revive,stylecheck

import (
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func SetVaultLengthCounter(
	ctx sdk.Context,
	vaultkeeper vaultkeeper.Keeper,
) {
	var count uint64
	appExtendedPairVaultData, found := vaultkeeper.GetAppMappingData(ctx, 2)
	if found {
		for _, data := range appExtendedPairVaultData {
			count += uint64(len(data.VaultIds))
		}
	}
	vaultkeeper.SetLengthOfVault(ctx, count)
}

func CreateUpgradeHandlerV5Beta(
	mm *module.Manager,
	configurator module.Configurator,
	lk lendkeeper.Keeper,
	liqk liquidationkeeper.Keeper,
	vaultkeeper vaultkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		SetVaultLengthCounter(ctx, vaultkeeper)
		err := FuncMigrateLiquidatedBorrow(ctx, lk, liqk)
		if err != nil {
			return nil, err
		}
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
