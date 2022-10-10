package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/auction/types"
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
	err = k.PlaceSurplusAuctionBid(ctx, msg.AppId, msg.AuctionMappingId, msg.AuctionId, bidder, msg.Amount)
	if err != nil {
		return nil, err
	}
	ctx.GasMeter().ConsumeGas(types.SurplusBidGas, "SurplusBidGas")
	return &types.MsgPlaceSurplusBidResponse{}, nil
}

func (k msgServer) MsgPlaceDebtBid(goCtx context.Context, msg *types.MsgPlaceDebtBidRequest) (*types.MsgPlaceDebtBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	bidder, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return nil, err
	}
	err = k.PlaceDebtAuctionBid(ctx, msg.AppId, msg.AuctionMappingId, msg.AuctionId, bidder, msg.Bid, msg.ExpectedUserToken)
	if err != nil {
		return nil, err
	}
	ctx.GasMeter().ConsumeGas(types.DebtBidGas, "DebtBidGas")
	return &types.MsgPlaceDebtBidResponse{}, nil
}

func (k msgServer) MsgPlaceDutchBid(goCtx context.Context, msg *types.MsgPlaceDutchBidRequest) (*types.MsgPlaceDutchBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	bidder, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return nil, err
	}
	err = k.PlaceDutchAuctionBid(ctx, msg.AppId, msg.AuctionMappingId, msg.AuctionId, bidder, msg.Amount, msg.Max)
	if err != nil {
		return nil, err
	}
	ctx.GasMeter().ConsumeGas(types.DutchBidGas, "DutchBidGas")
	return &types.MsgPlaceDutchBidResponse{}, nil
}
