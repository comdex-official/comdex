package collector

import (
	"github.com/comdex-official/comdex/x/collector/keeper"
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	k.SetParams(ctx, state.Params)

	for _, item := range state.NetFeeCollectedData {
		for _, initem := range item.AssetIdToFeeCollected {
			k.SetNetFeeCollectedData(ctx, item.AppId, initem.AssetId, initem.NetFeesCollected)
		}
	}

	for _, item := range state.AppIdToAssetCollectorMapping {
		k.SetAppidToAssetCollectorMapping(ctx, item)
	}

	for _, item := range state.CollectorLookup {
		k.SetCollectorLookupTable(ctx, item.AssetRateInfo...)
	}

	for _, item := range state.CollectorAuctionLookupTable {
		k.SetAuctionMappingForApp(ctx, item)
	}

	for _, item := range state.AppToDenomsMapping {
		k.SetAppToDenomsMapping(ctx, item.AppId, item)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	collectorAuctionLookupTable, _ := k.GetAllAuctionMappingForApp(ctx)
	return types.NewGenesisState(
		k.GetAllNetFeeCollectedData(ctx),
		k.GetAllAppidToAssetCollectorMapping(ctx),
		k.GetAllCollectorLookupTable(ctx),
		collectorAuctionLookupTable,
		k.GetAllAppToDenomsMapping(ctx),
		k.GetParams(ctx),
	)
}
