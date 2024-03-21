package keeper

import (
	"sort"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/cosmos/gogoproto/types"

	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	"github.com/comdex-official/comdex/x/vault/types"
)

func (k Keeper) SetLengthOfVault(ctx sdk.Context, length uint64) {
	var (
		store = k.Store(ctx)
		key   = types.VaultLengthPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: length,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetLengthOfVault(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.VaultLengthPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

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

// SetAppExtendedPairVaultMappingData Set AppExtendedPairVaultMapping to check the current status of the vault by extended pair vault id .
func (k Keeper) SetAppExtendedPairVaultMappingData(ctx sdk.Context, appExtendedPairVaultData types.AppExtendedPairVaultMappingData) {
	var (
		store = k.Store(ctx)
		key   = types.AppExtendedPairVaultMappingKey(appExtendedPairVaultData.AppId, appExtendedPairVaultData.ExtendedPairId)
		value = k.cdc.MustMarshal(&appExtendedPairVaultData)
	)

	store.Set(key, value)
}

// Get AppExtendedPairVaultMapping to check the current status of the vault by extended pair vault id
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

// Check AppExtendedPairVault Data,
// If exists fine --- go with the next steps from here
// else instantiate 1 and set it. and go for the next steps from here
// So best way will be to create a function which will first check if AppExtendedPairVault Data exists or not. If it does. then send counted value. else create a struct save it. and send counter value.
func (k Keeper) CheckAppExtendedPairVaultMapping(ctx sdk.Context, appMappingID uint64, extendedPairVaultID uint64) (mintedStatistics sdk.Int, lenVaults uint64) {
	appExtendedPairVaultData, found := k.GetAppExtendedPairVaultMappingData(ctx, appMappingID, extendedPairVaultID)
	if !found {
		// Initialising a new struct
		var newAppExtendedPairVault types.AppExtendedPairVaultMappingData
		newAppExtendedPairVault.AppId = appMappingID
		zeroVal := sdk.ZeroInt()
		newAppExtendedPairVault.ExtendedPairId = extendedPairVaultID
		newAppExtendedPairVault.CollateralLockedAmount = zeroVal
		newAppExtendedPairVault.TokenMintedAmount = zeroVal

		k.SetAppExtendedPairVaultMappingData(ctx, newAppExtendedPairVault)

		return newAppExtendedPairVault.TokenMintedAmount, 0
	}
	return appExtendedPairVaultData.TokenMintedAmount, uint64(len(appExtendedPairVaultData.VaultIds))
}

func (k Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, vaultData types.Vault) {
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, vaultData.AppId, vaultData.ExtendedPairVaultID)

	addedMintedData := appExtendedPairVaultData.TokenMintedAmount.Add(vaultData.AmountOut)
	addedCollateralData := appExtendedPairVaultData.CollateralLockedAmount.Add(vaultData.AmountIn)
	appExtendedPairVaultData.TokenMintedAmount = addedMintedData
	appExtendedPairVaultData.CollateralLockedAmount = addedCollateralData
	appExtendedPairVaultData.VaultIds = append(appExtendedPairVaultData.VaultIds, vaultData.Id)

	k.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
}

func (k Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreateStableMintVault(ctx sdk.Context, vaultData types.StableMintVault) {
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, vaultData.AppId, vaultData.ExtendedPairVaultID)
	addedMintedData := appExtendedPairVaultData.TokenMintedAmount.Add(vaultData.AmountOut)
	addedCollateralData := appExtendedPairVaultData.CollateralLockedAmount.Add(vaultData.AmountIn)
	appExtendedPairVaultData.TokenMintedAmount = addedMintedData
	appExtendedPairVaultData.CollateralLockedAmount = addedCollateralData
	appExtendedPairVaultData.VaultIds = append(appExtendedPairVaultData.VaultIds, vaultData.Id)

	k.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
}

// CalculateCollateralizationRatio Calculate Collaterlization Ratio .
func (k Keeper) CalculateCollateralizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error) {
	extendedPairVault, found := k.asset.GetPairsVault(ctx, extendedPairVaultID)
	if !found {
		return sdk.ZeroDec(), types.ErrorExtendedPairVaultDoesNotExists
	}
	pairData, found := k.asset.GetPair(ctx, extendedPairVault.PairId)
	if !found {
		return sdk.ZeroDec(), types.ErrorPairDoesNotExist
	}
	assetInData, found := k.asset.GetAsset(ctx, pairData.AssetIn)
	if !found {
		return sdk.ZeroDec(), types.ErrorAssetDoesNotExist
	}
	assetOutData, found := k.asset.GetAsset(ctx, pairData.AssetOut)
	if !found {
		return sdk.ZeroDec(), types.ErrorAssetDoesNotExist
	}
	esmStatus, found := k.esm.GetESMStatus(ctx, extendedPairVault.AppId)
	statusEsm := false
	if found {
		statusEsm = esmStatus.Status
	}
	var assetInTotalPrice sdk.Dec
	var err error
	// check to get calc asset price from esm
	if statusEsm && esmStatus.SnapshotStatus {
		price, found := k.esm.GetSnapshotOfPrices(ctx, extendedPairVault.AppId, assetInData.Id)
		if !found {
			return sdk.ZeroDec(), types.ErrorPriceDoesNotExist
		}
		numerator := sdk.NewDecFromInt(amountIn).Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(price)))
		denominator := sdk.NewDecFromInt(assetInData.Decimals)
		assetInTotalPrice = numerator.Quo(denominator)
	} else if !statusEsm {
		// calculating price of the asset_in
		assetInTotalPrice, err = k.oracle.CalcAssetPrice(ctx, assetInData.Id, amountIn)
		if err != nil {
			return sdk.ZeroDec(), err
		}
	}
	var assetOutTotalPrice sdk.Dec

	if extendedPairVault.AssetOutOraclePrice {
		// If oracle Price required for the assetOut
		if statusEsm && esmStatus.SnapshotStatus {
			price, found := k.esm.GetSnapshotOfPrices(ctx, extendedPairVault.AppId, assetOutData.Id)
			if !found {
				return sdk.ZeroDec(), types.ErrorPriceDoesNotExist
			}
			numerator := sdk.NewDecFromInt(amountOut).Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(price)))
			denominator := sdk.NewDecFromInt(assetOutData.Decimals)
			assetOutTotalPrice = numerator.Quo(denominator)
		} else {
			assetOutTotalPrice, err = k.oracle.CalcAssetPrice(ctx, assetOutData.Id, amountOut)
			if err != nil {
				return sdk.ZeroDec(), err
			}
		}
	} else {
		// If oracle Price is not required for the assetOut
		numerator := sdk.NewDecFromInt(amountOut).Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(extendedPairVault.AssetOutPrice)))
		denominator := sdk.NewDecFromInt(assetOutData.Decimals)
		assetOutTotalPrice = numerator.Quo(denominator)
	}

	if assetInTotalPrice.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), types.ErrorInvalidAmountIn
	}

	if assetOutTotalPrice.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), types.ErrorInvalidAmountOut
	}
	return assetInTotalPrice.Quo(assetOutTotalPrice), nil
}

