package reward

import (
	"math"
)

// CalculateRewardFactor returns the reward factor calculated using a logarmithic
// model based on miss counters. missCount is the current miss count, m is the
// maximum possible miss counts, and s is the smallest miss count in the period.
// If the logarimthic function returns NaN or Inf the Reward Factor returned will be 0.
// rewardFactor = 1 - logₘ₋ₛ₊₁(missCount - s + 1)
func CalculateRewardFactor(missCount, m, s int64) float64 {
	logBase := float64(m-s) + 1
	logKey := float64(missCount-s) + 1
	rewardFactor := 1 - (math.Log(logKey) / math.Log(logBase))
	if math.IsNaN(rewardFactor) || math.IsInf(rewardFactor, 0) {
		rewardFactor = 0
	}

	return rewardFactor
}
