package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"

	"github.com/comdex-official/comdex/x/cdp/types"
)

var (
	_ types.MsgServiceServer = (*msgServer)(nil)
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServiceServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k *msgServer) MsgCreate(c context.Context, msg *types.MsgCreateRequest) (*types.MsgCreateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if k.HasCDPForAddressByAssetPair(ctx, sender, msg.PairId) {
		return nil, types.ErrorCDPAlreadyExists
	}

	pair, found := k.GetAssetPair(ctx, msg.PairId)
	if !found {
		return nil, errors.Wrapf(types.ErrorAssetPairDoesNotExist, "%d", msg.PairId)
	}

	var (
		coins  = k.SpendableCoins(ctx, sender)
		amount = coins.AmountOf(pair.DenomIn)
	)

	if amount.LT(msg.AmountIn) {
		return nil, errors.Wrapf(types.ErrorInsufficientCollateral, "expected %s, have %s", msg.AmountIn, amount)
	}

	return &types.MsgCreateResponse{}, nil
}

func (k *msgServer) MsgDeposit(c context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	panic("implement me")
}

func (k *msgServer) MsgWithdraw(c context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	panic("implement me")
}

func (k *msgServer) MsgDraw(c context.Context, msg *types.MsgDrawRequest) (*types.MsgDrawResponse, error) {
	panic("implement me")
}

func (k *msgServer) MsgRepay(c context.Context, msg *types.MsgRepayRequest) (*types.MsgRepayResponse, error) {
	panic("implement me")
}

func (k *msgServer) MsgLiquidate(c context.Context, msg *types.MsgLiquidateRequest) (*types.MsgLiquidateResponse, error) {
	panic("implement me")
}
