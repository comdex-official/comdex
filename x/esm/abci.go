package esm

import (
	utils "github.com/comdex-official/comdex/types"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/esm/keeper"
	"github.com/comdex-official/comdex/x/esm/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {

	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		apps, found := k.GetApps(ctx)
		if !found {
			return assettypes.AppIdsDoesntExist
		}
		for _, v := range apps {
			esmStatus, found := k.GetESMStatus(ctx, v.Id)
			if !found {
				continue
			}
			if ctx.BlockTime().After(esmStatus.EndTime) && esmStatus.Status && !esmStatus.VaultRedemptionStatus {
				err := k.SetUpCollateralRedemptionForVault(ctx, esmStatus.AppId)
				if err != nil {
					continue
				}
			}
			if ctx.BlockTime().After(esmStatus.EndTime) && esmStatus.Status && !esmStatus.StableVaultRedemptionStatus {
				err := k.SetUpCollateralRedemptionForStableVault(ctx, esmStatus.AppId)
				if err != nil {
					continue
				}
			}
			if !esmStatus.SnapshotStatus && esmStatus.Status {
				err := k.SnapshotOfPrices(ctx, esmStatus)
				if err != nil {
					continue
				}
			}
		}
		return nil
	})
}
