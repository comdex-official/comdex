package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/esm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
var (
	_ types.MsgServer = (*msgServer)(nil)
)

type msgServer struct {
	keeper Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

func (m msgServer) Deposit(goCtx context.Context, deposit *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	appID := deposit.AppId

	if err := m.keeper.DepositESM(ctx, deposit.Depositor, appID, deposit.Amount); err != nil {
		return nil, err
	}

	return &types.MsgDepositResponse{}, nil
}

func (m msgServer) Execute(goCtx context.Context, execute *types.MsgExecute) (*types.MsgExecuteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	appID := execute.AppId

	if err := m.keeper.ExecuteESM(ctx, execute.Depositor, appID); err != nil {
		return nil, err
	}

	return &types.MsgExecuteResponse{}, nil
}
