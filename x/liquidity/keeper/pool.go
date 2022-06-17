package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
)

// getNextPoolIdWithUpdate increments pool id by one and set it.
func (k Keeper) getNextPoolIDWithUpdate(ctx sdk.Context, appID uint64) uint64 {
	id := k.GetLastPoolID(ctx, appID) + 1
	k.SetLastPoolID(ctx, appID, id)
	return id
}

// getNextDepositRequestIdWithUpdate increments the pool's last deposit request
// id and returns it.
func (k Keeper) getNextDepositRequestIDWithUpdate(ctx sdk.Context, pool types.Pool) uint64 {
	id := pool.LastDepositRequestId + 1
	pool.LastDepositRequestId = id
	k.SetPool(ctx, pool)
	return id
}

// getNextWithdrawRequestIdWithUpdate increments the pool's last withdraw
// request id and returns it.
func (k Keeper) getNextWithdrawRequestIDWithUpdate(ctx sdk.Context, pool types.Pool) uint64 {
	id := pool.LastWithdrawRequestId + 1
	pool.LastWithdrawRequestId = id
	k.SetPool(ctx, pool)
	return id
}

// GetPoolBalances returns the balances of the pool.
func (k Keeper) GetPoolBalances(ctx sdk.Context, pool types.Pool) (rx sdk.Coin, ry sdk.Coin) {
	reserveAddr := pool.GetReserveAddress()
	pair, _ := k.GetPair(ctx, pool.AppId, pool.PairId)
	spendable := k.bankKeeper.SpendableCoins(ctx, reserveAddr)
	rx = sdk.NewCoin(pair.QuoteCoinDenom, spendable.AmountOf(pair.QuoteCoinDenom))
	ry = sdk.NewCoin(pair.BaseCoinDenom, spendable.AmountOf(pair.BaseCoinDenom))
	return
}

// getPoolBalances returns the balances of the pool.
// It is used internally when caller already has types.Pair instance.
func (k Keeper) getPoolBalances(ctx sdk.Context, pool types.Pool, pair types.Pair) (rx sdk.Coin, ry sdk.Coin) {
	reserveAddr := pool.GetReserveAddress()
	spendable := k.bankKeeper.SpendableCoins(ctx, reserveAddr)
	rx = sdk.NewCoin(pair.QuoteCoinDenom, spendable.AmountOf(pair.QuoteCoinDenom))
	ry = sdk.NewCoin(pair.BaseCoinDenom, spendable.AmountOf(pair.BaseCoinDenom))
	return
}

// GetPoolCoinSupply returns total pool coin supply of the pool.
func (k Keeper) GetPoolCoinSupply(ctx sdk.Context, pool types.Pool) sdk.Int {
	return k.bankKeeper.GetSupply(ctx, pool.PoolCoinDenom).Amount
}

// MarkPoolAsDisabled marks a pool as disabled.
func (k Keeper) MarkPoolAsDisabled(ctx sdk.Context, pool types.Pool) {
	pool.Disabled = true
	k.SetPool(ctx, pool)
}

// ValidateMsgCreatePool validates types.MsgCreatePool.
func (k Keeper) ValidateMsgCreatePool(ctx sdk.Context, msg *types.MsgCreatePool) error {
	_, found := k.assetKeeper.GetApp(ctx, msg.AppId)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", msg.AppId)
	}

	pair, found := k.GetPair(ctx, msg.AppId, msg.PairId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	params, err := k.GetGenericParams(ctx, msg.AppId)
	if err != nil {
		return sdkerrors.Wrap(err, "params retreval failed")
	}
	for _, coin := range msg.DepositCoins {
		if coin.Denom != pair.BaseCoinDenom && coin.Denom != pair.QuoteCoinDenom {
			return sdkerrors.Wrapf(types.ErrInvalidCoinDenom, "coin denom %s is not in the pair", coin.Denom)
		}
		minDepositCoin := sdk.NewCoin(coin.Denom, params.MinInitialDepositAmount)
		if coin.IsLT(minDepositCoin) {
			return sdkerrors.Wrapf(
				types.ErrInsufficientDepositAmount, "%s is smaller than %s", coin, minDepositCoin)
		}
	}

	// Check if there is a pool in the pair.
	// Creating multiple pools within the same pair is disallowed, but it will be allowed in v2.
	duplicate := false
	_ = k.IteratePoolsByPair(ctx, pair.AppId, pair.Id, func(pool types.Pool) (stop bool, err error) {
		if !pool.Disabled {
			duplicate = true
			return true, nil
		}
		return false, nil
	})
	if duplicate {
		return types.ErrPoolAlreadyExists
	}

	return nil
}

