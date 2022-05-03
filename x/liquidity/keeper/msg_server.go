package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidity/types"
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

// CreatePair defines a method to create a pair.
func (m msgServer) CreatePair(goCtx context.Context, msg *types.MsgCreatePair) (*types.MsgCreatePairResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.CreatePair(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCreatePairResponse{}, nil
}

// CreatePool defines a method to create a liquidity pool.
func (m msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.CreatePool(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCreatePoolResponse{}, nil
}

// Deposit defines a method to deposit coins to the pool.
func (m msgServer) Deposit(goCtx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.Deposit(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgDepositResponse{}, nil
}

// Withdraw defines a method to withdraw pool coin from the pool.
func (m msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.Withdraw(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgWithdrawResponse{}, nil
}

// LimitOrder defines a method to making a limit order.
func (m msgServer) LimitOrder(goCtx context.Context, msg *types.MsgLimitOrder) (*types.MsgLimitOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.LimitOrder(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgLimitOrderResponse{}, nil
}

// MarketOrder defines a method to making a market order.
func (m msgServer) MarketOrder(goCtx context.Context, msg *types.MsgMarketOrder) (*types.MsgMarketOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.MarketOrder(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgMarketOrderResponse{}, nil
}

// CancelOrder defines a method to cancel an order.
func (m msgServer) CancelOrder(goCtx context.Context, msg *types.MsgCancelOrder) (*types.MsgCancelOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.CancelOrder(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCancelOrderResponse{}, nil
}

// CancelAllOrders defines a method to cancel all orders.
func (m msgServer) CancelAllOrders(goCtx context.Context, msg *types.MsgCancelAllOrders) (*types.MsgCancelAllOrdersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.CancelAllOrders(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCancelAllOrdersResponse{}, nil
}

func (k msgServer) BondPoolTokens(goCtx context.Context, msg *types.MsgBondPoolTokens) (*types.MsgBondPoolTokensResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	//IF user has provided liquidity to any given pool
	userPoolsData, found := k.GetIndividualUserPoolsData(ctx, sdk.AccAddress(msg.UserAddress))
	if !found {
		return nil, types.ErrUserNotHavingLiquidityInPools
	}
	//If user has provided liquidity to the pool he has entered
	poolExists := k.GetUserPoolsContributionData(userPoolsData, msg.PoolId)
	if !poolExists {
		return nil, types.ErrUserNotHavingLiquidityInCurrentPool
	}
	for _, pool := range userPoolsData.UserPoolWiseData {

		if pool.PoolId == msg.PoolId {
			//Check if that amount of  unbonded token exists ,
			//If exists , then bond user tokens.
			//Remove from unbonded section
			if pool.UnbondedPoolCoin.GTE(msg.PoolCoin.Amount) {
				updatedBondedTokens := msg.PoolCoin.Amount
				updatedUnBondedTokens := pool.UnbondedPoolCoin.Sub(msg.PoolCoin.Amount)
				pool.BondedPoolCoin = &updatedBondedTokens
				pool.UnbondedPoolCoin = &updatedUnBondedTokens

			} else {
				return nil, types.ErrNotEnoughCoinsForBonding
			}

		} else {
			continue
		}
	}
	k.SetIndividualUserPoolsData(ctx, userPoolsData)

	return &types.MsgBondPoolTokensResponse{}, nil
}

func (k msgServer) UnbondPoolTokens(goCtx context.Context, msg *types.MsgUnbondPoolTokens) (*types.MsgUnbondPoolTokensResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	//If user has provided liquidity to any pool
	userPoolsData, found := k.GetIndividualUserPoolsData(ctx, sdk.AccAddress(msg.UserAddress))
	if !found {
		return nil, types.ErrUserNotHavingLiquidityInPools
	}
	//If user has provided liquidity to the pool he has entered
	poolExists := k.GetUserPoolsContributionData(userPoolsData, msg.PoolId)
	if !poolExists {
		return nil, types.ErrUserNotHavingLiquidityInCurrentPool
	}
	for _, pool := range userPoolsData.UserPoolWiseData {

		if pool.PoolId == msg.PoolId {
			//Check if that amount of  bonded token exists,
			//If exists , then unbond user tokens.
			//Remove from bonded section
			if pool.BondedPoolCoin.GTE(msg.PoolCoin.Amount) {
				var userUnbondingTokens types.UserPoolUnbondingTokens
				userUnbondingTokens.IsUnbondingPoolCoin = &msg.PoolCoin.Amount
				//Check for mistakes in these values
				userUnbondingTokens.UnbondingStartTime = ctx.BlockTime()                            //Check for current second value
				userUnbondingTokens.UnbondingEndTime = k.CalculateUnbondingEndTime(ctx,ctx.BlockTime()) //Ending Time after the unbonding time will get over
				updatedBondedTokens := pool.BondedPoolCoin.Sub(msg.PoolCoin.Amount)
				pool.BondedPoolCoin = &updatedBondedTokens
				pool.UserPoolUnbondingTokens = append(pool.UserPoolUnbondingTokens, &userUnbondingTokens)

			} else {
				return nil, types.ErrNotEnoughCoinsForUnBonding
			}

		} else {
			continue
		}
	}

	k.SetIndividualUserPoolsData(ctx, userPoolsData)

	return &types.MsgUnbondPoolTokensResponse{}, nil
}
