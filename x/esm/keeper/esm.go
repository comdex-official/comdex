package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/app/wasm/bindings"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/esm/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
)

func (k Keeper) SetESMTriggerParams(ctx sdk.Context, esmTriggerParams types.ESMTriggerParams) {
	var (
		store = k.Store(ctx)
		key   = types.ESMTriggerParamsKey(esmTriggerParams.AppId)
		value = k.cdc.MustMarshal(&esmTriggerParams)
	)

	store.Set(key, value)
}

func (k Keeper) GetESMTriggerParams(ctx sdk.Context, id uint64) (esmTriggerParams types.ESMTriggerParams, found bool) {
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

func (k Keeper) GetAllESMTriggerParams(ctx sdk.Context) (eSMTriggerParams []types.ESMTriggerParams) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.ESMTriggerParamsKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var esm types.ESMTriggerParams
		k.cdc.MustUnmarshal(iter.Value(), &esm)
		eSMTriggerParams = append(eSMTriggerParams, esm)
	}
	return eSMTriggerParams
}

func (k Keeper) SetCurrentDepositStats(ctx sdk.Context, depositStats types.CurrentDepositStats) {
	var (
		store = k.Store(ctx)
		key   = types.CurrentDepositStatsKey(depositStats.AppId)
		value = k.cdc.MustMarshal(&depositStats)
	)

	store.Set(key, value)
}

func (k Keeper) GetCurrentDepositStats(ctx sdk.Context, id uint64) (depositStats types.CurrentDepositStats, found bool) {
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

func (k Keeper) GetAllCurrentDepositStats(ctx sdk.Context) (currentDepositStats []types.CurrentDepositStats) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.CurrentDepositStatsPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var esm types.CurrentDepositStats
		k.cdc.MustUnmarshal(iter.Value(), &esm)
		currentDepositStats = append(currentDepositStats, esm)
	}
	return currentDepositStats
}

func (k Keeper) SetESMStatus(ctx sdk.Context, esmStatus types.ESMStatus) {
	var (
		store = k.Store(ctx)
		key   = types.ESMStatusKey(esmStatus.AppId)
		value = k.cdc.MustMarshal(&esmStatus)
	)

	store.Set(key, value)
}

func (k Keeper) GetESMStatus(ctx sdk.Context, id uint64) (esmStatus types.ESMStatus, found bool) {
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

func (k Keeper) GetAllESMStatus(ctx sdk.Context) (eSMStatus []types.ESMStatus) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.ESMStatusPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var esm types.ESMStatus
		k.cdc.MustUnmarshal(iter.Value(), &esm)
		eSMStatus = append(eSMStatus, esm)
	}
	return eSMStatus
}

func (k Keeper) GetUserDepositByApp(ctx sdk.Context, address string, appID uint64) (userDeposits types.UsersDepositMapping, found bool) {
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

func (k Keeper) SetUserDepositByApp(ctx sdk.Context, userDeposits types.UsersDepositMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserDepositByAppKey(userDeposits.Depositor, userDeposits.AppId)
		value = k.cdc.MustMarshal(&userDeposits)
	)
	store.Set(key, value)
}

func (k Keeper) GetAllUserDepositByApp(ctx sdk.Context) (usersDepositMapping []types.UsersDepositMapping) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.UserDepositByAppPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var esm types.UsersDepositMapping
		k.cdc.MustUnmarshal(iter.Value(), &esm)
		usersDepositMapping = append(usersDepositMapping, esm)
	}
	return usersDepositMapping
}

