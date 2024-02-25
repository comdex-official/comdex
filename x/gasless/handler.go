package gasless

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/pkg/errors"

	"github.com/comdex-official/comdex/x/gasless/keeper"
	"github.com/comdex-official/comdex/x/gasless/types"
)

// NewHandler returns a new msg handler.
func NewHandler(k keeper.Keeper) sdk.Handler {
	// msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func NewGaslessProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		default:
			return errors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
		}
	}
}
