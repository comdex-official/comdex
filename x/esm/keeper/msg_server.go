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

func (m msgServer) DepositESM(goCtx context.Context, deposit *types.MsgDepositESM) (*types.MsgDepositESMResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	appID := deposit.AppId

	if err := m.keeper.DepositESM(ctx, deposit.Depositor, appID, deposit.Amount); err != nil {
		return nil, err
	}

	return &types.MsgDepositESMResponse{}, nil
}

func (m msgServer) ExecuteESM(goCtx context.Context, execute *types.MsgExecuteESM) (*types.MsgExecuteESMResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	appID := execute.AppId

	if err := m.keeper.ExecuteESM(ctx, execute.Depositor, appID); err != nil {
		return nil, err
	}

	return &types.MsgExecuteESMResponse{}, nil
}
