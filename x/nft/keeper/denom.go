package keeper

import (
	"github.com/comdex-official/comdex/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HasDenomID(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenomID(id))
}

func (k Keeper) HasDenomSymbol( ctx sdk.Context, symbol string) bool {
	if len(symbol) == 0 {
		return false
	}
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenomSymbol(symbol))
}

func (k Keeper) SetDenom(ctx sdk.Context, denom types.Denom) error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&denom)
	store.Set(types.KeyDenomID(denom.Id), bz)
	if len(denom.Symbol) > 0 {
		store.Set(types.KeyDenomSymbol(denom.Symbol), []byte(denom.Id))
	}
	return nil
}

func (k Keeper) setDenomOwner(ctx sdk.Context, denomId string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	bz := types.MustMarshalDenomID(k.cdc, denomId)
	store.Set((types.KeyDenomCreator(owner, denomId)), bz)
}