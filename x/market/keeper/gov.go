package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/market/types"
)

func (k *Keeper) HandleUpdateAdminProposal(ctx sdk.Context, prop *types.UpdateAdminProposal) error {
	params := k.GetParams(ctx)

	k.SetParams(ctx, params)

	return nil
}
