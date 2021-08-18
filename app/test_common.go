package app

import (
	cdpkeeper "github.com/comdex-official/comdex/x/cdp/keeper"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	ibctransferkeeper "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/keeper"
	ibckeeper "github.com/cosmos/cosmos-sdk/x/ibc/core/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	tmdb "github.com/tendermint/tm-db"
	"math/rand"
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
func (tApp TestApp) GetStakingKeeper() stakingkeeper.Keeper { return tApp.stakingKeeper }
func (tApp TestApp) GetCapabilityKeeper() *capabilitykeeper.Keeper { return tApp.capabilityKeeper }
func (tApp TestApp) GetSlashingKeeper() slashingkeeper.Keeper { return tApp.slashingKeeper }
func (tApp TestApp) GetMintKeeper() mintkeeper.Keeper { return tApp.mintKeeper }
func (tApp TestApp) GetDistrKeeper() distrkeeper.Keeper { return tApp.distrKeeper }
func (tApp TestApp) GetGovKeeper() govkeeper.Keeper { return tApp.govKeeper }
func (tApp TestApp) GetCrisisKeeper() crisiskeeper.Keeper { return tApp.crisisKeeper }
func (tApp TestApp) GetUpgradeKeeper() upgradekeeper.Keeper { return tApp.upgradeKeeper }
func (tApp TestApp) GetParamsKeeper() paramskeeper.Keeper { return tApp.paramsKeeper }
func (tApp TestApp) GetIbcKeeper() *ibckeeper.Keeper { return tApp.ibcKeeper }
func (tApp TestApp) GetEvidenceKeeper() evidencekeeper.Keeper { return tApp.evidenceKeeper }
func (tApp TestApp) GetIbcTransferKeeper() ibctransferkeeper.Keeper { return tApp.ibcTransferKeeper }