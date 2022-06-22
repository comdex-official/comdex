package keeper

import (
	"fmt"
	"time"

	"github.com/comdex-official/comdex/x/rewards/types"
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

func (k Keeper) OraclePrice(ctx sdk.Context, denom string) (uint64, bool) {
	asset, found := k.asset.GetAssetForDenom(ctx, denom)
	if !found {
		return 0, false
	}

	price, found := k.marketKeeper.GetPriceForAsset(ctx, asset.Id)
	if !found {
		return 0, false
	}
	return price, true
}

func (k Keeper) ValidateIfOraclePricesExists(ctx sdk.Context, appID, pairID uint64) error {
	pair, found := k.liquidityKeeper.GetPair(ctx, appID, pairID)
	if !found {
		return sdkerrors.Wrapf(types.ErrPairNotExists, "pair does not exists for given pool id")
	}

	_, baseCoinPriceFound := k.OraclePrice(ctx, pair.BaseCoinDenom)
	_, quoteCoinPricefound := k.OraclePrice(ctx, pair.QuoteCoinDenom)
	if !(baseCoinPriceFound || quoteCoinPricefound) {
		return sdkerrors.Wrapf(types.ErrPriceNotFound, "oracle price required for atleast %s or %s but not found", pair.QuoteCoinDenom, pair.BaseCoinDenom)
	}

	return nil
}

// ValidateMsgCreateGaugeLiquidityMetaData validates types.MsgCreateGauge_LiquidityMetaData.
func (k Keeper) ValidateMsgCreateGaugeLiquidityMetaData(ctx sdk.Context, appID uint64, kind *types.MsgCreateGauge_LiquidityMetaData, forSwapFee bool) error {
	_, found := k.asset.GetApp(ctx, appID)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", appID)
	}

	pool, found := k.liquidityKeeper.GetPool(ctx, appID, kind.LiquidityMetaData.PoolId)
	if !found {
		return types.ErrInvalidPoolID
	}

	if !forSwapFee {
		err := k.ValidateIfOraclePricesExists(ctx, appID, pool.PairId)
		if err != nil {
			return err
		}
	}

	childPoolIds := kind.LiquidityMetaData.ChildPoolIds
	for _, poolID := range childPoolIds {
		_, found := k.liquidityKeeper.GetPool(ctx, appID, poolID)
		if !found {
			return sdkerrors.Wrap(types.ErrInvalidPoolID, fmt.Sprintf("invalid child pool id : %d", poolID))
		}
	}

	return nil
}

// NewGauge returns the new Gauge object.
func (k Keeper) NewGauge(ctx sdk.Context, msg *types.MsgCreateGauge, forSwapFee bool) (types.Gauge, error) {
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
		ForSwapFee:        false,
		Kind:              nil,
		AppId:             msg.AppId,
	}

	switch extraData := msg.Kind.(type) {
	case *types.MsgCreateGauge_LiquidityMetaData:

		err := k.ValidateMsgCreateGaugeLiquidityMetaData(ctx, msg.AppId, extraData, forSwapFee)
		if err != nil {
			return types.Gauge{}, err
		}

		liquidityGaugeType := &types.Gauge_LiquidityMetaData{
			LiquidityMetaData: &types.LiquidtyGaugeMetaData{
				PoolId:       extraData.LiquidityMetaData.PoolId,
				IsMasterPool: extraData.LiquidityMetaData.IsMasterPool,
				ChildPoolIds: extraData.LiquidityMetaData.ChildPoolIds,
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

func (k Keeper) CreateNewGauge(ctx sdk.Context, msg *types.MsgCreateGauge, forSwapFee bool) error {
	newGauge, err := k.NewGauge(ctx, msg, forSwapFee)
	if err != nil {
		return err
	}
	newGauge.ForSwapFee = forSwapFee

	gaugeIdsByTriggerDuration, err := k.GetUpdatedGaugeIdsByTriggerDurationObj(ctx, newGauge.TriggerDuration, newGauge.Id)
	if err != nil {
		return err
	}

	from, _ := sdk.AccAddressFromBech32(newGauge.From)
	err = k.bank.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(newGauge.DepositAmount))
	if err != nil {
		return err
	}

	_, found := k.GetEpochInfoByDuration(ctx, newGauge.TriggerDuration)
	if !found {
		newEpochInfo := k.NewEpochInfo(ctx, newGauge.TriggerDuration)
		k.SetEpochInfoByDuration(ctx, newEpochInfo)
	}

	k.SetGaugeID(ctx, newGauge.Id)
	k.SetGauge(ctx, newGauge)
	k.SetGaugeIdsByTriggerDuration(ctx, gaugeIdsByTriggerDuration)
	return nil
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
		if !gauge.ForSwapFee {
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
		} else {
			// distribution method for swap fee
			poolID := gauge.GetLiquidityMetaData().PoolId

			ongoingEpochCount := gauge.TriggeredCount + 1
			coinToDistribute := gauge.DepositAmount

			// distributing swap fees amount which was accumulated at epoch-1
			if coinToDistribute.IsPositive() {
				coinsDistributed, err := k.BeginRewardDistributions(ctx, gauge, coinToDistribute, ongoingEpochCount, triggerDuration)
				if err != nil {
					logger.Info(fmt.Sprintf("error occurred while reward distribution in BeginRewardDistributions, err : %s", err))
					continue
				}

				gauge.DepositAmount = gauge.DepositAmount.Sub(coinsDistributed)

				// in case of swap fee distribution denom change in params
				if gauge.DistributedAmount.Denom == coinsDistributed.Denom {
					gauge.DistributedAmount.Amount = gauge.DistributedAmount.Amount.Add(coinsDistributed.Amount)
				} else {
					gauge.DistributedAmount = coinsDistributed
				}
			}

			// transferring swap fees amount of current epoch from swap fee collector address to rewards module
			receivedAmount, err := k.liquidityKeeper.TransferFundsForSwapFeeDistribution(ctx, gauge.AppId, poolID)
			if err != nil {
				logger.Info(fmt.Sprintf("error occurred while swap fee fund transfer, err : %s", err))
				continue
			}
			// in case of swap fee distribution denom change in params
			if gauge.DepositAmount.Denom == receivedAmount.Denom {
				gauge.DepositAmount = gauge.DepositAmount.Add(receivedAmount)
			} else {
				gauge.DepositAmount = receivedAmount
			}
			gauge.TriggeredCount = ongoingEpochCount
			k.SetGauge(ctx, gauge)
		}
	}

	return nil
}
