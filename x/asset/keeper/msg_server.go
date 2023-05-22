package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// AddAsset defines a method to add asset.
func (m msgServer) AddAsset(goCtx context.Context, msg *types.MsgAddAsset) (*types.MsgAddAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.AddAsset(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgAddAssetResponse{}, nil
}