func (k Keeper) AddESMTriggerParamsForApp(ctx sdk.Context, addESMTriggerParams *bindings.MsgAddESMTriggerParams) error {
	var debtRates []types.DebtAssetsRates
	for i := range addESMTriggerParams.AssetID {
		var debtRate types.DebtAssetsRates
		debtRate.AssetID = addESMTriggerParams.AssetID[i]
		debtRate.Rates = addESMTriggerParams.Rates[i]
		debtRates = append(debtRates, debtRate)
	}
	esmTriggerParams := types.ESMTriggerParams{
		AppId:         addESMTriggerParams.AppID,
		TargetValue:   addESMTriggerParams.TargetValue,
		CoolOffPeriod: addESMTriggerParams.CoolOffPeriod,
		AssetsRates:   debtRates,
	}
	k.SetESMTriggerParams(ctx, esmTriggerParams)

	return nil
}

func (k Keeper) WasmAddESMTriggerParamsQuery(ctx sdk.Context, appID uint64) (bool, string) {
	_, found := k.asset.GetApp(ctx, appID)
	if !found {
		return false, types.ErrAppIDDoesNotExists.Error()
	}
	return true, ""
}

func (k Keeper) SetDataAfterCoolOff(ctx sdk.Context, esmDataAfterCoolOff types.DataAfterCoolOff) {
	var (
		store = k.Store(ctx)
		key   = types.ESMDataAfterCoolOff(esmDataAfterCoolOff.AppId)
		value = k.cdc.MustMarshal(&esmDataAfterCoolOff)
	)

	store.Set(key, value)
}

func (k Keeper) GetDataAfterCoolOff(ctx sdk.Context, id uint64) (esmDataAfterCoolOff types.DataAfterCoolOff, found bool) {
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

func (k Keeper) GetAllDataAfterCoolOff(ctx sdk.Context) (dataAfterCoolOff []types.DataAfterCoolOff) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.ESMDataAfterCoolOffPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var esm types.DataAfterCoolOff
		k.cdc.MustUnmarshal(iter.Value(), &esm)
		dataAfterCoolOff = append(dataAfterCoolOff, esm)
	}
	return dataAfterCoolOff
}

