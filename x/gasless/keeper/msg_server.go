package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/gasless/types"
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

// CreateGasProvider defines a method to create a new gas provider
func (m msgServer) CreateGasProvider(goCtx context.Context, msg *types.MsgCreateGasProvider) (*types.MsgCreateGasProviderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.CreateGasProvider(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCreateGasProviderResponse{}, nil
}
