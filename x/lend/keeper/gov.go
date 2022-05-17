package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleAddWhitelistedPairsRecords(ctx sdk.Context, p *types.LendPairsProposal) error {
	return k.AddLendPairsRecords(ctx, p.Pairs...)
}

func (k Keeper) HandleUpdateWhitelistedPairRecords(ctx sdk.Context, p *types.UpdatePairProposal) error {
	return k.UpdateLendPairRecords(ctx, p.Pair)
}