func (k Keeper) SetUpCollateralRedemptionForVault(ctx sdk.Context, appID uint64) error {
	totalVaults := k.vault.GetVaults(ctx)
	esmStatus, found := k.GetESMStatus(ctx, appID)
	if !found {
		return types.ErrESMParamsNotFound
	}
	esmData, _ := k.GetESMTriggerParams(ctx, appID)
	for _, data := range totalVaults {
		if data.AppId == appID {
			extendedPairVault, found := k.asset.GetPairsVault(ctx, data.ExtendedPairVaultID)
			if !found {
				return vaulttypes.ErrorExtendedPairVaultDoesNotExists
			}
			pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
			if !found {
				return assettypes.ErrorPairDoesNotExist
			}
			assetInData, found := k.asset.GetAsset(ctx, pairData.AssetIn)
			if !found {
				return assettypes.ErrorAssetDoesNotExist
			}
			rateIn, found := k.GetSnapshotOfPrices(ctx, appID, pairData.AssetIn)
			if !found {
				continue
			}
			assetOutData, found := k.asset.GetAsset(ctx, pairData.AssetOut)
			if !found {
				return assettypes.ErrorAssetDoesNotExist
			}
			rateOut, _ := k.GetSnapshotOfPrices(ctx, appID, pairData.AssetOut)
			for _, data := range esmData.AssetsRates {
				if pairData.AssetOut == data.AssetID {
					rateOut = data.Rates
					break
				}
			}
			coolOffData, found := k.GetDataAfterCoolOff(ctx, appID)
			if !found {
				coolOffData.AppId = appID
				var itemc types.AssetToAmount
				var itemd types.AssetToAmount

				itemc.AppId = appID
				itemc.AssetID = assetInData.Id
				itemc.Amount = data.AmountIn
				itemc.IsCollateral = true
				coolOffData.CollateralTotalAmount = k.CalcDollarValueOfToken(ctx, rateIn, data.AmountIn, assetInData.Decimals)
				itemc.Share = sdk.OneDec()

				err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, data.AmountIn)))
				if err != nil {
					return err
				}
				k.SetAssetToAmount(ctx, itemc)

				itemd.AppId = appID
				itemd.AssetID = assetOutData.Id
				itemd.Amount = data.AmountOut
				itemd.IsCollateral = false
				coolOffData.DebtTotalAmount = k.CalcDollarValueOfToken(ctx, rateOut, data.AmountOut, assetOutData.Decimals)
				itemd.Share = sdk.OneDec()
				k.SetAssetToAmount(ctx, itemd)

				// debt token worth ????
				// itemd.DebtTokenWorth = coolOffData.CollateralTotalAmount.Mul(itemd.Share).Quo(sdk.NewDecFromInt(itemd.Amount))

				k.SetDataAfterCoolOff(ctx, coolOffData)
			} else {
				coolOffData.CollateralTotalAmount = coolOffData.CollateralTotalAmount.Add(k.CalcDollarValueOfToken(ctx, rateIn, data.AmountIn, assetInData.Decimals))
				coolOffData.DebtTotalAmount = coolOffData.DebtTotalAmount.Add(k.CalcDollarValueOfToken(ctx, rateOut, data.AmountOut, assetOutData.Decimals))
				assetToAmtInData, found := k.GetAssetToAmount(ctx, appID, assetInData.Id)
				if !found {
					assetToAmtInData.AppId = appID
					assetToAmtInData.AssetID = assetInData.Id
					assetToAmtInData.Amount = data.AmountIn
					assetToAmtInData.IsCollateral = true
				} else {
					assetToAmtInData.Amount = assetToAmtInData.Amount.Add(data.AmountIn)
				}
				assetToAmtOutData, found := k.GetAssetToAmount(ctx, appID, assetOutData.Id)
				if !found {
					assetToAmtOutData.AppId = appID
					assetToAmtOutData.AssetID = assetOutData.Id
					assetToAmtOutData.Amount = data.AmountOut
					assetToAmtOutData.IsCollateral = false
				} else {
					assetToAmtOutData.Amount = assetToAmtOutData.Amount.Add(data.AmountOut)
				}
				err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, data.AmountIn)))
				if err != nil {
					return err
				}
				k.SetDataAfterCoolOff(ctx, coolOffData)
				k.SetAssetToAmount(ctx, assetToAmtInData)
				k.SetAssetToAmount(ctx, assetToAmtOutData)
			}
			k.vault.DeleteVault(ctx, data.Id)
			k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, data.ExtendedPairVaultID, data.Id, data.AppId)
		}
	}
	coolOffData, found := k.GetDataAfterCoolOff(ctx, appID)
	allAssetToAmtData := k.GetAllAssetToAmount(ctx, appID)
	for _, amt := range allAssetToAmtData {
		assetData, found := k.asset.GetAsset(ctx, amt.AssetID)
		if !found {
			return assettypes.ErrorAssetDoesNotExist
		}
		rate, _ := k.GetSnapshotOfPrices(ctx, appID, amt.AssetID)
		for _, data := range esmData.AssetsRates {
			if amt.AssetID == data.AssetID {
				rate = data.Rates
				break
			}
		}
		amtDValue := k.CalcDollarValueOfToken(ctx, rate, amt.Amount, assetData.Decimals)
		if amt.IsCollateral {
			amt.Share = amtDValue.Quo(coolOffData.CollateralTotalAmount)
		} else {
			amt.Share = amtDValue.Quo(coolOffData.DebtTotalAmount)
			debtDValue := amt.Share.Mul(coolOffData.CollateralTotalAmount)
			// amt.DebtTokenWorth = debtDValue.Quo(sdk.Dec(amt.Amount))
			denominator := sdk.NewDecFromInt(sdk.NewIntFromUint64(uint64(assetData.Decimals)))
			numerator := sdk.NewDecFromInt(amt.Amount).Quo(denominator)
			amt.DebtTokenWorth = debtDValue.Quo(numerator)
		}
		k.SetAssetToAmount(ctx, amt)
	}
	esmStatus.VaultRedemptionStatus = true
	k.SetESMStatus(ctx, esmStatus)

	return nil
}

