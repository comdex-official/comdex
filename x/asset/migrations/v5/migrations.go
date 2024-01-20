package v5

import (
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/cosmos/cosmos-sdk/codec"
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

func MigrateValueApps(store storetypes.KVStore, cdc codec.BinaryCodec) error {
	app1 := assettypes.AppData{
		Id:               1,
		Name:             "CSWAP",
		ShortName:        "cswap",
		MinGovDeposit:    sdkmath.ZeroInt(),
		GovTimeInSeconds: 300,
		GenesisToken:     nil,
	}
	key1 := assettypes.AppKey(app1.Id)
	store.Delete(key1)
	SetApp(store, cdc, app1)

	genesisToken := assettypes.MintGenesisToken{
		AssetId:       9,
		GenesisSupply: sdkmath.NewIntFromUint64(1000000000000000),
		IsGovToken:    true,
		Recipient:     "comdex1unvvj23q89dlgh82rdtk5su7akdl5932reqarg",
	}
	var gToken []assettypes.MintGenesisToken
	gToken = append(gToken, genesisToken)
	app2 := assettypes.AppData{
		Id:               2,
		Name:             "HARBOR",
		ShortName:        "hbr",
		MinGovDeposit:    sdkmath.NewIntFromUint64(10000000),
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
		MinGovDeposit:    sdkmath.ZeroInt(),
		GovTimeInSeconds: 0,
		GenesisToken:     nil,
	}
	key3 := assettypes.AppKey(app3.Id)
	store.Delete(key3)
	SetApp(store, cdc, app3)
	return nil
}

func MigrateValueAsset(store storetypes.KVStore, cdc codec.BinaryCodec) error {
	asset1 := assettypes.Asset{
		Id:                    1,
		Name:                  "ATOM",
		Denom:                 "ibc/2E5D0AC026AC1AFA65A23023BA4F24BB8DDF94F118EDC0BAD6F625BFC557CDED",
		Decimals:              sdkmath.NewInt(1000000),
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
		Decimals:              sdkmath.NewInt(1000000),
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
		Decimals:              sdkmath.NewInt(1000000),
		IsOnChain:             false,
		IsOraclePriceRequired: true,
	}
	key3 := assettypes.AssetKey(asset3.Id)
	store.Delete(key3)
	SetAsset(store, cdc, asset3)

	asset4 := assettypes.Asset{
		Id:                    4,
		Name:                  "OSMO",
		Denom:                 "ibc/868AF0A32D53849B6093348F5A47BB969A98E71A3F0CD2D3BE406EA25DA7F836",
		Decimals:              sdkmath.NewInt(1000000),
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
		Decimals:              sdkmath.NewInt(1000000),
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
		Decimals:              sdkmath.NewInt(1000000),
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
		Decimals:              sdkmath.NewInt(1000000),
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
		Decimals:              sdkmath.NewInt(1000000),
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
		Decimals:              sdkmath.NewInt(1000000),
		IsOnChain:             true,
		IsOraclePriceRequired: false,
	}
	key9 := assettypes.AssetKey(asset9.Id)
	store.Delete(key9)
	SetAsset(store, cdc, asset9)

	asset10 := assettypes.Asset{
		Id:                    10,
		Name:                  "USDC",
		Denom:                 "ibc/EF8A76D0FD3F3F45D8DB7FEBFCF921206DF58CA41493ED16D69BF7B4E061C60C",
		Decimals:              sdkmath.NewInt(1000000),
		IsOnChain:             false,
		IsOraclePriceRequired: true,
	}
	key10 := assettypes.AssetKey(asset10.Id)
	store.Delete(key10)
	SetAsset(store, cdc, asset10)

	asset11 := assettypes.Asset{
		Id:                    11,
		Name:                  "WETH",
		Denom:                 "ibc/A99459944FD67B5711735B4B4D3FE30BA45328E94D437C78E47CA8DEFA781E49",
		Decimals:              sdkmath.NewInt(1000000000000000000),
		IsOnChain:             false,
		IsOraclePriceRequired: true,
	}
	key11 := assettypes.AssetKey(asset11.Id)
	store.Delete(key11)
	SetAsset(store, cdc, asset11)

	return nil
}

func SetApp(store storetypes.KVStore, cdc codec.BinaryCodec, app assettypes.AppData) {
	var (
		key   = assettypes.AppKey(app.Id)
		value = cdc.MustMarshal(&app)
	)

	store.Set(key, value)
}

func SetAsset(store storetypes.KVStore, cdc codec.BinaryCodec, asset assettypes.Asset) {
	var (
		key   = assettypes.AssetKey(asset.Id)
		value = cdc.MustMarshal(&asset)
	)

	store.Set(key, value)
}
