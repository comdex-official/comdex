package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/rewards/types"
)

func (k Keeper) HandleNewMintRewardsProposal(ctx sdk.Context, p *types.NewMintRewardsProposal) error {
	return k.AddNewMintingRewards(ctx, *p.MintRewards)
}
