package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/lend/types"
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
	ctx.GasMeter().ConsumeGas(types.LendGas, "LendGas")

	if err := m.keeper.LendAsset(ctx, lend.Lender, lend.AssetId, lend.Amount, lend.PoolId, lend.AppId); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLend,
			sdk.NewAttribute(types.AttributeKeyCreator, lend.Lender),
			sdk.NewAttribute(types.AttributeKeyAssetID, strconv.FormatUint(lend.AssetId, 10)),
			sdk.NewAttribute(types.AttributeKeyAmountIn, lend.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(lend.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyAppID, strconv.FormatUint(lend.AppId, 10)),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	})
	return &types.MsgLendResponse{}, nil
}

func (m msgServer) Withdraw(goCtx context.Context, withdraw *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.GasMeter().ConsumeGas(types.WithdrawGas, "WithdrawGas")

	lendID := withdraw.LendId

	if err := m.keeper.WithdrawAsset(ctx, withdraw.Lender, lendID, withdraw.Amount); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdraw,
			sdk.NewAttribute(types.AttributeKeyCreator, withdraw.Lender),
			sdk.NewAttribute(types.AttributeKeyLendID, strconv.FormatUint(withdraw.LendId, 10)),
			sdk.NewAttribute(types.AttributeKeyAmountOut, withdraw.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	})
	return &types.MsgWithdrawResponse{}, nil
}

func (m msgServer) Deposit(goCtx context.Context, deposit *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.GasMeter().ConsumeGas(types.DepositGas, "DepositGas")

	lendID := deposit.LendId

	if err := m.keeper.DepositAsset(ctx, deposit.Lender, lendID, deposit.Amount); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeposit,
			sdk.NewAttribute(types.AttributeKeyCreator, deposit.Lender),
			sdk.NewAttribute(types.AttributeKeyLendID, strconv.FormatUint(deposit.LendId, 10)),
			sdk.NewAttribute(types.AttributeKeyAmountIn, deposit.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	})
	return &types.MsgDepositResponse{}, nil
}

func (m msgServer) CloseLend(goCtx context.Context, lend *types.MsgCloseLend) (*types.MsgCloseLendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.GasMeter().ConsumeGas(types.CloseLendGas, "CloseLendGas")

	lendID := lend.LendId

	if err := m.keeper.CloseLend(ctx, lend.Lender, lendID); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClose,
			sdk.NewAttribute(types.AttributeKeyCreator, lend.Lender),
			sdk.NewAttribute(types.AttributeKeyLendID, strconv.FormatUint(lend.LendId, 10)),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	})
	return &types.MsgCloseLendResponse{}, nil
}

func (m msgServer) Borrow(goCtx context.Context, borrow *types.MsgBorrow) (*types.MsgBorrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.GasMeter().ConsumeGas(types.BorrowAssetGas, "BorrowAssetGas")

	if err := m.keeper.BorrowAsset(ctx, borrow.Borrower, borrow.LendId, borrow.PairId, borrow.IsStableBorrow, borrow.AmountIn, borrow.AmountOut); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBorrow,
			sdk.NewAttribute(types.AttributeKeyCreator, borrow.Borrower),
			sdk.NewAttribute(types.AttributeKeyLendID, strconv.FormatUint(borrow.LendId, 10)),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(borrow.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyIsStable, strconv.FormatBool(borrow.IsStableBorrow)),
			sdk.NewAttribute(types.AttributeKeyAmountIn, borrow.AmountIn.String()),
			sdk.NewAttribute(types.AttributeKeyAmountOut, borrow.AmountOut.String()),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	})
	return &types.MsgBorrowResponse{}, nil
}

func (m msgServer) Repay(goCtx context.Context, repay *types.MsgRepay) (*types.MsgRepayResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.GasMeter().ConsumeGas(types.RepayAssetGas, "RepayAssetGas")

	if err := m.keeper.RepayAsset(ctx, repay.BorrowId, repay.Borrower, repay.Amount); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRepay,
			sdk.NewAttribute(types.AttributeKeyCreator, repay.Borrower),
			sdk.NewAttribute(types.AttributeKeyBorrowID, strconv.FormatUint(repay.BorrowId, 10)),
			sdk.NewAttribute(types.AttributeKeyAmountIn, repay.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	})

	return &types.MsgRepayResponse{}, nil
}

