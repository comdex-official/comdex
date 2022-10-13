package keeper

import (
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleProposalAddLendRewards(ctx sdk.Context, p *types.AddLendExternalRewardsProposal) error {
	return k.AddLendExternalRewards(ctx, p.LendExternalRewards)
}
