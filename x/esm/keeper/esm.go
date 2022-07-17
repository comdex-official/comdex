package keeper

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/esm/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) AddESMTriggerParamsRecords(ctx sdk.Context, record types.ESMTriggerParams) error {
	_, found := k.GetESMTriggerParams(ctx, record.AppId)
	if found {
		return types.ErrorDuplicateESMTriggerParams
	}

	var (
		esmTriggerParams = types.ESMTriggerParams{
			AppId:         record.AppId,
			TargetValue:   record.TargetValue,
			CoolOffPeriod: record.CoolOffPeriod,
			AssetsRates:   record.AssetsRates,
		}
	)
	k.SetESMTriggerParams(ctx, esmTriggerParams)

	return nil
}

func (k *Keeper) SetESMTriggerParams(ctx sdk.Context, esmTriggerParams types.ESMTriggerParams) {
	var (
		store = k.Store(ctx)
		key   = types.ESMTriggerParamsKey(esmTriggerParams.AppId)
		value = k.cdc.MustMarshal(&esmTriggerParams)
	)

	store.Set(key, value)
}

func (k *Keeper) GetESMTriggerParams(ctx sdk.Context, id uint64) (esmTriggerParams types.ESMTriggerParams, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.ESMTriggerParamsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return esmTriggerParams, false
	}
	k.cdc.MustUnmarshal(value, &esmTriggerParams)
	return esmTriggerParams, true
}

func (k *Keeper) SetCurrentDepositStats(ctx sdk.Context, depositStats types.CurrentDepositStats) {
	var (
		store = k.Store(ctx)
		key   = types.CurrentDepositStatsKey(depositStats.AppId)
		value = k.cdc.MustMarshal(&depositStats)
	)

	store.Set(key, value)
}

func (k *Keeper) GetCurrentDepositStats(ctx sdk.Context, id uint64) (depositStats types.CurrentDepositStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.CurrentDepositStatsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return depositStats, false
	}

	k.cdc.MustUnmarshal(value, &depositStats)
	return depositStats, true
}

func (k *Keeper) SetESMStatus(ctx sdk.Context, esmStatus types.ESMStatus) {
	var (
		store = k.Store(ctx)
		key   = types.ESMStatusKey(esmStatus.AppId)
		value = k.cdc.MustMarshal(&esmStatus)
	)

	store.Set(key, value)
}

func (k *Keeper) GetESMStatus(ctx sdk.Context, id uint64) (esmStatus types.ESMStatus, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.ESMStatusKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return esmStatus, false
	}

	k.cdc.MustUnmarshal(value, &esmStatus)
	return esmStatus, true
}

func (k *Keeper) GetUserDepositByApp(ctx sdk.Context, address string, appID uint64) (userDeposits types.UsersDepositMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserDepositByAppKey(address, appID)
		value = store.Get(key)
	)
	if value == nil {
		return userDeposits, false
	}
	k.cdc.MustUnmarshal(value, &userDeposits)

	return userDeposits, true
}

func (k *Keeper) SetUserDepositByApp(ctx sdk.Context, userDeposits types.UsersDepositMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserDepositByAppKey(userDeposits.Depositor, userDeposits.AppId)
		value = k.cdc.MustMarshal(&userDeposits)
	)
	store.Set(key, value)
}

func (k *Keeper) SetESMMarketForAsset(ctx sdk.Context, esmMarket types.ESMMarketPrice) {
	var (
		store = k.Store(ctx)
		key   = types.ESMSPriceKey(esmMarket.AppId)
		value = k.cdc.MustMarshal(&esmMarket)
	)

	store.Set(key, value)
}

func (k *Keeper) GetESMMarketForAsset(ctx sdk.Context, id uint64) (esmMarket types.ESMMarketPrice, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.ESMSPriceKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return esmMarket, false
	}

	k.cdc.MustUnmarshal(value, &esmMarket)
	return esmMarket, true
}

func (k *Keeper) AddESMTriggerParamsForApp(ctx sdk.Context, addESMTriggerParams *bindings.MsgAddESMTriggerParams) error {

	var debtRates[] types.DebtAssetsRates
	for i := range addESMTriggerParams.AssetID{
		var debtRate types.DebtAssetsRates
		debtRate.AssetID = addESMTriggerParams.AssetID[i]
		debtRate.Rates = addESMTriggerParams.Rates[i]
		debtRates = append(debtRates, debtRate)
	}
	var (
		esmTriggerParams = types.ESMTriggerParams{
			AppId:         addESMTriggerParams.AppID,
			TargetValue:   addESMTriggerParams.TargetValue,
			CoolOffPeriod: addESMTriggerParams.CoolOffPeriod,
			AssetsRates:   debtRates,
		}
	)
	k.SetESMTriggerParams(ctx, esmTriggerParams)

	return nil
}

