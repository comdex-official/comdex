package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	if k.HasCDPForAddressByPair(ctx, from, msg.PairID) {
		return nil, types.ErrorDuplicateCDP
	}

	pair, found := k.GetPair(ctx, msg.PairID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	assetIn, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	var (
		balance = k.SpendableCoins(ctx, from)
		amount  = balance.AmountOf(assetIn.Denom)
	)

	if amount.LT(msg.AmountIn) {
		return nil, types.ErrorInsufficientAmount
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
