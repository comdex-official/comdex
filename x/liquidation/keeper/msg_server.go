package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/liquidation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) WhitelistApp(goCtx context.Context, id *types.WhitelistAppId) (*types.MsgWhitelistAppIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.WhitelistAppID(ctx, id.AppMappingId); err != nil {
		return nil, err
	}
	return &types.MsgWhitelistAppIdResponse{}, nil
}

func (m msgServer) RemoveWhitelistApp(goCtx context.Context, id *types.RemoveWhitelistAppId) (*types.MsgRemoveWhitelistAppIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.RemoveWhitelistAsset(ctx, id.AppMappingId); err != nil {
		return nil, err
	}
	return &types.MsgRemoveWhitelistAppIdResponse{}, nil
}
