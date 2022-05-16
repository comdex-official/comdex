package keeper

import (
	"strconv"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/vault/types"
)

func (k *Keeper) SetCounterID(ctx sdk.Context, AppId uint64, counter uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LookupTableVaultKey(AppId)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: counter + 1,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetCounterID(ctx sdk.Context, AppId uint64) (uint64, bool) {
	var (
		store = k.Store(ctx)
		key   = types.LookupTableVaultKey(AppId)
		value = store.Get(key)
	)

	if value == nil {
		return 0 , false
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)


	return id.GetValue(), true
}

func (k *Keeper) SetVault(ctx sdk.Context, vault types.Vault, sName string) {
	var (
		Id = sName + strconv.Itoa(int(vault.Id))
		store = k.Store(ctx)
		key   = types.VaultKey(Id)
		value = k.cdc.MustMarshal(&vault)
	)

	store.Set(key, value)
}

func (k *Keeper) GetVault(ctx sdk.Context, id string) (vault types.Vault, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.VaultKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return vault, false
	}

	k.cdc.MustUnmarshal(value, &vault)
	return vault, true
}

func (k *Keeper) DeleteVault(ctx sdk.Context, id string) {
	var (
		store = k.Store(ctx)
		key   = types.VaultKey(id)
	)

	store.Delete(key)
}

func (k *Keeper) GetVaults(ctx sdk.Context) (vaults []types.Vault) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.VaultKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var vault types.Vault
		k.cdc.MustUnmarshal(iter.Value(), &vault)
		vaults = append(vaults, vault)
	}

	return vaults
}

func (k *Keeper) SetVaultForAddressByPair(ctx sdk.Context, address sdk.AccAddress, ExtendedPairVaultID uint64, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.VaultForAddressByPair(address, ExtendedPairVaultID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) HasVaultForAddressByPair(ctx sdk.Context, address sdk.AccAddress, ExtendedPairVaultID uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.VaultForAddressByPair(address, ExtendedPairVaultID)
	)

	return store.Has(key)
}

func (k *Keeper) DeleteVaultForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.VaultForAddressByPair(address, pairID)
	)

	store.Delete(key)
}

func (k *Keeper) VerifyCollaterlizationRatio(
	ctx sdk.Context,
	amountIn sdk.Int,
	assetIn assettypes.Asset,
	amountOut sdk.Int,
	assetOut assettypes.Asset,
	liquidationRatio sdk.Dec,
) error {
	collaterlizationRatio, err := k.CalculateCollaterlizationRatio(ctx, amountIn, assetIn, amountOut, assetOut)
	if err != nil {
		return err
	}

	if collaterlizationRatio.LT(liquidationRatio) {
		return types.ErrorInvalidCollateralizationRatio
	}

	return nil
}

func (k *Keeper) CalculateCollaterlizationRatio(
	ctx sdk.Context,
	amountIn sdk.Int,
	assetIn assettypes.Asset,
	amountOut sdk.Int,
	assetOut assettypes.Asset,
) (sdk.Dec, error) {

	assetInPrice, found := k.GetPriceForAsset(ctx, assetIn.Id)
	if !found {
		return sdk.ZeroDec(), types.ErrorPriceDoesNotExist
	}

	assetOutPrice, found := k.GetPriceForAsset(ctx, assetOut.Id)
	if !found {
		return sdk.ZeroDec(), types.ErrorPriceDoesNotExist
	}

	totalIn := amountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).ToDec()
	if totalIn.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), types.ErrorInvalidAmountIn
	}

	totalOut := amountOut.Mul(sdk.NewIntFromUint64(assetOutPrice)).ToDec()
	if totalOut.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), types.ErrorInvalidAmountOut
	}

	return totalIn.Quo(totalOut), nil
}

func (k *Keeper) SetLookupTableVault(ctx sdk.Context, vault types.LookupTableVault, app_id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.CounterVaultKey(app_id)
		value = k.cdc.MustMarshal(&vault)
	)

	store.Set(key, value)
}

func (k *Keeper) GetLookupTableVault(ctx sdk.Context, AppId uint64) (vaults []types.LookupTableVault) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.VaultKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var vault types.LookupTableVault
		k.cdc.MustUnmarshal(iter.Value(), &vault)
		vaults = append(vaults, vault)
	}

	return vaults
}

func (k *Keeper) SetUserVaultIdMapping(ctx sdk.Context, vault types.UserVaultIdMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserVaultIdMappingKey(vault.Owner)
		value = k.cdc.MustMarshal(&vault)
	)

	store.Set(key, value)
}

func (k *Keeper) SetPairVaultMapping(ctx sdk.Context, vault types.PairToVaultMapping) {
	var (
		store = k.Store(ctx)
		key   = types.PairIDKey(vault.PairId)
		value = k.cdc.MustMarshal(&vault)
	)

	store.Set(key, value)
}

func (k *Keeper) SetExtendedPairMapping(ctx sdk.Context, vault types.AppExtendedPairVaultMapping) {
	var (
		store = k.Store(ctx)
		key   = types.ExtendedIDKey(vault.AppId)
		value = k.cdc.MustMarshal(&vault)
	)

	store.Set(key, value)
}

func (k *Keeper) SetExtendedVaultPairMapping(ctx sdk.Context, vault types.ExtendedPairVaultMapping) {
	var (
		store = k.Store(ctx)
		key   = types.ExtendedIDKey(vault.ExtendedPairId)
		value = k.cdc.MustMarshal(&vault)
	)

	store.Set(key, value)
}

func (k *Keeper) GetExtendedVaultPairMapping(ctx sdk.Context, vault types.ExtendedPairVaultMapping) (uint64, bool) {
	var (
		store = k.Store(ctx)
		key   = types.ExtendedIDKey(vault.ExtendedPairId)
		value = store.Get(key)
	)

	if value == nil {
		return 0 , false
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)


	return id.GetValue(), true
}


func (k *Keeper) SetTokenMintedID(ctx sdk.Context, ExtendedPairVaultID uint64, TokenMinted sdk.Int) {
	var (
		store = k.Store(ctx)
		key   = types.TokenMintedIDKey(ExtendedPairVaultID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value:TokenMinted.Uint64(),
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetTokenMintedID(ctx sdk.Context, ExtendedPairVaultID uint64) (sdk.Int, bool) {
	var (
		store = k.Store(ctx)
		key   = types.TokenMintedIDKey(ExtendedPairVaultID)
		value = store.Get(key)
	)

	if value == nil {
		return sdk.ZeroInt() , false
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)


	return sdk.NewIntFromUint64(id.GetValue()), true
}