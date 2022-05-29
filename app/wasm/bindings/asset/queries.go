package asset

import (
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
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

func (qp QueryPlugin) GetAppInfo(ctx sdk.Context, appMappingId uint64) (int64, int64, uint64, error) {
	MinGovDeposit, GovTimeInSeconds, AssetId, err := qp.assetKeeper.GetAppWasmQuery(ctx, appMappingId)
	if err != nil {
		return MinGovDeposit, GovTimeInSeconds, AssetId, nil
	}
	return MinGovDeposit, GovTimeInSeconds, AssetId, nil
}

func (qp QueryPlugin) GetAssetInfo(ctx sdk.Context, Id uint64) (string, error) {
	assetDenom := qp.assetKeeper.GetAssetDenom(ctx, Id)
	return assetDenom, nil
}
