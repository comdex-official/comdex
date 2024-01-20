package v9

import (
	"context"
	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func UpdateDenomAndAddAsset(ctx sdk.Context, assetKeeper assetkeeper.Keeper) {
	asset, found := assetKeeper.GetAsset(ctx, 17)
	if found {
		asset.Denom = "ibc/50EF138042B553362774A8A9DE967F610E52CAEB3BA864881C9A1436DED98075"
		assetKeeper.SetAsset(ctx, asset)
	}

	assetGDAI := assettypes.Asset{Name: "GDAI", Denom: "ibc/109DD45CF4093BEB472784A0C5B5F4643140900020B74B102B842A4BE2AE45DA", Decimals: sdkmath.NewInt(1000000000000000000), IsOnChain: false, IsOraclePriceRequired: true, IsCdpMintable: false}

	err := assetKeeper.AddAssetRecords(ctx, assetGDAI)
	if err != nil {
		return
	}
	getGDAI, found := assetKeeper.GetAssetForDenom(ctx, "ibc/109DD45CF4093BEB472784A0C5B5F4643140900020B74B102B842A4BE2AE45DA")
	if !found {
		return
	}
	var (
		id       = assetKeeper.GetPairID(ctx)
		pairGDAI = assettypes.Pair{
			Id:       id + 1,
			AssetIn:  getGDAI.Id,
			AssetOut: 3,
		}
	)

	assetKeeper.SetPairID(ctx, pairGDAI.Id)
	assetKeeper.SetPair(ctx, pairGDAI)

	getSTOSMO, found := assetKeeper.GetAssetForDenom(ctx, "ibc/CC482813CC038C614C2615A997621EA5E605ADCCD4040B83B0468BD72533A165")
	if !found {
		return
	}

	var (
		id2        = assetKeeper.GetPairID(ctx)
		pairSTOSMO = assettypes.Pair{
			Id:       id2 + 1,
			AssetIn:  getSTOSMO.Id,
			AssetOut: 3,
		}
	)

	assetKeeper.SetPairID(ctx, pairSTOSMO.Id)
	assetKeeper.SetPair(ctx, pairSTOSMO)
}

func CreateUpgradeHandlerV900(
	mm *module.Manager,
	configurator module.Configurator,
	assetKeeper assetkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		UpdateDenomAndAddAsset(ctx, assetKeeper)
		return vm, err
	}
}
