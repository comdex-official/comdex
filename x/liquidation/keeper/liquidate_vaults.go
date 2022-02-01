package keeper

import (
	"time"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) LiquidateVaults(ctx sdk.Context) error {
	vaults := k.GetVaults(ctx)
	for _, vault := range vaults {
		pair, found := k.GetPair(ctx, vault.PairID)
		if !found {
			continue
		}
		liquidationRatio := pair.LiquidationRatio
		assetIn, found := k.GetAsset(ctx, pair.AssetIn)
		if !found {
			continue
		}

		assetOut, found := k.GetAsset(ctx, pair.AssetOut)
		if !found {
			continue
		}
		collateralizationRatio, err := k.CalculateCollaterlizationRatio(ctx, vault.AmountIn, assetIn, vault.AmountOut, assetOut)
		if err != nil {
			continue
		}
		if sdk.Dec.LT(collateralizationRatio, liquidationRatio) {
			err := k.TransferCollateralCreateLockedVaultAndDeleteVault(ctx, vault, assetIn, assetOut, collateralizationRatio)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) UnLiquidateVaults(ctx sdk.Context) error {
	locked_vaults := k.GetLockedVaults(ctx)
	for _, locked_vault := range locked_vaults {
		pair, found := k.GetPair(ctx, locked_vault.PairId)
		if !found {
			continue
		}
		liquidationRatio := pair.LiquidationRatio
		assetIn, found := k.GetAsset(ctx, pair.AssetIn)
		if !found {
			continue
		}

		assetOut, found := k.GetAsset(ctx, pair.AssetOut)
		if !found {
			continue
		}
		collateralizationRatio, err := k.CalculateCollaterlizationRatio(ctx, locked_vault.AmountIn, assetIn, locked_vault.AmountOut, assetOut)
		if err != nil {
			continue
		}
		if sdk.Dec.GT(collateralizationRatio, liquidationRatio) && !locked_vault.IsAuctioned {
			err := k.TransferCollateralRecreateVaultAndDeleteLockedVault(ctx, locked_vault, assetIn, assetOut, collateralizationRatio)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) TransferCollateralCreateLockedVaultAndDeleteVault(
	ctx sdk.Context,
	vault vaulttypes.Vault,
	assetIn assettypes.Asset,
	assetOut assettypes.Asset,
	collateralizationRatio sdk.Dec,
) error {
	collateralAvailableInVaultModule := k.GetModAccountBalances(ctx, vaulttypes.ModuleName, assetIn.Denom)
	collateralCoins := sdk.NewCoin(assetIn.Denom, sdk.MinInt(vault.AmountIn, collateralAvailableInVaultModule))

	err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, liquidationtypes.ModuleName, sdk.NewCoins(collateralCoins))
	if err != nil {
		return err
	}

	locked_vault := liquidationtypes.LockedVault{
		Id:                   k.GetLockedVaultID(ctx) + 1,
		VaultId:              vault.ID,
		PairId:               vault.PairID,
		Owner:                vault.Owner,
		AmountIn:             vault.AmountIn,
		AmountOut:            vault.AmountOut,
		Initiator:            liquidationtypes.ModuleName,
		IsAuctioned:          false,
		CrAtLiquidation:      collateralizationRatio,
		LiquidationTimestamp: time.Now().UTC(),
	}
	k.SetLockedVaultID(ctx, locked_vault.Id)
	k.SetLockedVault(ctx, locked_vault)

	vaultOwner, _ := sdk.AccAddressFromBech32(vault.Owner)
	k.DeleteVault(ctx, vault.ID)
	k.DeleteVaultForAddressByPair(ctx, vaultOwner, vault.PairID)
	return nil
}

func (k Keeper) TransferCollateralRecreateVaultAndDeleteLockedVault(
	ctx sdk.Context,
	locked_vault liquidationtypes.LockedVault,
	assetIn assettypes.Asset,
	assetOut assettypes.Asset,
	collateralizationRatio sdk.Dec,
) error {
	collateralAvailableInLiquidationModule := k.GetModAccountBalances(ctx, liquidationtypes.ModuleName, assetIn.Denom)
	collateralCoins := sdk.NewCoin(assetIn.Denom, sdk.MinInt(locked_vault.AmountIn, collateralAvailableInLiquidationModule))

	err := k.bank.SendCoinsFromModuleToModule(ctx, liquidationtypes.ModuleName, vaulttypes.ModuleName, sdk.NewCoins(collateralCoins))
	if err != nil {
		return err
	}
	vault := vaulttypes.Vault{
		ID:        k.GetVaultID(ctx) + 1,
		PairID:    locked_vault.PairId,
		Owner:     locked_vault.Owner,
		AmountIn:  locked_vault.AmountIn,
		AmountOut: locked_vault.AmountOut,
	}
	vaultOwner, _ := sdk.AccAddressFromBech32(vault.Owner)
	k.SetVaultID(ctx, vault.ID)
	k.SetVault(ctx, vault)
	k.SetVaultForAddressByPair(ctx, vaultOwner, vault.PairID, vault.ID)
	k.DeleteLockedVault(ctx, locked_vault.Id)
	return nil
}

func (k Keeper) GetModAccountBalances(ctx sdk.Context, accountName string, denom string) sdk.Int {
	macc := k.GetModuleAccount(ctx, accountName)
	return k.GetBalance(ctx, macc.GetAddress(), denom).Amount
}

func (k *Keeper) GetLockedVaultID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = liquidationtypes.LockedVaultIdKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetLockedVaultID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = liquidationtypes.LockedVaultIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) SetLockedVault(ctx sdk.Context, locked_vault liquidationtypes.LockedVault) {
	var (
		store = k.Store(ctx)
		key   = liquidationtypes.LockedVaultKey(locked_vault.Id)
		value = k.cdc.MustMarshal(&locked_vault)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteLockedVault(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = liquidationtypes.LockedVaultKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) GetLockedVault(ctx sdk.Context, id uint64) (locked_vault liquidationtypes.LockedVault, found bool) {
	var (
		store = k.Store(ctx)
		key   = liquidationtypes.LockedVaultKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return locked_vault, false
	}

	k.cdc.MustUnmarshal(value, &locked_vault)
	return locked_vault, true
}

func (k *Keeper) GetLockedVaults(ctx sdk.Context) (locked_vaults []liquidationtypes.LockedVault) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, liquidationtypes.LockedVaultKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var locked_vault liquidationtypes.LockedVault
		k.cdc.MustUnmarshal(iter.Value(), &locked_vault)
		locked_vaults = append(locked_vaults, locked_vault)
	}

	return locked_vaults
}

func (k *Keeper) FlagLockedVaultAsAuctioned(ctx sdk.Context, id uint64) error {

	locked_vault, found := k.GetLockedVault(ctx, id)
	if !found {
		return liquidationtypes.LockedVaultDoesNotExist
	}
	locked_vault.IsAuctioned = true
	k.SetLockedVault(ctx, locked_vault)
	return nil
}

func (k *Keeper) UnflagLockedVaultAsAuctioned(ctx sdk.Context, id uint64) error {

	locked_vault, found := k.GetLockedVault(ctx, id)
	if !found {
		return liquidationtypes.LockedVaultDoesNotExist
	}
	locked_vault.IsAuctioned = false
	k.SetLockedVault(ctx, locked_vault)
	return nil
}
