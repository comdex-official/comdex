package keeper

import (
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	"github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) SetIDForVault(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.VaultIDPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetIDForVault(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.VaultIDPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetUserAppExtendedPairMappingData(ctx sdk.Context, mappingData types.OwnerAppExtendedPairVaultMappingData) {
	var (
		store = k.Store(ctx)
		key   = types.UserAppExtendedPairMappingKey(mappingData.Owner, mappingData.AppId, mappingData.ExtendedPairId)
		value = k.cdc.MustMarshal(&mappingData)
	)

	store.Set(key, value)
}

func (k Keeper) GetUserAppExtendedPairMappingData(ctx sdk.Context, address string, appID uint64, pairVaultID uint64) (mappingData types.OwnerAppExtendedPairVaultMappingData, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserAppExtendedPairMappingKey(address, appID, pairVaultID)
		value = store.Get(key)
	)

	if value == nil {
		return mappingData, false
	}

	k.cdc.MustUnmarshal(value, &mappingData)
	return mappingData, true
}

func (k Keeper) GetUserAppMappingData(ctx sdk.Context, address string, appID uint64) (mappingData []types.OwnerAppExtendedPairVaultMappingData, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserAppMappingKey(address, appID)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.OwnerAppExtendedPairVaultMappingData
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		mappingData = append(mappingData, mapData)
	}
	if mappingData == nil {
		return nil, false
	}

	return mappingData, true
}

func (k Keeper) DeleteUserVaultExtendedPairMapping(ctx sdk.Context, address string, appID uint64, pairVaultID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.UserAppExtendedPairMappingKey(address, appID, pairVaultID)
	)

	store.Delete(key)
}

func (k Keeper) GetUserMappingData(ctx sdk.Context, address string) (mappingData []types.OwnerAppExtendedPairVaultMappingData) {
	var (
		store = k.Store(ctx)
		key   = types.UserKey(address)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.OwnerAppExtendedPairVaultMappingData
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		mappingData = append(mappingData, mapData)
	}

	return mappingData
}

func (k Keeper) GetAllUserVaultExtendedPairMapping(ctx sdk.Context) (userVaultAssetData []types.OwnerAppExtendedPairVaultMappingData) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.UserVaultExtendedPairMappingKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var vault types.OwnerAppExtendedPairVaultMappingData
		k.cdc.MustUnmarshal(iter.Value(), &vault)
		userVaultAssetData = append(userVaultAssetData, vault)
	}
	return userVaultAssetData
}

// CheckUserAppToExtendedPairMapping Checking if for a certain user for the app type , whether there exists a certain asset or not and if it contains a locker id or not .
// func (k Keeper) CheckUserAppToExtendedPairMapping(ctx sdk.Context, userVaultAssetData types.UserVaultAssetMapping, extendedPairVaultID uint64, appMappingID uint64) (vaultID uint64, found bool) {
// 	for _, vaultAppMapping := range userVaultAssetData.UserVaultApp {
// 		if vaultAppMapping.AppId == appMappingID {
// 			for _, extendedPairToVaultIDMapping := range vaultAppMapping.UserExtendedPairVault {
// 				if extendedPairToVaultIDMapping.ExtendedPairId == extendedPairVaultID && extendedPairToVaultIDMapping.VaultId > 0 {
// 					vaultID = extendedPairToVaultIDMapping.VaultId
// 					return vaultID, true
// 				}
// 			}
// 		}
// 	}
// 	return vaultID, false
// }
// func (k Keeper) CheckUserToAppMapping(ctx sdk.Context, userVaultAssetData types.UserVaultAssetMapping, appMappingID uint64) (found bool) {
// 	for _, vaultAppMapping := range userVaultAssetData.UserVaultApp {
// 		if vaultAppMapping.AppId == appMappingID {
// 			return true
// 		}
// 	}
// 	return false
// }

// SetAppExtendedPairVaultMappingData Set AppExtendedPairVaultMapping to check the current status of the vault by extended pair vault id .
func (k Keeper) SetAppExtendedPairVaultMappingData(ctx sdk.Context, appExtendedPairVaultData types.AppExtendedPairVaultMappingData) error {
	var (
		store = k.Store(ctx)
		key   = types.AppExtendedPairVaultMappingKey(appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId)
		value = k.cdc.MustMarshal(&appExtendedPairVaultData)
	)

	store.Set(key, value)
	return nil
}

//Get AppExtendedPairVaultMapping to check the current status of the vault by extended pair vault id

