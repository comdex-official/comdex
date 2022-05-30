package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) MsgPlaceSurplusBid(goCtx context.Context, msg *types.MsgPlaceSurplusBidRequest) (*types.MsgPlaceSurplusBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	bidder, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return nil, err
	}
	err = k.PlaceSurplusBid(ctx, msg.AppId, msg.AuctionMappingId, msg.AuctionId, bidder, msg.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgPlaceSurplusBidResponse{}, nil
}

func (k msgServer) MsgPlaceDebtBid(goCtx context.Context, msg *types.MsgPlaceDebtBidRequest) (*types.MsgPlaceDebtBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	bidder, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return nil, err
	}
	err = k.PlaceDebtBid(ctx, msg.AppId, msg.AuctionMappingId, msg.AuctionId, bidder, msg.Bid, msg.ExpectedUserToken)
	if err != nil {
		return nil, err
	}
	return &types.MsgPlaceDebtBidResponse{}, nil
}

func (k msgServer) MsgPlaceDutchBid(goCtx context.Context, msg *types.MsgPlaceDutchBidRequest) (*types.MsgPlaceDutchBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	bidder, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return nil, err
	}
	err = k.PlaceDutchBid(ctx, msg.AppId, msg.AuctionMappingId, msg.AuctionId, bidder, msg.Amount, msg.Max)
	if err != nil {
		return nil, err
	}
	return &types.MsgPlaceDutchBidResponse{}, nil
}
