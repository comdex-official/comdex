package wasm

import (
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/comdex-official/comdex/app"
)

func SetupCustomApp() (*app.App, *sdk.Context) {
	osmosis, ctx := CreateTestInput()
	return osmosis, ctx
}

func CreateTestInput() (*app.App, *sdk.Context) {
	comdex := app.Setup(false)
	ctx := comdex.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "comdex-1", Time: time.Now().UTC()})
	return comdex, &ctx
}

func FundAccount(t *testing.T, ctx sdk.Context, comdex *app.App, acct sdk.AccAddress) {
	err := simapp.FundAccount(comdex.BankKeeper, ctx, acct, sdk.NewCoins(
		sdk.NewCoin("ucmdx", sdk.NewInt(10000000000)),
	))
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
	msg1 := []assetTypes.AppData{{
		Name:             "cswap",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900},
		{
			Name:             "commodo",
			ShortName:        "commodo",
			MinGovDeposit:    sdk.NewIntFromUint64(10000000),
			GovTimeInSeconds: 900},
	}
	_ = assetKeeper.AddAppRecords(*ctx, msg1...)

	msg2 := []assetTypes.Asset{
		{Name: "CMDX",
			Denom:     "ucmdx",
			Decimals:  1000000,
			IsOnChain: true}, {Name: "CMST",
			Denom:     "ucmst",
			Decimals:  1000000,
			IsOnChain: true}, {Name: "HARBOR",
			Denom:     "uharbor",
			Decimals:  1000000,
			IsOnChain: true},
	}
	_ = assetKeeper.AddAssetRecords(*ctx, msg2...)

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
		{"Add Pair 1: cmdx cmst",
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
