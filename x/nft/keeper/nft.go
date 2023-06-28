package keeper

import (
	"github.com/comdex-official/comdex/x/nft/exported"
	"github.com/comdex-official/comdex/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetNFT(ctx sdk.Context, denomID, nftID string) (nft exported.NFT, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyNFT(denomID, nftID))
	if bz == nil {
		return nil, sdkerrors.Wrapf(types.ErrUnknownCollection, "not found NFT: %s", denomID)
	}

	var NFT types.NFT
	k.cdc.MustUnmarshal(bz, &NFT)
	return NFT, nil
}

func (k Keeper) GetNFTs(ctx sdk.Context, denom string) (nfts []exported.NFT) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyNFT(denom, ""))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var NFT types.NFT
		k.cdc.MustUnmarshal(iterator.Value(), &NFT)
		nfts = append(nfts, NFT)
	}
	return nfts
}

func (k Keeper) GetOwnerNFTs(ctx sdk.Context, denom string, owner string) (nfts []*types.NFT) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyNFT(denom, ""))
	defer iterator.Close()
	var nftList []*types.NFT
	for ; iterator.Valid(); iterator.Next() {
		var NFT types.NFT
		k.cdc.MustUnmarshal(iterator.Value(), &NFT)
		if NFT.Owner == owner {
			nftList = append(nftList, &NFT)
		}
	}
	return nftList
}

func (k Keeper) Authorize(ctx sdk.Context, denomID, nftID string, owner sdk.AccAddress) (types.NFT, error) {
	nft, err := k.GetNFT(ctx, denomID, nftID)
	if err != nil {
		return types.NFT{}, err
	}

	if !owner.Equals(nft.GetOwner()) {
		return types.NFT{}, sdkerrors.Wrap(types.ErrUnauthorized, owner.String())
	}
	return nft.(types.NFT), nil
}

func (k Keeper) HasNFT(ctx sdk.Context, denomID, nftID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyNFT(denomID, nftID))
}

func (k Keeper) setNFT(ctx sdk.Context, denomID string, nft types.NFT) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&nft)
	store.Set(types.KeyNFT(denomID, nft.GetID()), bz)
}

func (k Keeper) deleteNFT(ctx sdk.Context, denomID string, nft exported.NFT) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyNFT(denomID, nft.GetID()))
}
