package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/cosmos/gogoproto/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidation/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
)

func (k Keeper) LiquidateVaults(ctx sdk.Context) error {
	appIds := k.GetAppIdsForLiquidation(ctx)
	params := k.GetParams(ctx)

	for i := range appIds {
		esmStatus, found := k.esm.GetESMStatus(ctx, appIds[i])
		status := false
		if found {
			status = esmStatus.Status
		}
		klwsParams, _ := k.esm.GetKillSwitchData(ctx, appIds[i])
		if klwsParams.BreakerEnable || status {
			ctx.Logger().Error("Kill Switch Or ESM is enabled For Liquidation, liquidate_vaults.go for AppID %d", appIds[i])
			continue
		}

		liquidationOffsetHolder, found := k.GetLiquidationOffsetHolder(ctx, appIds[i], types.VaultLiquidationsOffsetPrefix)
		if !found {
			liquidationOffsetHolder = types.NewLiquidationOffsetHolder(appIds[i], 0)
		}
		totalVaults := k.vault.GetVaults(ctx)
		lengthOfVaults := int(k.vault.GetLengthOfVault(ctx))
		//// get all vaults
		/// range over those vaults
		//// for length of vaults use vault counter
		//// wen inside the vault slice check if the app_id matches with that of app_id[i]

		start, end := types.GetSliceStartEndForLiquidations(lengthOfVaults, int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
		if start == end {
			liquidationOffsetHolder.CurrentOffset = 0
			start, end = types.GetSliceStartEndForLiquidations(lengthOfVaults, int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
		}

		newVaults := totalVaults[start:end]
		for _, vault := range newVaults {
			_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
				if vault.AppId != appIds[i] {
					return fmt.Errorf("vault and app id mismatch in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
				}
				extPair, _ := k.asset.GetPairsVault(ctx, vault.ExtendedPairVaultID)
				pair, _ := k.asset.GetPair(ctx, extPair.PairId)
				assetIn, found := k.asset.GetAsset(ctx, pair.AssetIn)
				if !found {
					return fmt.Errorf("asset not found in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
				}
				totalRate, err := k.market.CalcAssetPrice(ctx, assetIn.Id, vault.AmountIn)
				if err != nil {
					return fmt.Errorf("error in CalcAssetPrice in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
				}
				totalIn := totalRate

				liqRatio := extPair.MinCr
				totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
				collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
				if err != nil {
					return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
				}
				if collateralizationRatio.LT(liqRatio) {
					// calculate interest and update vault
					totalDebt := vault.AmountOut.Add(vault.InterestAccumulated)
					err1 := k.rewards.CalculateVaultInterest(ctx, vault.AppId, vault.ExtendedPairVaultID, vault.Id, totalDebt, vault.BlockHeight, vault.BlockTime.Unix())
					if err1 != nil {
						return fmt.Errorf("error Calculating vault interest in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
					}
					vault, _ := k.vault.GetVault(ctx, vault.Id)
					totalFees := vault.InterestAccumulated.Add(vault.ClosingFeeAccumulated)
					totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
					collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
					if err != nil {
						return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
					}
					err = k.CreateLockedVault(ctx, vault, totalIn, collateralizationRatio, appIds[i], totalFees)
					if err != nil {
						return fmt.Errorf("error Creating Locked Vaults in Liquidation, liquidate_vaults.go for Vault %d", vault.Id)
					}
					k.vault.DeleteVault(ctx, vault.Id)
					var rewards rewardstypes.VaultInterestTracker
					rewards.AppMappingId = appIds[i]
					rewards.VaultId = vault.Id
					k.rewards.DeleteVaultInterestTracker(ctx, rewards)
					k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, vault.ExtendedPairVaultID, vault.Id, appIds[i])
				}
				return nil
			})
		}
		liquidationOffsetHolder.CurrentOffset = uint64(end)
		k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)
	}
	return nil
}

func (k Keeper) CreateLockedVault(ctx sdk.Context, vault vaulttypes.Vault, totalIn sdk.Dec, collateralizationRatio sdk.Dec, appID uint64, totalFees sdk.Int) error {
	lockedVaultID := k.GetLockedVaultID(ctx)

	value := types.LockedVault{
		LockedVaultId:           lockedVaultID + 1,
		AppId:                   appID,
		OriginalVaultId:         vault.Id,
		ExtendedPairId:          vault.ExtendedPairVaultID,
		Owner:                   vault.Owner,
		AmountIn:                vault.AmountIn,
		AmountOut:               vault.AmountOut,
		UpdatedAmountOut:        sdk.ZeroInt(),
		Initiator:               types.ModuleName,
		IsAuctionComplete:       false,
		IsAuctionInProgress:     false,
		CrAtLiquidation:         collateralizationRatio,
		CollateralToBeAuctioned: totalIn,
		LiquidationTimestamp:    ctx.BlockTime(),
		InterestAccumulated:     totalFees,
		Kind:                    nil,
	}

	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)
	length := k.vault.GetLengthOfVault(ctx)
	k.vault.SetLengthOfVault(ctx, length-1)
	err := k.auction.DutchActivator(ctx, value)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) GetModAccountBalances(ctx sdk.Context, accountName string, denom string) sdk.Int {
	macc := k.account.GetModuleAccount(ctx, accountName)
	return k.bank.GetBalance(ctx, macc.GetAddress(), denom).Amount
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
