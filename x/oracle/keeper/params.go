package keeper

import (
	"github.com/comdex-official/comdex/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VotePeriod returns the number of blocks during which voting takes place.
func (k Keeper) VotePeriod(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyVotePeriod, &res)
	return
}

// SetVotePeriod updates the number of blocks during which voting takes place.
func (k Keeper) SetVotePeriod(ctx sdk.Context, votePeriod uint64) {
	k.paramSpace.Set(ctx, types.KeyVotePeriod, votePeriod)
}

// VoteThreshold returns the minimum percentage of votes that must be received
// for a ballot to pass.
func (k Keeper) VoteThreshold(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyVoteThreshold, &res)
	return
}

// SetVoteThreshold updates the minimum percentage of votes that must be received
// for a ballot to pass.
func (k Keeper) SetVoteThreshold(ctx sdk.Context, voteThreshold sdk.Dec) {
	k.paramSpace.Set(ctx, types.KeyVoteThreshold, voteThreshold)
}

// RewardBand returns the ratio of allowable exchange rate error that a validator
// can be rewarded.
func (k Keeper) RewardBands(ctx sdk.Context) (res types.RewardBandList) {
	k.paramSpace.Get(ctx, types.KeyRewardBands, &res)
	return
}

// VoteThreshold updates the ratio of allowable exchange rate error that a validator
// can be rewarded.
func (k Keeper) SetRewardBand(ctx sdk.Context, rewardBands types.RewardBandList) {
	k.paramSpace.Set(ctx, types.KeyRewardBands, rewardBands)
}

// RewardDistributionWindow returns the number of vote periods during which
// seigniorage reward comes in and then is distributed.
func (k Keeper) RewardDistributionWindow(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyRewardDistributionWindow, &res)
	return
}

// SetRewardDistributionWindow updates the number of vote periods during which
// seigniorage reward comes in and then is distributed.
func (k Keeper) SetRewardDistributionWindow(ctx sdk.Context, rewardDistributionWindow uint64) {
	k.paramSpace.Set(ctx, types.KeyRewardDistributionWindow, rewardDistributionWindow)
}

// AcceptList returns the denom list that can be activated
func (k Keeper) AcceptList(ctx sdk.Context) (res types.DenomList) {
	k.paramSpace.Get(ctx, types.KeyAcceptList, &res)
	return
}

// SetAcceptList updates the accepted list of assets supported by the x/oracle
// module.
func (k Keeper) SetAcceptList(ctx sdk.Context, acceptList types.DenomList) {
	k.paramSpace.Set(ctx, types.KeyAcceptList, acceptList)
}

// MandatoryList returns the denom list that are mandatory
func (k Keeper) MandatoryList(ctx sdk.Context) (res types.DenomList) {
	k.paramSpace.Get(ctx, types.KeyMandatoryList, &res)
	return
}

// SetMandatoryList updates the mandatory list of assets supported by the x/oracle
// module.
func (k Keeper) SetMandatoryList(ctx sdk.Context, mandatoryList types.DenomList) {
	k.paramSpace.Set(ctx, types.KeyMandatoryList, mandatoryList)
}

// SlashFraction returns the oracle voting penalty rate.
func (k Keeper) SlashFraction(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeySlashFraction, &res)
	return
}

// SetSlashFraction updates the oracle voting penalty rate.
func (k Keeper) SetSlashFraction(ctx sdk.Context, slashFraction sdk.Dec) {
	k.paramSpace.Set(ctx, types.KeySlashFraction, slashFraction)
}

// SlashWindow returns the number of total blocks in a slash window.
func (k Keeper) SlashWindow(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeySlashWindow, &res)
	return
}

// SetSlashWindow updates the number of total blocks in a slash window.
func (k Keeper) SetSlashWindow(ctx sdk.Context, slashWindow uint64) {
	k.paramSpace.Set(ctx, types.KeySlashWindow, slashWindow)
}

// MinValidPerWindow returns the oracle slashing threshold.
func (k Keeper) MinValidPerWindow(ctx sdk.Context) (res sdk.Dec) {
	k.paramSpace.Get(ctx, types.KeyMinValidPerWindow, &res)
	return
}

// MinValidPerWindow updates the oracle slashing threshold.
func (k Keeper) SetMinValidPerWindow(ctx sdk.Context, minValidPerWindow sdk.Dec) {
	k.paramSpace.Set(ctx, types.KeyMinValidPerWindow, minValidPerWindow)
}

// GetParams returns the total set of oracle parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of oracle parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// HistoricStampPeriod returns the amount of blocks the oracle module waits
// before recording a new historic price.
func (k Keeper) HistoricStampPeriod(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyHistoricStampPeriod, &res)
	return
}

// SetHistoricStampPeriod updates the amount of blocks the oracle module waits
// before recording a new historic price.
func (k Keeper) SetHistoricStampPeriod(ctx sdk.Context, historicPriceStampPeriod uint64) {
	k.paramSpace.Set(ctx, types.KeyHistoricStampPeriod, historicPriceStampPeriod)
}

// MedianStampPeriod returns the amount blocks the oracle module waits between
// calculating a new median and standard deviation of that median.
func (k Keeper) MedianStampPeriod(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyMedianStampPeriod, &res)
	return
}

// SetMedianStampPeriod updates the amount blocks the oracle module waits between
// calculating a new median and standard deviation of that median.
func (k Keeper) SetMedianStampPeriod(ctx sdk.Context, medianStampPeriod uint64) {
	k.paramSpace.Set(ctx, types.KeyMedianStampPeriod, medianStampPeriod)
}

// MaximumPriceStamps returns the maximum amount of historic prices the oracle
// module will hold.
func (k Keeper) MaximumPriceStamps(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyMaximumPriceStamps, &res)
	return
}

// SetMaximumPriceStamps updates the the maximum amount of historic prices the
// oracle module will hold.
func (k Keeper) SetMaximumPriceStamps(ctx sdk.Context, maximumPriceStamps uint64) {
	k.paramSpace.Set(ctx, types.KeyMaximumPriceStamps, maximumPriceStamps)
}

// MaximumMedianStamps returns the maximum amount of medians the oracle module will
// hold.
func (k Keeper) MaximumMedianStamps(ctx sdk.Context) (res uint64) {
	k.paramSpace.Get(ctx, types.KeyMaximumMedianStamps, &res)
	return
}

// SetMaximumMedianStamps updates the the maximum amount of medians the oracle module will
// hold.
func (k Keeper) SetMaximumMedianStamps(ctx sdk.Context, maximumMedianStamps uint64) {
	k.paramSpace.Set(ctx, types.KeyMaximumMedianStamps, maximumMedianStamps)
}
