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
	assets := k.GetAssetsForOracle(ctx)
	for _, v := range apps {
		esmMarket, found := k.GetESMMarketForAsset(ctx, v.Id)
		if !found {
			return
		}
		esmStatus, found := k.GetESMStatus(ctx, v.Id)
		if !found {
			return
		}

		if !esmMarket.IsPriceSet && esmStatus.Status {
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

		if ctx.BlockTime().After(esmStatus.EndTime) && esmStatus.Status {
			err := k.SetUpCollateralRedemption(ctx, esmStatus.AppId)
			if err != nil {
				return
			}
		}
	}
}
