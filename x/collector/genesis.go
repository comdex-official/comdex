package collector

import (
	"github.com/comdex-official/comdex/x/collector/keeper"
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		err := k.SetAuctionMappingForApp(ctx, item)
		if err != nil {
			return
		}
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
