package v5

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	err := MigrateValueApps(store, cdc)
	if err != nil {
		return err
	}
	err = MigrateValueAsset(store, cdc)
	if err != nil {
		return err
	}
	return err
}

func MigrateValueApps(store sdk.KVStore, cdc codec.BinaryCodec) error {

	genesisToken := assettypes.MintGenesisToken{
		AssetId:       9,
		GenesisSupply: sdk.NewIntFromUint64(1000000000000000),
		IsGovToken:    true,
		Recipient:     "comdex1unvvj23q89dlgh82rdtk5su7akdl5932reqarg",
	}
	var gToken []assettypes.MintGenesisToken
	gToken = append(gToken, genesisToken)
	app2 := assettypes.AppData{
		Id:               1,
		Name:             "HARBOR",
		ShortName:        "hbr",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 300,
		GenesisToken:     gToken,
	}
	key2 := assettypes.AppKey(app2.Id)
	store.Delete(key2)
	SetApp(store, cdc, app2)

	app1 := assettypes.AppData{
		Id:               2,
		Name:             "CSWAP",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.ZeroInt(),
		GovTimeInSeconds: 300,
		GenesisToken:     nil,
	}
	key1 := assettypes.AppKey(app1.Id)
	store.Delete(key1)
	SetApp(store, cdc, app1)

	app3 := assettypes.AppData{
		Id:               3,
		Name:             "commodo",
		ShortName:        "cmdo",
		MinGovDeposit:    sdk.ZeroInt(),
		GovTimeInSeconds: 0,
		GenesisToken:     nil,
	}
	key3 := assettypes.AppKey(app3.Id)
	store.Delete(key3)
	SetApp(store, cdc, app3)
	return nil
}

func MigrateValueAsset(store sdk.KVStore, cdc codec.BinaryCodec) error {
	asset1 := assettypes.Asset{
		Id:                    1,
		Name:                  "ATOM",
		Denom:                 "uatom",
		Decimals:              sdk.NewInt(1000000),
		IsOnChain:             false,
		IsOraclePriceRequired: true,
	}
	key1 := assettypes.AssetKey(asset1.Id)
	store.Delete(key1)
	SetAsset(store, cdc, asset1)

	asset2 := assettypes.Asset{
		Id:                    2,
		Name:                  "CMDX",
		Denom:                 "ucmdx",
		Decimals:              sdk.NewInt(1000000),
		IsOnChain:             false,
		IsOraclePriceRequired: true,
	}
	key2 := assettypes.AssetKey(asset2.Id)
	store.Delete(key2)
	SetAsset(store, cdc, asset2)

	asset3 := assettypes.Asset{
		Id:                    3,
		Name:                  "CMST",
		Denom:                 "ucmst",
		Decimals:              sdk.NewInt(1000000),
		IsOnChain:             false,
		IsOraclePriceRequired: true,
	}
	key3 := assettypes.AssetKey(asset3.Id)
	store.Delete(key3)
	SetAsset(store, cdc, asset3)

	asset4 := assettypes.Asset{
		Id:                    4,
		Name:                  "OSMO",
		Denom:                 "uosmo",
		Decimals:              sdk.NewInt(1000000),
		IsOnChain:             false,
		IsOraclePriceRequired: true,
	}
	key4 := assettypes.AssetKey(asset4.Id)
	store.Delete(key4)
	SetAsset(store, cdc, asset4)

	asset5 := assettypes.Asset{
		Id:                    5,
		Name:                  "cATOM",
		Denom:                 "ucatom",
		Decimals:              sdk.NewInt(1000000),
		IsOnChain:             false,
		IsOraclePriceRequired: false,
	}
	key5 := assettypes.AssetKey(asset5.Id)
	store.Delete(key5)
	SetAsset(store, cdc, asset5)

	asset6 := assettypes.Asset{
		Id:                    6,
		Name:                  "cCMDX",
		Denom:                 "uccmdx",
		Decimals:              sdk.NewInt(1000000),
		IsOnChain:             false,
		IsOraclePriceRequired: false,
	}
	key6 := assettypes.AssetKey(asset6.Id)
	store.Delete(key6)
	SetAsset(store, cdc, asset6)

	asset7 := assettypes.Asset{
		Id:                    7,
		Name:                  "cCMST",
		Denom:                 "uccmst",
		Decimals:              sdk.NewInt(1000000),
		IsOnChain:             false,
		IsOraclePriceRequired: false,
	}
	key7 := assettypes.AssetKey(asset7.Id)
	store.Delete(key7)
	SetAsset(store, cdc, asset7)

	asset8 := assettypes.Asset{
		Id:                    8,
		Name:                  "cOSMO",
		Denom:                 "ucosmo",
		Decimals:              sdk.NewInt(1000000),
		IsOnChain:             false,
		IsOraclePriceRequired: false,
	}
	key8 := assettypes.AssetKey(asset8.Id)
	store.Delete(key8)
	SetAsset(store, cdc, asset8)

	asset9 := assettypes.Asset{
		Id:                    9,
		Name:                  "HARBOR",
		Denom:                 "uharbor",
		Decimals:              sdk.NewInt(1000000),
		IsOnChain:             true,
		IsOraclePriceRequired: false,
	}
	key9 := assettypes.AssetKey(asset9.Id)
	store.Delete(key9)
	SetAsset(store, cdc, asset9)

	return nil
}

func SetApp(store sdk.KVStore, cdc codec.BinaryCodec, app assettypes.AppData) {
	var (
		key   = assettypes.AppKey(app.Id)
		value = cdc.MustMarshal(&app)
	)

	store.Set(key, value)
}

func SetAsset(store sdk.KVStore, cdc codec.BinaryCodec, asset assettypes.Asset) {
	var (
		key   = assettypes.AssetKey(asset.Id)
		value = cdc.MustMarshal(&asset)
	)

	store.Set(key, value)
}
