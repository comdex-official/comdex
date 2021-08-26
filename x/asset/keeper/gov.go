package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) HandleAddPairProposal(ctx sdk.Context, prop *types.AddPairProposal) error {
	var (
		id   = k.GetPairID(ctx)
		pair = types.Pair{
			ID:               id + 1,
			AssetIn:          prop.AssetIn,
			AssetOut:         prop.AssetOut,
			LiquidationRatio: prop.LiquidationRatio,
		}
	)

	k.SetPair(ctx, pair)
	k.SetPairID(ctx, id+1)

	_ = ctx.EventManager().EmitTypedEvent(
		&types.EventAddPair{
			ID: pair.ID,
		},
	)

	return nil
}
