package keeper

import (
	"strconv"

	"github.com/comdex-official/comdex/x/liquidation/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) LiquidateVaults(ctx sdk.Context) error {
	appIds := k.GetAppIds(ctx).WhitelistedAppIds
	params := k.GetParams(ctx)

	for i := range appIds {
		esmStatus, found := k.GetESMStatus(ctx, appIds[i])
		status := false
		if found {
			status = esmStatus.Status
		}
		klwsParams, _ := k.GetKillSwitchData(ctx, appIds[i])
		if klwsParams.BreakerEnable || status {
			continue
		}

		liquidationOffsetHolder, found := k.GetLiquidationOffsetHolder(ctx, appIds[i], types.VaultLiquidationsOffsetPrefix)
		if !found {
			liquidationOffsetHolder = types.NewLiquidationOffsetHolder(appIds[i], 0)
			k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)
		}

		vaultsMap, _ := k.GetAppMappingData(ctx, appIds[i])

		vaults := vaultsMap
		for j := range vaults {
			vaultIds := vaults[j].VaultIds
			start, end := types.GetSliceStartEndForLiquidations(len(vaultIds), int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
			if start == end {
				liquidationOffsetHolder.CurrentOffset = 0
				start, end = types.GetSliceStartEndForLiquidations(len(vaultIds), int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
			}
			newVaultIDs := vaultIds[start:end]
			for l := range newVaultIDs {
				vault, found := k.GetVault(ctx, newVaultIDs[l])
				if !found {
					continue
				}

				extPair, _ := k.GetPairsVault(ctx, vault.ExtendedPairVaultID)

				liqRatio := extPair.MinCr
				totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
				collateralitzationRatio, err := k.CalculateCollaterlizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
				if err != nil {
					continue
				}

				if sdk.Dec.LT(collateralitzationRatio, liqRatio) {
					err := k.CreateLockedVault(ctx, vault, collateralitzationRatio, appIds[i])
					if err != nil {
						continue
					}
					k.DeleteVault(ctx, vault.Id)
					k.DeleteAddressFromAppExtendedPairVaultMapping(ctx, vault.ExtendedPairVaultID, vault.Id, appIds[i])
				}
			}
		}
	}
	return nil
}

func (k Keeper) CreateLockedVault(ctx sdk.Context, vault vaulttypes.Vault, collateralizationRatio sdk.Dec, appID uint64) error {
	lockedVaultID := k.GetLockedVaultID(ctx)

	var value = types.LockedVault{
		LockedVaultId:                lockedVaultID + 1,
		AppId:                        appID,
		AppVaultTypeId:               strconv.FormatUint(appID, 10),
		OriginalVaultId:              vault.Id,
		ExtendedPairId:               vault.ExtendedPairVaultID,
		Owner:                        vault.Owner,
		AmountIn:                     vault.AmountIn,
		AmountOut:                    vault.AmountOut,
		UpdatedAmountOut:             vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated),
		Initiator:                    types.ModuleName,
		IsAuctionComplete:            false,
		IsAuctionInProgress:          false,
		CrAtLiquidation:              collateralizationRatio,
		CurrentCollaterlisationRatio: collateralizationRatio,
		CollateralToBeAuctioned:      sdk.ZeroDec(),
		LiquidationTimestamp:         ctx.BlockTime(),
		SellOffHistory:               nil,
		InterestAccumulated:          vault.InterestAccumulated,
		Kind:                         nil,
	}

	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)
	return nil
}

func (k Keeper) SetLockedVaultByAppID(ctx sdk.Context, msg types.LockedVaultToAppMapping) {
	var (
		store = k.Store(ctx)
		key   = types.AppIDLockedVaultMappingKey(msg.AppId)
		value = k.cdc.MustMarshal(&msg)
	)

	store.Set(key, value)
}

func (k Keeper) GetLockedVaultByAppID(ctx sdk.Context, appMappingID uint64) (msg types.LockedVaultToAppMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AppIDLockedVaultMappingKey(appMappingID)
		value = store.Get(key)
	)

	if value == nil {
		return msg, false
	}

	k.cdc.MustUnmarshal(value, &msg)
	return msg, true
}

func (k Keeper) GetAllLockedVaultByAppID(ctx sdk.Context) (lockedVaultToAppMapping []types.LockedVaultToAppMapping) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AppIDLockedVaultMappingKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var vault types.LockedVaultToAppMapping
		k.cdc.MustUnmarshal(iter.Value(), &vault)
		lockedVaultToAppMapping = append(lockedVaultToAppMapping, vault)
	}
	return lockedVaultToAppMapping
}

