package keeper

import (
	"fmt"
	"time"

	"github.com/comdex-official/comdex/x/incentives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetRewardDistributionData(
	ctx sdk.Context, gauge types.Gauge, coinToDistribute sdk.Coin, epochCount uint64, epochDuration time.Duration,
) []types.RewardDistributionDataCollector {
	gaugeType := gauge.Kind

	var rewardDistributionData = []types.RewardDistributionDataCollector{}

	switch gaugeExtraData := gaugeType.(type) {
	case *types.Gauge_LiquidityMetaData:
		rewardDistributionData = k.liquidityKeeper.GetFarmingRewardsData(ctx, *gaugeExtraData.LiquidityMetaData)
	}

	return rewardDistributionData
}

func (k Keeper) doDistributionSends(ctx sdk.Context, distrs *types.DistributionInfo) error {
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
		logger.Info(fmt.Sprintf("error occured while reward distribution, err : %s", err))
	}
	logger.Info("Finished sending, now creating reward distribution events")
	for id := 0; id < numIDs; id++ {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.TypeEvtDistribution,
				sdk.NewAttribute(types.AttributeReceiver, distrs.Addresses[id].String()),
				sdk.NewAttribute(types.AttributeAmount, distrs.Coins[id].String()),
			),
		})
	}
	logger.Info(fmt.Sprintf("Finished Distributing to %d users", numIDs))
	return nil
}

func (k Keeper) BeginRewardDistributions(
	ctx sdk.Context, gauge types.Gauge, coinToDistribute sdk.Coin, epochCount uint64, epochDuration time.Duration,
) (sdk.Coin, error) {
	rewardDistributionData := k.GetRewardDistributionData(ctx, gauge, coinToDistribute, epochCount, epochDuration)

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

	err := k.doDistributionSends(ctx, &newDistributionInfo)
	if err != nil {
		return sdk.NewCoin(coinToDistribute.Denom, sdk.NewInt(0)), err
	}

	return totalDistributionCoinsCalculated, nil
}
