package wasm

import (
	"fmt"
	"os"
	"testing"
	"time"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/comdex-official/comdex/app/wasm/bindings"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	tokenmintTypes "github.com/comdex-official/comdex/x/tokenmint/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	// simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/tokenmint/keeper"
)

// func SetupCustomApp(t *testing.T) (*app.App, *sdk.Context) {
// 	comdex, ctx := CreateTestInput(t)
// 	return comdex, ctx
// }

func SetupCustomApp(t *testing.T, addr sdk.AccAddress) (*app.App, *sdk.Context) {
	comdexApp, ctx := CreateTestInput(t)
	wasmKeeper := comdexApp.WasmKeeper

	storeReflectCode(t, *ctx, comdexApp, addr)

	cInfo := wasmKeeper.GetCodeInfo(*ctx, 1)
	require.NotNil(t, cInfo)

	return comdexApp, ctx
}

func storeReflectCode(t *testing.T, ctx sdk.Context, comdexApp *app.App, addr sdk.AccAddress) uint64 {
	wasmCode, err := os.ReadFile("../testdata/token_reflect.wasm")
	require.NoError(t, err)

	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(comdexApp.WasmKeeper)
	codeID, _, err := contractKeeper.Create(ctx, addr, wasmCode, nil)
	require.NoError(t, err)

	return codeID
}

func instantiateReflectContract(t *testing.T, ctx sdk.Context, comdexApp *app.App, funder sdk.AccAddress) sdk.AccAddress {
	initMsgBz := []byte("{}")
	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(comdexApp.WasmKeeper)
	codeID := uint64(1)
	addr, _, err := contractKeeper.Instantiate(ctx, codeID, funder, funder, initMsgBz, "demo contract", nil)
	require.NoError(t, err)

	return addr
}

func CreateTestInput(t *testing.T) (*app.App, *sdk.Context) {
	comdex := app.Setup(t, false)
	ctx := comdex.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "comdex-1", Time: time.Now().UTC()})
	return comdex, &ctx
}

func fundAccount(t *testing.T, ctx sdk.Context, comdexapp *app.App, addr sdk.AccAddress, coins sdk.Coins) {
	err := banktestutil.FundAccount(
		comdexapp.BankKeeper,
		ctx,
		addr,
		coins,
	)
	require.NoError(t, err)
}

// we need to make this deterministic (same every test run), as content might affect gas costs
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func RandomAccountAddress() sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress() string {
	return RandomAccountAddress().String()
}

func AddAppAsset(app *app.App, ctx1 sdk.Context) {
	assetKeeper, ctx := &app.AssetKeeper, &ctx1
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	genesisSupply := sdk.NewIntFromUint64(9000000)
	msg1 := assetTypes.AppData{
		Name:             "cswap",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
		GenesisToken: []assetTypes.MintGenesisToken{
			{
				3,
				genesisSupply,
				true,
				userAddress,
			},
		},
	}
	_ = assetKeeper.AddAppRecords(*ctx, msg1)

	msg2 := assetTypes.Asset{
		Name:          "CMDX",
		Denom:         "ucmdx",
		Decimals:      sdk.NewInt(1000000),
		IsOnChain:     true,
		IsCdpMintable: true,
	}
	_ = assetKeeper.AddAssetRecords(*ctx, msg2)

	msg3 := assetTypes.Asset{
		Name:          "CMST",
		Denom:         "ucmst",
		Decimals:      sdk.NewInt(1000000),
		IsOnChain:     true,
		IsCdpMintable: true,
	}
	_ = assetKeeper.AddAssetRecords(*ctx, msg3)

	msg4 := assetTypes.Asset{
		Name:          "HARBOR",
		Denom:         "uharbor",
		Decimals:      sdk.NewInt(1000000),
		IsOnChain:     true,
		IsCdpMintable: true,
	}
	_ = assetKeeper.AddAssetRecords(*ctx, msg4)
}

func AddPair(app *app.App, ctx1 sdk.Context) {
	AddAppAsset(app, ctx1)
	assetKeeper, ctx := &app.AssetKeeper, &ctx1
	for _, tc := range []struct {
		name                   string
		pair                   assetTypes.Pair
		symbol1                string
		symbol2                string
		isErrorExpectedForPair bool
		pairID                 uint64
	}{
		{
			"Add Pair 1: cmdx cmst",
			assetTypes.Pair{
				AssetIn:  1,
				AssetOut: 2,
			},
			"ucmdx",
			"ucmst",
			false,
			1,
		},
	} {
		_ = assetKeeper.AddPairsRecords(*ctx, tc.pair)
	}
}