// func (k Keeper) SetUpCollateralRedemptionForStableVault(ctx sdk.Context, appID uint64) error {
// 	totalStableVaults := k.vault.GetStableMintVaults(ctx)
// 	esmStatus, found := k.GetESMStatus(ctx, appID)
// 	if !found {
// 		return types.ErrESMParamsNotFound
// 	}
// 	esmData, _ := k.GetESMTriggerParams(ctx, appID)
// 	for _, data := range totalStableVaults {
// 		if data.AppId == appID {
// 			extendedPairVault, found := k.asset.GetPairsVault(ctx, data.ExtendedPairVaultID)
// 			if !found {
// 				return vaulttypes.ErrorExtendedPairVaultDoesNotExists
// 			}
// 			pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
// 			if !found {
// 				return assettypes.ErrorPairDoesNotExist
// 			}
// 			assetInData, found := k.asset.GetAsset(ctx, pairData.AssetIn)
// 			if !found {
// 				return assettypes.ErrorAssetDoesNotExist
// 			}
// 			rateIn, found := k.GetSnapshotOfPrices(ctx, appID, pairData.AssetIn)
// 			if !found {
// 				continue
// 			}
// 			assetOutData, found := k.asset.GetAsset(ctx, pairData.AssetOut)
// 			if !found {
// 				return assettypes.ErrorAssetDoesNotExist
// 			}
// 			rateOut, found := k.GetSnapshotOfPrices(ctx, appID, pairData.AssetOut)
// 			if !found {
// 				for _, data := range esmData.AssetsRates {
// 					if pairData.AssetOut == data.AssetID {
// 						rateOut = data.Rates
// 						break
// 					}
// 					rateOut = 0
// 				}
// 			}
// 			if rateOut == 0 {
// 				continue
// 			}
// 			coolOffData, found := k.GetDataAfterCoolOff(ctx, appID)
// 			if !found {
// 				coolOffData.AppId = appID
// 				var itemc types.AssetToAmount
// 				var itemd types.AssetToAmount

// 				itemc.AssetID = assetInData.Id
// 				itemc.Amount = data.AmountIn
// 				x := sdk.NewIntFromUint64(rateIn).Mul(data.AmountIn)
// 				coolOffData.CollateralTotalAmount = sdk.NewDecFromInt(x)
// 				itemc.Share = sdk.NewDecFromInt(x).Quo(coolOffData.CollateralTotalAmount)

// 				err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, data.AmountIn)))
// 				if err != nil {
// 					return err
// 				}
// 				coolOffData.CollateralAsset = append(coolOffData.CollateralAsset, itemc)

// 				itemd.AssetID = assetOutData.Id
// 				itemd.Amount = data.AmountOut
// 				y := sdk.NewIntFromUint64(rateOut).Mul(data.AmountOut)
// 				coolOffData.DebtTotalAmount = sdk.NewDecFromInt(y)
// 				itemd.Share = sdk.NewDecFromInt(y).Quo(coolOffData.DebtTotalAmount)
// 				itemd.DebtTokenWorth = coolOffData.CollateralTotalAmount.Mul(itemd.Share).Quo(sdk.NewDecFromInt(itemd.Amount))
// 				coolOffData.DebtAsset = append(coolOffData.DebtAsset, itemd)

// 				k.SetDataAfterCoolOff(ctx, coolOffData)
// 			} else {
// 				count := 0
// 				for i, indata := range coolOffData.CollateralAsset {
// 					if indata.AssetID == assetInData.Id {
// 						count++
// 						indata.Amount = indata.Amount.Add(data.AmountIn)
// 						x := sdk.NewIntFromUint64(rateIn).Mul(data.AmountIn)
// 						coolOffData.CollateralTotalAmount = coolOffData.CollateralTotalAmount.Add(sdk.NewDecFromInt(x))
// 						indata.Share = sdk.NewDecFromInt(x).Quo(coolOffData.CollateralTotalAmount)

