package lend

import (
	"fmt"
	"github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/pkg/errors"
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

		case *types.MsgFundModuleAccounts:
			res, err := server.FundModuleAccounts(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func NewLendPairHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.LendPairsProposal:
			return handleAddWhitelistedPairsProposal(ctx, k, c)
		case *types.UpdatePairProposal:
			return handleUpdateWhitelistedPairProposal(ctx, k, c)

		default:
			return errors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
		}
	}
}

func handleAddWhitelistedPairsProposal(ctx sdk.Context, k keeper.Keeper, p *types.LendPairsProposal) error {
	return k.HandleAddWhitelistedPairsRecords(ctx, p)
}

func handleUpdateWhitelistedPairProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdatePairProposal) error {
	return k.HandleUpdateWhitelistedPairRecords(ctx, p)
}
