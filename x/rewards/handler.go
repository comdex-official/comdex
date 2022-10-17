package rewards

import (
	"fmt"
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

		case *types.ActivateExternalRewardsLend:
			res, err := server.ExternalRewardsLend(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
