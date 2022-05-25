package keeper

import (
	"fmt"
	"time"

	"github.com/comdex-official/comdex/x/incentives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetRewardDistributionData returns the reward distribution data.
func (k Keeper) GetRewardDistributionData(
	ctx sdk.Context, gauge types.Gauge, coinToDistribute sdk.Coin, epochCount uint64, epochDuration time.Duration,
) ([]types.RewardDistributionDataCollector, error) {
	gaugeType := gauge.Kind

	var rewardDistributionData = []types.RewardDistributionDataCollector{}
	var err error = nil

	switch gaugeExtraData := gaugeType.(type) {
	case *types.Gauge_LiquidityMetaData:
		rewardDistributionData, err = k.liquidityKeeper.GetFarmingRewardsData(ctx, coinToDistribute, *gaugeExtraData.LiquidityMetaData)
	}

	return rewardDistributionData, err
}

// doDistributionSends initiates reward transfer from module account to wallet address.
func (k Keeper) doDistributionSends(ctx sdk.Context, gaugeTypeID uint64, distrs *types.DistributionInfo) {
	logger := k.Logger(ctx)
	numIDs := len(distrs.Addresses)
	logger.Info(fmt.Sprintf("Beginning reward distribution to %d users", numIDs))

	for i, address := range distrs.Addresses {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.ModuleName,
			address,
			distrs.Coins[i],
		)
		if err != nil {
			logger.Info(fmt.Sprintf("error occurred while reward distribution, err : %v", err))
		}
	}
	logger.Info("Finished sending, now creating reward distribution events")
	for id := 0; id < numIDs; id++ {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.TypeEvtDistribution,
				sdk.NewAttribute(types.AttributeReceiver, distrs.Addresses[id].String()),
				sdk.NewAttribute(types.AttributeAmount, distrs.Coins[id].String()),
				sdk.NewAttribute(types.AttributeGaugeTypeID, fmt.Sprintf("%v", gaugeTypeID)),
			),
		})
	}
	logger.Info(fmt.Sprintf("Finished distributing to %d users", numIDs))
}

// BeginRewardDistributions collects the reward data and initiates reward distribution.
func (k Keeper) BeginRewardDistributions(
	ctx sdk.Context, gauge types.Gauge, coinToDistribute sdk.Coin, epochCount uint64, epochDuration time.Duration,
) (sdk.Coin, error) {
	rewardDistributionData, err := k.GetRewardDistributionData(ctx, gauge, coinToDistribute, epochCount, epochDuration)
	if err != nil {
		return sdk.NewCoin(coinToDistribute.Denom, sdk.NewInt(0)), err
	}

	newDistributionInfo := types.DistributionInfo{
		Addresses: []sdk.AccAddress{},
		Coins:     []sdk.Coins{},
	}

	totalDistributionCoinsCalculated := sdk.NewCoin(coinToDistribute.Denom, sdk.NewInt(0))

	for _, distrData := range rewardDistributionData {
		newDistributionInfo.Addresses = append(newDistributionInfo.Addresses, distrData.RewardReceiver)
		newDistributionInfo.Coins = append(newDistributionInfo.Coins, sdk.NewCoins(distrData.RewardCoin))
		totalDistributionCoinsCalculated.Amount = totalDistributionCoinsCalculated.Amount.Add(distrData.RewardCoin.Amount)
	}

	if totalDistributionCoinsCalculated.Amount.GT(coinToDistribute.Amount) {
		return sdk.NewCoin(coinToDistribute.Denom, sdk.NewInt(0)), types.ErrInvalidCalculatedAMount
	}

	k.doDistributionSends(ctx, gauge.GaugeTypeId, &newDistributionInfo)

	return totalDistributionCoinsCalculated, nil
}
