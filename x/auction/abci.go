package auction

import (
	"fmt"

	"github.com/comdex-official/comdex/x/auction/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker compounds the debt in outstanding cdps and liquidates cdps that are below the required collateralization ratio
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	fmt.Println("Liquidation Here")
	vaults := k.GetVaults(ctx)
	for _, vault := range vaults {
		pair, found := k.GetPair(ctx, vault.PairID)
		if !found {
			continue
		}
		liquidationRatio := pair.LiquidationRatio
		assetIn, found := k.GetAsset(ctx, pair.AssetIn)
		if !found {
			continue
		}

		assetOut, found := k.GetAsset(ctx, pair.AssetOut)
		if !found {
			continue
		}
		collateralizationRatio, err := k.CalculateCollaterlizationRatio(ctx, vault.AmountIn, assetIn, vault.AmountOut, assetOut)
		if err != nil {
			continue
		}
		if sdk.Dec.LT(collateralizationRatio, liquidationRatio) {
			fmt.Println("Liquidating Vault Id : ", vault.ID)
			k.SeizeCollateral(ctx, vault, assetIn, assetOut, collateralizationRatio)
		}
	}
	fmt.Println("All Auctions Below")
	auctions := k.GetAuctions(ctx)
	for _, auction := range auctions {
		fmt.Println("ID : ", auction.ID)
		fmt.Println("AuctionedVault : ", auction.AuctionedVault)
		fmt.Println("Lot : ", auction.Lot)
		fmt.Println("Bidder : ", auction.Bidder)
		fmt.Println("Bid : ", auction.Bid)
		fmt.Println("EndTime : ", auction.EndTime)
		fmt.Println("MaxEndTime : ", auction.MaxEndTime)
		fmt.Println("MaxBid : ", auction.MaxBid)
	}
}
