package keeper

import (
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleProposalLookupTableParams(ctx sdk.Context, p *types.LookupTableParams) error {
	return k.SetCollectorLookupTable(ctx, p.LookupTableData...)
}

func (k Keeper) HandleProposalLookupAppToAuction(ctx sdk.Context, p *types.AuctionControlByAppIdProposal) error {
	return k.SetAuctionMappingForApp(ctx, p.AppIdToAuctionLookup...)
}