func (k Keeper) WasmAddESMTriggerParamsQuery(ctx sdk.Context, appID uint64) (bool, string) {
	_, found := k.GetApp(ctx, appID)
	if !found {
		return false, types.ErrAppIDDoesNotExists.Error()
	}
	return true, ""
}

func (k *Keeper) SetDataAfterCoolOff(ctx sdk.Context, esmDataAfterCoolOff types.DataAfterCoolOff) {
	var (
		store = k.Store(ctx)
		key   = types.ESMDataAfterCoolOff(esmDataAfterCoolOff.AppId)
		value = k.cdc.MustMarshal(&esmDataAfterCoolOff)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDataAfterCoolOff(ctx sdk.Context, id uint64) (esmDataAfterCoolOff types.DataAfterCoolOff, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.ESMDataAfterCoolOff(id)
		value = store.Get(key)
	)

	if value == nil {
		return esmDataAfterCoolOff, false
	}
	k.cdc.MustUnmarshal(value, &esmDataAfterCoolOff)
	return esmDataAfterCoolOff, true
}

func (k *Keeper) SetUpCollateralRedemptionForVault(ctx sdk.Context, appId uint64) error {
	totalVaults := k.GetVaults(ctx)
	for _, data := range totalVaults {
		if data.AppId == appId {
			extendedPairVault, found := k.GetPairsVault(ctx, data.ExtendedPairVaultID)
			if !found {
				return vaulttypes.ErrorExtendedPairVaultDoesNotExists
			}
			pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
			if !found {
				return assettypes.ErrorPairDoesNotExist
			}
			assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
			if !found {
				return assettypes.ErrorAssetDoesNotExist
			}
			assetOutData, found := k.GetAsset(ctx, pairData.AssetOut)
			if !found {
				return assettypes.ErrorAssetDoesNotExist
			}
			coolOffData, found := k.GetDataAfterCoolOff(ctx, appId)
			if !found {
				coolOffData.AppId = appId
				var item types.AssetToAmount

				item.AssetID = assetInData.Id
				item.Amount = data.AmountIn
				err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, data.AmountIn)))
				if err != nil {
					return err
				}
				coolOffData.CollateralAsset = append(coolOffData.CollateralAsset, item)

				item.AssetID = assetOutData.Id
				item.Amount = data.AmountOut

				coolOffData.DebtAsset = append(coolOffData.DebtAsset, item)

				k.SetDataAfterCoolOff(ctx, coolOffData)
			} else {
				var count = 0
				for i, indata := range coolOffData.CollateralAsset {
					if indata.AssetID == assetInData.Id {
						count++
						indata.Amount = indata.Amount.Add(data.AmountIn)
						err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, data.AmountIn)))
						if err != nil {
							return err
						}
						coolOffData.CollateralAsset = append(coolOffData.CollateralAsset[:i],coolOffData.CollateralAsset[i+1:]...)
						coolOffData.CollateralAsset = append(coolOffData.CollateralAsset, indata)
						break
					}
				}
				if count == 0 {
					var item types.AssetToAmount

					item.AssetID = assetInData.Id
					item.Amount = data.AmountIn

					err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, data.AmountIn)))
					if err != nil {
						return err
					}
					coolOffData.CollateralAsset = append(coolOffData.CollateralAsset, item)
					count = 0
				}

				for i, indata := range coolOffData.DebtAsset {
					if indata.AssetID == assetOutData.Id {
						count++
						indata.Amount = indata.Amount.Add(data.AmountOut)
						coolOffData.DebtAsset = append(coolOffData.DebtAsset[:i],coolOffData.DebtAsset[i+1:]...)
						coolOffData.DebtAsset = append(coolOffData.DebtAsset, indata)
						break
					}
				}
				if count == 0 {
					var item types.AssetToAmount

					item.AssetID = assetOutData.Id
					item.Amount = data.AmountOut
					coolOffData.DebtAsset = append(coolOffData.DebtAsset, item)
					count = 0
				}
				k.SetDataAfterCoolOff(ctx, coolOffData)
			}

			k.DeleteVault(ctx, data.Id)
			k.DeleteAddressFromAppExtendedPairVaultMapping(ctx, data.ExtendedPairVaultID, data.Id, data.AppId)
			esmStatus, found := k.GetESMStatus(ctx, data.AppId)
			esmStatus.VaultRedemptionStatus = true
			k.SetESMStatus(ctx, esmStatus)
		}
	}
	return nil
}