// 						err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, data.AmountIn)))
// 						if err != nil {
// 							return err
// 						}
// 						coolOffData.CollateralAsset = append(coolOffData.CollateralAsset[:i], coolOffData.CollateralAsset[i+1:]...)
// 						coolOffData.CollateralAsset = append(coolOffData.CollateralAsset, indata)
// 						break
// 					}
// 				}
// 				if count == 0 {
// 					var item types.AssetToAmount

// 					item.AssetID = assetInData.Id
// 					item.Amount = data.AmountIn
// 					x := sdk.NewIntFromUint64(rateIn).Mul(data.AmountIn)
// 					coolOffData.CollateralTotalAmount = coolOffData.CollateralTotalAmount.Add(sdk.NewDecFromInt(x))
// 					item.Share = sdk.NewDecFromInt(x).Quo(coolOffData.CollateralTotalAmount)

// 					err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetInData.Denom, data.AmountIn)))
// 					if err != nil {
// 						return err
// 					}
// 					coolOffData.CollateralAsset = append(coolOffData.CollateralAsset, item)
// 				}
// 				count = 0
// 				for i, indatadebt := range coolOffData.DebtAsset {
// 					if indatadebt.AssetID == assetOutData.Id {
// 						count++
// 						indatadebt.Amount = indatadebt.Amount.Add(data.AmountOut)
// 						x := sdk.NewIntFromUint64(rateOut).Mul(data.AmountOut)
// 						coolOffData.DebtTotalAmount = coolOffData.DebtTotalAmount.Add(sdk.NewDecFromInt(x))
// 						indatadebt.Share = sdk.NewDecFromInt(x).Quo(coolOffData.DebtTotalAmount)
// 						indatadebt.DebtTokenWorth = coolOffData.CollateralTotalAmount.Mul(indatadebt.Share).Quo(sdk.NewDecFromInt(indatadebt.Amount))
// 						coolOffData.DebtAsset = append(coolOffData.DebtAsset[:i], coolOffData.DebtAsset[i+1:]...)
// 						coolOffData.DebtAsset = append(coolOffData.DebtAsset, indatadebt)
// 						break
// 					}
// 				}
// 				if count == 0 {
// 					var item types.AssetToAmount

// 					item.AssetID = assetOutData.Id
// 					item.Amount = data.AmountOut
// 					x := sdk.NewIntFromUint64(rateOut).Mul(data.AmountOut)
// 					coolOffData.DebtTotalAmount = coolOffData.DebtTotalAmount.Add(sdk.NewDecFromInt(x))
// 					item.Share = sdk.NewDecFromInt(x).Quo(coolOffData.DebtTotalAmount)
// 					item.DebtTokenWorth = coolOffData.CollateralTotalAmount.Mul(item.Share).Quo(sdk.NewDecFromInt(item.Amount))
// 					coolOffData.DebtAsset = append(coolOffData.DebtAsset, item)
// 				}
// 				k.SetDataAfterCoolOff(ctx, coolOffData)
// 			}
// 		}
// 	}
// 	netFee, found1 := k.collector.GetAppNetFeeCollectedData(ctx, appID)
// 	coolOffData, found := k.GetDataAfterCoolOff(ctx, appID)
// 	if !found {
// 		return nil
// 	}
// 	if found1 {
// 		for _, data := range netFee {
// 			for i, indatanet := range coolOffData.DebtAsset {
// 				if data.AssetId == indatanet.AssetID {
// 					rateOut, found := k.GetSnapshotOfPrices(ctx, appID, data.AssetId)
// 					if !found {
// 						for _, datan := range esmData.AssetsRates {
// 							if data.AssetId == datan.AssetID {
// 								rateOut = datan.Rates
// 								break
// 							}
// 							rateOut = 0
// 						}
// 					}
// 					indatanet.Amount = indatanet.Amount.Sub(data.NetFeesCollected)
// 					x, _ := k.CalcDollarValueOfToken(ctx, data.AssetId, rateOut, indatanet.Amount)
// 					// x := sdk.NewIntFromUint64(rateOut).Mul(indatanet.Amount)
// 					coolOffData.DebtTotalAmount = coolOffData.DebtTotalAmount.Sub(x)
// 					indatanet.Share = x.Quo(coolOffData.DebtTotalAmount)
// 					indatanet.DebtTokenWorth = coolOffData.CollateralTotalAmount.Mul(indatanet.Share).Quo(sdk.NewDecFromInt(indatanet.Amount))
// 					coolOffData.DebtAsset = append(coolOffData.DebtAsset[:i], coolOffData.DebtAsset[i+1:]...)
// 					coolOffData.DebtAsset = append(coolOffData.DebtAsset, indatanet)
// 				}
// 			}
// 		}
// 	}
// 	esmStatus.StableVaultRedemptionStatus = true
// 	k.SetESMStatus(ctx, esmStatus)
// 	k.SetDataAfterCoolOff(ctx, coolOffData)
// 	return nil
// }

