package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/types/decmath"
	"github.com/comdex-official/comdex/types/genmap"
	"github.com/comdex-official/comdex/types/reward"
	"github.com/comdex-official/comdex/x/oracle/types"
)

// prependCmdxIfUnique pushes `ucmdx` denom to the front of the list, if it is not yet included.
func prependCmdxIfUnique(voteTargets []string) []string {
	if genmap.Contains(types.CmdxDenom, voteTargets) {
		return voteTargets
	}
	rewardDenoms := make([]string, len(voteTargets)+1)
	rewardDenoms[0] = types.CmdxDenom
	copy(rewardDenoms[1:], voteTargets)
	return rewardDenoms
}

// smallestMissCountInBallot iterates through a given list of Claims and returns the smallest
// misscount in that list
func (k Keeper) smallestMissCountInBallot(ctx sdk.Context, ballotWinners []types.Claim) int64 {
	missCount := k.GetMissCounter(ctx, ballotWinners[0].Recipient)
	for _, winner := range ballotWinners[1:] {
		count := k.GetMissCounter(ctx, winner.Recipient)
		if count < missCount {
			missCount = count
		}
	}

	return int64(missCount)
}

// RewardBallotWinners is executed at the end of every voting period, where we
// give out a portion of seigniorage reward(reward-weight) to the oracle voters
// that voted correctly.
func (k Keeper) RewardBallotWinners(
	ctx sdk.Context,
	votePeriod int64,
	rewardDistributionWindow int64,
	voteTargets []string,
	ballotWinners []types.Claim,
) {
	if len(ballotWinners) == 0 {
		return
	}

	distributionRatio := sdk.NewDec(votePeriod).QuoInt64(rewardDistributionWindow)
	var periodRewards sdk.DecCoins
	rewardDenoms := prependCmdxIfUnique(voteTargets)
	for _, denom := range rewardDenoms {
		rewardPool := k.GetRewardPool(ctx, denom)

		// return if there's no rewards to give out
		if rewardPool.IsZero() {
			continue
		}

		periodRewards = periodRewards.Add(sdk.NewDecCoinFromDec(
			denom,
			sdk.NewDecFromInt(rewardPool.Amount).Mul(distributionRatio),
		))
	}

	// distribute rewards
	var distributedReward sdk.Coins

	smallestMissCount := k.smallestMissCountInBallot(ctx, ballotWinners)
	for _, winner := range ballotWinners {
		receiverVal := k.StakingKeeper.Validator(ctx, winner.Recipient)
		// in case absence of the validator, we just skip distribution
		if receiverVal == nil {
			continue
		}

		missCount := int64(k.GetMissCounter(ctx, winner.Recipient))
		maxMissCount := int64(len(voteTargets) * (int(k.SlashWindow(ctx) / k.VotePeriod(ctx))))
		rewardFactor := reward.CalculateRewardFactor(
			missCount,
			maxMissCount,
			smallestMissCount,
		)
		rewardDec, err := decmath.NewDecFromFloat(rewardFactor)
		if err != nil {
			k.Logger(ctx).With(err).Error("unable to calculate validator reward factor!")
			return
		}
		ballotLength := int64(len(ballotWinners))

		rewardCoins, _ := periodRewards.MulDec(rewardDec.QuoInt64(
			ballotLength)).TruncateDecimal()
		if rewardCoins.IsZero() {
			continue
		}

		k.distrKeeper.AllocateTokensToValidator(ctx, receiverVal, sdk.NewDecCoinsFromCoins(rewardCoins...))
		distributedReward = distributedReward.Add(rewardCoins...)
	}

	// move distributed reward to distribution module
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distrName, distributedReward)
	if err != nil {
		panic(fmt.Errorf("failed to send coins to distribution module %w", err))
	}
}
