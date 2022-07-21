package auction

import (
	"fmt"
	"github.com/comdex-official/comdex/x/auction/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	auctionMapData, auctionMappingFound := k.GetAllAuctionMappingForApp(ctx)
	if auctionMappingFound {
		for _, data := range auctionMapData {
			for _, inData := range data.AssetIdToAuctionLookup {
				killSwitchParams, _ := k.GetKillSwitchData(ctx, data.AppId)
				esmStatus, found := k.GetESMStatus(ctx, data.AppId)
				status := false
				if found {
					status = esmStatus.Status
				}
				err1 := k.SurplusActivator(ctx, data, inData, killSwitchParams, status)
				if err1 != nil {
					_ = fmt.Errorf("error in surplus activator")
					return
				}
				err2 := k.DebtActivator(ctx, data, inData, killSwitchParams, status)
				if err2 != nil {
					_ = fmt.Errorf("error in debt activator")
					return
				}
			}
		}
	}

	lockedVaults := k.GetLockedVaults(ctx)

	if len(lockedVaults) > 0 {
		err3 := k.DutchActivator(ctx, lockedVaults)
		if err3 != nil {
			_ = fmt.Errorf("error in dutch activator")
			return
		}

		err5 := k.LendDutchActivator(ctx, lockedVaults)
		if err5 != nil {
			_ = fmt.Errorf("error in lend dutch activator")
			return
		}
	}

	apps, appsFound := k.GetApps(ctx)

	if appsFound {
		for _, app := range apps {
			err4 := k.RestartDutch(ctx, app.Id)
			if err4 != nil {
				_ = fmt.Errorf("error in restart dutch activator")
				return
			}

			err6 := k.RestartLendDutch(ctx, app.Id)
			if err6 != nil {
				_ = fmt.Errorf("error in restart lend dutch activator")
				return
			}
		}
	}
}
