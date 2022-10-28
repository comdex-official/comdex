package v5_0_0 //nolint:revive,stylecheck

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v5_0_0.beta
func CreateUpgradeHandlerV5Beta(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// This change is only for testnet upgrade
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

//func SetVaultLengthCounter(
//	ctx sdk.Context,
//	vaultkeeper vaultkeeper.Keeper,
//) {
//	var count uint64
//	appExtendedPairVaultData, found := vaultkeeper.GetAppMappingData(ctx, 2)
//	if found {
//		for _, data := range appExtendedPairVaultData {
//			count += uint64(len(data.VaultIds))
//		}
//	}
//	vaultkeeper.SetLengthOfVault(ctx, count)
//}
