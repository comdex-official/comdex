package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) HandleUpdateAdminProposal(ctx sdk.Context, p *types.UpdateAdminProposal) error {
	params := k.GetParams(ctx)

	params.Admin = p.Address
	k.SetParams(ctx, params)

	return nil
}
