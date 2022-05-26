package keeper

import (
	"github.com/comdex-official/comdex/x/esm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleProposalToggleEsm(ctx sdk.Context, p *types.ToggleEsmProposal) error {
	return k.SetTriggerEsm(ctx, p.EsmActive)
}