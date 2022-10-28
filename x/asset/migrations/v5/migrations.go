package v5

import (
	"fmt"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	fmt.Println("MigrateStore")
	return MigrateValueApps(store, cdc)
}

func MigrateValueApps(store sdk.KVStore, cdc codec.BinaryCodec) error {

	app1 := assettypes.AppData{
		Id:               1,
		Name:             "CSWAP",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.ZeroInt(),
		GovTimeInSeconds: 300,
		GenesisToken:     nil,
	}
	key1 := assettypes.AppKey(app1.Id)
	store.Delete(key1)
	SetApp(store, cdc, app1)

	genesisToken := assettypes.MintGenesisToken{
		AssetId:       9,
		GenesisSupply: sdk.NewIntFromUint64(1000000000000000),
		IsGovToken:    true,
		Recipient:     "comdex1unvvj23q89dlgh82rdtk5su7akdl5932reqarg",
	}
	var gToken []assettypes.MintGenesisToken
	gToken = append(gToken, genesisToken)
	app2 := assettypes.AppData{
		Id:               2,
		Name:             "HARBOR",
		ShortName:        "hbr",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 300,
		GenesisToken:     gToken,
	}
	key2 := assettypes.AppKey(app2.Id)
	store.Delete(key2)
	SetApp(store, cdc, app2)

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
	SetAppID(store, cdc, 3)
	return nil
}

func SetApp(store sdk.KVStore, cdc codec.BinaryCodec, app assettypes.AppData) {
	var (
		key   = assettypes.AppKey(app.Id)
		value = cdc.MustMarshal(&app)
	)

	store.Set(key, value)
}

func SetAppID(store sdk.KVStore, cdc codec.BinaryCodec, id uint64) {
	var (
		key   = assettypes.AppIDKey
		value = cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}
