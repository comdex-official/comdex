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
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (ms msgServer) CreateCDP(ctx context.Context, msg *types.MsgCreateCDPRequest) (*types.MsgCreateCDPResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) Deposit(ctx context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) Withdraw(ctx context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) DrawDebt(ctx context.Context, msg *types.MsgDrawDebtRequest) (*types.MsgDrawDebtResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) RepayDebt(ctx context.Context, msg *types.MsgRepayDebtRequest) (*types.MsgRepayDebtResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) Liquidate(ctx context.Context, msg *types.MsgLiquidateRequest) (*types.MsgLiquidateResponse, error) {
	//TODO
	return nil, nil
}
