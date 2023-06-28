package keeper

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HasDenomID(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenomID(id))
}

func (k Keeper) HasDenomSymbol(ctx sdk.Context, symbol string) bool {
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

func (k Keeper) GetDenom(ctx sdk.Context, id string) (denom types.Denom, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyDenomID(id))
	if bz == nil || len(bz) == 0 {
		return denom, sdkerrors.Wrapf(types.ErrInvalidDenom, "not found denomID: %s", id)
	}

	k.cdc.MustUnmarshal(bz, &denom)
	return denom, nil
}

func (k Keeper) GetDenoms(ctx sdk.Context) (denoms []types.Denom) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyDenomID(""))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var denom types.Denom
		k.cdc.MustUnmarshal(iterator.Value(), &denom)
		denoms = append(denoms, denom)
	}
	return denoms
}

func (k Keeper) GetDenomsByOwner(ctx sdk.Context, owner sdk.AccAddress) (denoms []types.Denom) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyDenomCreator(owner, ""))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		denomId := types.MustUnMarshalDenomID(k.cdc, iterator.Value())
		denom, _ := k.GetDenom(ctx, denomId)
		denoms = append(denoms, denom)
	}
	return denoms
}

func (k Keeper) AuthorizeDenomCreator(ctx sdk.Context, id string, creator sdk.AccAddress) (types.Denom, error) {
	denom, err := k.GetDenom(ctx, id)
	if err != nil {
		return types.Denom{}, err
	}

	if creator.String() != denom.Creator {
		return types.Denom{}, sdkerrors.Wrap(types.ErrUnauthorized, creator.String())
	}
	return denom, nil
}

func (k Keeper) HasPermissionToMint(ctx sdk.Context, denomID string, sender sdk.AccAddress) bool {
	denom, err := k.GetDenom(ctx, denomID)
	if err != nil {
		return false
	}

	if sender.String() == denom.Creator {
		return true
	}
	return false
}

func (k Keeper) deleteDenomOwner(ctx sdk.Context, denomID string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyDenomCreator(owner, denomID))
}

func (k Keeper) setDenomOwner(ctx sdk.Context, denomId string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	bz := types.MustMarshalDenomID(k.cdc, denomId)
	store.Set(types.KeyDenomCreator(owner, denomId), bz)
}

func (k Keeper) swapDenomOwner(ctx sdk.Context, denomID string, srcOwner, dstOwner sdk.AccAddress) {
	k.deleteDenomOwner(ctx, denomID, srcOwner)
	k.setDenomOwner(ctx, denomID, dstOwner)
}
