package keeper

import (
	"context"
	"fmt"

	"github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServiceServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServiceServer = msgServer{}

func (k msgServer) MsgPlaceBid(goCtx context.Context, msg *types.MsgPlaceBidRequest) (*types.MsgPlaceBidResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)

	bidder, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return nil, err
	}
	fmt.Println(bidder)
	fmt.Println(msg.Amount)
	fmt.Println(msg.AuctionId)
	// err = k.keeper.PlaceBid(ctx, msg.AuctionId, bidder, msg.Amount)
	// if err != nil {
	// 	return nil, err
	// }
	return &types.MsgPlaceBidResponse{}, nil
}
