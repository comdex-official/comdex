package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/cdp/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServiceServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServiceServer = msgServer{}

func (ms msgServer) MsgCreateCDP(ctx context.Context, msg *types.MsgCreateCDPRequest) (*types.MsgCreateCDPResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) MsgDeposit(ctx context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) MsgWithdraw(ctx context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) MsgDrawDebt(ctx context.Context, msg *types.MsgDrawDebtRequest) (*types.MsgDrawDebtResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) MsgRepayDebt(ctx context.Context, msg *types.MsgRepayDebtRequest) (*types.MsgRepayDebtResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) MsgLiquidate(ctx context.Context, msg *types.MsgLiquidateRequest) (*types.MsgLiquidateResponse, error) {
	//TODO
	return nil, nil
}
