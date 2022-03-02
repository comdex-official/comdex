package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ types.MsgServiceServer = (*msgServer)(nil)
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServiceServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k *msgServer) MsgAddAsset(c context.Context, msg *types.MsgAddAssetRequest) (*types.MsgAddAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	/*if msg.From != k.Admin(ctx) {
		return nil, types.ErrorUnauthorized
	}*/

	if k.HasAssetForDenom(ctx, msg.Denom) {
		return nil, types.ErrorDuplicateAsset
	}

	var (
		id    = k.GetAssetID(ctx)
		asset = types.Asset{
			Id:       id + 1,
			Name:     msg.Name,
			Denom:    msg.Denom,
			Decimals: msg.Decimals,
		}
	)

	k.SetAssetID(ctx, asset.Id)
	k.SetAsset(ctx, asset)
	k.SetAssetForDenom(ctx, asset.Denom, asset.Id)

	return &types.MsgAddAssetResponse{}, nil
}

func (k *msgServer) MsgUpdateAsset(c context.Context, msg *types.MsgUpdateAssetRequest) (*types.MsgUpdateAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
/*	if msg.From != k.Admin(ctx) {
		return nil, types.ErrorUnauthorized
	}
*/
	asset, found := k.GetAsset(ctx, msg.Id)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	if msg.Name != "" {
		asset.Name = msg.Name
	}
	if msg.Denom != "" {
		if k.HasAssetForDenom(ctx, msg.Denom) {
			return nil, types.ErrorDuplicateAsset
		}

		asset.Denom = msg.Denom

		k.DeleteAssetForDenom(ctx, asset.Denom)
		k.SetAssetForDenom(ctx, asset.Denom, asset.Id)
	}
	if msg.Decimals >= 0 {
		asset.Decimals = msg.Decimals
	}

	k.SetAsset(ctx, asset)
	return &types.MsgUpdateAssetResponse{}, nil
}

func (k *msgServer) MsgAddPair(c context.Context, msg *types.MsgAddPairRequest) (*types.MsgAddPairResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
/*	if msg.From != k.Admin(ctx) {
		return nil, types.ErrorUnauthorized
	}
*/
	if !k.HasAsset(ctx, msg.AssetIn) {
		return nil, types.ErrorAssetDoesNotExist
	}
	if !k.HasAsset(ctx, msg.AssetOut) {
		return nil, types.ErrorAssetDoesNotExist
	}

	var (
		id   = k.GetPairID(ctx)
		pair = types.Pair{
			Id:               id + 1,
			AssetIn:          msg.AssetIn,
			AssetOut:         msg.AssetOut,
			LiquidationRatio: msg.LiquidationRatio,
		}
	)

	k.SetPairID(ctx, pair.Id)
	k.SetPair(ctx, pair)

	return &types.MsgAddPairResponse{}, nil
}

func (k *msgServer) MsgUpdatePair(c context.Context, msg *types.MsgUpdatePairRequest) (*types.MsgUpdatePairResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
/*	if msg.From != k.Admin(ctx) {
		return nil, types.ErrorUnauthorized
	}
*/
	pair, found := k.GetPair(ctx, msg.Id)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	if !msg.LiquidationRatio.IsZero() {
		pair.LiquidationRatio = msg.LiquidationRatio
	}

	k.SetPair(ctx, pair)
	return &types.MsgUpdatePairResponse{}, nil
}
