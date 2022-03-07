package keeper

import (
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleNewMintingRewardsProposal(ctx sdk.Context, proposalData *types.NewMintRewards) error {
	return k.AddNewMintingRecords(ctx, *proposalData)
}
