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

func (q *queryServer) QueryAllTokenMintedForAllApps(c context.Context, req *types.QueryAllTokenMintedForAllAppsRequest) (*types.QueryAllTokenMintedForAllAppsResponse, error) {
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	totalMintedData := q.GetTotalTokenMinted(ctx)
	return &types.QueryAllTokenMintedForAllAppsResponse{
		TokenMint: totalMintedData,
	}, nil
}

func (q *queryServer) QueryTokenMintedByApp(c context.Context, req *types.QueryTokenMintedByAppRequest) (*types.QueryTokenMintedByAppResponse, error) {
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

func (q *queryServer) QueryTokenMintedByAppAndAsset(c context.Context, req *types.QueryTokenMintedByAppAndAssetRequest) (*types.QueryTokenMintedByAppAndAssetResponse, error) {
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
