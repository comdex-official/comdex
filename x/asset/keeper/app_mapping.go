package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) GetAppID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.AppIDkey
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
		key   = types.AppIDkey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetApp(ctx sdk.Context, app types.AppMapping) {
	var (
		store = k.Store(ctx)
		key   = types.AppKey(app.Id)
		value = k.cdc.MustMarshal(&app)
	)

	store.Set(key, value)
}

func (k *Keeper) GetApp(ctx sdk.Context, id uint64) (app types.AppMapping, found bool) {
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

func (k *Keeper) GetApps(ctx sdk.Context) (apps []types.AppMapping, found bool) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AppKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var app types.AppMapping
		k.cdc.MustUnmarshal(iter.Value(), &app)
		apps = append(apps, app)
	}
	if apps == nil {
		return nil, false
	}

	return apps, true
}

func (k *Keeper) GetMintGenesisTokenData(ctx sdk.Context, appId, assetId uint64) (types.MintGenesisToken, error) {
	appsData, found := k.GetApp(ctx, appId)
	var minted types.MintGenesisToken
	if !found {
		return minted, types.AppIdsDoesntExist
	}

	for _, data := range appsData.MintGenesisToken {
		if data.AssetId == assetId {
			minted = *data
		}
	}
	return minted, nil
}

func (k *Keeper) CheckIfAssetIsaddedToAppMapping(ctx sdk.Context, assetId uint64) bool {
	apps, _ := k.GetApps(ctx)
	for _, data := range apps {
		for _, inData := range data.MintGenesisToken {
			if inData.AssetId == assetId {
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

func (k *Keeper) AddAppMappingRecords(ctx sdk.Context, records ...types.AppMapping) error {
	for _, msg := range records {
		if k.HasAppForShortName(ctx, msg.ShortName) {
			return types.ErrorDuplicateApp
		}
		if k.HasAppForName(ctx, msg.ShortName) {
			return types.ErrorDuplicateApp
		}

		var (
			id  = k.GetAppID(ctx)
			app = types.AppMapping{
				Id:               id + 1,
				Name:             msg.Name,
				ShortName:        msg.ShortName,
				MintGenesisToken: msg.MintGenesisToken,
			}
		)

		k.SetAppID(ctx, app.Id)
		k.SetApp(ctx, app)
		k.SetAppForShortName(ctx, app.ShortName, app.Id)
		k.SetAppForName(ctx, app.Name, app.Id)

	}

	return nil
}

func (k *Keeper) AddAssetMappingRecords(ctx sdk.Context, records ...types.AppMapping) error {
	for _, msg := range records {

		appdata, found := k.GetApp(ctx, msg.Id)
		if !found {
			return types.AppIdsDoesntExist
		}
		var mintGenesis []*types.MintGenesisToken

		for _, data := range msg.MintGenesisToken {
			if !k.HasAsset(ctx, data.AssetId) {
				return types.ErrorAssetDoesNotExist
			}
			found := k.CheckIfAssetIsaddedToAppMapping(ctx, data.AssetId)
			if !found {
				return types.ErrorAssetAlreadyExistinApp
			}
		}
		for _, data := range msg.MintGenesisToken {
			var minter types.MintGenesisToken
			minter.AssetId = data.AssetId
			minter.GenesisSupply = data.GenesisSupply
			minter.IsgovToken = data.IsgovToken
			minter.Sender = data.Sender
			mintGenesis = append(mintGenesis, &minter)
			fmt.Println("data", data)
			fmt.Println("data", data)
			fmt.Println("data sender", data.Sender)

		}

		var (
			app = types.AppMapping{
				Id:               msg.Id,
				Name:             appdata.Name,
				ShortName:        appdata.ShortName,
				MintGenesisToken: mintGenesis,
			}
		)

		k.SetApp(ctx, app)

	}

	return nil
}
