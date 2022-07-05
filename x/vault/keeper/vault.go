package keeper

import (
	"strconv"
	"github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) SetUserVaultExtendedPairMapping(ctx sdk.Context, userVaultAssetData types.UserVaultAssetMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserVaultExtendedPairMappingKey(userVaultAssetData.Owner)
		value = k.cdc.MustMarshal(&userVaultAssetData)
	)

	store.Set(key, value)
}

func (k *Keeper) GetUserVaultExtendedPairMapping(ctx sdk.Context, address string) (userVaultAssetData types.UserVaultAssetMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserVaultExtendedPairMappingKey(address)
		value = store.Get(key)
	)

	if value == nil {
		return userVaultAssetData, false
	}

	k.cdc.MustUnmarshal(value, &userVaultAssetData)
	return userVaultAssetData, true
}

// CheckUserAppToExtendedPairMapping Checking if for a certain user for the app type , whether there exists a certain asset or not and if it contains a locker id or not .
func (k *Keeper) CheckUserAppToExtendedPairMapping(ctx sdk.Context, userVaultAssetData types.UserVaultAssetMapping, extendedPairVaultID uint64, appMappingID uint64) (vaultID string, found bool) {
	for _, vaultAppMapping := range userVaultAssetData.UserVaultApp {
		if vaultAppMapping.AppId == appMappingID {
			for _, extendedPairToVaultIDMapping := range vaultAppMapping.UserExtendedPairVault {
				if extendedPairToVaultIDMapping.ExtendedPairId == extendedPairVaultID && len(extendedPairToVaultIDMapping.VaultId) > 0 {
					vaultID = extendedPairToVaultIDMapping.VaultId
					return vaultID, true
				}
			}
		}
	}
	return vaultID, false
}
func (k *Keeper) CheckUserToAppMapping(ctx sdk.Context, userVaultAssetData types.UserVaultAssetMapping, appMappingID uint64) (found bool) {
	for _, vaultAppMapping := range userVaultAssetData.UserVaultApp {
		if vaultAppMapping.AppId == appMappingID {
			return true
		}
	}
	return false
}

// SetAppExtendedPairVaultMapping Set AppExtendedPairVaultMapping to check the current status of the vault by extended pair vault id .
func (k *Keeper) SetAppExtendedPairVaultMapping(ctx sdk.Context, appExtendedPairVaultData types.AppExtendedPairVaultMapping) error {
	var (
		store = k.Store(ctx)
		key   = types.AppExtendedPairVaultMappingKey(appExtendedPairVaultData.AppId)
		value = k.cdc.MustMarshal(&appExtendedPairVaultData)
	)

	store.Set(key, value)
	return nil
}

//Get AppExtendedPairVaultMapping to check the current status of the vault by extended pair vault id

func (k *Keeper) GetAppExtendedPairVaultMapping(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData types.AppExtendedPairVaultMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AppExtendedPairVaultMappingKey(appMappingID)
		value = store.Get(key)
	)

	if value == nil {
		return appExtendedPairVaultData, false
	}

	k.cdc.MustUnmarshal(value, &appExtendedPairVaultData)
	return appExtendedPairVaultData, true
}

//Check AppExtendedPairVault Data,
//If exists fine --- go with the next steps from here
//else instantiate 1 and set it. and go for the next steps from here
//So best way will be to create a function which will first check if AppExtendedPairVault Data exists or not. If it does. then send counted value. else create a struct save it. and send counter value.

