package keeper

import (
	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) IsWhitelistedAsset(ctx sdk.Context, tokenDenom string) bool {
	//store := ctx.KVStore(k.storeKey)
	//key := types.CreateRegisteredTokenKey(tokenDenom)
	//
	//return store.Has(key)
	return true
}

func (k *Keeper) HasAssetForDenom(ctx sdk.Context, denom string) bool {
	var (
		store = k.Store(ctx)
		key   = types.AssetForDenomKey(denom)
	)

	return store.Has(key)
}
