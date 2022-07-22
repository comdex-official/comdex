package esm

import (
	"github.com/comdex-official/comdex/x/esm/keeper"
	"github.com/comdex-official/comdex/x/esm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	apps, found := k.GetApps(ctx)
	if !found {
		return
	}
	for _, v := range apps {
		esmStatus, found := k.GetESMStatus(ctx, v.Id)
		if !found {
			return
		}
		if ctx.BlockTime().After(esmStatus.EndTime) && esmStatus.Status && !esmStatus.VaultRedemptionStatus {
			err := k.SetUpCollateralRedemptionForVault(ctx, esmStatus.AppId)
			if err != nil {
				return
			}
		}
		if ctx.BlockTime().After(esmStatus.EndTime) && esmStatus.Status && !esmStatus.StableVaultRedemptionStatus {
			err := k.SetUpCollateralRedemptionForStableVault(ctx, esmStatus.AppId)
			if err != nil {
				return
			}
		}
		esmMarket, found := k.GetESMMarketForAsset(ctx, v.Id)
		if found {
			return
		}
		if !esmMarket.IsPriceSet && esmStatus.Status {
			assets := k.GetAssetsForOracle(ctx)
			var markets []types.Market
			for _, a := range assets {
				price, found := k.GetPriceForAsset(ctx, a.Id)
				if !found {
					return
				}
				market := types.Market{
					AssetID: a.Id,
					Rates:   price,
				}
				markets = append(markets, market)
			}
			em := types.ESMMarketPrice{
				AppId:      v.Id,
				IsPriceSet: true,
				Market:     markets,
			}
			k.SetESMMarketForAsset(ctx, em)
		}
	}
}