func (k Keeper) CreateLockedVaultHistory(ctx sdk.Context, lockedVault types.LockedVault) error {
	lockedVaultID := k.GetLockedVaultIDHistory(ctx)
	k.SetLockedVaultHistory(ctx, lockedVault, lockedVaultID)
	k.SetLockedVaultIDHistory(ctx, lockedVaultID+1)

	return nil
}

func (k Keeper) UpdateLockedVaults(ctx sdk.Context) error {
	appIds := k.GetAppIds(ctx).WhitelistedAppIds
	for _, v := range appIds {
		newUpdatedVaults := k.GetLockedVaults(ctx)
		for _, lockedVault := range newUpdatedVaults {
			if lockedVault.AppId == v {
				ExtPair, _ := k.GetPairsVault(ctx, lockedVault.ExtendedPairId)
				if (!lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete) || (lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.LTE(ExtPair.MinCr)) {
					pair, _ := k.GetPair(ctx, ExtPair.PairId)
					assetIn, found := k.GetAsset(ctx, pair.AssetIn)
					if !found {
						continue
					}

					collateralizationRatio, err := k.CalculateCollaterlizationRatio(ctx, ExtPair.PairId, lockedVault.AmountIn, lockedVault.UpdatedAmountOut)
					if err != nil {
						continue
					}

					assetInPrice, found := k.GetPriceForAsset(ctx, assetIn.Id)
					if !found {
						continue
					}

					totalIn := lockedVault.AmountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).ToDec()
					updatedLockedVault := lockedVault
					updatedLockedVault.CurrentCollaterlisationRatio = collateralizationRatio
					updatedLockedVault.CollateralToBeAuctioned = totalIn
					k.SetLockedVault(ctx, updatedLockedVault)
					//k.UpdateLockedVaultsAppMapping(ctx, updatedLockedVault)

				}
			}
		}
	}
	return nil
}

func (k Keeper) GetModAccountBalances(ctx sdk.Context, accountName string, denom string) sdk.Int {
	macc := k.GetModuleAccount(ctx, accountName)
	return k.GetBalance(ctx, macc.GetAddress(), denom).Amount
}

func (k Keeper) GetLockedVaultIDbyApp(ctx sdk.Context, appID uint64) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.AppLockedVaultMappingKey(appID)
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) GetLockedVaultIDHistory(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKeyHistory
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetLockedVaultID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k Keeper) GetLockedVaultID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}
func (k Keeper) SetLockedVaultIDHistory(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKeyHistory
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k Keeper) SetLockedVault(ctx sdk.Context, lockedVault types.LockedVault) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(lockedVault.LockedVaultId)
		value = k.cdc.MustMarshal(&lockedVault)
	)
	store.Set(key, value)
}

func (k Keeper) SetLockedVaultHistory(ctx sdk.Context, lockedVault types.LockedVault, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultHistoryKey(id)
		value = k.cdc.MustMarshal(&lockedVault)
	)
	store.Set(key, value)
}

func (k Keeper) DeleteLockedVault(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(id)
	)
	store.Delete(key)
}

func (k Keeper) GetLockedVault(ctx sdk.Context, id uint64) (lockedVault types.LockedVault, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return lockedVault, false
	}

	k.cdc.MustUnmarshal(value, &lockedVault)
	return lockedVault, true
}

func (k Keeper) GetLockedVaults(ctx sdk.Context) (lockedVaults []types.LockedVault) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LockedVaultKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var lockedVault types.LockedVault
		k.cdc.MustUnmarshal(iter.Value(), &lockedVault)
		lockedVaults = append(lockedVaults, lockedVault)
	}

	return lockedVaults
}

func (k Keeper) SetFlagIsAuctionInProgress(ctx sdk.Context, id uint64, flag bool) error {
	lockedVault, found := k.GetLockedVault(ctx, id)
	if !found {
		return types.LockedVaultDoesNotExist
	}
	lockedVault.IsAuctionInProgress = flag
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

func (k Keeper) SetFlagIsAuctionComplete(ctx sdk.Context, id uint64, flag bool) error {
	lockedVault, found := k.GetLockedVault(ctx, id)
	if !found {
		return types.LockedVaultDoesNotExist
	}
	lockedVault.IsAuctionComplete = flag
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

func (k Keeper) SetAppID(ctx sdk.Context, AppIds types.WhitelistedAppIds) {
	var (
		store = k.Store(ctx)
		key   = types.AppIdsKeyPrefix
		value = k.cdc.MustMarshal(&AppIds)
	)

	store.Set(key, value)
}

func (k Keeper) GetAppIds(ctx sdk.Context) (appIds types.WhitelistedAppIds) {
	var (
		store = k.Store(ctx)
		key   = types.AppIdsKeyPrefix
		value = store.Get(key)
	)

	if value == nil {
		return appIds
	}

	k.cdc.MustUnmarshal(value, &appIds)
	return appIds
}
