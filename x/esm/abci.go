package esm

import (
	assettypes "github.com/petrichormoney/petri/x/asset/types"
	"github.com/petrichormoney/petri/x/esm/expected"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	utils "github.com/petrichormoney/petri/types"
	"github.com/petrichormoney/petri/x/esm/keeper"
	"github.com/petrichormoney/petri/x/esm/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper, assetKeeper expected.AssetKeeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		apps, found := assetKeeper.GetApps(ctx)
		if !found {
			return assettypes.AppIdsDoesntExist
		}
		for _, app := range apps {
			esmStatus, found := k.GetESMStatus(ctx, app.Id)
			if !found {
				continue
			}
			if found && esmStatus.Status {
				// Should check if price exists in the band or not. else should skip---- k.market.isPriceValidationActive
				// to add this check at all abci as well where price is important
				if !esmStatus.SnapshotStatus {
					err := k.SnapshotOfPrices(ctx, esmStatus)
					if err != nil {
						continue
					}
				}
				if ctx.BlockTime().After(esmStatus.EndTime) && esmStatus.SnapshotStatus {
					esmData, _ := k.GetESMTriggerParams(ctx, esmStatus.AppId)
					if !esmStatus.VaultRedemptionStatus {
						err := k.SetUpCollateralRedemptionForVault(ctx, esmStatus.AppId, esmData)
						if err != nil {
							continue
						}
					}
					if !esmStatus.StableVaultRedemptionStatus {
						err := k.SetUpCollateralRedemptionForStableVault(ctx, esmStatus.AppId, esmData)
						if err != nil {
							continue
						}
					}

					if !esmStatus.CollectorTransaction {
						err := k.SetUpDebtRedemptionForCollector(ctx, esmStatus.AppId)
						if err != nil {
							continue
						}
					}
					if !esmStatus.ShareCalculation && esmStatus.VaultRedemptionStatus && esmStatus.StableVaultRedemptionStatus && esmStatus.CollectorTransaction {
						err := k.SetUpShareCalculation(ctx, esmStatus.AppId)
						if err != nil {
							continue
						}
					}
				}
			}
		}
		return nil
	})
}

//Collector
//Stable Vault
//Vault
//Share Calculation
//------------------------
//Bool Value for all these four components
