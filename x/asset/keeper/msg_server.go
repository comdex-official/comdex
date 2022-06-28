package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ types.MsgServer = (*msgServer)(nil)
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k *msgServer) MsgAddAsset(c context.Context, msg *types.MsgAddAssetRequest) (*types.MsgAddAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

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
	return k.NewAddPair(ctx, msg)
}
