package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) Admin(ctx sdk.Context) (s string) {
	k.params.Get(ctx, types.KeyAdmin, &s)
	return
}

func (k *Keeper) IBCPort(ctx sdk.Context) (s string) {
	k.params.Get(ctx, types.KeyIBCPort, &s)
	return
}

func (k *Keeper) IBCVersion(ctx sdk.Context) (s string) {
	k.params.Get(ctx, types.KeyIBCVersion, &s)
	return
}

func (k *Keeper) OracleAskCount(ctx sdk.Context) (i uint64) {
	k.params.Get(ctx, types.KeyOracleAskCount, &i)
	return
}

func (k *Keeper) OracleMinCount(ctx sdk.Context) (i uint64) {
	k.params.Get(ctx, types.KeyOracleMinCount, &i)
	return
}

func (k *Keeper) OracleMultiplier(ctx sdk.Context) (i uint64) {
	k.params.Get(ctx, types.KeyOracleMultiplier, &i)
	return
}
