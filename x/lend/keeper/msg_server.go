package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/lend/types"
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

func (m msgServer) Lend(goCtx context.Context, lend *types.MsgLend) (*types.MsgLendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lenderAddr, err := sdk.AccAddressFromBech32(lend.Lender)
	if err != nil {
		return nil, err
	}

	if err := m.keeper.LendAsset(ctx, lenderAddr, lend.PairId, lend.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLoanAsset,
			sdk.NewAttribute(types.EventAttrLender, lenderAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, lend.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventAttrModule),
			sdk.NewAttribute(sdk.AttributeKeySender, lenderAddr.String()),
		),
	})

	return &types.MsgLendResponse{}, nil
}

func (m msgServer) Withdraw(goCtx context.Context, withdraw *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lenderAddr, err := sdk.AccAddressFromBech32(withdraw.Lender)
	if err != nil {
		return nil, err
	}

	if err := m.keeper.WithdrawAsset(ctx, lenderAddr, withdraw.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawLoanedAsset,
			sdk.NewAttribute(types.EventAttrLender, lenderAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, withdraw.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventAttrModule),
			sdk.NewAttribute(sdk.AttributeKeySender, lenderAddr.String()),
		),
	})

	return &types.MsgWithdrawResponse{}, nil
}

func (m msgServer) FundModuleAccounts(goCtx context.Context, accounts *types.MsgFundModuleAccounts) (*types.MsgFundModuleAccountsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lenderAddr, err := sdk.AccAddressFromBech32(accounts.Lender)
	if err != nil {
		return nil, err
	}

	if err := m.keeper.FundModAcc(ctx, accounts.ModuleName, lenderAddr, accounts.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLoanAsset,
			sdk.NewAttribute(types.EventAttrLender, lenderAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, accounts.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventAttrModule),
			sdk.NewAttribute(sdk.AttributeKeySender, lenderAddr.String()),
		),
	})

	return &types.MsgFundModuleAccountsResponse{}, nil
}
