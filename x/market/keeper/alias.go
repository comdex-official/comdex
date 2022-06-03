package keeper

import (
	"github.com/comdex-official/comdex/x/asset/types"
	bandoraclemoduletypes "github.com/comdex-official/comdex/x/bandoracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
)

func (k *Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scoped.AuthenticateCapability(ctx, cap, name)
}

func (k *Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scoped.ClaimCapability(ctx, cap, name)
}

func (k *Keeper) HasAsset(ctx sdk.Context, id uint64) bool {
	return k.assetKeeper.HasAsset(ctx, id)
}

func (k *Keeper) GetAssetsForOracle(ctx sdk.Context) (assets []types.Asset) {
	return k.assetKeeper.GetAssetsForOracle(ctx)
}

func (k *Keeper) GetLastFetchPriceID(ctx sdk.Context) int64 {
	return k.bandoraclekeeper.GetLastFetchPriceID(ctx)
}

func (k *Keeper) GetLastBlockheight(ctx sdk.Context) int64 {
	return k.bandoraclekeeper.GetLastBlockheight(ctx)
}

func (k *Keeper) GetFetchPriceMsg(ctx sdk.Context) bandoraclemoduletypes.MsgFetchPriceData {
	return k.bandoraclekeeper.GetFetchPriceMsg(ctx)
}
