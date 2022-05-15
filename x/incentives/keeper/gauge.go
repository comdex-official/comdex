package keeper

import (
	"fmt"

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
	return nil
}

func (k Keeper) ValidateMsgCreateGauge_LiquidityMetaData(ctx sdk.Context, kind *types.MsgCreateGauge_LiquidityMetaData) error {
	if !kind.LiquidityMetaData.StartTime.After(ctx.BlockTime()) {
		return types.ErrLiquidityGaugeStartTime
	}
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
		GaugeTypeId:       msg.GaugeTypeId,
		TriggerDuration:   msg.TriggerDuration,
		DepositAmount:     msg.DepositAmount,
		TotalTriggers:     msg.TotalTriggers,
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
				StartTime:    extraData.LiquidityMetaData.StartTime,
				LockDuration: extraData.LiquidityMetaData.LockDuration,
			},
		}
		newGauge.Kind = liquidityGaugeType
	}
	return newGauge, nil
}
