package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/tokenmint/types"
)

var (
	_ types.QueryServer = (*QueryServer)(nil)
)

type QueryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &QueryServer{
		Keeper: k,
	}
}

func (q *QueryServer) QueryAllTokenMintedForAllProducts(c context.Context, req *types.QueryAllTokenMintedForAllProductsRequest) (*types.QueryAllTokenMintedForAllProductsResponse, error) {
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	totalMintedData := q.GetTotalTokenMinted(ctx)
	return &types.QueryAllTokenMintedForAllAppsResponse{
		TokenMint: totalMintedData,
	}, nil
}

func (q *QueryServer) QueryTokenMintedByProduct(c context.Context, req *types.QueryTokenMintedByProductRequest) (*types.QueryTokenMintedByProductResponse, error) {
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	tokenMint, found := q.GetTokenMint(ctx, req.AppId)
	if !found {
		return nil, types.ErrorMintDataNotFound
	}

	return &types.QueryTokenMintedByAppResponse{
		TokenMint: tokenMint,
	}, nil
}

func (q *QueryServer) QueryTokenMintedByProductAndAsset(c context.Context, req *types.QueryTokenMintedByProductAndAssetRequest) (*types.QueryTokenMintedByProductAndAssetResponse, error) {
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	tokenMint, found := q.GetAssetDataInTokenMintByApp(ctx, req.AppId, req.AssetId)
	if !found {
		return nil, types.ErrorMintDataNotFound
	}

	return &types.QueryTokenMintedByAppAndAssetResponse{
		MintedTokens: tokenMint,
	}, nil
}
