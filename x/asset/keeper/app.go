package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) GetAppID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.AppIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetAppID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.AppIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetApp(ctx sdk.Context, app types.AppData) {
	var (
		store = k.Store(ctx)
		key   = types.AppKey(app.Id)
		value = k.cdc.MustMarshal(&app)
	)

	store.Set(key, value)
}

func (k *Keeper) GetApp(ctx sdk.Context, id uint64) (app types.AppData, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AppKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return app, false
	}

	k.cdc.MustUnmarshal(value, &app)
	return app, true
}
func (k *Keeper) GetAppWasmQuery(ctx sdk.Context, id uint64) (int64, int64, uint64, error) {
	appData, _ := k.GetApp(ctx, id)
	minGovDeposit := appData.MinGovDeposit.Int64()
	var assetID uint64
	gen := appData.GenesisToken
	govTimeInSeconds := int64(appData.GovTimeInSeconds)
	for _, v := range gen {
		if v.IsGovToken {
			assetID = v.AssetId
		}
	}
	return minGovDeposit, govTimeInSeconds, assetID, nil
}

func (k *Keeper) GetApps(ctx sdk.Context) (apps []types.AppData, found bool) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AppKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var app types.AppData
		k.cdc.MustUnmarshal(iter.Value(), &app)
		apps = append(apps, app)
	}
	if apps == nil {
		return nil, false
	}

	return apps, true
}

func (k *Keeper) GetMintGenesisTokenData(ctx sdk.Context, appID, assetID uint64) (mintData types.MintGenesisToken, found bool) {
	appsData, _ := k.GetApp(ctx, appID)

	for _, data := range appsData.GenesisToken {
		if data.AssetId == assetID {
			return data, true
		}
	}
	return mintData, false
}

func (k *Keeper) CheckIfAssetIsAddedToAppMapping(ctx sdk.Context, assetID uint64) bool {
	apps, _ := k.GetApps(ctx)
	for _, data := range apps {
		for _, inData := range data.GenesisToken {
			if inData.AssetId == assetID {
				return false
			}
		}
	}
	return true
}

func (k *Keeper) SetAppForShortName(ctx sdk.Context, shortName string, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForShortNameKey(shortName)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}
func (k *Keeper) SetAppForName(ctx sdk.Context, Name string, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForNameKey(Name)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) HasAppForShortName(ctx sdk.Context, shortName string) bool {
	var (
		store = k.Store(ctx)
		key   = types.AssetForShortNameKey(shortName)
	)

	return store.Has(key)
}

func (k *Keeper) HasAppForName(ctx sdk.Context, Name string) bool {
	var (
		store = k.Store(ctx)
		key   = types.AssetForNameKey(Name)
	)

	return store.Has(key)
}

func (k *Keeper) SetGenesisTokenForApp(ctx sdk.Context, appID uint64, assetID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.GenesisForApp(appID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: assetID,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetGenesisTokenForApp(ctx sdk.Context, appID uint64) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.GenesisForApp(appID)
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) AddAppRecords(ctx sdk.Context, records ...types.AppData) error {
	for _, msg := range records {
		if k.HasAppForShortName(ctx, msg.ShortName) {
			return types.ErrorDuplicateApp
		}
		if k.HasAppForName(ctx, msg.ShortName) {
			return types.ErrorDuplicateApp
		}

		var (
			id  = k.GetAppID(ctx)
			app = types.AppData{
				Id:               id + 1,
				Name:             msg.Name,
				ShortName:        msg.ShortName,
				MinGovDeposit:    msg.MinGovDeposit,
				GovTimeInSeconds: msg.GovTimeInSeconds,
				GenesisToken:     msg.GenesisToken,
			}
		)

		k.SetAppID(ctx, app.Id)
		k.SetApp(ctx, app)
		k.SetAppForShortName(ctx, app.ShortName, app.Id)
		k.SetAppForName(ctx, app.Name, app.Id)
	}

	return nil
}

func (k *Keeper) UpdateGovTimeInApp(ctx sdk.Context, msg types.AppAndGovTime) error {
	appDetails, found := k.GetApp(ctx, msg.AppId)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	appDetails.GovTimeInSeconds = msg.GovTimeInSeconds

	k.SetApp(ctx, appDetails)
	return nil
}

func (k *Keeper) AddAssetInAppRecords(ctx sdk.Context, records ...types.AppData) error {
	for _, msg := range records {
		appdata, found := k.GetApp(ctx, msg.Id)
		if !found {
			return types.AppIdsDoesntExist
		}
		// var mintGenesis []types.MintGenesisToken

		for _, data := range msg.GenesisToken {
			assetData, found := k.GetAsset(ctx, data.AssetId)
			if !found {
				return types.ErrorAssetDoesNotExist
			}
			if !assetData.IsOnChain {
				return types.ErrorAssetIsOffChain
			}
			hasAsset := k.GetGenesisTokenForApp(ctx, msg.Id)
			if hasAsset != 0 && data.IsGovToken {
				return types.ErrorGenesisTokenExistForApp
			}

			if data.IsGovToken && appdata.MinGovDeposit.Equal(sdk.ZeroInt()) {
				return types.ErrorMinGovDepositIsZero
			}

			checkFound := k.CheckIfAssetIsAddedToAppMapping(ctx, data.AssetId)
			if !checkFound {
				return types.ErrorAssetAlreadyExistingApp
			}
			if hasAsset == 0 && data.IsGovToken {
				k.SetGenesisTokenForApp(ctx, msg.Id, data.AssetId)
			}
			appdata.GenesisToken = append(appdata.GenesisToken, data)
		}
		var (
			app = types.AppData{
				Id:               msg.Id,
				Name:             appdata.Name,
				ShortName:        appdata.ShortName,
				MinGovDeposit:    appdata.MinGovDeposit,
				GovTimeInSeconds: appdata.GovTimeInSeconds,
				GenesisToken:     appdata.GenesisToken,
			}
		)
		k.SetApp(ctx, app)
	}
	return nil
}
