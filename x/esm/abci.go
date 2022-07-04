package esm

import (
	"github.com/comdex-official/comdex/x/esm/keeper"
	"github.com/comdex-official/comdex/x/esm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	apps, _ := k.GetApps(ctx)
	assets := k.GetAssetsForOracle(ctx)
	for _, v := range apps {
		esmMarket, _ := k.GetESMMarketForAsset(ctx, v.Id)
		esmStatus, _ := k.GetESMStatus(ctx, v.Id)

		if !esmMarket.IsPriceSet && esmStatus.Status {
			var markets []types.Market
			for _, a := range assets {
				price, _ := k.GetPriceForAsset(ctx, a.Id)
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
