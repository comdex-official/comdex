package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DistributionInfo struct {
	Addresses []sdk.AccAddress
	Coins     []sdk.Coins
}

type RewardDistributionDataCollector struct {
	RewardReceiver sdk.AccAddress
	RewardCoin     sdk.Coin
}
