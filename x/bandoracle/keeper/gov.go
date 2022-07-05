package keeper

import (
	"github.com/comdex-official/comdex/x/bandoracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleProposalFetchPrice(ctx sdk.Context, p *types.FetchPriceProposal) error {
	return k.AddFetchPriceRecords(ctx, p.FetchPrice)
}
