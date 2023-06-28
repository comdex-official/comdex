package keeper

import (
	"github.com/comdex-official/comdex/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetOwner gets all the ID collections owned by an address and denom ID

func (k Keeper) deleteOwner(ctx sdk.Context, denomID, nftId string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyOwner(owner, denomID, nftId))
}

func (k Keeper) setOwner(ctx sdk.Context,
	denomId, nftId string,
	owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	bz := types.MustMarshalNFTID(k.cdc, nftId)
	store.Set(types.KeyOwner(owner, denomId, nftId), bz)
}

func (k Keeper) swapOwner(ctx sdk.Context, denomID, tokenID string, srcOwner, dstOwner sdk.AccAddress) {
	k.deleteOwner(ctx, denomID, tokenID, srcOwner)
	k.setOwner(ctx, denomID, tokenID, dstOwner)
}
