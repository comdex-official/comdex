package keeper

import (
<<<<<<< HEAD
	"fmt"

=======
>>>>>>> 394fefd550310c2dae2b25c9005c3e4843bd47ff
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

func (k *Keeper) GetMintGenesisTokenData(ctx sdk.Context, appId, assetId uint64) (mintData types.MintGenesisToken, found bool) {
	appsData, _ := k.GetApp(ctx, appId)

	for _, data := range appsData.MintGenesisToken {
		if data.AssetId == assetId {
			return *data, true
		}
	}
	return mintData, false

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

func (k *Keeper) SetGenesisTokenForApp(ctx sdk.Context, appId uint64, assetId uint64) {
	var (
		store = k.Store(ctx)
		key   = types.GensisForApp(appId)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: assetId,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetGenesisTokenForApp(ctx sdk.Context, appId uint64) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.GensisForApp(appId)
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
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

			assetData, found := k.GetAsset(ctx, data.AssetId)
			if !found{
				return types.ErrorAssetDoesNotExist
			}
			if assetData.IsOnchain{
				return types.ErrorAssetIsOnChain
			}
			hasAsset := k.GetGenesisTokenForApp(ctx, msg.Id)
			if hasAsset != 0{
				return types.ErrorGenesisTokenExistForApp
			}

			checkfound := k.CheckIfAssetIsaddedToAppMapping(ctx, data.AssetId)
			if !checkfound {
				return types.ErrorAssetAlreadyExistinApp
			}
			if data.IsgovToken {
				k.SetGenesisTokenForApp(ctx, msg.Id ,data.AssetId)
			}
			if data.GenesisSupply.IsZero() {
				return types.ErrorGenesisCantBeZero
			}
		}
		for _, data := range msg.MintGenesisToken {
			var minter types.MintGenesisToken
			minter.AssetId = data.AssetId
			minter.GenesisSupply = data.GenesisSupply
			minter.IsgovToken = data.IsgovToken
			minter.Sender = data.Sender
			mintGenesis = append(mintGenesis, &minter)

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