// CreatePool handles types.MsgCreatePool and creates a pool.
func (k Keeper) CreatePool(ctx sdk.Context, msg *types.MsgCreatePool) (types.Pool, error) {
	if err := k.ValidateMsgCreatePool(ctx, msg); err != nil {
		return types.Pool{}, err
	}

	params, err := k.GetGenericParams(ctx, msg.AppId)
	if err != nil {
		return types.Pool{}, sdkerrors.Wrap(err, "params retreval failed")
	}
	pair, _ := k.GetPair(ctx, msg.AppId, msg.PairId)

	// Create and save the new pool object.
	poolID := k.getNextPoolIDWithUpdate(ctx, msg.AppId)
	pool := types.NewPool(msg.AppId, poolID, pair.Id)
	k.SetPool(ctx, pool)
	k.SetPoolByReserveIndex(ctx, pool)
	k.SetPoolsByPairIndex(ctx, pool)

	// Send deposit coins to the pool's reserve account.
	creator := msg.GetCreator()
	if err := k.bankKeeper.SendCoins(ctx, creator, pool.GetReserveAddress(), msg.DepositCoins); err != nil {
		return types.Pool{}, err
	}

	// Send the pool creation fee to the fee collector.
	feeCollectorAddr, _ := sdk.AccAddressFromBech32(params.FeeCollectorAddress)
	if err := k.bankKeeper.SendCoins(ctx, creator, feeCollectorAddr, params.PoolCreationFee); err != nil {
		return types.Pool{}, sdkerrors.Wrap(err, "insufficient pool creation fee")
	}

	// Mint and send pool coin to the creator.
	// Minting pool coin amount is calculated based on two coins' amount.
	// Minimum minting amount is params.MinInitialPoolCoinSupply.
	ps := sdk.MaxInt(
		amm.InitialPoolCoinSupply(msg.DepositCoins[0].Amount, msg.DepositCoins[1].Amount),
		params.MinInitialPoolCoinSupply,
	)
	poolCoin := sdk.NewCoin(pool.PoolCoinDenom, ps)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(poolCoin)); err != nil {
		return types.Pool{}, err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, sdk.NewCoins(poolCoin)); err != nil {
		return types.Pool{}, err
	}

	newGauge := rewardstypes.NewMsgCreateGauge(
		msg.AppId,
		pair.GetSwapFeeCollectorAddress(),
		ctx.BlockTime(),
		rewardstypes.LiquidityGaugeTypeID,
		types.DefaultSwapFeeDistributionDuration,
		sdk.NewCoin(params.SwapFeeDistrDenom, sdk.NewInt(0)),
		1,
	)
	newGauge.Kind = &rewardstypes.MsgCreateGauge_LiquidityMetaData{
		LiquidityMetaData: &rewardstypes.LiquidtyGaugeMetaData{
			PoolId:       pool.Id,
			IsMasterPool: false,
			ChildPoolIds: []uint64{},
		},
	}
	err = k.rewardsKeeper.CreateNewGauge(ctx, newGauge, true)
	if err != nil {
		return types.Pool{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(msg.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyDepositCoins, msg.DepositCoins.String()),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(pool.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyReserveAddress, pool.ReserveAddress),
			sdk.NewAttribute(types.AttributeKeySwapFeeCollectorAddress, pair.SwapFeeCollectorAddress),
			sdk.NewAttribute(types.AttributeKeyMintedPoolCoin, poolCoin.String()),
		),
	})

	return pool, nil
}

// ValidateMsgDeposit validates types.MsgDeposit.
func (k Keeper) ValidateMsgDeposit(ctx sdk.Context, msg *types.MsgDeposit) error {
	_, found := k.assetKeeper.GetApp(ctx, msg.AppId)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", msg.AppId)
	}

	pool, found := k.GetPool(ctx, msg.AppId, msg.PoolId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pool %d not found", msg.PoolId)
	}
	if pool.Disabled {
		return types.ErrDisabledPool
	}

	pair, _ := k.GetPair(ctx, msg.AppId, pool.PairId)

	for _, coin := range msg.DepositCoins {
		if coin.Denom != pair.BaseCoinDenom && coin.Denom != pair.QuoteCoinDenom {
			return sdkerrors.Wrapf(types.ErrInvalidCoinDenom, "coin denom %s is not in the pair", coin.Denom)
		}
	}

	rx, ry := k.getPoolBalances(ctx, pool, pair)
	if rx.Amount.Add(msg.DepositCoins.AmountOf(rx.Denom)).GT(amm.MaxCoinAmount) {
		return types.ErrTooLargePool
	}
	if ry.Amount.Add(msg.DepositCoins.AmountOf(ry.Denom)).GT(amm.MaxCoinAmount) {
		return types.ErrTooLargePool
	}

	return nil
}