func (k Keeper) VerifyCollaterlizationRatio(
	ctx sdk.Context,
	extendedPairVaultID uint64,
	amountIn sdk.Int,
	amountOut sdk.Int,
	minCrRequired sdk.Dec,
	statusEsm bool,
) error {
	collaterlizationRatio, err := k.CalculateCollateralizationRatio(ctx, extendedPairVaultID, amountIn, amountOut)
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
func (k Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdk.Int, changeType bool) {
	// if Change type true = Add to collateral Locked
	// If change type false = Subtract from the collateral Locked
	appExtendedPairVaultData, found := k.GetAppExtendedPairVaultMappingData(ctx, appMappingID, extendedPairID)
	if found {
		if changeType {
			updatedVal := appExtendedPairVaultData.CollateralLockedAmount.Add(amount)
			appExtendedPairVaultData.CollateralLockedAmount = updatedVal
		} else {
			updatedVal := appExtendedPairVaultData.CollateralLockedAmount.Sub(amount)
			appExtendedPairVaultData.CollateralLockedAmount = updatedVal
		}

		k.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
	}
}

// UpdateTokenMintedAmountLockerMapping For updating token stats of minted .
func (k Keeper) UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdk.Int, changeType bool) {
	// if Change type true = Add to token Locked
	// If change type false = Subtract from the token Locked
	appExtendedPairVaultData, found := k.GetAppExtendedPairVaultMappingData(ctx, appMappingID, extendedPairID)
	if found {
		if changeType {
			updatedVal := appExtendedPairVaultData.TokenMintedAmount.Add(amount)
			appExtendedPairVaultData.TokenMintedAmount = updatedVal
		} else {
			updatedVal := appExtendedPairVaultData.TokenMintedAmount.Sub(amount)
			appExtendedPairVaultData.TokenMintedAmount = updatedVal
		}
		k.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
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

func (k Keeper) DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID uint64, appMappingID uint64) {
	appExtendedPairVaultData, found := k.GetAppExtendedPairVaultMappingData(ctx, appMappingID, extendedPairID)

	if found {
		lengthOfVaults := len(appExtendedPairVaultData.VaultIds)

		dataIndex := sort.Search(lengthOfVaults, func(i int) bool { return appExtendedPairVaultData.VaultIds[i] >= userVaultID })

		if dataIndex < lengthOfVaults && appExtendedPairVaultData.VaultIds[dataIndex] == userVaultID {
			appExtendedPairVaultData.VaultIds = append(appExtendedPairVaultData.VaultIds[:dataIndex], appExtendedPairVaultData.VaultIds[dataIndex+1:]...)
			k.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
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
	appMapping, _ := k.asset.GetApp(ctx, AppID)
	extendedPairVault, _ := k.asset.GetPairsVault(ctx, ExtendedPairVaultID)
	//checking if vault exists
	userData, userExists := k.GetUserAppExtendedPairMappingData(ctx, From, AppID, ExtendedPairVaultID)
	if userExists {

		vaultData, _ := k.GetVault(ctx, userData.VaultId)
		vaultData.AmountIn = vaultData.AmountIn.Add(AmountIn)
		vaultData.AmountOut = vaultData.AmountOut.Add(AmountOut)
		k.SetVault(ctx, vaultData)

		return nil
	}

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
	length := k.GetLengthOfVault(ctx)
	k.SetLengthOfVault(ctx, length+1)

	var mappingData types.OwnerAppExtendedPairVaultMappingData
	mappingData.Owner = From
	mappingData.AppId = AppID
	mappingData.ExtendedPairId = ExtendedPairVaultID
	mappingData.VaultId = newVault.Id
	appExtendedPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, AppID, ExtendedPairVaultID)
	appExtendedPairVaultData.VaultIds = append(appExtendedPairVaultData.VaultIds, newVault.Id)
	k.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
	k.SetUserAppExtendedPairMappingData(ctx, mappingData)

	return nil
}

func (k Keeper) calculateUserToken(userVault types.Vault, amountIn sdk.Int) (userToken sdk.Int) {
	nume := userVault.AmountOut.Mul(amountIn)
	deno := userVault.AmountIn
	userToken = nume.Quo(deno)

	return userToken
}

func (k Keeper) WasmMsgAddEmissionRewards(ctx sdk.Context, appID uint64, amount sdk.Int, extPair []uint64, votingRatio []sdk.Int) error {
	var assetID uint64
	var vaultsData types.AppExtendedPairVaultMappingData

	totalVote := sdk.ZeroInt()
	app, _ := k.asset.GetApp(ctx, appID)
	govToken := app.GenesisToken
	for _, v := range govToken {
		if v.IsGovToken {
			assetID = v.AssetId
		}
	}
	asset, _ := k.asset.GetAsset(ctx, assetID)
	if amount.GT(sdk.ZeroInt()) {
		err := k.bank.MintCoins(ctx, tokenminttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Denom, amount)))
		if err != nil {
			return err
		}
	}
	k.tokenmint.UpdateAssetDataInTokenMintByApp(ctx, appID, assetID, true, amount)
	for i := range votingRatio {
		totalVote = totalVote.Add(votingRatio[i])
	}
	for j, extP := range extPair {
		extPairVaultMappingData, found := k.GetAppExtendedPairVaultMappingData(ctx, appID, extP)
		individualVote := votingRatio[j]
		votingR := sdk.NewDecFromInt(individualVote).Quo(sdk.NewDecFromInt(totalVote))
		shareByExtPair := votingR.Mul(sdk.NewDecFromInt(amount))
		if !found || extPairVaultMappingData.TokenMintedAmount.IsZero() {
			continue
		}
		perUserShareByAmtDec := shareByExtPair.Quo(sdk.NewDecFromInt(extPairVaultMappingData.TokenMintedAmount))
		vaultsData, _ = k.GetAppExtendedPairVaultMappingData(ctx, appID, extP)

		for _, vaultID := range vaultsData.VaultIds {
			vault, _ := k.GetVault(ctx, vaultID)
			amt := sdk.NewDecFromInt(vault.AmountOut).Mul(perUserShareByAmtDec)
			addr, _ := sdk.AccAddressFromBech32(vault.Owner)
			if amt.GT(sdk.NewDec(0)) {
				err := k.bank.SendCoinsFromModuleToAccount(ctx, tokenminttypes.ModuleName, addr, sdk.NewCoins(sdk.NewCoin(asset.Denom, amt.TruncateInt())))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (k Keeper) GetAmountOfOtherToken(ctx sdk.Context, id1 uint64, rate1 sdk.Dec, amt1 sdk.Int, id2 uint64, rate2 sdk.Dec) (sdk.Dec, sdk.Int, error) {
	asset1, found := k.asset.GetAsset(ctx, id1)
	if !found {
		return sdk.ZeroDec(), sdk.ZeroInt(), assettypes.ErrorAssetDoesNotExist
	}
	asset2, found := k.asset.GetAsset(ctx, id2)
	if !found {
		return sdk.ZeroDec(), sdk.ZeroInt(), assettypes.ErrorAssetDoesNotExist
	}

	numerator := sdk.NewDecFromInt(amt1).Mul(rate1) //rate urate 1000000
	denominator := sdk.NewDecFromInt(asset1.Decimals)
	t1dAmount := numerator.Quo(denominator)

	newAmount := t1dAmount.Quo(rate2)
	tokenAmount := newAmount.Mul(sdk.NewDecFromInt(asset2.Decimals))
	// return sdk.Int(tokenAmount), nil
	return t1dAmount, tokenAmount.TruncateInt(), nil
}

func (k Keeper) SetStableMintVaultRewards(ctx sdk.Context, stableMintVaultRewards types.StableMintVaultRewards) {
	var (
		store = k.Store(ctx)
		key   = types.StableMintVaultRewardsKey(stableMintVaultRewards.AppId, stableMintVaultRewards.User, stableMintVaultRewards.BlockHeight)
		value = k.cdc.MustMarshal(&stableMintVaultRewards)
	)

	store.Set(key, value)
}

func (k Keeper) GetStableMintVaultRewards(ctx sdk.Context, stableMintVaultRewards types.StableMintVaultRewards) (mappingData types.StableMintVaultRewards, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.StableMintVaultRewardsKey(stableMintVaultRewards.AppId, stableMintVaultRewards.User, stableMintVaultRewards.BlockHeight)
		value = store.Get(key)
	)

	if value == nil {
		return mappingData, false
	}

	k.cdc.MustUnmarshal(value, &mappingData)
	return mappingData, true
}

func (k Keeper) DeleteStableMintVaultRewards(ctx sdk.Context, stableMintVaultRewards types.StableMintVaultRewards) {
	var (
		store = k.Store(ctx)
		key   = types.StableMintVaultRewardsKey(stableMintVaultRewards.AppId, stableMintVaultRewards.User, stableMintVaultRewards.BlockHeight)
	)

	store.Delete(key)
}

func (k Keeper) GetStableMintVaultUserRewards(ctx sdk.Context, appID uint64, user string) (mappingData []types.StableMintVaultRewards, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.StableMintRewardsKey(appID, user)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.StableMintVaultRewards
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		mappingData = append(mappingData, mapData)
	}
	if mappingData == nil {
		return nil, false
	}

	return mappingData, true
}

func (k Keeper) GetStableMintVaultRewardsByApp(ctx sdk.Context, appID uint64) (mappingData []types.StableMintVaultRewards, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.StableMintRewardsAppKey(appID)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.StableMintVaultRewards
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		mappingData = append(mappingData, mapData)
	}
	if mappingData == nil {
		return nil, false
	}

	return mappingData, true
}

func (k Keeper) DeleteUserStableRewardEntries(ctx sdk.Context, appID uint64, user string, quanitity sdk.Int) {
	stableVaultRewards, found := k.GetStableMintVaultUserRewards(ctx, appID, user)
	if found {
		for _, userRewards := range stableVaultRewards {
			if userRewards.Amount.GTE(quanitity) {
				userRewards.Amount = userRewards.Amount.Sub(quanitity)
				k.SetStableMintVaultRewards(ctx, userRewards)
				break
			} else if userRewards.Amount.LT(quanitity) {
				quanitity = quanitity.Sub(userRewards.Amount)
				k.DeleteStableMintVaultRewards(ctx, userRewards)
			}
		}
	}
}

func (k Keeper) GetStableMintVaultRewardsOfAllApps(ctx sdk.Context) (mappingData []types.StableMintVaultRewards) {
	var (
		store = k.Store(ctx)
		key   = types.StableVaultRewardsKeyPrefix
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.StableMintVaultRewards
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		mappingData = append(mappingData, mapData)
	}
	if mappingData == nil {
		return nil
	}

	return mappingData
}

func (k Keeper) SetWithdrawStableMintControl(ctx sdk.Context, control bool) {
	var (
		store = k.Store(ctx)
		key   = types.StableVaultControlKeyPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.BoolValue{
				Value: control,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetWithdrawStableMintControl(ctx sdk.Context) bool {
	var (
		store = k.Store(ctx)
		key   = types.StableVaultControlKeyPrefix
		value = store.Get(key)
	)

	if value == nil {
		return false
	}

	var id protobuftypes.BoolValue
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}