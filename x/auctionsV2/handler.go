package auctionsV2

import (
	"github.com/comdex-official/comdex/x/auctionsV2/keeper"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	server := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgPlaceMarketBidRequest:
			res, err := server.MsgPlaceMarketBid(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgDepositLimitBidRequest:
			res, err := server.MsgDepositLimitBid(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgCancelLimitBidRequest:
			res, err := server.MsgCancelLimitBid(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgWithdrawLimitBidRequest:
			res, err := server.MsgWithdrawLimitBid(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "%T", msg)
		}
	}
}

func NewAuctionsV2Handler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.DutchAutoBidParamsProposal:
			return handleAddAuctionParamsProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
		}
	}
}

func handleAddAuctionParamsProposal(ctx sdk.Context, k keeper.Keeper, p *types.DutchAutoBidParamsProposal) error {
	return k.HandleAuctionParamsProposal(ctx, p)
}
