package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/asset/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	var (
		assetID uint64 = 0
		pairID  uint64 = 0
	)

	k.SetParams(ctx, state.Params)

	for _, item := range state.Assets {
		if item.Id > assetID {
			assetID = item.Id
		}

		k.SetAsset(ctx, item)
	}

	for _, item := range state.Pairs {
		if item.Id > assetID {
			pairID = item.Id
		}

		k.SetPair(ctx, item)
	}

	k.SetAssetID(ctx, assetID)
	k.SetPairID(ctx, pairID)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetAssets(ctx),
		k.GetPairs(ctx),
		k.GetParams(ctx),
	)
}
