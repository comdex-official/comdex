package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetReserveFunds(ctx sdk.Context, denom string) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	key := types.ReserveFundsKey(denom)
	amount := sdk.ZeroInt()

	if bz := store.Get(key); bz != nil {
		err := amount.Unmarshal(bz)
		if err != nil {
			panic(err)
		}
	}

	if amount.IsNegative() {
		panic("negative reserve amount detected")
	}

	//return amount
	return sdk.NewInt(0)
}
