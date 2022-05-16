package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) GetCollateralAmount(ctx sdk.Context, borrowerAddr sdk.AccAddress, denom string) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	collateral := sdk.NewCoin(denom, sdk.ZeroInt())
	key := types.CreateCollateralAmountKey(borrowerAddr, denom)

	if bz := store.Get(key); bz != nil {
		err := collateral.Amount.Unmarshal(bz)
		if err != nil {
			panic(err)
		}
	}

	return collateral
}

func (k Keeper) setCollateralAmount(ctx sdk.Context, borrowerAddr sdk.AccAddress, collateral sdk.Coin) error {
	if !collateral.IsValid() {
		return sdkerrors.Wrap(types.ErrInvalidAsset, collateral.String())
	}

	if borrowerAddr.Empty() {
		return types.ErrEmptyAddress
	}

	bz, err := collateral.Amount.Marshal()
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	key := types.CreateCollateralAmountKey(borrowerAddr, collateral.Denom)

	if collateral.Amount.IsZero() {
		store.Delete(key)
	} else {
		store.Set(key, bz)
	}
	return nil
}

func (k *Keeper) SetLendID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetLendID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LendIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetLend(ctx sdk.Context, lend types.Lend_Asset) {
	var (
		store = k.Store(ctx)
		key   = types.LendKey(lend.ID)
		value = k.cdc.MustMarshal(&lend)
	)

	store.Set(key, value)
}

func (k *Keeper) GetLend(ctx sdk.Context, id uint64) (asset types.Lend_Asset, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return asset, false
	}

	k.cdc.MustUnmarshal(value, &asset)
	return asset, true
}

func (k *Keeper) GetLends(ctx sdk.Context) (lends []types.Lend_Asset) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LendKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var lend types.Lend_Asset
		k.cdc.MustUnmarshal(iter.Value(), &lend)
		lends = append(lends, lend)
	}

	return lends
}

func (k *Keeper) SetLendForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendForAddressByPair(address, pairID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) HasLendForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.LendForAddressByPair(address, pairID)
	)

	return store.Has(key)
}