// Deposit handles types.MsgDeposit and stores the request.
func (k Keeper) Deposit(ctx sdk.Context, msg *types.MsgDeposit) (types.DepositRequest, error) {
	if err := k.ValidateMsgDeposit(ctx, msg); err != nil {
		return types.DepositRequest{}, err
	}

	params, err := k.GetGenericParams(ctx, msg.AppId)
	if err != nil {
		return types.DepositRequest{}, sdkerrors.Wrap(err, "params retreval failed")
	}

	if err := k.bankKeeper.SendCoins(ctx, msg.GetDepositor(), types.GlobalEscrowAddress, msg.DepositCoins); err != nil {
		return types.DepositRequest{}, err
	}

	pool, _ := k.GetPool(ctx, msg.AppId, msg.PoolId)
	requestID := k.getNextDepositRequestIDWithUpdate(ctx, pool)
	req := types.NewDepositRequest(msg, pool, requestID, ctx.BlockHeight())
	k.SetDepositRequest(ctx, req)
	k.SetDepositRequestIndex(ctx, req)

	ctx.GasMeter().ConsumeGas(params.DepositExtraGas, "DepositExtraGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeposit,
			sdk.NewAttribute(types.AttributeKeyDepositor, msg.Depositor),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(pool.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyDepositCoins, msg.DepositCoins.String()),
			sdk.NewAttribute(types.AttributeKeyRequestID, strconv.FormatUint(req.Id, 10)),
		),
	})

	return req, nil
}

// ValidateMsgWithdraw validates types.MsgWithdraw.
func (k Keeper) ValidateMsgWithdraw(ctx sdk.Context, msg *types.MsgWithdraw) error {
	_, found := k.assetKeeper.GetApp(ctx, msg.AppId)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", msg.AppId)
	}

	pool, found := k.GetPool(ctx, msg.AppId, msg.PoolId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pool %d not found", msg.PoolId)
	}
	if pool.Disabled {
		return types.ErrDisabledPool
	}

	if msg.PoolCoin.Denom != pool.PoolCoinDenom {
		return types.ErrWrongPoolCoinDenom
	}

	return nil
}

// Withdraw handles types.MsgWithdraw and stores the request.
func (k Keeper) Withdraw(ctx sdk.Context, msg *types.MsgWithdraw) (types.WithdrawRequest, error) {
	if err := k.ValidateMsgWithdraw(ctx, msg); err != nil {
		return types.WithdrawRequest{}, err
	}

	params, err := k.GetGenericParams(ctx, msg.AppId)
	if err != nil {
		return types.WithdrawRequest{}, sdkerrors.Wrap(err, "params retreval failed")
	}

	pool, _ := k.GetPool(ctx, msg.AppId, msg.PoolId)
	if err := k.bankKeeper.SendCoins(ctx, msg.GetWithdrawer(), types.GlobalEscrowAddress, sdk.NewCoins(msg.PoolCoin)); err != nil {
		return types.WithdrawRequest{}, err
	}

	requestID := k.getNextWithdrawRequestIDWithUpdate(ctx, pool)
	req := types.NewWithdrawRequest(msg, requestID, ctx.BlockHeight())
	k.SetWithdrawRequest(ctx, req)
	k.SetWithdrawRequestIndex(ctx, req)

	ctx.GasMeter().ConsumeGas(params.WithdrawExtraGas, "WithdrawExtraGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdraw,
			sdk.NewAttribute(types.AttributeKeyWithdrawer, msg.Withdrawer),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(pool.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyPoolCoin, msg.PoolCoin.String()),
			sdk.NewAttribute(types.AttributeKeyRequestID, strconv.FormatUint(req.Id, 10)),
		),
	})

	return req, nil
}