func (k Keeper) GetAppExtendedPairVaultMappingData(ctx sdk.Context, appMappingID uint64, pairVaultID uint64) (appExtendedPairVaultData types.AppExtendedPairVaultMappingData, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AppExtendedPairVaultMappingKey(appMappingID, pairVaultID)
		value = store.Get(key)
	)

	if value == nil {
		return appExtendedPairVaultData, false
	}

	k.cdc.MustUnmarshal(value, &appExtendedPairVaultData)
	return appExtendedPairVaultData, true
}

func (k Keeper) GetAppMappingData(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData []types.AppExtendedPairVaultMappingData, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AppMappingKey(appMappingID)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.AppExtendedPairVaultMappingData
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		appExtendedPairVaultData = append(appExtendedPairVaultData, mapData)
	}
	if appExtendedPairVaultData == nil {
		return nil, false
	}

	return appExtendedPairVaultData, true
}

func (k Keeper) GetAllAppExtendedPairVaultMapping(ctx sdk.Context) (appExtendedPairVaultData []types.AppExtendedPairVaultMappingData) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AppExtendedPairVaultMappingKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var vault types.AppExtendedPairVaultMappingData
		k.cdc.MustUnmarshal(iter.Value(), &vault)
		appExtendedPairVaultData = append(appExtendedPairVaultData, vault)
	}
	return appExtendedPairVaultData
}

//Check AppExtendedPairVault Data,
//If exists fine --- go with the next steps from here
//else instantiate 1 and set it. and go for the next steps from here
//So best way will be to create a function which will first check if AppExtendedPairVault Data exists or not. If it does. then send counted value. else create a struct save it. and send counter value.

func (k Keeper) CheckAppExtendedPairVaultMapping(ctx sdk.Context, appMappingID uint64, extendedPairVaultID uint64) (mintedStatistics sdk.Int, lenVaults uint64) {
	appExtendedPairVaultData, found := k.GetAppExtendedPairVaultMappingData(ctx, appMappingID, extendedPairVaultID)
	if !found {
		//Initialising a new struct
		var newAppExtendedPairVault types.AppExtendedPairVaultMappingData
		newAppExtendedPairVault.AppId = appMappingID
		zeroVal := sdk.ZeroInt()
		newAppExtendedPairVault.ExtendedPairId = extendedPairVaultID
		newAppExtendedPairVault.CollateralLockedAmount = zeroVal
		newAppExtendedPairVault.TokenMintedAmount = zeroVal

		err := k.SetAppExtendedPairVaultMappingData(ctx, newAppExtendedPairVault)
		if err != nil {
			return sdk.Int{}, 0
		}

		return newAppExtendedPairVault.TokenMintedAmount, 0
	}

	//Check the Zero Value once
	zeroVal := sdk.ZeroInt()
	var newExtendedPairVault types.AppExtendedPairVaultMappingData
	newExtendedPairVault.ExtendedPairId = extendedPairVaultID
	newExtendedPairVault.CollateralLockedAmount = zeroVal
	newExtendedPairVault.TokenMintedAmount = zeroVal

	err := k.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
	if err != nil {
		return sdk.Int{}, 0
	}

	return newExtendedPairVault.TokenMintedAmount, 0
}

func (k Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, vaultData types.Vault) {
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, vaultData.AppId, vaultData.ExtendedPairVaultID)

	addedMintedData := appExtendedPairVaultData.TokenMintedAmount.Add(vaultData.AmountOut)
	addedCollateralData := appExtendedPairVaultData.CollateralLockedAmount.Add(vaultData.AmountIn)
	appExtendedPairVaultData.TokenMintedAmount = addedMintedData
	appExtendedPairVaultData.CollateralLockedAmount = addedCollateralData
	appExtendedPairVaultData.VaultIds = append(appExtendedPairVaultData.VaultIds, vaultData.Id)

	err := k.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
	if err != nil {
		return
	}
}

func (k Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreateStableMintVault(ctx sdk.Context, vaultData types.StableMintVault) {
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, vaultData.AppId, vaultData.ExtendedPairVaultID)

	addedMintedData := appExtendedPairVaultData.TokenMintedAmount.Add(vaultData.AmountOut)
	addedCollateralData := appExtendedPairVaultData.CollateralLockedAmount.Add(vaultData.AmountIn)
	appExtendedPairVaultData.TokenMintedAmount = addedMintedData
	appExtendedPairVaultData.CollateralLockedAmount = addedCollateralData
	appExtendedPairVaultData.VaultIds = append(appExtendedPairVaultData.VaultIds, vaultData.Id)

	err := k.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
	if err != nil {
		return
	}
}