func (k Keeper) SnapshotOfPrices(ctx sdk.Context, esmStatus types.ESMStatus) error {
	assets := k.asset.GetAssets(ctx)
	for _, a := range assets {
		if a.IsOraclePriceRequired {
			price, found := k.market.GetTwa(ctx, a.Id)
			// not checking is price active as this fn is called at the end of protocol and active relayer service is not certain
			// so, we are not implementing is price active field.
			if !found {
				continue
			}
			k.SetSnapshotOfPrices(ctx, esmStatus.AppId, a.Id, price.Twa)
		}
	}
	esmStatus.SnapshotStatus = true
	k.SetESMStatus(ctx, esmStatus)
	return nil
}

func (k Keeper) SetSnapshotOfPrices(ctx sdk.Context, appID, assetID, price uint64) {
	var (
		store = k.Store(ctx)
		key   = types.SnapshotTypeKey(appID, assetID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: price,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetSnapshotOfPrices(ctx sdk.Context, appID, assetID uint64) (price uint64, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.SnapshotTypeKey(appID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return price, false
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)
	return id.GetValue(), true
}

func (k Keeper) SetAssetToAmount(ctx sdk.Context, assetToAmount types.AssetToAmount) {
	var (
		store = k.Store(ctx)
		key   = types.AssetToAmountKey(assetToAmount.AppId, assetToAmount.AssetID)
		value = k.cdc.MustMarshal(&assetToAmount)
	)

	store.Set(key, value)
}

func (k Keeper) GetAssetToAmount(ctx sdk.Context, appID, assetId uint64) (assetToAmount types.AssetToAmount, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetToAmountKey(appID, assetId)
		value = store.Get(key)
	)

	if value == nil {
		return assetToAmount, false
	}
	k.cdc.MustUnmarshal(value, &assetToAmount)
	return assetToAmount, true
}

func (k Keeper) GetAllAssetToAmount(ctx sdk.Context, appID uint64) (assetToAmount []types.AssetToAmount) {
	var (
		store = k.Store(ctx)
		key   = types.AppAssetToAmountKey(appID)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var esm types.AssetToAmount
		k.cdc.MustUnmarshal(iter.Value(), &esm)
		assetToAmount = append(assetToAmount, esm)
	}
	return assetToAmount
}

func (k Keeper) CalcDollarValueOfToken(ctx sdk.Context, rate uint64, amt sdk.Int, decimals int64) (price sdk.Dec) {
	numerator := sdk.NewDecFromInt(amt).Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(rate)))
	denominator := sdk.NewDecFromInt(sdk.NewIntFromUint64(uint64(decimals)))
	return numerator.Quo(denominator)
}
