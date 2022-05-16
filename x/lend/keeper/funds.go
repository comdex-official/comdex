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
	return sdk.NewInt(100)
}

// setReserveFunds sets the amount reserved of a specified token.
func (k Keeper) setReserveFunds(ctx sdk.Context, coin sdk.Coin) error {
	if err := coin.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	reserveKey := types.ReserveFundsKey(coin.Denom)

	// save the new reserve funds
	bz, err := coin.Amount.Marshal()
	if err != nil {
		return err
	}

	store.Set(reserveKey, bz)
	return nil
}