// CalculateCollaterlizationRatio Calculate Collaterlization Ratio .
func (k Keeper) CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error) {
	extendedPairVault, found := k.GetPairsVault(ctx, extendedPairVaultID)
	if !found {
		return sdk.ZeroDec(), types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return sdk.ZeroDec(), types.ErrorPairDoesNotExist
	}
	assetInData, found := k.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return sdk.ZeroDec(), types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return sdk.ZeroDec(), types.ErrorAssetDoesNotExist
	}

	assetInPrice, found := k.GetPriceForAsset(ctx, assetInData.Id)
	if !found {
		return sdk.ZeroDec(), types.ErrorPriceDoesNotExist
	}
	esmStatus, found := k.GetESMStatus(ctx, extendedPairVault.AppId)
	statusEsm := false
	if found {
		statusEsm = esmStatus.Status
	}
	if statusEsm && esmStatus.SnapshotStatus {
		price, found := k.GetSnapshotOfPrices(ctx, extendedPairVault.AppId, assetInData.Id)
		if !found {
			return sdk.ZeroDec(), types.ErrorPriceDoesNotExist
		}
		assetInPrice = price
	}
	var assetOutPrice uint64

	if extendedPairVault.AssetOutOraclePrice {
		//If oracle Price required for the assetOut
		if statusEsm && esmStatus.SnapshotStatus {
			price, found := k.GetSnapshotOfPrices(ctx, extendedPairVault.AppId, assetOutData.Id)
			if !found {
				return sdk.ZeroDec(), types.ErrorPriceDoesNotExist
			}
			assetOutPrice = price
		} else {
			assetOutPrice, found = k.GetPriceForAsset(ctx, assetOutData.Id)
			if !found {
				return sdk.ZeroDec(), types.ErrorPriceDoesNotExist
			}
		}
	} else {
		//If oracle Price is not required for the assetOut
		assetOutPrice = extendedPairVault.AssetOutPrice
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

func (k Keeper) VerifyCollaterlizationRatio(
	ctx sdk.Context,
	extendedPairVaultID uint64,
	amountIn sdk.Int,
	amountOut sdk.Int,
	minCrRequired sdk.Dec,
	statusEsm bool,
) error {
	collaterlizationRatio, err := k.CalculateCollaterlizationRatio(ctx, extendedPairVaultID, amountIn, amountOut)
	if err != nil {
		return err
	}
	if collaterlizationRatio.LT(minCrRequired) && !statusEsm {
		return types.ErrorInvalidCollateralizationRatio
	} else if collaterlizationRatio.LT(sdk.MustNewDecFromStr("1")) && statusEsm {
		return types.ErrorInvalidCollateralizationRatio
	}
	return nil
}

// func (k Keeper) SetVaultID(ctx sdk.Context, id uint64) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.VaultIDKey
// 		value = k.cdc.MustMarshal(
// 			&protobuftypes.UInt64Value{
// 				Value: id,
// 			},
// 		)
// 	)

// 	store.Set(key, value)
// }

// func (k Keeper) GetVaultID(ctx sdk.Context) uint64 {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.VaultIDKey
// 		value = store.Get(key)
// 	)

// 	if value == nil {
// 		return 0
// 	}

// 	var id protobuftypes.UInt64Value
// 	k.cdc.MustUnmarshal(value, &id)

// 	return id.GetValue()
// }

func (k Keeper) SetVault(ctx sdk.Context, vault types.Vault) {
	var (
		store = k.Store(ctx)
		key   = types.VaultKey(vault.Id)
		value = k.cdc.MustMarshal(&vault)
	)
	store.Set(key, value)
}

func (k Keeper) GetVault(ctx sdk.Context, id uint64) (vault types.Vault, found bool) {
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

// UpdateCollateralLockedAmountLockerMapping For updating token stats of collateral .
func (k Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, vaultLookupData types.AppExtendedPairVaultMappingData, amount sdk.Int, changeType bool) {
	//if Change type true = Add to collateral Locked
	//If change type false = Subtract from the collateral Locked

	if changeType {
		updatedVal := vaultLookupData.CollateralLockedAmount.Add(amount)
		vaultLookupData.CollateralLockedAmount = updatedVal
	} else {
		updatedVal := vaultLookupData.CollateralLockedAmount.Sub(amount)
		vaultLookupData.CollateralLockedAmount = updatedVal
	}

	err := k.SetAppExtendedPairVaultMappingData(ctx, vaultLookupData)
	if err != nil {
		return
	}
}

// UpdateTokenMintedAmountLockerMapping For updating token stats of minted .
func (k Keeper) UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, vaultLookupData types.AppExtendedPairVaultMappingData, amount sdk.Int, changeType bool) {
	//if Change type true = Add to token Locked
	//If change type false = Subtract from the token Locked

	if changeType {
		updatedVal := vaultLookupData.TokenMintedAmount.Add(amount)
		vaultLookupData.TokenMintedAmount = updatedVal
	} else {
		updatedVal := vaultLookupData.TokenMintedAmount.Sub(amount)
		vaultLookupData.TokenMintedAmount = updatedVal
	}
	err := k.SetAppExtendedPairVaultMappingData(ctx, vaultLookupData)
	if err != nil {
		return
	}
}

func (k Keeper) DeleteVault(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.VaultKey(id)
	)

	store.Delete(key)
}

