package keeper

import (
	"context"
	"fmt"
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

	if err := m.keeper.LendAsset(ctx, lend.Lender, lend.AssetId, lend.Amount, lend.PoolId); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLoanAsset,
			sdk.NewAttribute(types.EventAttrLender, lend.Lender),
			sdk.NewAttribute(sdk.AttributeKeyAmount, lend.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventAttrModule),
			sdk.NewAttribute(sdk.AttributeKeySender, lend.Lender),
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

	lendId := withdraw.LendId

	if err := m.keeper.WithdrawAsset(ctx, withdraw.Lender, lendId, withdraw.Amount); err != nil {
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

func (m msgServer) Deposit(goCtx context.Context, deposit *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lenderAddr, err := sdk.AccAddressFromBech32(deposit.Lender)
	if err != nil {
		return nil, err
	}

	lendId := deposit.LendId

	if err := m.keeper.DepositAsset(ctx, deposit.Lender, lendId, deposit.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawLoanedAsset,
			sdk.NewAttribute(types.EventAttrLender, lenderAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, deposit.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventAttrModule),
			sdk.NewAttribute(sdk.AttributeKeySender, lenderAddr.String()),
		),
	})

	return &types.MsgDepositResponse{}, nil
}

func (m msgServer) CloseLend(goCtx context.Context, lend *types.MsgCloseLend) (*types.MsgCloseLendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lenderAddr, err := sdk.AccAddressFromBech32(lend.Lender)
	if err != nil {
		return nil, err
	}

	lendId := lend.LendId

	if err := m.keeper.CloseLend(ctx, lend.Lender, lendId); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawLoanedAsset,
			sdk.NewAttribute(types.EventAttrLender, lenderAddr.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventAttrModule),
			sdk.NewAttribute(sdk.AttributeKeySender, lenderAddr.String()),
		),
	})

	return &types.MsgCloseLendResponse{}, nil
}

func (m msgServer) Borrow(goCtx context.Context, borrow *types.MsgBorrow) (*types.MsgBorrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	borrowerAddr, err := sdk.AccAddressFromBech32(borrow.Borrower)
	if err != nil {
		return nil, err
	}

	if err := m.keeper.BorrowAsset(ctx, borrow.Borrower, borrow.LendId, borrow.PairId, borrow.IsStableBorrow, borrow.AmountIn, borrow.AmountOut); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBorrowAsset,
			sdk.NewAttribute(types.EventAttrBorrower, borrowerAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, borrow.AmountIn.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, borrow.AmountOut.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventAttrModule),
			sdk.NewAttribute(sdk.AttributeKeySender, borrowerAddr.String()),
		),
	})

	return &types.MsgBorrowResponse{}, nil
}

func (m msgServer) Repay(goCtx context.Context, repay *types.MsgRepay) (*types.MsgRepayResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	borrowerAddr, err := sdk.AccAddressFromBech32(repay.Borrower)
	if err != nil {
		return nil, err
	}

	if err := m.keeper.RepayAsset(ctx, repay.BorrowId, repay.Borrower, repay.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRepayBorrowedAsset,
			sdk.NewAttribute(types.EventAttrBorrower, borrowerAddr.String()),
			sdk.NewAttribute(types.EventAttrAttempted, repay.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventAttrModule),
			sdk.NewAttribute(sdk.AttributeKeySender, borrowerAddr.String()),
		),
	})

	return &types.MsgRepayResponse{}, nil
}

func (m msgServer) DepositBorrow(goCtx context.Context, borrow *types.MsgDepositBorrow) (*types.MsgDepositBorrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	borrowerAddr, err := sdk.AccAddressFromBech32(borrow.Borrower)
	if err != nil {
		return nil, err
	}

	if err := m.keeper.DepositBorrowAsset(ctx, borrow.BorrowId, borrow.Borrower, borrow.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRepayBorrowedAsset,
			sdk.NewAttribute(types.EventAttrBorrower, borrowerAddr.String()),
			sdk.NewAttribute(types.EventAttrAttempted, borrow.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventAttrModule),
			sdk.NewAttribute(sdk.AttributeKeySender, borrowerAddr.String()),
		),
	})

	return &types.MsgDepositBorrowResponse{}, nil
}

func (m msgServer) Draw(goCtx context.Context, draw *types.MsgDraw) (*types.MsgDrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	borrowerAddr, err := sdk.AccAddressFromBech32(draw.Borrower)
	if err != nil {
		return nil, err
	}

	err = m.keeper.DrawAsset(ctx, draw.BorrowId, draw.Borrower, draw.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRepayBorrowedAsset,
			sdk.NewAttribute(types.EventAttrBorrower, borrowerAddr.String()),
			sdk.NewAttribute(types.EventAttrAttempted, draw.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventAttrModule),
			sdk.NewAttribute(sdk.AttributeKeySender, borrowerAddr.String()),
		),
	})

	return &types.MsgDrawResponse{}, nil
}

func (m msgServer) CloseBorrow(goCtx context.Context, borrow *types.MsgCloseBorrow) (*types.MsgCloseBorrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lenderAddr, err := sdk.AccAddressFromBech32(borrow.Borrower)
	if err != nil {
		return nil, err
	}

	BorrowId := borrow.BorrowId

	if err := m.keeper.CloseBorrow(ctx, borrow.Borrower, BorrowId); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawLoanedAsset,
			sdk.NewAttribute(types.EventAttrLender, lenderAddr.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.EventAttrModule),
			sdk.NewAttribute(sdk.AttributeKeySender, lenderAddr.String()),
		),
	})

	return &types.MsgCloseBorrowResponse{}, nil
}

func (m msgServer) FundModuleAccounts(goCtx context.Context, accounts *types.MsgFundModuleAccounts) (*types.MsgFundModuleAccountsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lenderAddr, err := sdk.AccAddressFromBech32(accounts.Lender)
	if err != nil {
		return nil, err
	}
	fmt.Println("goose......")

	if err := m.keeper.FundModAcc(ctx, accounts.ModuleName, accounts.AssetId, lenderAddr, accounts.Amount); err != nil {
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
