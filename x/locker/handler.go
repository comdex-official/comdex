package locker

import (
	"github.com/comdex-official/comdex/x/locker/keeper"
	"github.com/comdex-official/comdex/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	server := keeper.NewMsgServer(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateLockerRequest:
			res, err := server.MsgCreateLocker(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDepositAssetRequest:
			res, err := server.MsgDepositAsset(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgWithdrawAssetRequest:
			res, err := server.MsgWithdrawAsset(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCloseLockerRequest:
			res, err := server.MsgCloseLocker(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(types.ErrorUnknownMsgType, "%T", msg)
		}
	}
}
