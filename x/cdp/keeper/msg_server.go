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

	err = ms.Keeper.AddCdp(ctx, sender, msg.Collateral, msg.Debt, msg.CollateralType)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventCreateCDP{
			Sender:         sender.String(),
			CollateralType: msg.CollateralType,
		},
	)
	return &types.MsgCreateCDPResponse{}, nil
}

func (ms msgServer) MsgDepositCollateral(context context.Context, msg *types.MsgDepositCollateralRequest) (*types.MsgDepositCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	err = ms.DepositCollateral(ctx, owner, msg.Collateral, msg.CollateralType)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventDepositCollateral{
			Owner:          owner.String(),
			CollateralType: msg.CollateralType,
			Collateral:     msg.Collateral,
		},
	)

	return &types.MsgDepositCollateralResponse{}, nil
}

func (ms msgServer) MsgWithdrawCollateral(context context.Context, msg *types.MsgWithdrawCollateralRequest) (*types.MsgWithdrawCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	err = ms.WithdrawCollateral(ctx, owner, msg.Collateral, msg.CollateralType)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventWithdrawCollateral{
			Owner:          owner.String(),
			CollateralType: msg.CollateralType,
			Collateral:     msg.Collateral,
		},
	)

	return &types.MsgWithdrawCollateralResponse{}, nil
}

func (ms msgServer) MsgDrawDebt(context context.Context, msg *types.MsgDrawDebtRequest) (*types.MsgDrawDebtResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	err = ms.DrawDebt(ctx, owner, msg.CollateralType, msg.Debt)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventDrawDebt{
			Owner:          owner.String(),
			CollateralType: msg.CollateralType,
			Debt:           msg.Debt,
		},
	)

	return &types.MsgDrawDebtResponse{}, nil
}

func (ms msgServer) MsgRepayDebt(context context.Context, msg *types.MsgRepayDebtRequest) (*types.MsgRepayDebtResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	err = ms.RepayDebt(ctx, owner, msg.CollateralType, msg.Debt)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventRepayDebt{
			Owner:          owner.String(),
			CollateralType: msg.CollateralType,
			Debt:           msg.Debt,
		},
	)

	return &types.MsgRepayDebtResponse{}, nil
}

func (ms msgServer) MsgLiquidateCDP(context context.Context, msg *types.MsgLiquidateCDPRequest) (*types.MsgLiquidateCDPResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	err = ms.AttemptLiquidation(ctx, owner, msg.CollateralType)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventLiquidateCDP{
			Owner:          owner.String(),
			CollateralType: msg.CollateralType,
		},
	)

	return &types.MsgLiquidateCDPResponse{}, nil
}
