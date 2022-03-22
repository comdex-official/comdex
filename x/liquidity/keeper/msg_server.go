package keeper

// DONTCOVER

// Although written in msg_server_test.go, it is approached at the keeper level rather than at the msgServer level
// so is not included in the coverage.

import (
	"context"
	"fmt"
	"strconv"


	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the distribution MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// Message server, handler for CreatePool msg
func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.GetCircuitBreakerEnabled(ctx) {
		return nil, types.ErrCircuitBreakerEnabled
	}

	pool, err := k.Keeper.CreatePool(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(pool.Id, 10)),
			sdk.NewAttribute(types.AttributeValuePoolTypeId, fmt.Sprintf("%d", msg.PoolTypeId)),
			sdk.NewAttribute(types.AttributeValuePoolName, pool.Name()),
			sdk.NewAttribute(types.AttributeValueReserveAccount, pool.ReserveAccountAddress),
			sdk.NewAttribute(types.AttributeValueDepositCoins, msg.DepositCoins.String()),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, pool.PoolCoinDenom),
		),
	})

	return &types.MsgCreatePoolResponse{}, nil
}

// Message server, handler for MsgDepositWithinBatch
func (k msgServer) DepositWithinBatch(goCtx context.Context, msg *types.MsgDepositWithinBatch) (*types.MsgDepositWithinBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.GetCircuitBreakerEnabled(ctx) {
		return nil, types.ErrCircuitBreakerEnabled
	}

	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}

	batchMsg, err := k.Keeper.DepositWithinBatch(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeDepositWithinBatch,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(poolBatch.Index, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueDepositCoins, batchMsg.Msg.DepositCoins.String()),
		),
	})

	return &types.MsgDepositWithinBatchResponse{}, nil
}

// Message server, handler for MsgWithdrawWithinBatch
func (k msgServer) WithdrawWithinBatch(goCtx context.Context, msg *types.MsgWithdrawWithinBatch) (*types.MsgWithdrawWithinBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}

	batchMsg, err := k.Keeper.WithdrawWithinBatch(ctx, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeWithdrawWithinBatch,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(poolBatch.Index, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, batchMsg.Msg.PoolCoin.Denom),
			sdk.NewAttribute(types.AttributeValuePoolCoinAmount, batchMsg.Msg.PoolCoin.Amount.String()),
		),
	})

	return &types.MsgWithdrawWithinBatchResponse{}, nil
}

// Message server, handler for MsgSwapWithinBatch
func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwapWithinBatch) (*types.MsgSwapWithinBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.GetCircuitBreakerEnabled(ctx) {
		return nil, types.ErrCircuitBreakerEnabled
	}

	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}

	batchMsg, err := k.Keeper.SwapWithinBatch(ctx, msg, types.CancelOrderLifeSpan)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeSwapWithinBatch,
			sdk.NewAttribute(types.AttributeValuePoolId, strconv.FormatUint(batchMsg.Msg.PoolId, 10)),
			sdk.NewAttribute(types.AttributeValueBatchIndex, strconv.FormatUint(poolBatch.Index, 10)),
			sdk.NewAttribute(types.AttributeValueMsgIndex, strconv.FormatUint(batchMsg.MsgIndex, 10)),
			sdk.NewAttribute(types.AttributeValueSwapTypeId, strconv.FormatUint(uint64(batchMsg.Msg.SwapTypeId), 10)),
			sdk.NewAttribute(types.AttributeValueOfferCoinDenom, batchMsg.Msg.OfferCoin.Denom),
			sdk.NewAttribute(types.AttributeValueOfferCoinAmount, batchMsg.Msg.OfferCoin.Amount.String()),
			sdk.NewAttribute(types.AttributeValueOfferCoinFeeAmount, batchMsg.Msg.OfferCoinFee.Amount.String()),
			sdk.NewAttribute(types.AttributeValueDemandCoinDenom, batchMsg.Msg.DemandCoinDenom),
			sdk.NewAttribute(types.AttributeValueOrderPrice, batchMsg.Msg.OrderPrice.String()),
		),
	})

	return &types.MsgSwapWithinBatchResponse{}, nil
}

func (k msgServer) BondPoolTokens(goCtx context.Context, msg *types.MsgBondPoolTokens) (*types.MsgBondPoolTokensResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.GetCircuitBreakerEnabled(ctx) {
		return nil, types.ErrCircuitBreakerEnabled
	}
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
	if k.GetCircuitBreakerEnabled(ctx) {
		return nil, types.ErrCircuitBreakerEnabled
	}
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
				userUnbondingTokens.UnbondingStartTime = float32(ctx.BlockTime().Second())//Check for current second value
				userUnbondingTokens.UnbondingEndTime =float32(k.CalculateUnbondingEndTime(int64(ctx.BlockTime().Second())))//Ending Time after the unbonding time will get over
				updatedBondedTokens := pool.BondedPoolCoin.Sub(msg.PoolCoin.Amount)
				pool.BondedPoolCoin=&updatedBondedTokens
				pool.UserPoolUnbondingTokens=append(pool.UserPoolUnbondingTokens, &userUnbondingTokens)

			} else {
				return nil, types.ErrNotEnoughCoinsForUnBonding
			}

		} else {
			continue
		}
	}

	fmt.Print("------------------------------------------------------")
	fmt.Print("------------------------------------------------------")
	fmt.Print("------------------------------------------------------")
	fmt.Print("------------------------------------------------------")
	fmt.Print("------------------------------------------------------")
	fmt.Print("------------------------------------------------------")
	fmt.Print("Checking updated data----1", userPoolsData)
	fmt.Print("------------------------------------------------------")
	fmt.Print("------------------------------------------------------")
	fmt.Print("------------------------------------------------------")
	fmt.Print("------------------------------------------------------")
	fmt.Print("------------------------------------------------------")
	fmt.Print("------------------------------------------------------")
	k.SetIndividualUserPoolsData(ctx, userPoolsData)


	return &types.MsgUnbondPoolTokensResponse{}, nil
}
