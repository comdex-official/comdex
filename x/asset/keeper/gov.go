package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) HandleAddPoolProposal(ctx sdk.Context, p *types.AddPoolProposal) error {
	var (
		count = k.GetCount(ctx)
		pool  = types.Pool{
			Id:               count + 1,
			DenomIn:          p.DenomIn,
			DenomOut:         p.DenomOut,
			LiquidationRatio: p.LiquidationRatio,
		}
	)

	k.SetPool(ctx, pool)
	k.SetCount(ctx, count+1)

	_ = ctx.EventManager().EmitTypedEvent(
		&types.EventAddPool{
			Id: pool.Id,
		},
	)

	return nil
}
