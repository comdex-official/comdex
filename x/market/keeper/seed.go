package keeper

import (
	"fmt"
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/comdex-official/comdex/x/market/types"
)

func (k Keeper) GetRandomSeed(ctx sdk.Context, height int64) []byte {
	store := ctx.KVStore(k.key)

	random := store.Get(types.GetRandomKey(height))

	return random
}

func (k Keeper) SetRandomSeed(ctx sdk.Context, random []byte) {
	store := ctx.KVStore(k.key)

	ctx.Logger().Info(fmt.Sprintf("Setting random: %s", hex.EncodeToString(random)))

	store.Set(types.GetRandomKey(ctx.BlockHeight()), random)
}
