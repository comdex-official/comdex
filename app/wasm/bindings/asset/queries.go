package asset

import (
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QueryPlugin struct {
	assetKeeper *assetKeeper.Keeper
}

func NewQueryPlugin(
	assetKeeper *assetKeeper.Keeper,
) *QueryPlugin {
	return &QueryPlugin{
		assetKeeper: assetKeeper,
	}
}

func (qp QueryPlugin) GetAppInfo(ctx sdk.Context, appMappingId uint64) (assetTypes.AppMapping, error) {
	appData, err := qp.assetKeeper.GetApp(ctx, appMappingId)
	if err != true {
		return appData, nil
	}
	return appData, nil
}
