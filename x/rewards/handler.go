package rewards

import (
	"fmt"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/x/rewards/keeper"
	"github.com/comdex-official/comdex/x/rewards/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	server := keeper.NewMsgServerImpl(k)
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateGauge:
			res, err := server.CreateGauge(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.ActivateExternalRewardsLockers:
			res, err := server.ExternalRewardsLockers(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.ActivateExternalRewardsVault:
			res, err := server.ExternalRewardsVault(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func NewAddRewardsProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddLendExternalRewardsProposal:
			return handleAddLendRewardsProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
		}
	}
}

func handleAddLendRewardsProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddLendExternalRewardsProposal) error {
	return k.HandleProposalAddLendRewards(ctx, p)
}
