package cdp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"

	"github.com/comdex-official/comdex/x/cdp/keeper"
	"github.com/comdex-official/comdex/x/cdp/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	server := keeper.NewMsgServiceServer(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateRequest:
			res, err := server.MsgCreate(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDepositRequest:
			res, err := server.MsgDeposit(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgWithdrawRequest:
			res, err := server.MsgWithdraw(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDrawRequest:
			res, err := server.MsgDraw(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRepayRequest:
			res, err := server.MsgRepay(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgLiquidateRequest:
			res, err := server.MsgLiquidate(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, errors.Wrapf(types.ErrorUnknownMsgType, "%T", msg)
		}
	}
}