func (k *Keeper) SetUpCollateralRedemptionForStableVault(ctx sdk.Context, appId uint64) error {
	totalStableVaults := k.GetStableMintVaults(ctx)
	for _, data := range totalStableVaults {
		if data.AppId == appId {
			extendedPairVault, found := k.GetPairsVault(ctx, data.ExtendedPairVaultID)
			if !found {
				return vaulttypes.ErrorExtendedPairVaultDoesNotExists
			}
			pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
			if !found {
				return assettypes.ErrorPairDoesNotExist
			}
			assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
			if !found {
				return assettypes.ErrorAssetDoesNotExist
			}
			assetOutData, found := k.GetAsset(ctx, pairData.AssetOut)
			if !found {
				return assettypes.ErrorAssetDoesNotExist
			}
			coolOffData, found := k.GetDataAfterCoolOff(ctx, appId)
			if !found {
				coolOffData.AppId = appId
				var item types.AssetToAmount

				item.AssetID = assetInData.Id
				item.Amount = data.AmountIn
				err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, data.AmountIn)))
				if err != nil {
					return err
				}
				coolOffData.CollateralAsset = append(coolOffData.CollateralAsset, item)

				item.AssetID = assetOutData.Id
				item.Amount = data.AmountOut

				coolOffData.DebtAsset = append(coolOffData.DebtAsset, item)

				k.SetDataAfterCoolOff(ctx, coolOffData)
			} else {
				var count = 0
				for i, indata := range coolOffData.CollateralAsset {
					if indata.AssetID == assetInData.Id {
						count++
						indata.Amount = indata.Amount.Add(data.AmountIn)
						err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, data.AmountIn)))
						if err != nil {
							return err
						}
						coolOffData.CollateralAsset = append(coolOffData.CollateralAsset[:i],coolOffData.CollateralAsset[i+1:]...)
						coolOffData.CollateralAsset = append(coolOffData.CollateralAsset, indata)
						break
					}
				}
				if count == 0 {
					var item types.AssetToAmount

					item.AssetID = assetInData.Id
					item.Amount = data.AmountIn

					err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, data.AmountIn)))
					if err != nil {
						return err
					}
					coolOffData.CollateralAsset = append(coolOffData.CollateralAsset, item)
					count = 0
				}

				for i, indata := range coolOffData.DebtAsset {
					if indata.AssetID == assetOutData.Id {
						count++
						indata.Amount = indata.Amount.Add(data.AmountOut)
						coolOffData.DebtAsset = append(coolOffData.DebtAsset[:i],coolOffData.DebtAsset[i+1:]...)
						coolOffData.DebtAsset = append(coolOffData.DebtAsset, indata)
						break
					}
				}
				if count == 0 {
					var item types.AssetToAmount

					item.AssetID = assetOutData.Id
					item.Amount = data.AmountOut
					coolOffData.DebtAsset = append(coolOffData.DebtAsset, item)
					count = 0
				}
				k.SetDataAfterCoolOff(ctx, coolOffData)
			}

			k.DeleteVault(ctx, data.Id)
			k.DeleteAddressFromAppExtendedPairVaultMapping(ctx, data.ExtendedPairVaultID, data.Id, data.AppId)
			esmStatus, found := k.GetESMStatus(ctx, data.AppId)
			esmStatus.StableVaultRedemptionStatus = true
			k.SetESMStatus(ctx, esmStatus)
		}
	}
	netFee, found := k.GetNetFeeCollectedData(ctx, appId)
	coolOffData, found := k.GetDataAfterCoolOff(ctx, appId)
	if !found{
		return nil
	}else {
		for _, data := range netFee.AssetIdToFeeCollected{
			for _, indata := range coolOffData.DebtAsset{
				if data.AssetId == indata.AssetID{
					indata.Amount = indata.Amount.Sub(data.NetFeesCollected)
					coolOffData.DebtAsset = append(coolOffData.DebtAsset, indata)
				}
			}
		}
		k.SetDataAfterCoolOff(ctx, coolOffData)
	}
	return nil
}

func (k *Keeper) SetAssetToAmountValue(ctx sdk.Context, assetToAmountValue types.AssetToAmountValue) {
	var (
		store = k.Store(ctx)
		key   = types.AssetToAmountValueKey(assetToAmountValue.AppId, assetToAmountValue.AssetID)
		value = k.cdc.MustMarshal(&assetToAmountValue)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAssetToAmountValue(ctx sdk.Context, appID, assetID uint64) (assetToAmountValue types.AssetToAmountValue, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetToAmountValueKey(appID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return assetToAmountValue, false
	}

	k.cdc.MustUnmarshal(value, &assetToAmountValue)
	return assetToAmountValue, true
}

func (k *Keeper) SetAppToAmtValue(ctx sdk.Context, appToAmt types.AppToAmountValue) {
	var (
		store = k.Store(ctx)
		key   = types.AppToAmountValueKey(appToAmt.AppId)
		value = k.cdc.MustMarshal(&appToAmt)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAppToAmtValue(ctx sdk.Context, id uint64) (appToAmt types.AppToAmountValue, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AppToAmountValueKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return appToAmt, false
	}

	k.cdc.MustUnmarshal(value, &appToAmt)
	return appToAmt, true
}
