package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ types.QueryServer = &Keeper{}

type QueryServer struct {
	Keeper
}

func (k Keeper) Collection(ctx context.Context, request *types.QueryCollectionRequest) (*types.QueryCollectionResponse, error) {
	return nil, nil
}

func (k Keeper) Denom(ctx context.Context, request *types.QueryDenomRequest) (*types.QueryDenomResponse, error) {
	return nil, nil
}

func (k Keeper) Denoms(ctx context.Context, request *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	return nil, nil
}

func (k Keeper) NFT(c context.Context, req *types.QueryNFTRequest) (*types.QueryNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	nft, err := k.GetNFT(ctx, req.DenomId, req.Id)
	if err != nil {
		return nil, err
	}
	NFT, ok := nft.(types.NFT)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrInvalidNFT, "invalid nft for denom-id %s for nft-id %s", req.DenomId, req.Id)
	}
	return &types.QueryNFTResponse{NFT: &NFT}, nil
}

func (k Keeper) OwnerNFTs(c context.Context, req *types.QueryOwnerNFTsRequest) (*types.QueryOwnerNFTsResponse, error) {

	ctx := sdk.UnwrapSDKContext(c)
	ownerNfts := k.GetOwnerNFTs(ctx, req.DenomId, req.Owner)

	var ownerNFTsCollection types.OwnerNFTCollection

	ownerNFTsCollection = types.OwnerNFTCollection{
		Denom: types.Denom{},
		Nfts:  ownerNfts,
	}

	return &types.QueryOwnerNFTsResponse{
		Owner:       req.Owner,
		Collections: ownerNFTsCollection,
		Pagination:  nil,
	}, nil
}

func (k Keeper) Supply(ctx context.Context, request *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	return nil, nil
}
