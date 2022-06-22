package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/tokenmint/types"
)

var (
	_ types.QueryServer = (*queryServer)(nil)
)

type queryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServer{
		Keeper: k,
	}
}

func (q *queryServer) QueryAllTokenMintedForAllProducts(c context.Context, req *types.QueryAllTokenMintedForAllProductsRequest) (*types.QueryAllTokenMintedForAllProductsResponse, error) {
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	totalMintedData := q.GetTotalTokenMinted(ctx)
	return &types.QueryAllTokenMintedForAllProductsResponse{
		TokenMint: totalMintedData,
	}, nil
}

func (q *queryServer) QueryTokenMintedByProduct(c context.Context, req *types.QueryTokenMintedByProductRequest) (*types.QueryTokenMintedByProductResponse, error) {
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	tokenMint, found := q.GetTokenMint(ctx, req.AppId)
	if !found {
		return nil, types.ErrorMintDataNotFound
	}

	return &types.QueryTokenMintedByProductResponse{
		TokenMint: tokenMint,
	}, nil
}

func (q *queryServer) QueryTokenMintedByProductAndAsset(c context.Context, req *types.QueryTokenMintedByProductAndAssetRequest) (*types.QueryTokenMintedByProductAndAssetResponse, error) {
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	tokenMint, found := q.GetAssetDataInTokenMintByApp(ctx, req.AppId, req.AssetId)
	if !found {
		return nil, types.ErrorMintDataNotFound
	}

	return &types.QueryTokenMintedByProductAndAssetResponse{
		MintedTokens: tokenMint,
	}, nil
}
