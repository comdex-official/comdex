package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) Deposit(goCtx context.Context, deposit *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.keeper.Deposit(ctx, deposit.Amount, deposit.AppId, deposit.Addr); err != nil {
		return nil, err
	}

	return &types.MsgDepositResponse{}, nil
}

func (m msgServer) Refund(goCtx context.Context, refund *types.MsgRefund) (*types.MsgRefundResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	counter := m.keeper.GetRefundCounterStatus(ctx)
	if counter != 0 {
		return nil, types.ErrorRefundCompleted
	}

	if err := m.keeper.Refund(ctx); err != nil {
		return nil, err
	}

	m.keeper.SetRefundCounterStatus(ctx, counter+1)

	return &types.MsgRefundResponse{}, nil
}

func (m msgServer) UpdateDebtParams(goCtx context.Context, msg *types.MsgUpdateDebtParams) (*types.MsgUpdateDebtParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.keeper.UpdateDebtParams(ctx, msg.AppId, msg.AssetId, msg.Slots, msg.DebtThreshold, msg.LotSize, msg.DebtLotSize, msg.IsDebtAuction, msg.Addr); err != nil {
		return nil, err
	}

	return &types.MsgUpdateDebtParamsResponse{}, nil
}