package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DistributionInfo defines the type for reward distribution.
type DistributionInfo struct {
	Addresses []sdk.AccAddress
	Coins     []sdk.Coins
}

// RewardDistributionDataCollector defines the type for the data collection of reward distribution.
type RewardDistributionDataCollector struct {
	RewardReceiver sdk.AccAddress
	RewardCoin     sdk.Coin
}
