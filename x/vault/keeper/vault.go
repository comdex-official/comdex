package keeper

import (

	// assettypes "github.com/comdex-official/comdex/x/asset/types"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	// protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/vault/types"
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

//Checking if for a certain user for the app type , whether there exists a certain asset or not and if it contains a locker id or not
func (k *Keeper) CheckUserAppToExtendedPairMapping(ctx sdk.Context, userVaultAssetData types.UserVaultAssetMapping, extendedPairVaultID uint64, appMappingId uint64) (vault_id string, found bool) {

	for _, vault_app_mapping := range userVaultAssetData.UserVaultApp {

		if vault_app_mapping.AppMappingId == appMappingId {
			for _, extendedPair_to_vaultId_mapping := range vault_app_mapping.UserExtendedPairVault {

				if extendedPair_to_vaultId_mapping.ExtendedPairId == extendedPairVaultID && len(extendedPair_to_vaultId_mapping.VaultId) > 0 {

					vault_id = extendedPair_to_vaultId_mapping.VaultId
					return vault_id, true

				}

			}

		}

	}
	return vault_id, false

}
func (k *Keeper) CheckUserToAppMapping(ctx sdk.Context, userVaultAssetData types.UserVaultAssetMapping, appMappingId uint64) (found bool) {
	for _, vault_app_mapping := range userVaultAssetData.UserVaultApp {

		if vault_app_mapping.AppMappingId == appMappingId {
			return true
		}
	}
	return false

}

//Set AppExtendedPairVaultMapping to check the current status of the vault by extended pair vault id
func (k *Keeper) SetAppExtendedPairVaultMapping(ctx sdk.Context, appExtendedPairVaultData types.AppExtendedPairVaultMapping) {

	var (
		store = k.Store(ctx)
		key   = types.AppExtendedPairVaultMappingKey(appExtendedPairVaultData.AppMappingId)
		value = k.cdc.MustMarshal(&appExtendedPairVaultData)
	)

	store.Set(key, value)

}

//Get AppExtendedPairVaultMapping to check the current status of the vault by extended pair vault id

func (k *Keeper) GetAppExtendedPairVaultMapping(ctx sdk.Context, appMappingId uint64) (appExtendedPairVaultData types.AppExtendedPairVaultMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AppExtendedPairVaultMappingKey(appMappingId)
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

func (k *Keeper) CheckAppExtendedPairVaultMapping(ctx sdk.Context, appMappingId uint64, extendedPairVaultId uint64) (counter uint64, minted_stastics sdk.Int, lenVaults uint64) {

	app_extended_pair_vault_data, found := k.GetAppExtendedPairVaultMapping(ctx, appMappingId)
	if !found {

		//Initialising a new struct
		var newAppExtendedPairVault types.AppExtendedPairVaultMapping
		var newExtendedPairVault types.ExtendedPairVaultMapping
		newAppExtendedPairVault.AppMappingId = appMappingId
		newAppExtendedPairVault.Counter = 0
		zero_val := sdk.ZeroInt()
		newExtendedPairVault.ExtendedPairId = extendedPairVaultId
		newExtendedPairVault.CollateralLockedAmount = &zero_val
		newExtendedPairVault.TokenMintedAmount = &zero_val
		newAppExtendedPairVault.ExtendedPairVaults = append(newAppExtendedPairVault.ExtendedPairVaults, &newExtendedPairVault)
		k.SetAppExtendedPairVaultMapping(ctx, newAppExtendedPairVault)

		return newAppExtendedPairVault.Counter, *newExtendedPairVault.TokenMintedAmount, 0

	} else {

		for _, extendedPairVaultData := range app_extended_pair_vault_data.ExtendedPairVaults {

			if extendedPairVaultData.ExtendedPairId == extendedPairVaultId {

				lenOfVaults := len(app_extended_pair_vault_data.ExtendedPairVaults)

				return app_extended_pair_vault_data.Counter, *extendedPairVaultData.TokenMintedAmount, uint64(lenOfVaults)
			}

		}
		//Check the Zero Value once
		zero_val := sdk.ZeroInt()
		var newExtendedPairVault types.ExtendedPairVaultMapping
		newExtendedPairVault.ExtendedPairId = extendedPairVaultId
		newExtendedPairVault.CollateralLockedAmount = &zero_val
		newExtendedPairVault.TokenMintedAmount = &zero_val
		app_extended_pair_vault_data.ExtendedPairVaults = append(app_extended_pair_vault_data.ExtendedPairVaults, &newExtendedPairVault)
		k.SetAppExtendedPairVaultMapping(ctx, app_extended_pair_vault_data)

		return app_extended_pair_vault_data.Counter, *newExtendedPairVault.TokenMintedAmount, 0

	}

}

func (k *Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, counter uint64, vaultData types.Vault) {

	app_extended_pair_vault_data, _ := k.GetAppExtendedPairVaultMapping(ctx, vaultData.AppMappingId)

	app_extended_pair_vault_data.Counter = counter

	for _, appData := range app_extended_pair_vault_data.ExtendedPairVaults {

		if appData.ExtendedPairId == vaultData.ExtendedPairVaultID {

			addedMintedData := appData.TokenMintedAmount.Add(vaultData.AmountOut)
			addedCollateralData := appData.CollateralLockedAmount.Add(vaultData.AmountIn)
			appData.TokenMintedAmount = &addedMintedData
			appData.CollateralLockedAmount = &addedCollateralData
			appData.VaultIds = append(appData.VaultIds, vaultData.Id)

		}

	}
	k.SetAppExtendedPairVaultMapping(ctx, app_extended_pair_vault_data)

}

func (k *Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreateStableMintVault(ctx sdk.Context, counter uint64, vaultData types.StableMintVault) {

	app_extended_pair_vault_data, _ := k.GetAppExtendedPairVaultMapping(ctx, vaultData.AppMappingId)

	app_extended_pair_vault_data.Counter = counter

	for _, appData := range app_extended_pair_vault_data.ExtendedPairVaults {

		if appData.ExtendedPairId == vaultData.ExtendedPairVaultID {

			addedMintedData := appData.TokenMintedAmount.Add(vaultData.AmountOut)
			addedCollateralData := appData.CollateralLockedAmount.Add(vaultData.AmountIn)
			appData.TokenMintedAmount = &addedMintedData
			appData.CollateralLockedAmount = &addedCollateralData
			appData.VaultIds = append(appData.VaultIds, vaultData.Id)

		}

	}
	k.SetAppExtendedPairVaultMapping(ctx, app_extended_pair_vault_data)

}

//Calculate Collaterlization Ratio
func (k *Keeper) CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultId uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error) {

	extended_pair_vault, found := k.GetPairsVault(ctx, extendedPairVaultId)
	if !found {
		return sdk.ZeroDec(), types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.GetPair(ctx, extended_pair_vault.PairId)
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
	fmt.Println(assetInPrice, "price of asset ")
	var assetOutPrice uint64

	if extended_pair_vault.AssetOutOraclePrice {
		fmt.Println(extended_pair_vault.AssetOutOraclePrice, "value bool price required")
		//If oracle Price required for the assetOut
		assetOutPrice, found = k.GetPriceForAsset(ctx, assetOutData.Id)
		fmt.Println(assetOutPrice, "should be what is set dollar ")

		if !found {
			return sdk.ZeroDec(), types.ErrorPriceDoesNotExist
		}
	} else {
		//If oracle Price is not required for the assetOut
		assetOutPrice = extended_pair_vault.AsssetOutPrice.BigInt().Uint64()

	}

	totalIn := amountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).ToDec()
	if totalIn.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), types.ErrorInvalidAmountIn
	}

	totalOut := amountOut.Mul(sdk.NewIntFromUint64(assetOutPrice)).ToDec()
	if totalOut.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), types.ErrorInvalidAmountOut
	}
	fmt.Println(amountIn, "amountIn")
	fmt.Println(amountOut, "amountout")
	fmt.Println(totalIn.Quo(totalOut))
	fmt.Println(totalIn)
	fmt.Println(totalOut)

	return totalIn.Quo(totalOut), nil

}