func (k Keeper) GetVaults(ctx sdk.Context) (vaults []types.Vault) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.VaultKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var vault types.Vault
		k.cdc.MustUnmarshal(iter.Value(), &vault)
		vaults = append(vaults, vault)
	}
	return vaults
}

// func (k Keeper) UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairID uint64, userAddress string, appMappingID uint64) {
// 	userData, found := k.GetUserVaultExtendedPairMapping(ctx, userAddress)

// 	var dataIndex int
// 	if found {
// 		for _, appData := range userData.UserVaultApp {
// 			if appData.AppId == appMappingID {
// 				for index, extendedPairData := range appData.UserExtendedPairVault {
// 					if extendedPairData.ExtendedPairId == extendedPairID {
// 						dataIndex = index
// 					}
// 				}
// 				appData.UserExtendedPairVault = append(appData.UserExtendedPairVault[:dataIndex], appData.UserExtendedPairVault[dataIndex+1:]...)
// 				break
// 			}
// 		}
// 		k.SetUserVaultExtendedPairMapping(ctx, userData)
// 	}
// }

func (k Keeper) DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID uint64, appMappingID uint64) {
	appExtendedPairVaultData, found := k.GetAppExtendedPairVaultMappingData(ctx, appMappingID, extendedPairID)

	var dataIndex int
	if found {
		for index, vaultID := range appExtendedPairVaultData.VaultIds {
			if vaultID == userVaultID {
				dataIndex = index
			}
		}
		appExtendedPairVaultData.VaultIds = append(appExtendedPairVaultData.VaultIds[:dataIndex], appExtendedPairVaultData.VaultIds[dataIndex+1:]...)

		err := k.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
		if err != nil {
			return
		}
	}
}

func (k Keeper) SetIDForStableVault(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.StableVaultIDPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetIDForStableVault(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.StableVaultIDPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetStableMintVault(ctx sdk.Context, stableVault types.StableMintVault) {
	var (
		store = k.Store(ctx)
		key   = types.StableMintVaultKey(stableVault.Id)
		value = k.cdc.MustMarshal(&stableVault)
	)

	store.Set(key, value)
}

func (k Keeper) GetStableMintVault(ctx sdk.Context, id uint64) (stableVault types.StableMintVault, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.StableMintVaultKey(id)
		value = store.Get(key)
	)
	if value == nil {
		return stableVault, false
	}

	k.cdc.MustUnmarshal(value, &stableVault)
	return stableVault, true
}

func (k Keeper) GetStableMintVaults(ctx sdk.Context) (stableVaults []types.StableMintVault) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.StableMintVaultKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var stableVault types.StableMintVault
		k.cdc.MustUnmarshal(iter.Value(), &stableVault)
		stableVaults = append(stableVaults, stableVault)
	}
	return stableVaults
}

