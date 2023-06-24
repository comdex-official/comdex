package nft

import (
	"github.com/comdex-official/comdex/x/nft/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k keeper.Keeper) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "%s", msg)
		}
	}
}
