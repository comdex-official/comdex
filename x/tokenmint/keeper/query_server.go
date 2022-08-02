package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/tokenmint/types"
)

var (
	_ types.QueryServer = QueryServer{}
)

type QueryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &QueryServer{
		Keeper: k,
	}
}

func (q QueryServer) QueryAllTokenMintedForAllApps(c context.Context, req *types.QueryAllTokenMintedForAllAppsRequest) (*types.QueryAllTokenMintedForAllAppsResponse, error) {
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	totalMintedData := q.GetTotalTokenMinted(ctx)
	return &types.QueryAllTokenMintedForAllAppsResponse{
		TokenMint: totalMintedData,
	}, nil
}

func (q QueryServer) QueryTokenMintedByApp(c context.Context, req *types.QueryTokenMintedByAppRequest) (*types.QueryTokenMintedByAppResponse, error) {
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

func (q QueryServer) QueryTokenMintedByAppAndAsset(c context.Context, req *types.QueryTokenMintedByAppAndAssetRequest) (*types.QueryTokenMintedByAppAndAssetResponse, error) {
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