func (k Keeper) CreateNewVault(ctx sdk.Context, From string, AppID uint64, ExtendedPairVaultID uint64, AmountIn sdk.Int, AmountOut sdk.Int) error {
	appMapping, _ := k.GetApp(ctx, AppID)
	extendedPairVault, _ := k.GetPairsVault(ctx, ExtendedPairVaultID)

	zeroVal := sdk.ZeroInt()
	oldID := k.GetIDForVault(ctx)
	var newVault types.Vault
	newID := oldID + 1
	newVault.Id = newID
	newVault.AmountIn = AmountIn

	newVault.ClosingFeeAccumulated = zeroVal
	newVault.AmountOut = AmountOut
	newVault.AppId = appMapping.Id
	newVault.InterestAccumulated = zeroVal
	newVault.Owner = From
	newVault.CreatedAt = ctx.BlockTime()
	newVault.ExtendedPairVaultID = extendedPairVault.Id
	k.SetVault(ctx, newVault)
	k.SetIDForVault(ctx, newID)

	//get extendedpair lookup table data
	// push the new vault id in that extended pair of the app
	// update the counter of the extendedpair lookup tabe
	////// ///////

	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, AppID, ExtendedPairVaultID)

	appExtendedPairVaultData.VaultIds = append(appExtendedPairVaultData.VaultIds, newVault.Id)

	err := k.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
	if err != nil {
		return err
	}

	//////////////////
	// k.UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx, updated_counter, new_vault)

	// userVaultExtendedPairMappingData, _ := k.GetUserVaultExtendedPairMapping(ctx,AppID, From, ExtendedPairVaultID)

	//So only need to add the locker id with asset
	// var userExtendedPairData types.ExtendedPairToVaultMapping
	// userExtendedPairData.VaultId = newVault.Id
	// userExtendedPairData.ExtendedPairId = newVault.ExtendedPairVaultID

	// for _, appData := range userVaultExtendedPairMappingData.UserVaultApp {
	// 	if appData.AppId == appMapping.Id {

	// 		appData.UserExtendedPairVault = append(appData.UserExtendedPairVault, &userExtendedPairData)
	// 	}

	// }
	var mappingData types.OwnerAppExtendedPairVaultMappingData
	mappingData.Owner = From
	mappingData.AppId = AppID
	mappingData.ExtendedPairId = ExtendedPairVaultID
	mappingData.VaultId = newVault.Id

	k.SetUserAppExtendedPairMappingData(ctx, mappingData)

	return nil
}

func (k Keeper) calculateUserToken(ctx sdk.Context, userVault types.Vault, cr sdk.Dec, amountIn sdk.Int, assetInPrice, assetOutPrice uint64) (userToken sdk.Int) {
	nume := userVault.AmountIn.Add(amountIn)
	newNum := nume.Mul(sdk.NewIntFromUint64(assetInPrice))
	denom := sdk.NewIntFromUint64(assetOutPrice).Mul(cr.TruncateInt())
	newVal := newNum.Quo(denom)
	userToken = newVal.Sub(userVault.AmountOut)

	return userToken
}
func (k Keeper) WasmMsgAddEmissionRewards(ctx sdk.Context, appID uint64, amount sdk.Int, extPair []uint64, votingRatio []sdk.Int) error {
	var assetID uint64
	var perUserShareByAmt sdk.Int
	var vaultsData types.AppExtendedPairVaultMappingData

	totalVote := sdk.ZeroInt()
	app, _ := k.GetApp(ctx, appID)
	govToken := app.GenesisToken
	for _, v := range govToken {
		if v.IsGovToken {
			assetID = v.AssetId
		}
	}
	asset, _ := k.GetAsset(ctx, assetID)
	if amount.GT(sdk.ZeroInt()) {
		err := k.MintCoin(ctx, tokenminttypes.ModuleName, sdk.NewCoin(asset.Denom, amount))
		if err != nil {
			return err
		}
	}
	k.UpdateAssetDataInTokenMintByApp(ctx, appID, assetID, true, amount)
	for i, _ := range votingRatio {
		totalVote = totalVote.Add(votingRatio[i])
	}
	for j, extP := range extPair {
		extPairVaultMappingData, _ := k.GetAppExtendedPairVaultMappingData(ctx, appID, extP)
		individualVote := votingRatio[j]
		votingR := individualVote.Quo(totalVote)
		shareByExtPair := votingR.Mul(amount)
		perUserShareByAmt = shareByExtPair.Quo(extPairVaultMappingData.TokenMintedAmount)
		vaultsData, _ = k.GetAppExtendedPairVaultMappingData(ctx, appID, extP)

		for _, vaultID := range vaultsData.VaultIds {
			vault, _ := k.GetVault(ctx, vaultID)
			amt := vault.AmountOut.Mul(perUserShareByAmt)
			addr, _ := sdk.AccAddressFromBech32(vault.Owner)
			if amt.GT(sdk.ZeroInt()) {
				err := k.SendCoinFromModuleToAccount(ctx, tokenminttypes.ModuleName, addr, sdk.NewCoin(asset.Denom, amt))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
