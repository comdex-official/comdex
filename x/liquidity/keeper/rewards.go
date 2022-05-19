package keeper

import (
	incentivestypes "github.com/comdex-official/comdex/x/incentives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetFarmingRewardsData(ctx sdk.Context, liquidityGaugeData incentivestypes.LiquidtyGaugeMetaData) []incentivestypes.RewardDistributionDataCollector {
	return []incentivestypes.RewardDistributionDataCollector{}
}
