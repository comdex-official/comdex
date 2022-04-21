package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) IsWhitelistedAsset(ctx sdk.Context, tokenDenom string) bool {
	//store := ctx.KVStore(k.storeKey)
	//key := types.CreateRegisteredTokenKey(tokenDenom)
	//
	//return store.Has(key)
	return true
}
