package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/comdex-official/comdex/x/nft/types"
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

func (m msgServer) CreateDenom( goCtx context.Context, msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.CreateDenom(ctx, msg.Id, msg.Symbol, msg.Name, msg.Schema, sender, msg.Description, msg.PreviewURI,); err != nil {
		return nil, err
	}

	return &types.MsgCreateDenomResponse{}, nil
}
