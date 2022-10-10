package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	"github.com/comdex-official/comdex/x/asset/types"
	bandoraclemoduletypes "github.com/comdex-official/comdex/x/bandoracle/types"
)

func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scoped.AuthenticateCapability(ctx, cap, name)
}

func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scoped.ClaimCapability(ctx, cap, name)
}

func (k Keeper) HasAsset(ctx sdk.Context, id uint64) bool {
	return k.assetKeeper.HasAsset(ctx, id)
}

func (k Keeper) GetAssets(ctx sdk.Context) (assets []types.Asset) {
	return k.assetKeeper.GetAssets(ctx)
}

func (k Keeper) GetAsset(ctx sdk.Context, id uint64) (asset types.Asset, found bool) {
	return k.assetKeeper.GetAsset(ctx, id)
}

func (k Keeper) GetLastFetchPriceID(ctx sdk.Context) int64 {
	return k.bandoraclekeeper.GetLastFetchPriceID(ctx)
}

func (k Keeper) GetLastBlockheight(ctx sdk.Context) int64 {
	return k.bandoraclekeeper.GetLastBlockHeight(ctx)
}

func (k Keeper) GetFetchPriceMsg(ctx sdk.Context) bandoraclemoduletypes.MsgFetchPriceData {
	return k.bandoraclekeeper.GetFetchPriceMsg(ctx)
}

func (k Keeper) GetFetchPriceResult(ctx sdk.Context, id bandoraclemoduletypes.OracleRequestID) (bandoraclemoduletypes.FetchPriceResult, error) {
	return k.bandoraclekeeper.GetFetchPriceResult(ctx, id)
}

func (k Keeper) GetCheckFlag(ctx sdk.Context) bool {
	return k.bandoraclekeeper.GetCheckFlag(ctx)
}

func (k Keeper) SetCheckFlag(ctx sdk.Context, flag bool) {
	k.bandoraclekeeper.SetCheckFlag(ctx, flag)
}

func (k Keeper) GetOracleValidationResult(ctx sdk.Context) bool {
	return k.bandoraclekeeper.GetOracleValidationResult(ctx)
}
