package lend

import (
	"fmt"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/pkg/errors"

	"github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {

	server := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgLend:
			res, err := server.Lend(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgWithdraw:
			res, err := server.Withdraw(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgBorrow:
			res, err := server.Borrow(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgRepay:
			res, err := server.Repay(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func NewUpdateAssetProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddWhitelistedAssetsProposal:
			return handleAddAssetProposal(ctx, k, c)
		case *types.UpdateWhitelistedAssetProposal:
			return handleUpdateAssetProposal(ctx, k, c)
		case *types.AddWhitelistedPairsProposal:
			return handleAddPairsProposal(ctx, k, c)
		case *types.UpdateWhitelistedPairProposal:
			return handleUpdatePairProposal(ctx, k, c)

		default:
			return errors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
		}
	}
}

func handleAddAssetProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddWhitelistedAssetsProposal) error {
	return k.HandleProposalAddAsset(ctx, p)
}

func handleUpdateAssetProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdateWhitelistedAssetProposal) error {
	return k.HandleProposalUpdateAsset(ctx, p)
}

func handleAddPairsProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddWhitelistedPairsProposal) error {
	return k.HandleProposalAddPairs(ctx, p)
}

func handleUpdatePairProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdateWhitelistedPairProposal) error {
	return k.HandleProposalUpdatePair(ctx, p)
}
