package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) HandleUpdateLiquidationRatio(ctx sdk.Context, prop *types.UpdateLiquidationRatioProposal) error {

	err := k.UpdateLiquidationRatio(ctx, prop)

	return err
}
