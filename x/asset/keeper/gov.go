package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) HandleUpdateLiquidationRatioProposal(ctx sdk.Context, prop *types.UpdateLiquidationRatioProposal) error {
	params := k.GetParams(ctx)

	params.LiquidationRatio = prop.LiquidationRatio
	k.SetParams(ctx, params)

	return nil
}
