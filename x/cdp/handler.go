package cdp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/x/cdp/keeper"
	"github.com/comdex-official/comdex/x/cdp/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateCDPRequest:
			result, err := msgServer.MsgCreateCDP(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, result, err)
		case *types.MsgDepositCollateralRequest:
			result, err := msgServer.MsgDepositCollateral(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, result, err)
		case *types.MsgWithdrawCollateralRequest:
			result, err := msgServer.MsgWithdrawCollateral(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, result, err)
		case *types.MsgDrawDebtRequest:
			result, err := msgServer.MsgDrawDebt(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, result, err)
		case *types.MsgRepayDebtRequest:
			result, err := msgServer.MsgRepayDebt(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, result, err)
		case *types.MsgCloseCDPRequest:
			result, err := msgServer.MsgCloseCDP(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, result, err)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
