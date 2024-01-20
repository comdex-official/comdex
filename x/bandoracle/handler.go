package bandoracle

import (
	"fmt"

	bam "github.com/cosmos/cosmos-sdk/baseapp"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/x/bandoracle/keeper"
	"github.com/comdex-official/comdex/x/bandoracle/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) bam.MsgServiceHandler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func NewFetchPriceHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.FetchPriceProposal:
			return handleFetchPriceProposal(ctx, k, c)

		default:
			return errors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
		}
	}
}

func handleFetchPriceProposal(ctx sdk.Context, k keeper.Keeper, p *types.FetchPriceProposal) error {
	return k.HandleProposalFetchPrice(ctx, p)
}