func (m msgServer) DepositBorrow(goCtx context.Context, borrow *types.MsgDepositBorrow) (*types.MsgDepositBorrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.GasMeter().ConsumeGas(types.DepositBorrowAssetGas, "DepositBorrowAssetGas")

	if err := m.keeper.DepositBorrowAsset(ctx, borrow.BorrowId, borrow.Borrower, borrow.Amount); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDepositBorrow,
			sdk.NewAttribute(types.AttributeKeyCreator, borrow.Borrower),
			sdk.NewAttribute(types.AttributeKeyBorrowID, strconv.FormatUint(borrow.BorrowId, 10)),
			sdk.NewAttribute(types.AttributeKeyAmountIn, borrow.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	})

	return &types.MsgDepositBorrowResponse{}, nil
}

func (m msgServer) Draw(goCtx context.Context, draw *types.MsgDraw) (*types.MsgDrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.GasMeter().ConsumeGas(types.DrawAssetGas, "DrawAssetGas")

	err := m.keeper.DrawAsset(ctx, draw.BorrowId, draw.Borrower, draw.Amount)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDraw,
			sdk.NewAttribute(types.AttributeKeyCreator, draw.Borrower),
			sdk.NewAttribute(types.AttributeKeyBorrowID, strconv.FormatUint(draw.BorrowId, 10)),
			sdk.NewAttribute(types.AttributeKeyAmountOut, draw.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	})
	return &types.MsgDrawResponse{}, nil
}

func (m msgServer) CloseBorrow(goCtx context.Context, borrow *types.MsgCloseBorrow) (*types.MsgCloseBorrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.GasMeter().ConsumeGas(types.CloseBorrowAssetGas, "CloseBorrowAssetGas")

	borrowID := borrow.BorrowId

	if err := m.keeper.CloseBorrow(ctx, borrow.Borrower, borrowID); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCloseBorrow,
			sdk.NewAttribute(types.AttributeKeyCreator, borrow.Borrower),
			sdk.NewAttribute(types.AttributeKeyBorrowID, strconv.FormatUint(borrow.BorrowId, 10)),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	})
	return &types.MsgCloseBorrowResponse{}, nil
}

func (m msgServer) BorrowAlternate(goCtx context.Context, alternate *types.MsgBorrowAlternate) (*types.MsgBorrowAlternateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.GasMeter().ConsumeGas(types.BorrowAssetAlternateGas, "BorrowAssetAlternateGas")

	if err := m.keeper.BorrowAlternate(ctx, alternate.Lender, alternate.AssetId, alternate.PoolId, alternate.AmountIn, alternate.PairId, alternate.IsStableBorrow, alternate.AmountOut, alternate.AppId); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBorrowAlternate,
			sdk.NewAttribute(types.AttributeKeyCreator, alternate.Lender),
			sdk.NewAttribute(types.AttributeKeyAssetID, strconv.FormatUint(alternate.AssetId, 10)),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(alternate.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyAmountIn, alternate.AmountIn.String()),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(alternate.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyIsStable, strconv.FormatBool(alternate.IsStableBorrow)),
			sdk.NewAttribute(types.AttributeKeyAmountOut, alternate.AmountOut.String()),
			sdk.NewAttribute(types.AttributeKeyAppID, strconv.FormatUint(alternate.AppId, 10)),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	})
	return &types.MsgBorrowAlternateResponse{}, nil
}

func (m msgServer) FundModuleAccounts(goCtx context.Context, accounts *types.MsgFundModuleAccounts) (*types.MsgFundModuleAccountsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lenderAddr, err := sdk.AccAddressFromBech32(accounts.Lender)
	if err != nil {
		return nil, err
	}
	if err = m.keeper.FundModAcc(ctx, accounts.PoolId, accounts.AssetId, lenderAddr, accounts.Amount); err != nil {
		return nil, err
	}

	return &types.MsgFundModuleAccountsResponse{}, nil
}

func (m msgServer) CalculateInterestAndRewards(goCtx context.Context, rewards *types.MsgCalculateInterestAndRewards) (*types.MsgCalculateInterestAndRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.GasMeter().ConsumeGas(types.CalculateInterestAndRewardGas, "CalculateInterestAndRewardGas")

	if err := m.keeper.MsgCalculateInterestAndRewards(ctx, rewards.Borrower); err != nil {
		return nil, err
	}

	return &types.MsgCalculateInterestAndRewardsResponse{}, nil
}
