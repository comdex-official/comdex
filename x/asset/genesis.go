package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/petrichormoney/petri/x/asset/keeper"
	"github.com/petrichormoney/petri/x/asset/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	var (
		assetID        uint64
		pairID         uint64
		appID          uint64
		extendedPairID uint64
	)

	k.SetParams(ctx, state.Params)

	for _, item := range state.Assets {
		if item.Id > assetID {
			assetID = item.Id
		}
		k.SetAssetForDenom(ctx, item.Denom, item.Id)
		k.SetAssetForName(ctx, item.Name, item.Id)
		k.SetAsset(ctx, item)
	}

	for _, item := range state.Pairs {
		if item.Id > pairID {
			pairID = item.Id
		}

		k.SetPair(ctx, item)
	}

	for _, item := range state.AppData {
		if item.Id > appID {
			appID = item.Id
		}
		k.SetAppForShortName(ctx, item.ShortName, item.Id)
		k.SetAppForName(ctx, item.Name, item.Id)
		k.SetApp(ctx, item)
	}

	for _, item := range state.ExtendedPairVault {
		if item.Id > extendedPairID {
			extendedPairID = item.Id
		}

		k.SetPairsVault(ctx, item)
	}

	k.SetAssetID(ctx, assetID)
	k.SetPairID(ctx, pairID)
	k.SetAppID(ctx, appID)
	k.SetPairsVaultID(ctx, extendedPairID)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	apps, _ := k.GetApps(ctx)
	pairVaults, _ := k.GetPairsVaults(ctx)
	return types.NewGenesisState(
		k.GetAssets(ctx),
		k.GetPairs(ctx),
		apps,
		pairVaults,
		k.GetParams(),
	)
}
