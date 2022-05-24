package rewards

import (
	"github.com/comdex-official/comdex/x/rewards/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	rewards := k.GetRewards(ctx)
	for _, v := range rewards {
		appId := v.App_mapping_ID
		assetIds := v.Asset_ID
		err := k.Iterate(ctx, appId, assetIds)
		if err != nil {
			return
		}
	}
	AppIdsVault := k.GetAppIds(ctx).WhitelistedAppMappingIdsVaults
	for i := range AppIdsVault {
		err := k.IterateVaults(ctx, AppIdsVault[i])
		if err != nil {
			return
		}
	}
}