// ExecuteDepositRequest executes a deposit request.
func (k Keeper) ExecuteDepositRequest(ctx sdk.Context, req types.DepositRequest) error {
	pool, _ := k.GetPool(ctx, req.AppId, req.PoolId)
	if pool.Disabled {
		if err := k.FinishDepositRequest(ctx, req, types.RequestStatusFailed); err != nil {
			return fmt.Errorf("refund deposit request: %w", err)
		}
		return nil
	}

	pair, _ := k.GetPair(ctx, req.AppId, pool.PairId)
	rx, ry := k.getPoolBalances(ctx, pool, pair)
	ps := k.GetPoolCoinSupply(ctx, pool)
	ammPool := amm.NewBasicPool(rx.Amount, ry.Amount, ps)
	if ammPool.IsDepleted() {
		k.MarkPoolAsDisabled(ctx, pool)
		if err := k.FinishDepositRequest(ctx, req, types.RequestStatusFailed); err != nil {
			return err
		}
		return nil
	}

	ax, ay, pc := ammPool.Deposit(req.DepositCoins.AmountOf(pair.QuoteCoinDenom), req.DepositCoins.AmountOf(pair.BaseCoinDenom))

	if pc.IsZero() {
		if err := k.FinishDepositRequest(ctx, req, types.RequestStatusFailed); err != nil {
			return err
		}
		return nil
	}

	mintedPoolCoin := sdk.NewCoin(pool.PoolCoinDenom, pc)
	mintingCoins := sdk.NewCoins(mintedPoolCoin)

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintingCoins); err != nil {
		return err
	}

	acceptedCoins := sdk.NewCoins(sdk.NewCoin(pair.QuoteCoinDenom, ax), sdk.NewCoin(pair.BaseCoinDenom, ay))
	bulkOp := types.NewBulkSendCoinsOperation()
	bulkOp.QueueSendCoins(types.GlobalEscrowAddress, pool.GetReserveAddress(), acceptedCoins)
	bulkOp.QueueSendCoins(k.accountKeeper.GetModuleAddress(types.ModuleName), req.GetDepositor(), mintingCoins)
	if err := bulkOp.Run(ctx, k.bankKeeper); err != nil {
		return err
	}

	req.AcceptedCoins = acceptedCoins
	req.MintedPoolCoin = mintedPoolCoin
	if err := k.FinishDepositRequest(ctx, req, types.RequestStatusSucceeded); err != nil {
		return err
	}
	return nil
}

// FinishDepositRequest refunds unhandled deposit coins and set request status.
func (k Keeper) FinishDepositRequest(ctx sdk.Context, req types.DepositRequest, status types.RequestStatus) error {
	if req.Status != types.RequestStatusNotExecuted { // sanity check
		return nil
	}

	refundingCoins := req.DepositCoins.Sub(req.AcceptedCoins)
	if !refundingCoins.IsZero() {
		if err := k.bankKeeper.SendCoins(ctx, types.GlobalEscrowAddress, req.GetDepositor(), refundingCoins); err != nil {
			return err
		}
	}
	req.SetStatus(status)
	k.SetDepositRequest(ctx, req)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDepositResult,
			sdk.NewAttribute(types.AttributeKeyRequestID, strconv.FormatUint(req.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyDepositor, req.Depositor),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(req.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyDepositCoins, req.DepositCoins.String()),
			sdk.NewAttribute(types.AttributeKeyAcceptedCoins, req.AcceptedCoins.String()),
			sdk.NewAttribute(types.AttributeKeyRefundedCoins, refundingCoins.String()),
			sdk.NewAttribute(types.AttributeKeyMintedPoolCoin, req.MintedPoolCoin.String()),
			sdk.NewAttribute(types.AttributeKeyStatus, req.Status.String()),
		),
	})

	return nil
}

