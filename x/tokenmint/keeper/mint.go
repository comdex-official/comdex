package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/tokenmint/types"
)

func (k Keeper) SetTokenMint(ctx sdk.Context, appTokenMintData types.TokenMint) {
	var (
		store = k.Store(ctx)
		key   = types.TokenMintKey(appTokenMintData.AppId)
		value = k.cdc.MustMarshal(&appTokenMintData)
	)

	store.Set(key, value)
}

func (k Keeper) GetTokenMint(ctx sdk.Context, appMappingID uint64) (appTokenMintData types.TokenMint, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.TokenMintKey(appMappingID)
		value = store.Get(key)
	)

	if value == nil {
		return appTokenMintData, false
	}

	k.cdc.MustUnmarshal(value, &appTokenMintData)
	return appTokenMintData, true
}

func (k Keeper) GetTotalTokenMinted(ctx sdk.Context) (appTokenMintData []types.TokenMint) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.TokenMintKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var totalMinted types.TokenMint
		k.cdc.MustUnmarshal(iter.Value(), &totalMinted)
		appTokenMintData = append(appTokenMintData, totalMinted)
	}

	return appTokenMintData
}

func (k Keeper) GetAssetDataInTokenMintByApp(ctx sdk.Context, appMappingID uint64, assetID uint64) (tokenData types.MintedTokens, found bool) {
	mintData, found := k.GetTokenMint(ctx, appMappingID)
	if !found {
		return tokenData, false
	}

	for _, mintAssetData := range mintData.MintedTokens {
		if mintAssetData.AssetId == assetID {
			tokenData = *mintAssetData
			return tokenData, true
		}
	}
	return tokenData, false
}

func (k Keeper) GetAssetDataInTokenMintByAppSupply(ctx sdk.Context, appMappingID uint64, assetID uint64) (tokenDataSupply int64, found bool) {
	tokenData, found := k.GetAssetDataInTokenMintByApp(ctx, appMappingID, assetID)

	if !found {
		return 0, false
	}
	return tokenData.CurrentSupply.Int64(), found
}

func (k Keeper) MintNewTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, address string, amount sdk.Int) error {
	assetData, found := k.GetAsset(ctx, assetID)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	_, found = k.GetTokenMint(ctx, appMappingID)
	if !found {
		return types.ErrorAppMappingDoesNotExists
	}

	_, found = k.GetAssetDataInTokenMintByApp(ctx, appMappingID, assetID)
	if !found {
		return types.ErrorAssetNotWhiteListedForGenesisMinting
	}
	if err := k.MintCoin(ctx, types.ModuleName, sdk.NewCoin(assetData.Denom, amount)); err != nil {
		return err
	}
	userAddress, err := sdk.AccAddressFromBech32(address)

	if err != nil {
		return err
	}
	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, userAddress, sdk.NewCoin(assetData.Denom, amount)); err != nil {
		return err
	}
	k.UpdateAssetDataInTokenMintByApp(ctx, appMappingID, assetID, true, amount)

	return nil
}

func (k Keeper) BurnTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, amount sdk.Int) error {
	assetData, found := k.GetAsset(ctx, assetID)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	_, found = k.GetTokenMint(ctx, appMappingID)
	if !found {
		return types.ErrorAppMappingDoesNotExists
	}

	tokenData, found := k.GetAssetDataInTokenMintByApp(ctx, appMappingID, assetID)
	if !found {
		return types.ErrorAssetNotWhiteListedForGenesisMinting
	}
	if tokenData.CurrentSupply.Sub(amount).LTE(sdk.NewInt(0)) {
		return types.ErrorBurningMakesSupplyLessThanZero
	}
	if err := k.BurnCoin(ctx, types.ModuleName, sdk.NewCoin(assetData.Denom, amount)); err != nil {
		return err
	}
	// if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(address), sdk.NewCoin(assetData.Denom, amount)); err != nil {
	// 	return err
	// }
	k.UpdateAssetDataInTokenMintByApp(ctx, appMappingID, assetID, false, amount)

	return nil
}

func (k Keeper) UpdateAssetDataInTokenMintByApp(ctx sdk.Context, appMappingID uint64, assetID uint64, changeType bool, amount sdk.Int) {
	//ChangeType + == add to current supply
	//Change type - == reduce from supply
	mintData, found := k.GetTokenMint(ctx, appMappingID)

	if found {
		for _, mintAssetData := range mintData.MintedTokens {
			if mintAssetData.AssetId == assetID {
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
}

func (k Keeper) BurnGovTokensForApp(ctx sdk.Context, appMappingID uint64, from sdk.AccAddress, amount sdk.Coin) error {
	_, found := k.GetApp(ctx, appMappingID)
	if !found {
		return types.ErrorAppMappingDoesNotExists
	}

	err := k.BurnFrom(ctx, amount, from)
	if err != nil {
		return err
	}
	asset, _ := k.GetAssetForDenom(ctx, amount.Denom)

	k.UpdateAssetDataInTokenMintByApp(ctx, appMappingID, asset.Id, false, amount.Amount)

	return nil
}

func (k Keeper) BurnFrom(ctx sdk.Context, amount sdk.Coin, burnFrom sdk.AccAddress) error {
	err := k.SendCoinFromAccountToModule(ctx,
		burnFrom,
		types.ModuleName,
		amount)
	if err != nil {
		return err
	}
	return k.BurnCoin(ctx, types.ModuleName, amount)
}
