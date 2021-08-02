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
func NewMsgServerImpl(keeper Keeper) types.MsgServiceServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServiceServer = msgServer{}

func (ms msgServer) MsgCreateCDP(context context.Context, msg *types.MsgCreateCDPRequest) (*types.MsgCreateCDPResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)

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

func (ms msgServer) MsgDeposit(context context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	err := ms.DepositCollateral(ctx, msg.Sender, msg.Collateral, msg.CollateralType)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Collateral.Amount.String()),
		),
	)

	return &types.MsgDepositResponse{}, nil
}

func (ms msgServer) MsgWithdraw(context context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	err := ms.DepositCollateral(ctx, msg.Sender, msg.Collateral, msg.CollateralType)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Collateral.Amount.String()),
		),
	)

	return &types.MsgWithdrawResponse{}, nil
}

func (ms msgServer) MsgDrawDebt(context context.Context, msg *types.MsgDrawDebtRequest) (*types.MsgDrawDebtResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	err := ms.AddPrincipal(ctx, msg.Sender, msg.CollateralType, msg.Principal)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Principal.Amount.String()),
		),
	)

	return &types.MsgDrawDebtResponse{}, nil
}

func (ms msgServer) MsgRepayDebt(context context.Context, msg *types.MsgRepayDebtRequest) (*types.MsgRepayDebtResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	err := ms.RepayPrincipal(ctx, msg.Sender, msg.CollateralType, msg.Payment)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Payment.Amount.String()),
		),
	)

	return &types.MsgRepayDebtResponse{}, nil
}

func (ms msgServer) MsgLiquidate(context context.Context, msg *types.MsgLiquidateRequest) (*types.MsgLiquidateResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	err := ms.AttemptLiquidation(ctx, msg.Sender, msg.CollateralType)
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

	return &types.MsgLiquidateResponse{}, nil
}
