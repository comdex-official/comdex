package keeper

import (
	"fmt"
	"time"

	"github.com/comdex-official/comdex/x/incentives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateMsgCreateCreateGauge validates types.MsgCreateGauge.
func (k Keeper) ValidateMsgCreateCreateGauge(ctx sdk.Context, msg *types.MsgCreateGauge) error {
	isValidGaugeTypeID := false
	for _, gaugeTypeID := range types.ValidGaugeTypeIds {
		if gaugeTypeID == msg.GaugeTypeId {
			isValidGaugeTypeID = true
			break
		}
	}

	if !isValidGaugeTypeID {
		return types.ErrInvalidGaugeTypeID
	}

	if msg.TriggerDuration <= 0 {
		return types.ErrInvalidDuration
	}

	if msg.DepositAmount.Amount.IsNegative() || msg.DepositAmount.Amount.IsZero() {
		return types.ErrInvalidDepositAmount
	}

	if msg.DepositAmount.Amount.LT(sdk.NewIntFromUint64(msg.TotalTriggers)) {
		return types.ErrDepositSmallThanEpoch
	}

	if msg.StartTime.Before(ctx.BlockTime()) {
		return types.ErrInvalidGaugeStartTime
	}

	return nil
}

// ValidateMsgCreateGaugeLiquidityMetaData validates types.MsgCreateGauge_LiquidityMetaData.
func (k Keeper) ValidateMsgCreateGaugeLiquidityMetaData(ctx sdk.Context, kind *types.MsgCreateGauge_LiquidityMetaData) error {
	if kind.LiquidityMetaData.LockDuration <= 0 {
		return types.ErrInvalidDuration
	}

	_, found := k.liquidityKeeper.GetPool(ctx, kind.LiquidityMetaData.PoolId)
	if !found {
		return types.ErrInvalidPoolID
	}

	childPoolIds := kind.LiquidityMetaData.ChildPoolIds
	for _, poolID := range childPoolIds {
		_, found := k.liquidityKeeper.GetPool(ctx, poolID)
		if !found {
			return sdkerrors.Wrap(types.ErrInvalidPoolID, fmt.Sprintf("invalid child pool id : %d", poolID))
		}
	}

	return nil
}

// NewGauge returns the new Gauge object.
func (k Keeper) NewGauge(ctx sdk.Context, msg *types.MsgCreateGauge) (types.Gauge, error) {
	newGauge := types.Gauge{
		Id:                k.GetGaugeID(ctx) + 1,
		From:              msg.From,
		CreatedAt:         ctx.BlockTime(),
		StartTime:         msg.StartTime,
		GaugeTypeId:       msg.GaugeTypeId,
		TriggerDuration:   msg.TriggerDuration,
		DepositAmount:     msg.DepositAmount,
		TotalTriggers:     msg.TotalTriggers,
		TriggeredCount:    0,
		DistributedAmount: sdk.NewCoin(msg.DepositAmount.Denom, sdk.NewInt(0)),
		IsActive:          true,
		Kind:              nil,
	}

	switch extraData := msg.Kind.(type) {
	case *types.MsgCreateGauge_LiquidityMetaData:

		err := k.ValidateMsgCreateGaugeLiquidityMetaData(ctx, extraData)
		if err != nil {
			return types.Gauge{}, err
		}

		liquidityGaugeType := &types.Gauge_LiquidityMetaData{
			LiquidityMetaData: &types.LiquidtyGaugeMetaData{
				PoolId:       extraData.LiquidityMetaData.PoolId,
				IsMasterPool: extraData.LiquidityMetaData.IsMasterPool,
				ChildPoolIds: extraData.LiquidityMetaData.ChildPoolIds,
				LockDuration: extraData.LiquidityMetaData.LockDuration,
			},
		}
		newGauge.Kind = liquidityGaugeType
	}
	return newGauge, nil
}

// NewGaugeIdsByDuration return new GaugeByTriggerDuration.
func (k Keeper) NewGaugeIdsByDuration(ctx sdk.Context, duration time.Duration) types.GaugeByTriggerDuration {
	return types.GaugeByTriggerDuration{
		TriggerDuration: duration,
		GaugeIds:        []uint64{},
	}
}

// GetUpdatedGaugeIdsByTriggerDurationObj returns gauge id by duration.
func (k Keeper) GetUpdatedGaugeIdsByTriggerDurationObj(ctx sdk.Context, triggerDuration time.Duration, newGaugeID uint64) (types.GaugeByTriggerDuration, error) {
	gaugeIdsByTriggerDuration, found := k.GetGaugeIdsByTriggerDuration(ctx, triggerDuration)

	if !found {
		gaugeIdsByTriggerDuration = k.NewGaugeIdsByDuration(ctx, triggerDuration)
	}
	gaugeIDAlreadyExists := false

	for _, gaugeID := range gaugeIdsByTriggerDuration.GaugeIds {
		if gaugeID == newGaugeID {
			gaugeIDAlreadyExists = true
		}
	}

	if gaugeIDAlreadyExists {
		return types.GaugeByTriggerDuration{}, sdkerrors.Wrapf(types.ErrInvalidGaugeID, "gauge id already exists in map : %d", newGaugeID)
	}
	gaugeIdsByTriggerDuration.GaugeIds = append(gaugeIdsByTriggerDuration.GaugeIds, newGaugeID)
	return gaugeIdsByTriggerDuration, nil
}

// InitateGaugesForDuration triggers the gauge in the event of triggerDuration.
func (k Keeper) InitateGaugesForDuration(ctx sdk.Context, triggerDuration time.Duration) error {
	logger := k.Logger(ctx)
	gaugesForDuration, found := k.GetGaugeIdsByTriggerDuration(ctx, triggerDuration)
	if !found {
		return sdkerrors.Wrapf(types.ErrNoGaugeForDuration, "duration : %d", triggerDuration)
	}

	for _, gaugeID := range gaugesForDuration.GaugeIds {
		gauge, found := k.GetGaugeByID(ctx, gaugeID)
		if !found {
			continue
		}
		if ctx.BlockTime().Before(gauge.StartTime) || !gauge.IsActive {
			continue
		}

		if gauge.TriggeredCount == gauge.TotalTriggers {
			gauge.IsActive = false
			k.SetGauge(ctx, gauge)
			continue
		}

		depositAmountSplitsByEpochs := SplitTotalAmountPerEpoch(gauge.DepositAmount.Amount.Uint64(), gauge.TotalTriggers)
		if len(depositAmountSplitsByEpochs) <= int(gauge.TriggeredCount) {
			logger.Info("triggered counts are higher than total trigger splits, exceptions avoided")
			continue
		}
		amountToDistribute := depositAmountSplitsByEpochs[gauge.TriggeredCount]
		availableDeposits := gauge.DepositAmount.Amount.Sub(gauge.DistributedAmount.Amount)

		// just in case (exception handelled), but this will never pass
		if availableDeposits.LT(sdk.NewIntFromUint64(amountToDistribute)) {
			continue
		}

		ongoingEpochCount := gauge.TriggeredCount + 1
		coinToDistribute := sdk.NewCoin(gauge.DepositAmount.Denom, sdk.NewIntFromUint64(amountToDistribute))

		coinsDistributed, err := k.BeginRewardDistributions(ctx, gauge, coinToDistribute, ongoingEpochCount, triggerDuration)
		if err != nil {
			logger.Info(fmt.Sprintf("error occurred while reward distribution in BeginRewardDistributions, err : %s", err))
			continue
		}
		gauge.TriggeredCount = ongoingEpochCount
		gauge.DistributedAmount.Amount = gauge.DistributedAmount.Amount.Add(coinsDistributed.Amount)
		k.SetGauge(ctx, gauge)
	}

	return nil
}
