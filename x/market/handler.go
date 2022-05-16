package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/comdex-official/comdex/x/market/keeper"
	"github.com/comdex-official/comdex/x/market/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	server := keeper.NewMsgServiceServer(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgAddMarketRequest:
			res, err := server.MsgAddMarket(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateMarketRequest:
			res, err := server.MsgUpdateMarket(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveMarketForAssetRequest:
			res, err := server.MsgRemoveMarketForAsset(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, errors.Wrapf(types.ErrorUnknownMsgType, "%T", msg)
		}
	}
}

func NewProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.UpdateAdminProposal:
			return k.HandleUpdateAdminProposal(ctx, c)
		default:
			return errors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
		}
	}
}
