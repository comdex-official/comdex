package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/market/types"
)


func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.params.SetParamSet(ctx, &params)
}

func (k *Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams()
}
