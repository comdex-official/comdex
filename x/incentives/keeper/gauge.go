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
	isValidGaugeTypeId := false
	for _, gaugeTypeId := range types.ValidGaugeTypeIds {
		if gaugeTypeId == msg.GaugeTypeId {
			isValidGaugeTypeId = true
			break
		}
	}

	if !isValidGaugeTypeId {
		return types.ErrInvalidGaugeTypeId
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

func (k Keeper) ValidateMsgCreateGauge_LiquidityMetaData(ctx sdk.Context, kind *types.MsgCreateGauge_LiquidityMetaData) error {
	if kind.LiquidityMetaData.LockDuration <= 0 {
		return types.ErrInvalidDuration
	}

	_, found := k.liquidityKeeper.GetPool(ctx, kind.LiquidityMetaData.PoolId)
	if !found {
		return types.ErrInvalidPoolId
	}

	childPoolIds := kind.LiquidityMetaData.ChildPoolIds
	for _, poolId := range childPoolIds {
		_, found := k.liquidityKeeper.GetPool(ctx, poolId)
		if !found {
			return sdkerrors.Wrap(types.ErrInvalidPoolId, fmt.Sprintf("invalid child pool id : %d", poolId))
		}
	}

	return nil
}

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

		err := k.ValidateMsgCreateGauge_LiquidityMetaData(ctx, extraData)
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

func (k Keeper) NewGaugeIdsByDuration(ctx sdk.Context, duration time.Duration) types.GaugeByTriggerDuration {
	return types.GaugeByTriggerDuration{
		TriggerDuration: duration,
		GaugeIds:        []uint64{},
	}
}

func (k Keeper) GetUpdatedGaugeIdsByTriggerDurationObj(ctx sdk.Context, triggerDuration time.Duration, newGaugeId uint64) (types.GaugeByTriggerDuration, error) {

	gaugeIdsByTriggerDuration, found := k.GetGaugeIdsByTriggerDuration(ctx, triggerDuration)

	if !found {
		gaugeIdsByTriggerDuration = k.NewGaugeIdsByDuration(ctx, triggerDuration)
	}
	gaugeIdAlreadyExists := false

	for _, gaugeId := range gaugeIdsByTriggerDuration.GaugeIds {
		if gaugeId == newGaugeId {
			gaugeIdAlreadyExists = true
		}
	}

	if gaugeIdAlreadyExists {
		return types.GaugeByTriggerDuration{}, sdkerrors.Wrapf(types.ErrInvalidGaugeId, "gauge id already exists in map : %d", newGaugeId)
	}
	gaugeIdsByTriggerDuration.GaugeIds = append(gaugeIdsByTriggerDuration.GaugeIds, newGaugeId)
	return gaugeIdsByTriggerDuration, nil
}

func (k Keeper) InitateGaugesForDuration(ctx sdk.Context, triggerDuration time.Duration) error {
	logger := k.Logger(ctx)
	gaugesForDuration, found := k.GetGaugeIdsByTriggerDuration(ctx, triggerDuration)
	if !found {
		return sdkerrors.Wrapf(types.ErrNoGaugeForDuration, "duration : %d", triggerDuration)
	}

	for _, gaugeId := range gaugesForDuration.GaugeIds {
		gauge, found := k.GetGaugeById(ctx, gaugeId)
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
			logger.Info(fmt.Sprintf("error occured while reward distribution in BeginRewardDistributions, err : %s", err))
			continue
		}
		gauge.TriggeredCount = ongoingEpochCount
		gauge.DistributedAmount.Amount = gauge.DistributedAmount.Amount.Add(coinsDistributed.Amount)
		k.SetGauge(ctx, gauge)
	}

	return nil
}
