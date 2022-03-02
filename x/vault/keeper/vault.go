package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/vault/types"
)

func (k *Keeper) SetID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.IDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.IDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetVault(ctx sdk.Context, vault types.Vault) {
	var (
		store = k.Store(ctx)
		key   = types.VaultKey(vault.ID)
		value = k.cdc.MustMarshal(&vault)
	)

	store.Set(key, value)
}

func (k *Keeper) GetVault(ctx sdk.Context, id uint64) (vault types.Vault, found bool) {
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

func (k *Keeper) DeleteVault(ctx sdk.Context, id uint64) {
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

func (k *Keeper) SetVaultForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.VaultForAddressByPair(address, pairID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetUserVaults(ctx sdk.Context, address string) (userVaults types.UserVaultIdMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserVaultsForAddressKey(address)
		value = store.Get(key)
	)
	if value == nil {
		return userVaults, false
	}
	k.cdc.MustUnmarshal(value, &userVaults)

	return userVaults, true
}

func (k *Keeper) SetUserVaults(ctx sdk.Context, userVaults types.UserVaultIdMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserVaultsForAddressKey(userVaults.Owner)
		value = k.cdc.MustMarshal(&userVaults)
	)
	store.Set(key, value)
}

func (k *Keeper) HasVaultForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.VaultForAddressByPair(address, pairID)
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
		return sdk.ZeroDec(), types.ErrorPriceInDoesNotExist
	}

	assetOutPrice, found := k.GetPriceForAsset(ctx, assetOut.Id)
	if !found {
		return sdk.ZeroDec(), types.ErrorPriceOutDoesNotExist
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

func (k *Keeper) GetAllCAssetMintRecords(ctx sdk.Context) (mintRecords []*types.CAssetsMintStatistics) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.CAssetMintStatisticsKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var mintRecord types.CAssetsMintStatistics
		k.cdc.MustUnmarshal(iter.Value(), &mintRecord)
		mintRecords = append(mintRecords, &mintRecord)
	}

	return mintRecords
}

func (k *Keeper) GetCAssetMintRecords(ctx sdk.Context, collateralDenom string) (mintRecords types.CAssetsMintStatistics, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.CAssetMintRecordsKey(collateralDenom)
		value = store.Get(key)
	)
	if value == nil {
		return mintRecords, false
	}
	k.cdc.MustUnmarshal(value, &mintRecords)

	return mintRecords, true
}

func (k *Keeper) SetCAssetMintRecords(ctx sdk.Context, mintRecords types.CAssetsMintStatistics) {
	var (
		store = k.Store(ctx)
		key   = types.CAssetMintRecordsKey(mintRecords.CollateralDenom)
		value = k.cdc.MustMarshal(&mintRecords)
	)
	store.Set(key, value)
}

func (k *Keeper) MintCAssets(
	ctx sdk.Context,
	moduleName string,
	collateralDenom string,
	denom string,
	amount sdk.Int,
) error {
	if err := k.MintCoin(ctx, moduleName, sdk.NewCoin(denom, amount)); err != nil {
		return err
	}
	mintRecords, found := k.GetCAssetMintRecords(ctx, collateralDenom)
	if !found {
		mintRecords = types.CAssetsMintStatistics{
			CollateralDenom: collateralDenom,
			MintedAssets:    map[string]uint64{},
		}
	}
	mintRecords.MintedAssets[denom] += amount.Uint64()
	k.SetCAssetMintRecords(ctx, mintRecords)
	return nil
}

func (k *Keeper) BurnCAssets(
	ctx sdk.Context,
	moduleName string,
	collateralDenom string,
	denom string,
	amount sdk.Int,
) error {
	mintRecords, found := k.GetCAssetMintRecords(ctx, collateralDenom)
	if !found {
		return types.ErrorCAssetRecordDoesNotExist
	}
	if mintRecords.MintedAssets[denom] < amount.Uint64() {
		return types.ErrorEnoughCAssetsNotMinted
	}
	mintRecords.MintedAssets[denom] -= amount.Uint64()
	k.SetCAssetMintRecords(ctx, mintRecords)

	if err := k.BurnCoin(ctx, moduleName, sdk.NewCoin(denom, amount)); err != nil {
		return err
	}
	return nil
}

func (k *Keeper) UpdateUserVaultIdMapping(
	ctx sdk.Context,
	vaultOwner string,
	vaultId uint64,
	isInsert bool,
) error {

	userVaults, found := k.GetUserVaults(ctx, vaultOwner)

	if !found && isInsert {
		userVaults = types.UserVaultIdMapping{
			Owner:    vaultOwner,
			VaultIds: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorVaultOwnerNotFound
	}

	if isInsert {
		userVaults.VaultIds = append(userVaults.VaultIds, vaultId)
	} else {
		for index, id := range userVaults.VaultIds {
			if id == vaultId {
				userVaults.VaultIds = append(userVaults.VaultIds[:index], userVaults.VaultIds[index+1:]...)
				break
			}
		}
	}

	k.SetUserVaults(ctx, userVaults)
	return nil
}