func AddCollectorLookuptable(app *app.App, ctx1 sdk.Context) {
	AddAppAsset(app, ctx1)
	collectorKeeper, ctx := &app.CollectorKeeper, &ctx1
	for _, tc := range []struct {
		name string
		msg  bindings.MsgSetCollectorLookupTable
	}{
		{
			"Wasm Add MsgSetCollectorLookupTable AppID 1 CollectorAssetID 2",
			bindings.MsgSetCollectorLookupTable{
				AppID:            1,
				CollectorAssetID: 2,
				SecondaryAssetID: 3,
				SurplusThreshold: sdk.NewInt(10000000),
				DebtThreshold:    sdk.NewInt(5000000),
				LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
				LotSize:          sdk.NewInt(2000000),
				BidFactor:        sdk.MustNewDecFromStr("0.01"),
				DebtLotSize:      sdk.NewInt(2000000),
			},
		},
	} {
		_ = collectorKeeper.WasmSetCollectorLookupTable(*ctx, &tc.msg)
	}
}

func AddExtendedPairVault(app *app.App, ctx1 sdk.Context) {
	AddAppAsset(app, ctx1)
	assetKeeper, ctx := &app.AssetKeeper, &ctx1
	for _, tc := range []struct {
		name string
		msg  bindings.MsgAddExtendedPairsVault
	}{
		{
			"Add Extended Pair Vault : cmdx cmst",

			bindings.MsgAddExtendedPairsVault{
				AppID:               1,
				PairID:              1,
				StabilityFee:        sdk.MustNewDecFromStr("0.01"),
				ClosingFee:          sdk.MustNewDecFromStr("0"),
				LiquidationPenalty:  sdk.MustNewDecFromStr("0.12"),
				DrawDownFee:         sdk.MustNewDecFromStr("0.01"),
				IsVaultActive:       true,
				DebtCeiling:         sdk.NewInt(1000000000000),
				DebtFloor:           sdk.NewInt(1000000),
				IsStableMintVault:   false,
				MinCr:               sdk.MustNewDecFromStr("1.5"),
				PairName:            "CMDX-A",
				AssetOutOraclePrice: true,
				AssetOutPrice:       1000000,
				MinUsdValueLeft:     1000000,
			},
		},
	} {
		_ = assetKeeper.WasmAddExtendedPairsVaultRecords(*ctx, &tc.msg)
	}
}

func WhitelistAppIDLiquidation(app *app.App, ctx1 sdk.Context) {
	AddAppAsset(app, ctx1)
	liquidationKeeper, ctx := &app.LiquidationKeeper, &ctx1
	for _, tc := range []struct {
		name string
		msg  bindings.MsgWhitelistAppIDLiquidation
	}{
		{
			"Whitelist AppID Liquidation",

			bindings.MsgWhitelistAppIDLiquidation{
				AppID: 1,
			},
		},
	} {
		_ = liquidationKeeper.WasmWhitelistAppIDLiquidation(*ctx, tc.msg.AppID)
	}
}

func MsgMintNewTokens(app *app.App, ctx1 sdk.Context) {
	AddAppAsset(app, ctx1)
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	tokenmintKeeper, ctx := &app.TokenmintKeeper, &ctx1
	wctx := sdk.WrapSDKContext(*ctx)

	server := keeper.NewMsgServer(*tokenmintKeeper)
	for _, tc := range []struct {
		name          string
		msg           tokenmintTypes.MsgMintNewTokensRequest
		expectedError bool
	}{
		{
			"Mint New Tokens : App ID : 1, Asset ID : 3",
			tokenmintTypes.MsgMintNewTokensRequest{
				From:    userAddress,
				AppId:   1,
				AssetId: 3,
			},
			false,
		},
	} {
		_, err := server.MsgMintNewTokens(wctx, &tc.msg)
		fmt.Println(err)
	}
}

func FundAccountFunc(bankKeeper bankkeeper.Keeper, ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}
