package auction

import (
	"fmt"
	"github.com/comdex-official/comdex/x/auction/expected"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/auction/keeper"
	"github.com/comdex-official/comdex/x/auction/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper, assetKeeper expected.AssetKeeper, collectorKeeper expected.CollectorKeeper, esmKeeper expected.EsmKeeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	auctionMapData, auctionMappingFound := collectorKeeper.GetAllAuctionMappingForApp(ctx)
	if auctionMappingFound {
		for _, data := range auctionMapData {
			killSwitchParams, _ := esmKeeper.GetKillSwitchData(ctx, data.AppId)
			esmStatus, found := esmKeeper.GetESMStatus(ctx, data.AppId)
			status := false
			if found {
				status = esmStatus.Status
			}
			_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
				err1 := k.SurplusActivator(ctx, data, killSwitchParams, status)
				if err1 != nil {
					ctx.EventManager().EmitEvent(
						sdk.NewEvent(
							types.EventTypeSurplusActivator,
							sdk.NewAttribute(types.DataAppID, fmt.Sprintf("%d", data.AppId)),
							sdk.NewAttribute(types.DataAssetID, fmt.Sprintf("%d", data.AssetId)),
							sdk.NewAttribute(types.DataAssetOutOraclePrice, fmt.Sprintf("%t", data.AssetOutOraclePrice)),
							sdk.NewAttribute(types.DataAssetOutPrice, fmt.Sprintf("%d", data.AssetOutPrice)),
							sdk.NewAttribute(types.DatIsAuctionActive, fmt.Sprintf("%t", data.IsAuctionActive)),
							sdk.NewAttribute(types.DataIsDebtAuction, fmt.Sprintf("%t", data.IsDebtAuction)),
							sdk.NewAttribute(types.DataIsDistributor, fmt.Sprintf("%t", data.IsDistributor)),
							sdk.NewAttribute(types.DataIsSurplusAuction, fmt.Sprintf("%t", data.IsSurplusAuction)),
							sdk.NewAttribute(types.KillSwitchParamsBreakerEnabled, fmt.Sprintf("%t", killSwitchParams.BreakerEnable)),
							sdk.NewAttribute(types.Status, fmt.Sprintf("%t", status)),
						),
					)
					ctx.Logger().Error("error in surplus activator")
					return err1
				}
				return nil
			})
			_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
				err2 := k.DebtActivator(ctx, data, killSwitchParams, status)
				if err2 != nil {
					ctx.EventManager().EmitEvent(
						sdk.NewEvent(
							types.EventTypeDebtActivator,
							sdk.NewAttribute(types.DataAppID, fmt.Sprintf("%d", data.AppId)),
							sdk.NewAttribute(types.DataAssetID, fmt.Sprintf("%d", data.AssetId)),
							sdk.NewAttribute(types.DataAssetOutOraclePrice, fmt.Sprintf("%t", data.AssetOutOraclePrice)),
							sdk.NewAttribute(types.DataAssetOutPrice, fmt.Sprintf("%d", data.AssetOutPrice)),
							sdk.NewAttribute(types.DatIsAuctionActive, fmt.Sprintf("%t", data.IsAuctionActive)),
							sdk.NewAttribute(types.DataIsDebtAuction, fmt.Sprintf("%t", data.IsDebtAuction)),
							sdk.NewAttribute(types.DataIsDistributor, fmt.Sprintf("%t", data.IsDistributor)),
							sdk.NewAttribute(types.DataIsSurplusAuction, fmt.Sprintf("%t", data.IsSurplusAuction)),
							sdk.NewAttribute(types.KillSwitchParamsBreakerEnabled, fmt.Sprintf("%t", killSwitchParams.BreakerEnable)),
							sdk.NewAttribute(types.Status, fmt.Sprintf("%t", status)),
						),
					)
					ctx.Logger().Error("error in debt activator")
					return err2
				}
				return nil
			})
		}
	}

	apps, appsFound := assetKeeper.GetApps(ctx)

	if appsFound {
		for _, app := range apps {
			err4 := k.RestartDutch(ctx, app.Id)
			if err4 != nil {
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeRestartDutch,
						sdk.NewAttribute(types.DataAppID, fmt.Sprintf("%d", app.Id)),
					),
				)
				ctx.Logger().Error("error in restart dutch activator")
			}

			err6 := k.RestartLendDutch(ctx, app.Id)
			if err6 != nil {
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeRestartLendDutch,
						sdk.NewAttribute(types.DataAppID, fmt.Sprintf("%d", app.Id)),
					),
				)
				ctx.Logger().Error("error in restart lend dutch activator")
			}
		}
	}
}
