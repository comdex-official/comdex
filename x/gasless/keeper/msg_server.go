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

// AuthorizeActors defines a method to update the actors in gas provider
func (m msgServer) AuthorizeActors(goCtx context.Context, msg *types.MsgAuthorizeActors) (*types.MsgAuthorizeActorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.AuthorizeActors(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgAuthorizeActorsResponse{}, nil
}

// UpdateGasProviderStatus defines a method to update the active status of gas provider
func (m msgServer) UpdateGasProviderStatus(goCtx context.Context, msg *types.MsgUpdateGasProviderStatus) (*types.MsgUpdateGasProviderStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.UpdateGasProviderStatus(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgUpdateGasProviderStatusResponse{}, nil
}

// UpdateGasProviderConfigs defines a method to update a gas provider
func (m msgServer) UpdateGasProviderConfigs(goCtx context.Context, msg *types.MsgUpdateGasProviderConfig) (*types.MsgUpdateGasProviderConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.UpdateGasProviderConfig(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgUpdateGasProviderConfigResponse{}, nil
}

// BlockConsumer defines a method to block a gas consumer
func (m msgServer) BlockConsumer(goCtx context.Context, msg *types.MsgBlockConsumer) (*types.MsgBlockConsumerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.BlockConsumer(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgBlockConsumerResponse{}, nil
}

// UnblockConsumer defines a method to unblock a consumer
func (m msgServer) UnblockConsumer(goCtx context.Context, msg *types.MsgUnblockConsumer) (*types.MsgUnblockConsumerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.UnblockConsumer(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgUnblockConsumerResponse{}, nil
}

// UpdateGasConsumerLimit defines a method to increase consumption limit for a consumer
func (m msgServer) UpdateGasConsumerLimit(goCtx context.Context, msg *types.MsgUpdateGasConsumerLimit) (*types.MsgUpdateGasConsumerLimitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.UpdateGasConsumerLimit(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgUpdateGasConsumerLimitResponse{}, nil
}
