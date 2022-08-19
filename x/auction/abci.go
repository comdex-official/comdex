package auction

import (
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/auction/keeper"
	"github.com/comdex-official/comdex/x/auction/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {

	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		auctionMapData, auctionMappingFound := k.GetAllAuctionMappingForApp(ctx)
		if auctionMappingFound {
			for _, data := range auctionMapData {
				killSwitchParams, _ := k.GetKillSwitchData(ctx, data.AppId)
				esmStatus, found := k.GetESMStatus(ctx, data.AppId)
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
				err3 := k.DistributorActivator(ctx, data, killSwitchParams, status)
				if err3 != nil {
					ctx.Logger().Error("error in distributor activator")
				}
			}
		}

		lockedVaults := k.GetLockedVaults(ctx)

		if len(lockedVaults) > 0 {
			err3 := k.DutchActivator(ctx, lockedVaults)
			if err3 != nil {
				ctx.Logger().Error("error in dutch activator")
			}

			err5 := k.LendDutchActivator(ctx, lockedVaults)
			if err5 != nil {
				ctx.Logger().Error("error in lend dutch activator")
			}
		}

		apps, appsFound := k.GetApps(ctx)

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
