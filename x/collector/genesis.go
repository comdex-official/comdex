package collector

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/petrichormoney/petri/x/collector/keeper"
	"github.com/petrichormoney/petri/x/collector/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	k.SetParams(ctx, state.Params)

	for _, item := range state.NetFeeCollectedData {
		err := k.SetNetFeeCollectedData(ctx, item.AppId, item.AssetId, item.NetFeesCollected)
		if err != nil {
			return
		}
	}

	for _, item := range state.AppIdToAssetCollectorMapping {
		k.SetAppidToAssetCollectorMapping(ctx, item)
	}

	for _, item := range state.CollectorLookup {
		err := k.SetCollectorLookupTable(ctx, item)
		if err != nil {
			return
		}
	}

	for _, item := range state.CollectorAuctionLookupTable {
		k.SetGenAuctionMappingForApp(ctx, item)
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