// ExecuteWithdrawRequest executes a withdraw request.
func (k Keeper) ExecuteWithdrawRequest(ctx sdk.Context, req types.WithdrawRequest) error {
	params, err := k.GetGenericParams(ctx, req.AppId)
	if err != nil {
		return sdkerrors.Wrap(err, "params retreval failed")
	}

	pool, _ := k.GetPool(ctx, req.AppId, req.PoolId)
	if pool.Disabled {
		if err := k.FinishWithdrawRequest(ctx, req, types.RequestStatusFailed); err != nil {
			return err
		}
		return nil
	}

	pair, _ := k.GetPair(ctx, req.AppId, pool.PairId)
	rx, ry := k.getPoolBalances(ctx, pool, pair)
	ps := k.GetPoolCoinSupply(ctx, pool)
	ammPool := amm.NewBasicPool(rx.Amount, ry.Amount, ps)
	if ammPool.IsDepleted() {
		k.MarkPoolAsDisabled(ctx, pool)
		if err := k.FinishWithdrawRequest(ctx, req, types.RequestStatusFailed); err != nil {
			return err
		}
		return nil
	}

	x, y := ammPool.Withdraw(req.PoolCoin.Amount, params.WithdrawFeeRate)
	if x.IsZero() && y.IsZero() {
		if err := k.FinishWithdrawRequest(ctx, req, types.RequestStatusFailed); err != nil {
			return err
		}
		return nil
	}

	withdrawnCoins := sdk.NewCoins(sdk.NewCoin(pair.QuoteCoinDenom, x), sdk.NewCoin(pair.BaseCoinDenom, y))
	burningCoins := sdk.NewCoins(req.PoolCoin)

	bulkOp := types.NewBulkSendCoinsOperation()
	bulkOp.QueueSendCoins(types.GlobalEscrowAddress, k.accountKeeper.GetModuleAddress(types.ModuleName), burningCoins)
	bulkOp.QueueSendCoins(pool.GetReserveAddress(), req.GetWithdrawer(), withdrawnCoins)
	if err := bulkOp.Run(ctx, k.bankKeeper); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burningCoins); err != nil {
		return err
	}

	// If the pool coin supply becomes 0, disable the pool.
	if req.PoolCoin.Amount.Equal(ps) {
		k.MarkPoolAsDisabled(ctx, pool)
	}

	req.WithdrawnCoins = withdrawnCoins
	if err := k.FinishWithdrawRequest(ctx, req, types.RequestStatusSucceeded); err != nil {
		return err
	}
	return nil
}

// FinishWithdrawRequest refunds unhandled pool coin and set request status.
func (k Keeper) FinishWithdrawRequest(ctx sdk.Context, req types.WithdrawRequest, status types.RequestStatus) error {
	if req.Status != types.RequestStatusNotExecuted { // sanity check
		return nil
	}

	var refundingCoins sdk.Coins
	if status == types.RequestStatusFailed {
		refundingCoins = sdk.NewCoins(req.PoolCoin)
		if err := k.bankKeeper.SendCoins(ctx, types.GlobalEscrowAddress, req.GetWithdrawer(), refundingCoins); err != nil {
			return err
		}
	}
	req.SetStatus(status)
	k.SetWithdrawRequest(ctx, req)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawalResult,
			sdk.NewAttribute(types.AttributeKeyRequestID, strconv.FormatUint(req.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyWithdrawer, req.Withdrawer),
			sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(req.PoolId, 10)),
			sdk.NewAttribute(types.AttributeKeyPoolCoin, req.PoolCoin.String()),
			sdk.NewAttribute(types.AttributeKeyRefundedCoins, refundingCoins.String()),
			sdk.NewAttribute(types.AttributeKeyWithdrawnCoins, req.WithdrawnCoins.String()),
			sdk.NewAttribute(types.AttributeKeyStatus, req.Status.String()),
		),
	})

	return nil
}

func (k Keeper) TransferFundsForSwapFeeDistribution(ctx sdk.Context, appID, poolId uint64) (sdk.Coin, error) {
	pool, found := k.GetPool(ctx, appID, poolId)
	if !found {
		return sdk.Coin{}, types.ErrInvalidPoolID
	}

	pair, found := k.GetPair(ctx, appID, pool.PairId)
	if !found {
		return sdk.Coin{}, types.ErrInvalidPairId
	}

	params, err := k.GetGenericParams(ctx, appID)
	if err != nil {
		return sdk.Coin{}, sdkerrors.Wrap(err, "params retreval failed")
	}

	availableBalance := k.bankKeeper.GetBalance(ctx, pair.GetSwapFeeCollectorAddress(), params.SwapFeeDistrDenom)

	burnAmount := availableBalance.Amount.ToDec().MulTruncate(params.SwapFeeBurnRate).TruncateInt()
	burnCoin := sdk.NewCoin(availableBalance.Denom, burnAmount)

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, pair.GetSwapFeeCollectorAddress(), types.ModuleName, sdk.NewCoins(burnCoin))
	if err != nil {
		return sdk.Coin{}, err
	}

	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burnCoin))
	if err != nil {
		return sdk.Coin{}, err
	}

	availableBalance.Amount = availableBalance.Amount.Sub(burnCoin.Amount)

	if availableBalance.IsPositive() {
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, pair.GetSwapFeeCollectorAddress(), rewardstypes.ModuleName, sdk.NewCoins(availableBalance))
		if err != nil {
			return sdk.Coin{}, err
		}
	} else {
		// negative amount handalling
		availableBalance.Amount = sdk.NewInt(0)
	}
	return availableBalance, nil
}
