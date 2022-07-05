package keeper

import (
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandelUpdateGenericParamsProposal(ctx sdk.Context, p *types.UpdateGenericParamsProposal) error {
	return k.UpdateGenericParams(ctx, p.AppId, p.Keys, p.Values)
}