func (k *Keeper) CheckAppExtendedPairVaultMapping(ctx sdk.Context, appMappingID uint64, extendedPairVaultID uint64) (counter uint64, mintedStatistics sdk.Int, lenVaults uint64) {
	appExtendedPairVaultData, found := k.GetAppExtendedPairVaultMapping(ctx, appMappingID)
	if !found {
		//Initialising a new struct
		var newAppExtendedPairVault types.AppExtendedPairVaultMapping
		var newExtendedPairVault types.ExtendedPairVaultMapping
		newAppExtendedPairVault.AppId = appMappingID
		newAppExtendedPairVault.Counter = 0
		zeroVal := sdk.ZeroInt()
		newExtendedPairVault.ExtendedPairId = extendedPairVaultID
		newExtendedPairVault.CollateralLockedAmount = zeroVal
		newExtendedPairVault.TokenMintedAmount = zeroVal
		newAppExtendedPairVault.ExtendedPairVaults = append(newAppExtendedPairVault.ExtendedPairVaults, &newExtendedPairVault)

		err := k.SetAppExtendedPairVaultMapping(ctx, newAppExtendedPairVault)
		if err != nil {
			return 0, sdk.Int{}, 0
		}

		return newAppExtendedPairVault.Counter, newExtendedPairVault.TokenMintedAmount, 0
	}

	for _, extendedPairVaultData := range appExtendedPairVaultData.ExtendedPairVaults {
		if extendedPairVaultData.ExtendedPairId == extendedPairVaultID {
			lenOfVaults := len(appExtendedPairVaultData.ExtendedPairVaults)
			return appExtendedPairVaultData.Counter, extendedPairVaultData.TokenMintedAmount, uint64(lenOfVaults)
		}
	}
	//Check the Zero Value once
	zeroVal := sdk.ZeroInt()
	var newExtendedPairVault types.ExtendedPairVaultMapping
	newExtendedPairVault.ExtendedPairId = extendedPairVaultID
	newExtendedPairVault.CollateralLockedAmount = zeroVal
	newExtendedPairVault.TokenMintedAmount = zeroVal
	appExtendedPairVaultData.ExtendedPairVaults = append(appExtendedPairVaultData.ExtendedPairVaults, &newExtendedPairVault)

	err := k.SetAppExtendedPairVaultMapping(ctx, appExtendedPairVaultData)
	if err != nil {
		return 0, sdk.Int{}, 0
	}

	return appExtendedPairVaultData.Counter, newExtendedPairVault.TokenMintedAmount, 0
}

func (k *Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, counter uint64, vaultData types.Vault) {
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMapping(ctx, vaultData.AppId)

	appExtendedPairVaultData.Counter = counter

	for _, appData := range appExtendedPairVaultData.ExtendedPairVaults {
		if appData.ExtendedPairId == vaultData.ExtendedPairVaultID {
			addedMintedData := appData.TokenMintedAmount.Add(vaultData.AmountOut)
			addedCollateralData := appData.CollateralLockedAmount.Add(vaultData.AmountIn)
			appData.TokenMintedAmount = addedMintedData
			appData.CollateralLockedAmount = addedCollateralData
			appData.VaultIds = append(appData.VaultIds, vaultData.Id)
		}
	}

	err := k.SetAppExtendedPairVaultMapping(ctx, appExtendedPairVaultData)
	if err != nil {
		return
	}
}

