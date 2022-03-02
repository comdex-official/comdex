package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

/*func (k *Keeper) Admin(ctx sdk.Context) (s string) {
	k.params.Get(ctx, types.KeyAdmin, &s)
	return
}*/

func (k *Keeper) LiqRatio(ctx sdk.Context) (s string) {
	k.params.Get(ctx, types.KeyLiqRatio, &s)
	return
}

func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.params.SetParamSet(ctx, &params)
}

func (k *Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams( k.LiqRatio(ctx) )
}
