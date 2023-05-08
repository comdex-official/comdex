package liquidationsV2

import (
	"github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	server := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgLiquidateInternalKeeperRequest:
			res, err := server.MsgLiquidateInternalKeeper(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(types.ErrorUnknownMsgType, "%T", msg)
		}
	}
}

func NewLiquidationsV2Handler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.WhitelistLiquidationProposal:
			return handleWhitelistLiquidationProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
		}
	}
}

func handleWhitelistLiquidationProposal(ctx sdk.Context, k keeper.Keeper, p *types.WhitelistLiquidationProposal) error {
	return k.HandleWhitelistLiquidationProposal(ctx, p)
}
