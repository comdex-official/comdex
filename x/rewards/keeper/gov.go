package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/rewards/types"
)

func (k Keeper) HandleNewMintRewardsProposal(ctx sdk.Context, p *types.NewMintRewardsProposal) error {
	return k.AddNewMintingRewards(ctx, *p.MintRewards)
}

func (k Keeper) HandleDisableMintRewardsProposal(ctx sdk.Context, p *types.DisbaleMintRewardsProposal) error {
	return k.DisableMintingReward(ctx, p.MintRewardId)
}
