package keeper

import (
	"github.com/comdex-official/comdex/x/liquidation/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) LiquidateVaults(ctx sdk.Context) error {
	appIds := k.GetAppIdsForLiquidation(ctx)
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
				vault, found := k.GetVault(ctx, vaultIds[l])
				if !found {
					continue
				}

				extPair, _ := k.GetPairsVault(ctx, vault.ExtendedPairVaultID)
				pair, _ := k.GetPair(ctx, extPair.PairId)
				assetIn, found := k.GetAsset(ctx, pair.AssetIn)
				if !found {
					continue
				}
				assetInPrice, found := k.GetPriceForAsset(ctx, assetIn.Id)
				if !found {
					continue
				}
				totalIn := vault.AmountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).ToDec()

				liqRatio := extPair.MinCr
				totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
				collateralitzationRatio, err := k.CalculateCollaterlizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
				if err != nil {
					continue
				}

				if sdk.Dec.LT(collateralitzationRatio, liqRatio) {
					//calculate interest and update vault
					totalDebt := vault.AmountOut.Add(vault.InterestAccumulated)
					err1 := k.CalculateVaultInterest(ctx, vault.AppId, vault.ExtendedPairVaultID, vault.Id, totalDebt, vault.BlockHeight, vault.BlockTime.Unix())
					if err1 != nil {
						continue
					}
					vault, _ := k.GetVault(ctx, vaultIds[l])

					err := k.CreateLockedVault(ctx, vault, totalIn, collateralitzationRatio, appIds[i])
					if err != nil {
						continue
					}
					k.DeleteVault(ctx, vault.Id)
					var rewards rewardstypes.VaultInterestTracker
					rewards.AppMappingId = appIds[i]
					rewards.VaultId = vaultIds[l]
					k.DeleteVaultInterestTracker(ctx, rewards)
					k.DeleteAddressFromAppExtendedPairVaultMapping(ctx, vault.ExtendedPairVaultID, vault.Id, appIds[i])
				}
			}
			liquidationOffsetHolder.CurrentOffset = uint64(end)
			k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)
		}
	}
	return nil
}

func (k Keeper) CreateLockedVault(ctx sdk.Context, vault vaulttypes.Vault, totalIn sdk.Dec, collateralizationRatio sdk.Dec, appID uint64) error {
	lockedVaultID := k.GetLockedVaultID(ctx)

	var value = types.LockedVault{
		LockedVaultId:           lockedVaultID + 1,
		AppId:                   appID,
		OriginalVaultId:         vault.Id,
		ExtendedPairId:          vault.ExtendedPairVaultID,
		Owner:                   vault.Owner,
		AmountIn:                vault.AmountIn,
		AmountOut:               vault.AmountOut,
		UpdatedAmountOut:        vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated),
		Initiator:               types.ModuleName,
		IsAuctionComplete:       false,
		IsAuctionInProgress:     false,
		CrAtLiquidation:         collateralizationRatio,
		CollateralToBeAuctioned: totalIn,
		LiquidationTimestamp:    ctx.BlockTime(),
		InterestAccumulated:     vault.InterestAccumulated,
		Kind:                    nil,
	}

	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)
	err := k.DutchActivator(ctx, value)
	if err != nil {
		ctx.Logger().Error("error in dutch activator")
	}
	return nil
}

func (k Keeper) GetModAccountBalances(ctx sdk.Context, accountName string, denom string) sdk.Int {
	macc := k.GetModuleAccount(ctx, accountName)
	return k.GetBalance(ctx, macc.GetAddress(), denom).Amount
}

// Locked vault history kv

func (k Keeper) CreateLockedVaultHistory(ctx sdk.Context, lockedVault types.LockedVault) error {
	lockedVaultID := k.GetLockedVaultIDHistory(ctx)
	k.SetLockedVaultHistory(ctx, lockedVault, lockedVaultID)
	k.SetLockedVaultIDHistory(ctx, lockedVaultID+1)

	return nil
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

func (k Keeper) SetLockedVaultHistory(ctx sdk.Context, lockedVault types.LockedVault, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultHistoryKey(lockedVault.AppId, id)
		value = k.cdc.MustMarshal(&lockedVault)
	)
	store.Set(key, value)
}

func (k Keeper) GetLockedVaultHistory(ctx sdk.Context, appID, id uint64) (lockedVault types.LockedVault, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultHistoryKey(appID, id)
		value = store.Get(key)
	)

	if value == nil {
		return lockedVault, false
	}

	k.cdc.MustUnmarshal(value, &lockedVault)
	return lockedVault, true

}

// locked vaults kvs

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

func (k Keeper) SetLockedVault(ctx sdk.Context, lockedVault types.LockedVault) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(lockedVault.AppId, lockedVault.LockedVaultId)
		value = k.cdc.MustMarshal(&lockedVault)
	)
	store.Set(key, value)
}

func (k Keeper) DeleteLockedVault(ctx sdk.Context, appID, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(appID, id)
	)
	store.Delete(key)
}

func (k Keeper) GetLockedVault(ctx sdk.Context, appID, id uint64) (lockedVault types.LockedVault, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(appID, id)
		value = store.Get(key)
	)

	if value == nil {
		return lockedVault, false
	}

	k.cdc.MustUnmarshal(value, &lockedVault)
	return lockedVault, true
}

func (k Keeper) GetLockedVaultByApp(ctx sdk.Context, appID uint64) (lockedVault []types.LockedVault) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKeyByApp(appID)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.LockedVault
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		lockedVault = append(lockedVault, mapData)
	}
	return lockedVault
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

// auction flags kvs

func (k Keeper) SetFlagIsAuctionInProgress(ctx sdk.Context, appID, id uint64, flag bool) error {
	lockedVault, found := k.GetLockedVault(ctx, appID, id)
	if !found {
		return types.LockedVaultDoesNotExist
	}
	lockedVault.IsAuctionInProgress = flag
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

func (k Keeper) SetFlagIsAuctionComplete(ctx sdk.Context, appID, id uint64, flag bool) error {
	lockedVault, found := k.GetLockedVault(ctx, appID, id)
	if !found {
		return types.LockedVaultDoesNotExist
	}
	lockedVault.IsAuctionComplete = flag
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

// whitlisted appIds kvs

func (k Keeper) SetAppIDForLiquidation(ctx sdk.Context, appID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistAppKeyByApp(appID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: appID,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetAppIDByAppForLiquidation(ctx sdk.Context, appID uint64) (uint64, bool) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistAppKeyByApp(appID)
		value = store.Get(key)
	)

	if value == nil {
		return 0, false
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue(), true
}

func (k Keeper) GetAppIdsForLiquidation(ctx sdk.Context) (appIds []uint64) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AppIdsKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {

		var app protobuftypes.UInt64Value
		k.cdc.MustUnmarshal(iter.Value(), &app)
		appIds = append(appIds, app.Value)
	}
	return appIds
}

func (k Keeper) DeleteAppID(ctx sdk.Context, appID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistAppKeyByApp(appID)
	)

	store.Delete(key)
}
