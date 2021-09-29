package app

import (
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"
	"github.com/cosmos/cosmos-sdk/simapp"
	tmdb "github.com/tendermint/tm-db"
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


func (tApp TestApp) GetVaultKeeper() vaultkeeper.Keeper { return tApp.vaultKeeper }
func (tApp TestApp) GetAssetKeeper() assetkeeper.Keeper { return tApp.assetKeeper }