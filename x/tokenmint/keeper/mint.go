package keeper

import (

	// assettypes "github.com/comdex-official/comdex/x/asset/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/tokenmint/types"
)

func (k *Keeper) SetTokenMint(ctx sdk.Context, appTokenMintData types.TokenMint) {

	var (
		store = k.Store(ctx)
		key   = types.TokenMintKey(appTokenMintData.AppMappingId)
		value = k.cdc.MustMarshal(&appTokenMintData)
	)

	store.Set(key, value)

}

func (k *Keeper) GetTokenMint(ctx sdk.Context, appMappingId uint64) (appTokenMintData types.TokenMint, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.TokenMintKey(appMappingId)
		value = store.Get(key)
	)

	if value == nil {
		return appTokenMintData, false
	}

	k.cdc.MustUnmarshal(value, &appTokenMintData)
	return appTokenMintData, true
}

func (k *Keeper) GetTotalTokenMinted(ctx sdk.Context) (appTokenMintData []types.TokenMint) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.TokenMintKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var totalMinted types.TokenMint
		k.cdc.MustUnmarshal(iter.Value(), &totalMinted)
		appTokenMintData = append(appTokenMintData, totalMinted)
	}

	return appTokenMintData
}

func (k *Keeper) GetAssetDataInTokenMintByApp(ctx sdk.Context, appMappingId uint64, assetId uint64) (tokenData types.MintedTokens, found bool) {

	mintData, found := k.GetTokenMint(ctx, appMappingId)

	for _, mintAssetData := range mintData.MintedTokens {

		if mintAssetData.AssetId == assetId {
			tokenData = *mintAssetData
			return tokenData, true
		}

	}
	return tokenData, false
}

func (k *Keeper) GetAssetDataInTokenMintByAppSupply(ctx sdk.Context, appMappingId uint64, assetId uint64) (tokenDataSupply int64, found bool) {
	tokenData, found := k.GetAssetDataInTokenMintByApp(ctx, appMappingId, assetId)
	return tokenData.CurrentSupply.Int64(), found
}

func (k *Keeper) MintNewTokensForApp(ctx sdk.Context, appMappingId uint64, assetId uint64, address string, amount sdk.Int) error {

	assetData, found := k.GetAsset(ctx, assetId)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	// appMappingData, found := k.GetApp(ctx, appMappingId)
	// if !found {
	// 	return types.ErrorAppMappingDoesNotExists
	// }
	// //Checking if asset exists in the app

	// _, found = k.GetMintGenesisTokenData(ctx, appMappingData.Id, assetData.Id)
	// if !found {
	// 	return types.ErrorAssetNotWhiteListedForGenesisMinting
	// }
	_, found = k.GetTokenMint(ctx, appMappingId)
	if !found {
		return types.ErrorAppMappingDoesNotExists
	}

	_, found = k.GetAssetDataInTokenMintByApp(ctx, appMappingId, assetId)
	if !found {
		return types.ErrorAssetNotWhiteListedForGenesisMinting

	}
	if err := k.MintCoin(ctx, types.ModuleName, sdk.NewCoin(assetData.Denom, amount)); err != nil {
		return err
	}
	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(address), sdk.NewCoin(assetData.Denom, amount)); err != nil {
		return err
	}
	k.UpdateAssetDataInTokenMintByApp(ctx, appMappingId, assetId, true, amount)

	return nil

}

func (k *Keeper) BurnTokensForApp(ctx sdk.Context, appMappingId uint64, assetId uint64, amount sdk.Int) error {

	assetData, found := k.GetAsset(ctx, assetId)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	// appMappingData, found := k.GetApp(ctx, appMappingId)
	// if !found {
	// 	return types.ErrorAppMappingDoesNotExists
	// }
	// //Checking if asset exists in the app

	// _, found = k.GetMintGenesisTokenData(ctx, appMappingData.Id, assetData.Id)
	// if !found {
	// 	return types.ErrorAssetNotWhiteListedForGenesisMinting
	// }
	_, found = k.GetTokenMint(ctx, appMappingId)
	if !found {
		return types.ErrorAppMappingDoesNotExists
	}

	tokenData, found := k.GetAssetDataInTokenMintByApp(ctx, appMappingId, assetId)
	if !found {
		return types.ErrorAssetNotWhiteListedForGenesisMinting

	}
	if tokenData.CurrentSupply.Sub(amount).LTE(sdk.NewInt(0)) {
		return types.ErrorBuringMakesSupplyLessThanZero

	}
	if err := k.BurnCoin(ctx, types.ModuleName, sdk.NewCoin(assetData.Denom, amount)); err != nil {
		return err
	}
	// if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(address), sdk.NewCoin(assetData.Denom, amount)); err != nil {
	// 	return err
	// }
	k.UpdateAssetDataInTokenMintByApp(ctx, appMappingId, assetId, false, amount)

	return nil

}

func (k *Keeper) UpdateAssetDataInTokenMintByApp(ctx sdk.Context, appMappingId uint64, assetId uint64, changeType bool, amount sdk.Int) {

	//ChangeType + == add to current supply
	//Change type - == reduce from supply
	mintData, _ := k.GetTokenMint(ctx, appMappingId)

	for _, mintAssetData := range mintData.MintedTokens {

		if mintAssetData.AssetId == assetId {

			if changeType {
				mintAssetData.CurrentSupply = mintAssetData.CurrentSupply.Add(amount)
				break

			} else {
				mintAssetData.CurrentSupply = mintAssetData.CurrentSupply.Sub(amount)
				break
			}

		}

	}
	k.SetTokenMint(ctx, mintData)

}
