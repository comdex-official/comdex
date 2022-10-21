package keeper

import (
	"context"

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

	if err := m.keeper.LendAsset(ctx, lend.Lender, lend.AssetId, lend.Amount, lend.PoolId, lend.AppId); err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.LendGas, "LendGas")

	return &types.MsgLendResponse{}, nil
}

func (m msgServer) Withdraw(goCtx context.Context, withdraw *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lendID := withdraw.LendId

	if err := m.keeper.WithdrawAsset(ctx, withdraw.Lender, lendID, withdraw.Amount); err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.WithdrawGas, "WithdrawGas")

	return &types.MsgWithdrawResponse{}, nil
}

func (m msgServer) Deposit(goCtx context.Context, deposit *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lendID := deposit.LendId

	if err := m.keeper.DepositAsset(ctx, deposit.Lender, lendID, deposit.Amount); err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.DepositGas, "DepositGas")

	return &types.MsgDepositResponse{}, nil
}

func (m msgServer) CloseLend(goCtx context.Context, lend *types.MsgCloseLend) (*types.MsgCloseLendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lendID := lend.LendId

	if err := m.keeper.CloseLend(ctx, lend.Lender, lendID); err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.CloseLendGas, "CloseLendGas")

	return &types.MsgCloseLendResponse{}, nil
}

func (m msgServer) Borrow(goCtx context.Context, borrow *types.MsgBorrow) (*types.MsgBorrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.keeper.BorrowAsset(ctx, borrow.Borrower, borrow.LendId, borrow.PairId, borrow.IsStableBorrow, borrow.AmountIn, borrow.AmountOut); err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.BorrowAssetGas, "BorrowAssetGas")

	return &types.MsgBorrowResponse{}, nil
}

func (m msgServer) Repay(goCtx context.Context, repay *types.MsgRepay) (*types.MsgRepayResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.keeper.RepayAsset(ctx, repay.BorrowId, repay.Borrower, repay.Amount); err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.RepayAssetGas, "RepayAssetGas")

	return &types.MsgRepayResponse{}, nil
}

func (m msgServer) DepositBorrow(goCtx context.Context, borrow *types.MsgDepositBorrow) (*types.MsgDepositBorrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.keeper.DepositBorrowAsset(ctx, borrow.BorrowId, borrow.Borrower, borrow.Amount); err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.DepositBorrowAssetGas, "DepositBorrowAssetGas")

	return &types.MsgDepositBorrowResponse{}, nil
}

func (m msgServer) Draw(goCtx context.Context, draw *types.MsgDraw) (*types.MsgDrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.keeper.DrawAsset(ctx, draw.BorrowId, draw.Borrower, draw.Amount)
	if err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.DrawAssetGas, "DrawAssetGas")

	return &types.MsgDrawResponse{}, nil
}

func (m msgServer) CloseBorrow(goCtx context.Context, borrow *types.MsgCloseBorrow) (*types.MsgCloseBorrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	borrowID := borrow.BorrowId

	if err := m.keeper.CloseBorrow(ctx, borrow.Borrower, borrowID); err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.CloseBorrowAssetGas, "CloseBorrowAssetGas")

	return &types.MsgCloseBorrowResponse{}, nil
}

func (m msgServer) BorrowAlternate(goCtx context.Context, alternate *types.MsgBorrowAlternate) (*types.MsgBorrowAlternateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.keeper.BorrowAlternate(ctx, alternate.Lender, alternate.AssetId, alternate.PoolId, alternate.AmountIn, alternate.PairId, alternate.IsStableBorrow, alternate.AmountOut, alternate.AppId); err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.BorrowAssetAlternateGas, "BorrowAssetAlternateGas")

	return &types.MsgBorrowAlternateResponse{}, nil
}

func (m msgServer) FundModuleAccounts(goCtx context.Context, accounts *types.MsgFundModuleAccounts) (*types.MsgFundModuleAccountsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lenderAddr, err := sdk.AccAddressFromBech32(accounts.Lender)
	if err != nil {
		return nil, err
	}

	if err = m.keeper.FundModAcc(ctx, accounts.ModuleName, accounts.AssetId, lenderAddr, accounts.Amount); err != nil {
		return nil, err
	}

	return &types.MsgFundModuleAccountsResponse{}, nil
}

func (m msgServer) CalculateInterestAndRewards(goCtx context.Context, rewards *types.MsgCalculateInterestAndRewards) (*types.MsgCalculateInterestAndRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.keeper.MsgCalculateInterestAndRewards(ctx, rewards.Borrower); err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.CalculateInterestAndRewardGas, "CalculateInterestAndRewardGas")

	return &types.MsgCalculateInterestAndRewardsResponse{}, nil
}
