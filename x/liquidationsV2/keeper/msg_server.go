package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) MsgLiquidateInternalKeeper(c context.Context, req *types.MsgLiquidateInternalKeeperRequest) (*types.MsgLiquidateInternalKeeperResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := m.keeper.MsgLiquidate(ctx, req.From, req.LiqType, req.Id); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLiquidateInternalKeeper,
			sdk.NewAttribute(types.AttributeKeyLiqType, strconv.FormatUint(req.LiqType, 10)),
			sdk.NewAttribute(types.AttributeKeyID, strconv.FormatUint(req.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, req.From),
		),
	})
	return &types.MsgLiquidateInternalKeeperResponse{}, nil
}
