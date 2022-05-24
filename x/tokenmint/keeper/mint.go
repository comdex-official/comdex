package keeper

import (

	// assettypes "github.com/comdex-official/comdex/x/asset/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/tokenmint/types"
)

func (k *Keeper) SetTokenMint(ctx sdk.Context, appTokenMintData types.TokenMint) {

	var (
		store = k.Store(ctx)
		key   = types.TokenMintKey(appTokenMintData.AppMappingId)
		value = k.cdc.MustMarshal(&appTokenMintData)
	)

	store.Set(key, value)

}

func (k *Keeper) GetTokenMint(ctx sdk.Context, appMappingId uint64) (appTokenMintData types.TokenMint, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.TokenMintKey(appMappingId)
		value = store.Get(key)
	)

	if value == nil {
		return appTokenMintData, false
	}

	k.cdc.MustUnmarshal(value, &appTokenMintData)
	return appTokenMintData, true
}
