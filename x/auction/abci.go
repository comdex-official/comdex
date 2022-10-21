package auction

import (
	"github.com/comdex-official/comdex/x/auction/expected"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/auction/keeper"
	"github.com/comdex-official/comdex/x/auction/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper, assetKeeper expected.AssetKeeper, collectorKeeper expected.CollectorKeeper, esmKeeper expected.EsmKeeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		auctionMapData, auctionMappingFound := collectorKeeper.GetAllAuctionMappingForApp(ctx)
		if auctionMappingFound {
			for _, data := range auctionMapData {
				killSwitchParams, _ := esmKeeper.GetKillSwitchData(ctx, data.AppId)
				esmStatus, found := esmKeeper.GetESMStatus(ctx, data.AppId)
				status := false
				if found {
					status = esmStatus.Status
				}
				err1 := k.SurplusActivator(ctx, data, killSwitchParams, status)
				if err1 != nil {
					ctx.Logger().Error("error in surplus activator")
				}
				err2 := k.DebtActivator(ctx, data, killSwitchParams, status)
				if err2 != nil {
					ctx.Logger().Error("error in debt activator")
				}
			}
		}

		apps, appsFound := assetKeeper.GetApps(ctx)

		if appsFound {
			for _, app := range apps {
				err4 := k.RestartDutch(ctx, app.Id)
				if err4 != nil {
					ctx.Logger().Error("error in restart dutch activator")
				}

				err6 := k.RestartLendDutch(ctx, app.Id)
				if err6 != nil {
					ctx.Logger().Error("error in restart lend dutch activator")
				}
			}
		}
		return nil
	})
}
