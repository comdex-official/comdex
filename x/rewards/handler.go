package rewards

import (
	"github.com/comdex-official/comdex/x/rewards/keeper"
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	server := keeper.NewMsgServiceServer(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgDepositMintingRewardAmountRequest:
			res, err := server.MsgDepositMintingRewards(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateMintRewardStartTimeRequest:
			res, err := server.MsgUpdateMintRewardStartTime(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(types.ErrorUnknownMsgType, "%T", msg)
		}
	}
}

func NewRewardsProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.NewMintRewardsProposal:
			return handleNewMintRewardsProposal(ctx, k, c)
		case *types.DisbaleMintRewardsProposal:
			return handleDisableMintRewardsProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized Minting rewards proposal content type: %T", c)
		}
	}
}

func handleNewMintRewardsProposal(ctx sdk.Context, k keeper.Keeper, p *types.NewMintRewardsProposal) error {
	return k.HandleNewMintRewardsProposal(ctx, p)
}

func handleDisableMintRewardsProposal(ctx sdk.Context, k keeper.Keeper, p *types.DisbaleMintRewardsProposal) error {
	return k.HandleDisableMintRewardsProposal(ctx, p)
}
