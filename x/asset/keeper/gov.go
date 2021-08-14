package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) HandleAddPairProposal(ctx sdk.Context, p *types.AddPairProposal) error {
	var (
		count = k.GetCount(ctx)
		pair  = types.Pair{
			Id:               count + 1,
			DenomIn:          p.DenomIn,
			DenomOut:         p.DenomOut,
			LiquidationRatio: p.LiquidationRatio,
		}
	)

	k.SetPair(ctx, pair)
	k.SetCount(ctx, count+1)

	_ = ctx.EventManager().EmitTypedEvent(
		&types.EventAddPair{
			Id: pair.Id,
		},
	)

	return nil
}