func (k *Keeper) VerifyCollaterlizationRatio(
	ctx sdk.Context,
	extendedPairVaultId uint64,
	amountIn sdk.Int,
	amountOut sdk.Int,
	minCrRequired sdk.Dec,
) error {
	collaterlizationRatio, err := k.CalculateCollaterlizationRatio(ctx, extendedPairVaultId, amountIn, amountOut)
	if err != nil {
		return err
	}

	if collaterlizationRatio.LT(minCrRequired) {
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

//For updating token stats of collateral
func (k *Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, valutLookupData types.AppExtendedPairVaultMapping, extendedPairId uint64, amount sdk.Int, changeType bool) {

	//if Change type true = Add to collateral Locked
	//If change type false = Substract from the collateral Locked

	for _, extendedPairData := range valutLookupData.ExtendedPairVaults {
		if extendedPairData.ExtendedPairId == extendedPairId {
			if changeType {
				updatedVal := extendedPairData.CollateralLockedAmount.Add(amount)
				extendedPairData.CollateralLockedAmount = &updatedVal
			} else {
				updatedVal := extendedPairData.CollateralLockedAmount.Sub(amount)
				extendedPairData.CollateralLockedAmount = &updatedVal
			}
		}
	}
	k.SetAppExtendedPairVaultMapping(ctx, valutLookupData)

}

//For updating token stats of minted
func (k *Keeper) UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, valutLookupData types.AppExtendedPairVaultMapping, extendedPairId uint64, amount sdk.Int, changeType bool) {

	//if Change type true = Add to token Locked
	//If change type false = Substract from the token Locked

	for _, extendedPairData := range valutLookupData.ExtendedPairVaults {
		if extendedPairData.ExtendedPairId == extendedPairId {
			if changeType {
				updatedVal := extendedPairData.TokenMintedAmount.Add(amount)
				extendedPairData.TokenMintedAmount = &updatedVal
			} else {
				updatedVal := extendedPairData.TokenMintedAmount.Sub(amount)
				extendedPairData.TokenMintedAmount = &updatedVal
			}
		}
	}
	k.SetAppExtendedPairVaultMapping(ctx, valutLookupData)

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

func (k *Keeper) UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairId uint64, userAddress string, appMappingId uint64) {

	userData, _ := k.GetUserVaultExtendedPairMapping(ctx, userAddress)

	var dataIndex int
	for _, appData := range userData.UserVaultApp {

		if appData.AppMappingId == appMappingId {

			for index, extendedPairData := range appData.UserExtendedPairVault {

				if extendedPairData.ExtendedPairId == extendedPairId {

					dataIndex = index
				}
			}
			a := appData.UserExtendedPairVault[1:dataIndex]
			b := appData.UserExtendedPairVault[dataIndex+1:]
			a = append(a, b...)
			appData.UserExtendedPairVault = a
			break
		}
	}

	k.SetUserVaultExtendedPairMapping(ctx, userData)

}

func (k *Keeper) DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairId uint64, userVaultId string, appMappingId uint64) {

	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMapping(ctx, appMappingId)

	var dataIndex int

	for _, appData := range appExtendedPairVaultData.ExtendedPairVaults {

		if appData.ExtendedPairId == extendedPairId {

			for index, vaultId := range appData.VaultIds {

				if vaultId == userVaultId {
					dataIndex = index

				}

			}
			a := appData.VaultIds[1:dataIndex]
			b := appData.VaultIds[dataIndex+1:]
			a = append(a, b...)
			appData.VaultIds = a

		}

	}
	k.SetAppExtendedPairVaultMapping(ctx, appExtendedPairVaultData)

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

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var stableVault types.StableMintVault
		k.cdc.MustUnmarshal(iter.Value(), &stableVault)
		stableVaults = append(stableVaults, stableVault)
	}

	return stableVaults
}
