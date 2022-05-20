package keeper

import (
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleProposalLookupTableParams(ctx sdk.Context, p *types.LookupTableParams) error {
	return k.SetCollectorLookupTable(ctx, p.LookupTableData...)
}
