package keeper

import (
	"github.com/comdex-official/comdex/x/nft/exported"
	"github.com/comdex-official/comdex/x/nft/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) error {
	denom := collection.Denom
	creator, err := sdk.AccAddressFromBech32(denom.Creator)
	if err != nil {
		return err
	}
	if err := k.CreateDenom(
		ctx,
		denom.Id,
		denom.Symbol,
		denom.Name,
		denom.Schema,
		creator,
		denom.Description,
		denom.PreviewURI,
	); err != nil {
		return err
	}

	for _, nft := range collection.NFTs {
		metadata := types.Metadata{
			Name:        nft.GetName(),
			Description: nft.GetDescription(),
			MediaURI:    nft.GetMediaURI(),
			PreviewURI:  nft.GetPreviewURI(),
		}

		if err := k.MintNFT(ctx,
			collection.Denom.Id,
			nft.GetID(),
			metadata,
			nft.GetData(),
			nft.IsTransferable(),
			nft.IsExtensible(),
			nft.IsNSFW(),
			nft.RoyaltyShare,
			creator,
			nft.GetOwner(),
		); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) GetCollection(ctx sdk.Context, denomID string) (types.Collection, error) {
	denom, err := k.GetDenom(ctx, denomID)
	if err != nil {
		return types.Collection{}, sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not existed ", denomID)
	}

	nfts := k.GetNFTs(ctx, denomID)
	return types.NewCollection(denom, nfts), nil
}

func (k Keeper) GetCollections(ctx sdk.Context) (cs []types.Collection) {
	for _, denom := range k.GetDenoms(ctx) {
		nfts := k.GetNFTs(ctx, denom.Id)
		cs = append(cs, types.NewCollection(denom, nfts))
	}
	return cs
}

func (k Keeper) GetPaginateCollection(ctx sdk.Context,
	request *types.QueryCollectionRequest, denomId string) (types.Collection, *query.PageResponse, error) {

	denom, err := k.GetDenom(ctx, denomId)
	if err != nil {
		return types.Collection{}, nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denomId %s not existed ", denomId)
	}
	var nfts []exported.NFT
	store := ctx.KVStore(k.storeKey)
	nftStore := prefix.NewStore(store, types.KeyNFT(denomId, ""))
	pagination, err := query.Paginate(nftStore, request.Pagination, func(key []byte, value []byte) error {
		var NFT types.NFT
		k.cdc.MustUnmarshal(value, &NFT)
		nfts = append(nfts, NFT)
		return nil
	})
	if err != nil {
		return types.Collection{}, nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}
	return types.NewCollection(denom, nfts), pagination, nil
}

func (k Keeper) GetTotalSupply(ctx sdk.Context, denomID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyCollection(denomID))
	if len(bz) == 0 {
		return 0
	}
	return types.MustUnMarshalSupply(k.cdc, bz)
}

func (k Keeper) GetTotalSupplyOfOwner(ctx sdk.Context, id string, owner sdk.AccAddress) (supply uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyOwner(owner, id, ""))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		supply++
	}
	return supply
}

func (k Keeper) increaseSupply(ctx sdk.Context, denomID string) {
	supply := k.GetTotalSupply(ctx, denomID)
	supply++

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeyCollection(denomID), bz)
}

func (k Keeper) decreaseSupply(ctx sdk.Context, denomID string) {
	supply := k.GetTotalSupply(ctx, denomID)
	supply--

	store := ctx.KVStore(k.storeKey)
	if supply == 0 {
		store.Delete(types.KeyCollection(denomID))
		return
	}

	bz := types.MustMarshalSupply(k.cdc, supply)
	store.Set(types.KeyCollection(denomID), bz)
}
