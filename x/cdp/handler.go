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
			result, err := msgServer.CreateCDP(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, result, err)
		case *types.MsgDepositRequest:
			result, err:= msgServer.Deposit(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, result, err)
		case *types.MsgWithdrawRequest:
			result, err:= msgServer.Withdraw(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, result, err)
		case *types.MsgDrawDebtRequest:
			result, err:= msgServer.DrawDebt(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, result, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
