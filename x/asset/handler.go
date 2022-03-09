package asset

import (
	"github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	server := keeper.NewMsgServiceServer(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgAddAssetRequest:
			res, err := server.MsgAddAsset(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateAssetRequest:
			res, err := server.MsgUpdateAsset(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddPairRequest:
			res, err := server.MsgAddPair(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdatePairRequest:
			res, err := server.MsgUpdatePair(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, errors.Wrapf(types.ErrorUnknownMsgType, "%T", msg)
		}
	}
}

func NewUpdateLiquidationRatioProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddAssetsProposal:
			return handleUpdateLiquidationRatioProposal(ctx, k, c)
		default:
			return errors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
		}
	}
}

func handleUpdateLiquidationRatioProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddAssetsProposal) error {
	return k.HandleProposalAddAsset(ctx, p)
}
