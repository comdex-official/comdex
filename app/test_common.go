package app

import (
	cdpkeeper "github.com/comdex-official/comdex/x/cdp/keeper"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmdb "github.com/tendermint/tm-db"
	"math/rand"
)

import (
	"github.com/tendermint/tendermint/libs/log"
)

type TestApp struct {
	App
}

func NewTestApp() TestApp {

	encoding := MakeEncodingConfig()
	db := tmdb.NewMemDB()
	app := New(log.NewNopLogger(),  db, nil, true, map[int64]bool{}, string(" "), 0, encoding, simapp.EmptyAppOptions{})
	return TestApp{App: *app}

}
func GeneratePrivKeyAddressPairs(n int) (keys []crypto.PrivKey, addrs []sdk.AccAddress) {
	r := rand.New(rand.NewSource(12345)) // make the generation deterministic
	keys = make([]crypto.PrivKey, n)
	addrs = make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		secret := make([]byte, 32)
		_, err := r.Read(secret)
		if err != nil {
			panic("Could not read randomness")
		}
		keys[i] = secp256k1.GenPrivKeySecp256k1(secret)
		addrs[i] = sdk.AccAddress(keys[i].PubKey().Address())
	}
	return

}


func (tApp TestApp) GetAccountKeeper() authkeeper.AccountKeeper { return tApp.accountKeeper }
func (tApp TestApp) GetCDPKeeper()      cdpkeeper.Keeper     { return tApp.cdpKeeper }
func (tApp TestApp) GetBankKeeper() bankkeeper.Keeper { return tApp.bankKeeper }
