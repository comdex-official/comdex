package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/nft/types"
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

func (m msgServer) CreateDenom(goCtx context.Context,
	msg *types.MsgCreateDenom) (*types.MsgCreateDenomResponse, error) {

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.CreateDenom(ctx,
		msg.Id,
		msg.Symbol,
		msg.Name,
		msg.Schema,
		sender,
		msg.Description,
		msg.PreviewURI,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventCreateDenom{
			Id:      msg.Id,
			Symbol:  msg.Symbol,
			Name:    msg.Name,
			Creator: msg.Sender,
		},
	)

	return &types.MsgCreateDenomResponse{}, nil
}

func (m msgServer) UpdateDenom(goCtx context.Context, msg *types.MsgUpdateDenom) (*types.MsgUpdateDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.UpdateDenom(ctx, msg.Id, msg.Name, msg.Description, msg.PreviewURI, sender)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventUpdateDenom{
			Id:      msg.Id,
			Name:    msg.Name,
			Creator: msg.Sender,
		},
	)

	return &types.MsgUpdateDenomResponse{}, nil
}

func (m msgServer) TransferDenom(goCtx context.Context, msg *types.MsgTransferDenom) (*types.MsgTransferDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.TransferDenomOwner(ctx, msg.Id, sender, recipient)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitTypedEvent(
		&types.EventTransferDenom{
			Id:        msg.Id,
			Sender:    msg.Sender,
			Recipient: msg.Recipient,
		},
	)

	return &types.MsgTransferDenomResponse{}, nil
}

func (m msgServer) MintNFT(goCtx context.Context,
	msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.MintNFT(ctx,
		msg.DenomId,
		msg.Id,
		msg.Metadata,
		msg.Data,
		msg.Transferable, // to mint a non-transferable NFT (optional, default is false) 
		msg.Extensible, // to mint an inextensible NFT (optional, default is false) 
		msg.Nsfw, // to mark the NFT as not safe for work (optional, default is false)
		msg.RoyaltyShare,
		sender,
		recipient,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventMintNFT{
			Id:      msg.Id,
			DenomId: msg.DenomId,
			URI:     msg.Metadata.MediaURI,
			Owner:   msg.Recipient,
		},
	)

	return &types.MsgMintNFTResponse{}, nil
}

func (m msgServer) TransferNFT(goCtx context.Context,
	msg *types.MsgTransferNFT) (*types.MsgTransferNFTResponse, error) {

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.TransferOwnership(ctx, msg.DenomId, msg.Id,
		sender,
		recipient,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventTransferNFT{
			Id:        msg.Id,
			DenomId:   msg.DenomId,
			Sender:    msg.Sender,
			Recipient: msg.Recipient,
		},
	)

	return &types.MsgTransferNFTResponse{}, nil
}

func (m msgServer) BurnNFT(goCtx context.Context,
	msg *types.MsgBurnNFT) (*types.MsgBurnNFTResponse, error) {

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.BurnNFT(ctx, msg.DenomId, msg.Id, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventBurnNFT{
			Id:      msg.Id,
			DenomId: msg.DenomId,
			Owner:   msg.Sender,
		},
	)

	return &types.MsgBurnNFTResponse{}, nil
}
