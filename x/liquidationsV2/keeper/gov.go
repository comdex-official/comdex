package keeper

import (
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleWhitelistLiquidationProposal(ctx sdk.Context, p *types.WhitelistLiquidationProposal) error {
	return k.WhitelistLiquidation(ctx, p.Whitelisting)
}