func (k *Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreateStableMintVault(ctx sdk.Context, counter uint64, vaultData types.StableMintVault) {
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMapping(ctx, vaultData.AppId)

	appExtendedPairVaultData.Counter = counter

	for _, appData := range appExtendedPairVaultData.ExtendedPairVaults {
		if appData.ExtendedPairId == vaultData.ExtendedPairVaultID {
			addedMintedData := appData.TokenMintedAmount.Add(vaultData.AmountOut)
			addedCollateralData := appData.CollateralLockedAmount.Add(vaultData.AmountIn)
			appData.TokenMintedAmount = addedMintedData
			appData.CollateralLockedAmount = addedCollateralData
			appData.VaultIds = append(appData.VaultIds, vaultData.Id)
		}
	}

	err := k.SetAppExtendedPairVaultMapping(ctx, appExtendedPairVaultData)
	if err != nil {
		return
	}
}

// CalculateCollaterlizationRatio Calculate Collaterlization Ratio .
func (k *Keeper) CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error) {
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
	esmStatus, found := k.GetESMStatus(ctx,extendedPairVault.AppId)
	statusEsm := false
	if found{
		statusEsm = esmStatus.Status
	}
	if statusEsm{
		marketStatus, found := k.GetESMMarketForAsset(ctx, extendedPairVault.AppId)
		if !found {
			return sdk.ZeroDec(), types.ErrorPriceDoesNotExist
		}
		if marketStatus.IsPriceSet{
			for _, data := range marketStatus.Market{
				if assetInData.Id == data.AssetID{
					assetInPrice = data.Rates
				}
			}
		}
	}
	var assetOutPrice uint64

	if extendedPairVault.AssetOutOraclePrice {
		//If oracle Price required for the assetOut
		if statusEsm{
			marketStatus, found := k.GetESMMarketForAsset(ctx, extendedPairVault.AppId)
			if !found {
				return sdk.ZeroDec(), types.ErrorPriceDoesNotExist
			}
			if marketStatus.IsPriceSet{
				for _, data := range marketStatus.Market{
					if assetOutData.Id == data.AssetID{
						assetOutPrice = data.Rates
					}
				}
			}
		}else{
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

func (k *Keeper) VerifyCollaterlizationRatio(
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
	if collaterlizationRatio.LT(minCrRequired) && !statusEsm{
		return types.ErrorInvalidCollateralizationRatio
	}else if collaterlizationRatio.LT(sdk.SmallestDec()) && statusEsm{
		return types.ErrorInvalidCollateralizationRatio
	}
	return nil
}

func (k *Keeper) SetVault(ctx sdk.Context, vault types.Vault) {
	var (
		store = k.Store(ctx)
		key   = types.VaultKey(vault.Id)
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

// UpdateCollateralLockedAmountLockerMapping For updating token stats of collateral .
func (k *Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, valutLookupData types.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool) {
	//if Change type true = Add to collateral Locked
	//If change type false = Subtract from the collateral Locked

	for _, extendedPairData := range valutLookupData.ExtendedPairVaults {
		if extendedPairData.ExtendedPairId == extendedPairID {
			if changeType {
				updatedVal := extendedPairData.CollateralLockedAmount.Add(amount)
				extendedPairData.CollateralLockedAmount = updatedVal
			} else {
				updatedVal := extendedPairData.CollateralLockedAmount.Sub(amount)
				extendedPairData.CollateralLockedAmount = updatedVal
			}
		}
	}
	err := k.SetAppExtendedPairVaultMapping(ctx, valutLookupData)
	if err != nil {
		return
	}
}

// UpdateTokenMintedAmountLockerMapping For updating token stats of minted .
func (k *Keeper) UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, valutLookupData types.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool) {
	//if Change type true = Add to token Locked
	//If change type false = Subtract from the token Locked

	for _, extendedPairData := range valutLookupData.ExtendedPairVaults {
		if extendedPairData.ExtendedPairId == extendedPairID {
			if changeType {
				updatedVal := extendedPairData.TokenMintedAmount.Add(amount)
				extendedPairData.TokenMintedAmount = updatedVal
			} else {
				updatedVal := extendedPairData.TokenMintedAmount.Sub(amount)
				extendedPairData.TokenMintedAmount = updatedVal
			}
		}
	}
	err := k.SetAppExtendedPairVaultMapping(ctx, valutLookupData)
	if err != nil {
		return
	}
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

func (k *Keeper) UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairID uint64, userAddress string, appMappingID uint64) {
	userData, found := k.GetUserVaultExtendedPairMapping(ctx, userAddress)

	var dataIndex int
	if found {
		for _, appData := range userData.UserVaultApp {
			if appData.AppId == appMappingID {
				for index, extendedPairData := range appData.UserExtendedPairVault {
					if extendedPairData.ExtendedPairId == extendedPairID {
						dataIndex = index
					}
				}
				appData.UserExtendedPairVault = append(appData.UserExtendedPairVault[:dataIndex], appData.UserExtendedPairVault[dataIndex+1:]...)
				break
			}
		}
		k.SetUserVaultExtendedPairMapping(ctx, userData)
	}
}

func (k *Keeper) DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID string, appMappingID uint64) {
	appExtendedPairVaultData, found := k.GetAppExtendedPairVaultMapping(ctx, appMappingID)

	var dataIndex int
	if found {
		for _, appData := range appExtendedPairVaultData.ExtendedPairVaults {
			if appData.ExtendedPairId == extendedPairID {
				for index, vaultID := range appData.VaultIds {
					if vaultID == userVaultID {
						dataIndex = index
					}
				}
				appData.VaultIds = append(appData.VaultIds[:dataIndex], appData.VaultIds[dataIndex+1:]...)
			}
		}
		err := k.SetAppExtendedPairVaultMapping(ctx, appExtendedPairVaultData)
		if err != nil {
			return
		}
	}
}

func (k *Keeper) SetStableMintVault(ctx sdk.Context, stableVault types.StableMintVault) {
	var (
		store = k.Store(ctx)
		key   = types.StableMintVaultKey(stableVault.Id)
		value = k.cdc.MustMarshal(&stableVault)
	)

	store.Set(key, value)
}

func (k *Keeper) GetStableMintVault(ctx sdk.Context, id string) (stableVault types.StableMintVault, found bool) {
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

func (k *Keeper) GetStableMintVaults(ctx sdk.Context) (stableVaults []types.StableMintVault) {
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

func (k *Keeper) CreateNewVault(ctx sdk.Context, From string, AppId uint64, ExtendedPairVaultID uint64, AmountIn sdk.Int, AmountOut sdk.Int) error {
	appMapping, _ := k.GetApp(ctx, AppId)
	extendedPairVault, _ := k.GetPairsVault(ctx, ExtendedPairVaultID)
	counterVal, _, _ := k.CheckAppExtendedPairVaultMapping(ctx, appMapping.Id, extendedPairVault.Id)

	zero_val := sdk.ZeroInt()
	var new_vault types.Vault
	updated_counter := counterVal + 1
	new_vault.Id = appMapping.ShortName + strconv.FormatUint(updated_counter, 10)
	new_vault.AmountIn = AmountIn

	new_vault.ClosingFeeAccumulated = zero_val
	new_vault.AmountOut = AmountOut
	new_vault.AppId = appMapping.Id
	new_vault.InterestAccumulated = zero_val
	new_vault.Owner = From
	new_vault.CreatedAt = ctx.BlockTime()
	new_vault.ExtendedPairVaultID = extendedPairVault.Id
	k.SetVault(ctx, new_vault)

	//get extendedpair lookup table data 
	// push the new vault id in that extended pair of the app
	// update the counter of the extendedpair lookup tabe
	////// ///////

	app_extended_pair_vault_data, _ := k.GetAppExtendedPairVaultMapping(ctx, AppId)

	app_extended_pair_vault_data.Counter = updated_counter

	for _, appData := range app_extended_pair_vault_data.ExtendedPairVaults {

  			if appData.ExtendedPairId == new_vault.ExtendedPairVaultID {

    		appData.VaultIds = append(appData.VaultIds, new_vault.Id)
  		}
	}
	
	k.SetAppExtendedPairVaultMapping(ctx, app_extended_pair_vault_data)

	//////////////////
	// k.UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx, updated_counter, new_vault)

	userVaultExtendedpairMappingData, _ := k.GetUserVaultExtendedPairMapping(ctx, From)

	//So only need to add the locker id with asset
	var userExtendedpairData types.ExtendedPairToVaultMapping
	userExtendedpairData.VaultId = new_vault.Id
	userExtendedpairData.ExtendedPairId = new_vault.ExtendedPairVaultID

	for _, appData := range userVaultExtendedpairMappingData.UserVaultApp {
		if appData.AppId == appMapping.Id {

			appData.UserExtendedPairVault = append(appData.UserExtendedPairVault, &userExtendedpairData)
		}

	}
	k.SetUserVaultExtendedPairMapping(ctx, userVaultExtendedpairMappingData)

	return nil
}