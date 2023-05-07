package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/auctionsV2/types"
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

func (k msgServer) MsgPlaceMarketBid(goCtx context.Context, msg *types.MsgPlaceMarketBidRequest) (*types.MsgPlaceMarketBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	bidder, err := sdk.AccAddressFromBech32(msg.Bidder)
	if err != nil {
		return nil, err
	}
	auctionData, err := k.GetAuction(ctx, msg.AuctionId)
	if err != nil {

		return nil, err
	}
	//From auction ID, checking the whether its an english or a dutch auction
	//If true triggering Dutch Auction Bid Request
	if auctionData.AuctionType {

		err = k.PlaceDutchAuctionBid(ctx, msg.AuctionId, bidder, msg.Amount, auctionData)
		if err != nil {
			return nil, err
		}

	} else {
		//Else ENGLISH - triggering English Auction Bid Request
		err = k.PlaceEnglishAuctionBid(ctx, msg.AuctionId, bidder, msg.Amount, auctionData)
		if err != nil {
			return nil, err
		}
	}

	// ctx.GasMeter().ConsumeGas(types.DutchBidGas, "DutchBidGas")
	return &types.MsgPlaceMarketBidResponse{}, nil
}
