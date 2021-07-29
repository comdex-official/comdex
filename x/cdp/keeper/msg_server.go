package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (ms msgServer) CreateCDP(goCtx context.Context, msg *types.MsgCreateCDP) (*types.MsgCreateCDPResponse, error) {
	//TODO
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = ms.Keeper.AddCdp(ctx, sender, msg.Collateral, msg.Principal, msg.CollateralType)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	)
	return &types.MsgCreateCDPResponse{}, nil
}

func (ms msgServer) Deposit(ctx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) Withdraw(ctx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) DrawDebt(ctx context.Context, msg *types.MsgDrawDebt) (*types.MsgDrawDebtResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) RepayDebt(ctx context.Context, msg *types.MsgRepayDebt) (*types.MsgRepayDebtResponse, error) {
	//TODO
	return nil, nil
}

func (ms msgServer) Liquidate(ctx context.Context, msg *types.MsgLiquidate) (*types.MsgLiquidateResponse, error) {
	//TODO
	return nil, nil
}
