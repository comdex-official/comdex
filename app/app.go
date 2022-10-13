package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"
	icahost "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/spf13/cast"

	wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfer "github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v3/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v3/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	"github.com/comdex-official/comdex/x/liquidation"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/comdex-official/comdex/x/asset"
	assetclient "github.com/comdex-official/comdex/x/asset/client"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/auction"
	auctionkeeper "github.com/comdex-official/comdex/x/auction/keeper"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/collector"
	collectorkeeper "github.com/comdex-official/comdex/x/collector/keeper"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	"github.com/comdex-official/comdex/x/esm"
	esmkeeper "github.com/comdex-official/comdex/x/esm/keeper"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	rewardsclient "github.com/comdex-official/comdex/x/rewards/client"

	"github.com/comdex-official/comdex/x/lend"
	lendclient "github.com/comdex-official/comdex/x/lend/client"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/comdex-official/comdex/x/locker"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"

	bandoraclemodule "github.com/comdex-official/comdex/x/bandoracle"
	bandoraclemoduleclient "github.com/comdex-official/comdex/x/bandoracle/client"
	bandoraclemodulekeeper "github.com/comdex-official/comdex/x/bandoracle/keeper"
	bandoraclemoduletypes "github.com/comdex-official/comdex/x/bandoracle/types"

	"github.com/comdex-official/comdex/x/market"
	marketkeeper "github.com/comdex-official/comdex/x/market/keeper"
	markettypes "github.com/comdex-official/comdex/x/market/types"

	"github.com/comdex-official/comdex/x/rewards"
	rewardskeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"

	"github.com/comdex-official/comdex/x/tokenmint"
	tokenmintkeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"

	"github.com/comdex-official/comdex/x/vault"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"

	"github.com/comdex-official/comdex/x/liquidity"
	liquidityclient "github.com/comdex-official/comdex/x/liquidity/client"
	liquiditykeeper "github.com/comdex-official/comdex/x/liquidity/keeper"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"

	cwasm "github.com/comdex-official/comdex/app/wasm"

	mv5 "github.com/comdex-official/comdex/app/upgrades/mainnet/v5"
	tv1_0_0 "github.com/comdex-official/comdex/app/upgrades/testnet/v1_0_0"
	tv2_0_0 "github.com/comdex-official/comdex/app/upgrades/testnet/v2_0_0"
	tv3_0_0 "github.com/comdex-official/comdex/app/upgrades/testnet/v3_0_0"
	tv4_0_0 "github.com/comdex-official/comdex/app/upgrades/testnet/v4_0_0"
)

const (
	AccountAddressPrefix = "comdex"
	Name                 = "comdex"
)

// GetWasmEnabledProposals parses the WasmProposalsEnabled / EnableSpecificWasmProposals values to
// produce a list of enabled proposals to pass into wasmd app.
func GetWasmEnabledProposals() []wasm.ProposalType {
	if EnableSpecificWasmProposals == "" {
		if WasmProposalsEnabled == "true" {
			return wasm.EnableAllProposals
		}
		return wasm.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificWasmProposals, ",")
	proposals, err := wasm.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
}

func GetGovProposalHandlers() []govclient.ProposalHandler {
	proposalHandlers := []govclient.ProposalHandler{
		bandoraclemoduleclient.AddFetchPriceHandler,
		lendclient.AddLendPairsHandler,
		lendclient.AddPoolHandler,
		lendclient.UpdateLendPairsHandler,
		lendclient.AddAssetToPairHandler,
		lendclient.AddAssetRatesParamsHandler,
		lendclient.AddAuctionParamsHandler,
		paramsclient.ProposalHandler,
		distrclient.ProposalHandler,
		upgradeclient.ProposalHandler,
		upgradeclient.CancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
	}
	proposalHandlers = append(proposalHandlers, wasmclient.ProposalHandlers...)
	proposalHandlers = append(proposalHandlers, assetclient.AddAssetsHandler...)
	proposalHandlers = append(proposalHandlers, rewardsclient.AddRewardsHandler...)
	proposalHandlers = append(proposalHandlers, liquidityclient.LiquidityProposalHandler...)
	return proposalHandlers
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
	// If EnableSpecificWasmProposals is "", and this is "true", then enable all x/wasm proposals.
	// If EnableSpecificWasmProposals is "", and this is not "true", then disable all x/wasm proposals.
	WasmProposalsEnabled = "true"
	// If set to non-empty string it must be comma-separated list of values that are all a subset
	// of "EnableAllProposals" (takes precedence over WasmProposalsEnabled)
	// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
	EnableSpecificWasmProposals = ""
	// use this for clarity in argument list
	EmptyWasmOpts []wasm.Option
	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(GetGovProposalHandlers()...),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		vault.AppModuleBasic{},
		asset.AppModuleBasic{},
		esm.AppModuleBasic{},
		lend.AppModuleBasic{},

		market.AppModuleBasic{},
		locker.AppModuleBasic{},
		bandoraclemodule.AppModuleBasic{},
		collector.AppModuleBasic{},
		liquidation.AppModuleBasic{},
		auction.AppModuleBasic{},
		tokenmint.AppModuleBasic{},
		wasm.AppModuleBasic{},
		liquidity.AppModuleBasic{},
		rewards.AppModuleBasic{},
		ica.AppModuleBasic{},
	)
)

var _ servertypes.Application = (*App)(nil)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)
}

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp

	amino *codec.LegacyAmino
	cdc   codec.Codec

	interfaceRegistry codectypes.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey
	mkeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper     authkeeper.AccountKeeper
	FeegrantKeeper    feegrantkeeper.Keeper
	BankKeeper        bankkeeper.Keeper
	AuthzKeeper       authzkeeper.Keeper
	CapabilityKeeper  *capabilitykeeper.Keeper
	StakingKeeper     stakingkeeper.Keeper
	SlashingKeeper    slashingkeeper.Keeper
	MintKeeper        mintkeeper.Keeper
	DistrKeeper       distrkeeper.Keeper
	GovKeeper         govkeeper.Keeper
	CrisisKeeper      crisiskeeper.Keeper
	UpgradeKeeper     upgradekeeper.Keeper
	ParamsKeeper      paramskeeper.Keeper
	IbcKeeper         *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	ICAHostKeeper     icahostkeeper.Keeper
	EvidenceKeeper    evidencekeeper.Keeper
	IbcTransferKeeper ibctransferkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper         capabilitykeeper.ScopedKeeper
	ScopedIBCTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedIBCOracleKeeper   capabilitykeeper.ScopedKeeper
	ScopedBandoracleKeeper  capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper     capabilitykeeper.ScopedKeeper

	BandoracleKeeper bandoraclemodulekeeper.Keeper
	AssetKeeper      assetkeeper.Keeper
	CollectorKeeper  collectorkeeper.Keeper
	VaultKeeper      vaultkeeper.Keeper

	MarketKeeper      marketkeeper.Keeper
	LiquidationKeeper liquidationkeeper.Keeper
	LockerKeeper      lockerkeeper.Keeper
	EsmKeeper         esmkeeper.Keeper
	LendKeeper        lendkeeper.Keeper
	ScopedWasmKeeper  capabilitykeeper.ScopedKeeper
	AuctionKeeper     auctionkeeper.Keeper
	TokenmintKeeper   tokenmintkeeper.Keeper
	LiquidityKeeper   liquiditykeeper.Keeper
	Rewardskeeper     rewardskeeper.Keeper

	WasmKeeper wasm.Keeper
	// the module manager
	mm *module.Manager
	// Module configurator
	configurator module.Configurator
}

// New returns a reference to an initialized App.
func New(
	logger log.Logger,
	db tmdb.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encoding EncodingConfig,
	appOptions servertypes.AppOptions,
	wasmEnabledProposals []wasm.ProposalType,
	wasmOpts []wasm.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	appCodec := encoding.Marshaler
	var (
		tkeys = sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
		mkeys = sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
		keys  = sdk.NewKVStoreKeys(
			authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
			minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
			govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, icahosttypes.StoreKey, upgradetypes.StoreKey,
			evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
			vaulttypes.StoreKey, assettypes.StoreKey, collectortypes.StoreKey, liquidationtypes.StoreKey,
			markettypes.StoreKey, bandoraclemoduletypes.StoreKey, lockertypes.StoreKey,
			wasm.StoreKey, authzkeeper.StoreKey, auctiontypes.StoreKey, tokenminttypes.StoreKey,
			rewardstypes.StoreKey, feegrant.StoreKey, liquiditytypes.StoreKey, esmtypes.ModuleName, lendtypes.StoreKey,
		)
	)

	baseApp := baseapp.NewBaseApp(Name, logger, db, encoding.TxConfig.TxDecoder(), baseAppOptions...)
	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetVersion(version.Version)
	baseApp.SetInterfaceRegistry(encoding.InterfaceRegistry)

	app := &App{
		BaseApp:           baseApp,
		amino:             encoding.Amino,
		cdc:               encoding.Marshaler,
		interfaceRegistry: encoding.InterfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		mkeys:             mkeys,
	}

	app.ParamsKeeper = paramskeeper.NewKeeper(
		app.cdc,
		app.amino,
		app.keys[paramstypes.StoreKey],
		app.tkeys[paramstypes.TStoreKey],
	)

	//nolint:godox  //TODO: refactor this code
	app.ParamsKeeper.Subspace(authtypes.ModuleName)
	app.ParamsKeeper.Subspace(banktypes.ModuleName)
	app.ParamsKeeper.Subspace(stakingtypes.ModuleName)
	app.ParamsKeeper.Subspace(minttypes.ModuleName)
	app.ParamsKeeper.Subspace(distrtypes.ModuleName)
	app.ParamsKeeper.Subspace(slashingtypes.ModuleName)
	app.ParamsKeeper.Subspace(govtypes.ModuleName).
		WithKeyTable(govtypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(crisistypes.ModuleName)
	app.ParamsKeeper.Subspace(ibctransfertypes.ModuleName)
	app.ParamsKeeper.Subspace(ibchost.ModuleName)
	app.ParamsKeeper.Subspace(icahosttypes.SubModuleName)
	app.ParamsKeeper.Subspace(vaulttypes.ModuleName)
	app.ParamsKeeper.Subspace(assettypes.ModuleName)
	app.ParamsKeeper.Subspace(collectortypes.ModuleName)
	app.ParamsKeeper.Subspace(esmtypes.ModuleName)
	app.ParamsKeeper.Subspace(lendtypes.ModuleName)
	app.ParamsKeeper.Subspace(markettypes.ModuleName)
	app.ParamsKeeper.Subspace(liquidationtypes.ModuleName)
	app.ParamsKeeper.Subspace(lockertypes.ModuleName)
	app.ParamsKeeper.Subspace(bandoraclemoduletypes.ModuleName)
	app.ParamsKeeper.Subspace(wasmtypes.ModuleName)
	app.ParamsKeeper.Subspace(auctiontypes.ModuleName)
	app.ParamsKeeper.Subspace(tokenminttypes.ModuleName)
	app.ParamsKeeper.Subspace(liquiditytypes.ModuleName)
	app.ParamsKeeper.Subspace(rewardstypes.ModuleName)

	// set the BaseApp's parameter store
	baseApp.SetParamStore(
		app.ParamsKeeper.
			Subspace(baseapp.Paramspace).
			WithKeyTable(paramskeeper.ConsensusParamsKeyTable()),
	)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		app.cdc,
		app.keys[capabilitytypes.StoreKey],
		app.mkeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	var (
		scopedIBCKeeper        = app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
		scopedTransferKeeper   = app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
		scopedIBCOracleKeeper  = app.CapabilityKeeper.ScopeToModule(markettypes.ModuleName)
		scopedWasmKeeper       = app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
		scopedICAHostKeeper    = app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
		scopedBandoracleKeeper = app.CapabilityKeeper.ScopeToModule(bandoraclemoduletypes.ModuleName)
	)

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		app.cdc,
		app.keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		app.ModuleAccountsPermissions(),
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		app.cdc,
		app.keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		app.ModuleAccountAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		app.cdc,
		app.keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)
	app.MintKeeper = mintkeeper.NewKeeper(
		app.cdc,
		app.keys[minttypes.StoreKey],
		app.GetSubspace(minttypes.ModuleName),
		&stakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		app.cdc,
		app.keys[distrtypes.StoreKey],
		app.GetSubspace(distrtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		authtypes.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		app.cdc,
		app.keys[slashingtypes.StoreKey],
		&stakingKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		app.cdc,
		baseApp.MsgServiceRouter(),
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		app.keys[upgradetypes.StoreKey],
		app.cdc,
		homePath,
		app.BaseApp,
	)
	// register the staking hooks
	// NOTE: StakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
		),
	)

	// Create IBC Keeper
	app.IbcKeeper = ibckeeper.NewKeeper(
		app.cdc,
		app.keys[ibchost.StoreKey],
		app.GetSubspace(ibchost.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, app.keys[icahosttypes.StoreKey],
		app.GetSubspace(icahosttypes.SubModuleName),
		app.IbcKeeper.ChannelKeeper,
		&app.IbcKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)

	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)
	app.AssetKeeper = assetkeeper.NewKeeper(
		app.cdc,
		app.keys[assettypes.StoreKey],
		app.GetSubspace(assettypes.ModuleName),
		&app.Rewardskeeper,
		&app.VaultKeeper,
		&app.BandoracleKeeper,
	)

	app.LendKeeper = lendkeeper.NewKeeper(
		app.cdc,
		app.keys[lendtypes.StoreKey],
		app.keys[lendtypes.StoreKey],
		app.GetSubspace(lendtypes.ModuleName),
		app.BankKeeper,
		app.AccountKeeper,
		&app.AssetKeeper,
		&app.MarketKeeper,
		&app.EsmKeeper,
	)

	app.EsmKeeper = esmkeeper.NewKeeper(
		app.cdc,
		app.keys[esmtypes.StoreKey],
		app.keys[esmtypes.StoreKey],
		app.GetSubspace(esmtypes.ModuleName),
		&app.AssetKeeper,
		&app.VaultKeeper,
		app.BankKeeper,
		&app.MarketKeeper,
		&app.TokenmintKeeper,
		&app.CollectorKeeper,
	)

	app.VaultKeeper = vaultkeeper.NewKeeper(
		app.cdc,
		app.keys[vaulttypes.StoreKey],
		app.BankKeeper,
		&app.AssetKeeper,
		&app.MarketKeeper,
		&app.CollectorKeeper,
		&app.EsmKeeper,
		&app.TokenmintKeeper,
		&app.Rewardskeeper,
	)

	app.TokenmintKeeper = tokenmintkeeper.NewKeeper(
		app.cdc,
		app.keys[tokenminttypes.StoreKey],
		app.BankKeeper,
		&app.AssetKeeper,
	)

	app.BandoracleKeeper = bandoraclemodulekeeper.NewKeeper(
		appCodec,
		keys[bandoraclemoduletypes.StoreKey],
		keys[bandoraclemoduletypes.MemStoreKey],
		app.GetSubspace(bandoraclemoduletypes.ModuleName),
		app.IbcKeeper.ChannelKeeper,
		&app.IbcKeeper.PortKeeper,
		scopedBandoracleKeeper,
		&app.MarketKeeper,
		app.AssetKeeper,
	)
	bandoracleModule := bandoraclemodule.NewAppModule(
		appCodec,
		app.BandoracleKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		app.ScopedBandoracleKeeper,
		&app.IbcKeeper.PortKeeper,
		app.IbcKeeper.ChannelKeeper,
	)

	app.MarketKeeper = marketkeeper.NewKeeper(
		app.cdc,
		app.keys[markettypes.StoreKey],
		app.GetSubspace(markettypes.ModuleName),
		scopedIBCOracleKeeper,
		app.AssetKeeper,
		&app.BandoracleKeeper,
	)

	app.LiquidationKeeper = liquidationkeeper.NewKeeper(
		app.cdc,
		keys[liquidationtypes.StoreKey],
		keys[liquidationtypes.MemStoreKey],
		app.GetSubspace(liquidationtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&app.AssetKeeper,
		&app.VaultKeeper,
		&app.MarketKeeper,
		&app.AuctionKeeper,
		&app.EsmKeeper,
		&app.Rewardskeeper,
		&app.LendKeeper,
	)

	app.AuctionKeeper = auctionkeeper.NewKeeper(
		app.cdc,
		keys[auctiontypes.StoreKey],
		keys[auctiontypes.MemStoreKey],
		app.GetSubspace(auctiontypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&app.MarketKeeper,
		&app.LiquidationKeeper,
		&app.AssetKeeper,
		&app.VaultKeeper,
		&app.CollectorKeeper,
		&app.TokenmintKeeper,
		&app.EsmKeeper,
		&app.LendKeeper,
	)

	app.CollectorKeeper = collectorkeeper.NewKeeper(
		app.cdc,
		app.keys[collectortypes.StoreKey],
		app.keys[collectortypes.MemStoreKey],
		&app.AssetKeeper,
		&app.AuctionKeeper,
		&app.LockerKeeper,
		&app.Rewardskeeper,
		app.GetSubspace(collectortypes.ModuleName),
		app.BankKeeper,
	)

	// Create Transfer Keepers
	app.IbcTransferKeeper = ibctransferkeeper.NewKeeper(
		app.cdc,
		app.keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IbcKeeper.ChannelKeeper,
		app.IbcKeeper.ChannelKeeper,
		&app.IbcKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)

	app.LockerKeeper = lockerkeeper.NewKeeper(
		app.cdc,
		app.keys[lockertypes.StoreKey],
		app.GetSubspace(lockertypes.ModuleName),
		app.BankKeeper,
		&app.AssetKeeper,
		&app.CollectorKeeper,
		&app.EsmKeeper,
		&app.Rewardskeeper,
	)

	app.LiquidityKeeper = liquiditykeeper.NewKeeper(
		app.cdc,
		app.keys[liquiditytypes.StoreKey],
		app.GetSubspace(liquiditytypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&app.AssetKeeper,
		&app.MarketKeeper,
		&app.Rewardskeeper,
	)

	app.Rewardskeeper = rewardskeeper.NewKeeper(
		app.cdc,
		app.keys[rewardstypes.StoreKey],
		app.keys[rewardstypes.MemStoreKey],
		app.GetSubspace(rewardstypes.ModuleName),
		&app.LockerKeeper,
		&app.CollectorKeeper,
		&app.VaultKeeper,
		&app.AssetKeeper,
		app.BankKeeper,
		app.LiquidityKeeper,
		&app.MarketKeeper,
		&app.EsmKeeper,
		&app.LendKeeper,
	)

	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOptions)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}
	supportedFeatures := "iterator,staking,stargate,comdex"

	wasmOpts = append(cwasm.RegisterCustomPlugins(&app.LockerKeeper, &app.TokenmintKeeper, &app.AssetKeeper, &app.Rewardskeeper, &app.CollectorKeeper, &app.LiquidationKeeper, &app.AuctionKeeper, &app.EsmKeeper, &app.VaultKeeper), wasmOpts...)

	app.WasmKeeper = wasmkeeper.NewKeeper(
		app.cdc,
		keys[wasmtypes.StoreKey],
		app.GetSubspace(wasmtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.DistrKeeper,
		app.IbcKeeper.ChannelKeeper,
		&app.IbcKeeper.PortKeeper,
		scopedWasmKeeper,
		app.IbcTransferKeeper,
		baseApp.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		wasmOpts...,
	)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(assettypes.RouterKey, asset.NewUpdateAssetProposalHandler(app.AssetKeeper)).
		AddRoute(rewardstypes.RouterKey, rewards.NewAddRewardsProposalHandler(app.Rewardskeeper)).
		AddRoute(lendtypes.RouterKey, lend.NewLendHandler(app.LendKeeper)).
		AddRoute(bandoraclemoduletypes.RouterKey, bandoraclemodule.NewFetchPriceHandler(app.BandoracleKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(app.IbcKeeper.ClientKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IbcKeeper.ClientKeeper)).
		AddRoute(liquiditytypes.RouterKey, liquidity.NewLiquidityProposalHandler(app.LiquidityKeeper))

	if len(wasmEnabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, wasmEnabledProposals))
	}

	app.GovKeeper = govkeeper.NewKeeper(
		app.cdc,
		app.keys[govtypes.StoreKey],
		app.GetSubspace(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		govRouter,
	)

	var (
		evidenceRouter      = evidencetypes.NewRouter()
		ibcRouter           = ibcporttypes.NewRouter()
		transferModule      = ibctransfer.NewAppModule(app.IbcTransferKeeper)
		transferIBCModule   = ibctransfer.NewIBCModule(app.IbcTransferKeeper)
		oracleModule        = market.NewAppModule(app.cdc, app.MarketKeeper)
		bandOracleIBCModule = bandoraclemodule.NewIBCModule(app.BandoracleKeeper)
	)

	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferIBCModule)
	ibcRouter.AddRoute(bandoraclemoduletypes.ModuleName, bandOracleIBCModule)
	ibcRouter.AddRoute(wasm.ModuleName, wasm.NewIBCHandler(app.WasmKeeper, app.IbcKeeper.ChannelKeeper))
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)
	app.IbcKeeper.SetRouter(ibcRouter)
	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	app.EvidenceKeeper = *evidencekeeper.NewKeeper(
		app.cdc,
		app.keys[evidencetypes.StoreKey],
		&app.StakingKeeper,
		app.SlashingKeeper,
	)
	app.EvidenceKeeper.SetRouter(evidenceRouter)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOptions.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx, encoding.TxConfig),
		auth.NewAppModule(app.cdc, app.AccountKeeper, nil),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(app.cdc, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(app.cdc, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(app.cdc, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(app.cdc, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(app.cdc, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(app.cdc, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(app.cdc, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		authzmodule.NewAppModule(app.cdc, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IbcKeeper),
		ica.NewAppModule(nil, &app.ICAHostKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		asset.NewAppModule(app.cdc, app.AssetKeeper),
		vault.NewAppModule(app.cdc, app.VaultKeeper),
		oracleModule,
		bandoracleModule,
		liquidation.NewAppModule(app.cdc, app.LiquidationKeeper, app.AccountKeeper, app.BankKeeper),
		locker.NewAppModule(app.cdc, app.LockerKeeper, app.AccountKeeper, app.BankKeeper),
		collector.NewAppModule(app.cdc, app.CollectorKeeper, app.AccountKeeper, app.BankKeeper),
		esm.NewAppModule(app.cdc, app.EsmKeeper, app.AccountKeeper, app.BankKeeper),
		lend.NewAppModule(app.cdc, app.LendKeeper, app.AccountKeeper, app.BankKeeper),
		wasm.NewAppModule(app.cdc, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		auction.NewAppModule(app.cdc, app.AuctionKeeper, app.AccountKeeper, app.BankKeeper),
		tokenmint.NewAppModule(app.cdc, app.TokenmintKeeper, app.AccountKeeper, app.BankKeeper),
		liquidity.NewAppModule(app.cdc, app.LiquidityKeeper, app.AccountKeeper, app.BankKeeper),
		rewards.NewAppModule(app.cdc, app.Rewardskeeper, app.AccountKeeper, app.BankKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName, minttypes.ModuleName, distrtypes.ModuleName, slashingtypes.ModuleName,
		evidencetypes.ModuleName, stakingtypes.ModuleName, ibchost.ModuleName, ibctransfertypes.ModuleName, icatypes.ModuleName,
		bandoraclemoduletypes.ModuleName, markettypes.ModuleName, lockertypes.ModuleName,
		crisistypes.ModuleName, genutiltypes.ModuleName, authtypes.ModuleName, capabilitytypes.ModuleName,
		authz.ModuleName, transferModule.Name(), assettypes.ModuleName, collectortypes.ModuleName, vaulttypes.ModuleName,
		liquidationtypes.ModuleName, auctiontypes.ModuleName, tokenminttypes.ModuleName,
		vesting.AppModuleBasic{}.Name(), paramstypes.ModuleName, wasmtypes.ModuleName, banktypes.ModuleName,
		govtypes.ModuleName, rewardstypes.ModuleName, liquiditytypes.ModuleName, lendtypes.ModuleName, esmtypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName, govtypes.ModuleName, stakingtypes.ModuleName,
		minttypes.ModuleName, bandoraclemoduletypes.ModuleName, markettypes.ModuleName, lockertypes.ModuleName,
		distrtypes.ModuleName, genutiltypes.ModuleName, vesting.AppModuleBasic{}.Name(), evidencetypes.ModuleName, ibchost.ModuleName,
		icatypes.ModuleName, vaulttypes.ModuleName, liquidationtypes.ModuleName, auctiontypes.ModuleName, tokenminttypes.ModuleName,
		wasmtypes.ModuleName, authtypes.ModuleName, slashingtypes.ModuleName, authz.ModuleName,
		paramstypes.ModuleName, capabilitytypes.ModuleName, upgradetypes.ModuleName, transferModule.Name(), lendtypes.ModuleName,
		assettypes.ModuleName, collectortypes.ModuleName, banktypes.ModuleName, rewardstypes.ModuleName, liquiditytypes.ModuleName, esmtypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		assettypes.ModuleName,
		collectortypes.ModuleName,
		esmtypes.ModuleName,
		lendtypes.ModuleName,
		vaulttypes.ModuleName,
		tokenminttypes.ModuleName,
		bandoraclemoduletypes.ModuleName,
		markettypes.ModuleName,
		liquidationtypes.ModuleName,
		auctiontypes.ModuleName,
		lockertypes.StoreKey,
		wasmtypes.ModuleName,
		authz.ModuleName,
		vesting.AppModuleBasic{}.Name(),
		upgradetypes.ModuleName,
		paramstypes.ModuleName,
		liquiditytypes.ModuleName,
		rewardstypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encoding.Amino)
	app.configurator = module.NewConfigurator(app.cdc, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)
	app.registerUpgradeHandlers()
	// initialize stores
	app.MountKVStores(app.keys)
	app.MountTransientStores(app.tkeys)
	app.MountMemoryStores(app.mkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeegrantKeeper,
				SignModeHandler: encoding.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			wasmConfig:        wasmConfig,
			txCounterStoreKey: app.GetKey(wasm.StoreKey),
			IBCChannelKeeper:  app.IbcKeeper,
		},
	)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}

		// Initialize and seal the capability keeper so all persistent capabilities
		// are loaded in-memory and prevent any further modules from creating scoped
		// sub-keepers.
		// This must be done during creation of baseapp rather than in InitChain so
		// that in-memory capabilities get regenerated on app restart.
		// Note that since this reads from the store, we can only perform it when
		// `loadLatest` is set to true.
		ctx := app.BaseApp.NewUncachedContext(true, tmprototypes.Header{})
		app.CapabilityKeeper.InitMemStore(ctx)
		app.CapabilityKeeper.Seal()
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedIBCTransferKeeper = scopedTransferKeeper
	app.ScopedIBCOracleKeeper = scopedIBCOracleKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper
	app.ScopedBandoracleKeeper = scopedBandoracleKeeper

	app.ScopedWasmKeeper = scopedWasmKeeper
	return app
}

// Name returns the name of the App
func (a *App) Name() string { return a.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (a *App) BeginBlocker(ctx sdk.Context, req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return a.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block.
func (a *App) EndBlocker(ctx sdk.Context, req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return a.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization.
func (a *App) InitChainer(ctx sdk.Context, req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	var state GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &state); err != nil {
		panic(err)
	}
	a.UpgradeKeeper.SetModuleVersionMap(ctx, a.mm.GetVersionMap())
	return a.mm.InitGenesis(ctx, a.cdc, state)
}

// LoadHeight loads a particular height.
func (a *App) LoadHeight(height int64) error {
	return a.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (a *App) ModuleAccountAddrs() map[string]bool {
	accounts := make(map[string]bool)

	names := make([]string, 0)
	for name := range a.ModuleAccountsPermissions() {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		accounts[authtypes.NewModuleAddress(name).String()] = true
	}

	return accounts
}

// LegacyAmino returns App's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (a *App) LegacyAmino() *codec.LegacyAmino {
	return a.amino
}

// AppCodec returns App's codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (a *App) AppCodec() codec.BinaryCodec {
	return a.cdc
}

// InterfaceRegistry returns Gaia's InterfaceRegistry.
func (a *App) InterfaceRegistry() codectypes.InterfaceRegistry {
	return a.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (a *App) GetKey(storeKey string) *sdk.KVStoreKey {
	return a.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (a *App) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return a.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (a *App) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return a.mkeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (a *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := a.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (a *App) RegisterAPIRoutes(server *api.Server, apiConfig serverconfig.APIConfig) {
	ctx := server.ClientCtx
	rpc.RegisterRoutes(ctx, server.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(ctx, server.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(ctx, server.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(ctx, server.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(ctx, server.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(ctx, server.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(ctx, server.Router)
	}
}

// RegisterSwaggerAPI registers swagger route with API Server.
func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticServer))
	rtr.PathPrefix("/swagger/").Handler(staticServer)
}

// RegisterTxService implements the Application.RegisterTxService method.
func (a *App) RegisterTxService(ctx client.Context) {
	authtx.RegisterTxService(a.BaseApp.GRPCQueryRouter(), ctx, a.BaseApp.Simulate, a.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (a *App) RegisterTendermintService(ctx client.Context) {
	tmservice.RegisterTendermintService(a.BaseApp.GRPCQueryRouter(), ctx, a.interfaceRegistry)
}

func (a *App) ModuleAccountsPermissions() map[string][]string {
	return map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		collectortypes.ModuleName:      {authtypes.Burner, authtypes.Staking},
		vaulttypes.ModuleName:          {authtypes.Minter, authtypes.Burner},
		lendtypes.ModuleName:           {authtypes.Minter, authtypes.Burner},
		tokenminttypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
		lendtypes.ModuleAcc1:           {authtypes.Minter, authtypes.Burner},
		lendtypes.ModuleAcc2:           {authtypes.Minter, authtypes.Burner},
		lendtypes.ModuleAcc3:           {authtypes.Minter, authtypes.Burner},
		liquidationtypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		auctiontypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		lockertypes.ModuleName:         {authtypes.Minter, authtypes.Burner},
		esmtypes.ModuleName:            {authtypes.Burner},
		wasm.ModuleName:                {authtypes.Burner},
		liquiditytypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
		rewardstypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:            nil,
	}
}

func (a *App) registerUpgradeHandlers() {
	a.UpgradeKeeper.SetUpgradeHandler(
		tv4_0_0.UpgradeNameV4_4_0,
		tv4_0_0.CreateUpgradeHandlerV440(a.mm, a.configurator, a.LendKeeper, a.LiquidationKeeper, a.AuctionKeeper),
	)

	// When a planned update height is reached, the old binary will panic
	// writing on disk the height and name of the update that triggered it
	// This will read that value, and execute the preparations for the upgrade.
	upgradeInfo, err := a.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	var storeUpgrades *storetypes.StoreUpgrades

	storeUpgrades = upgradeHandlers(upgradeInfo, a, storeUpgrades)

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		a.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}

func upgradeHandlers(upgradeInfo storetypes.UpgradeInfo, a *App, storeUpgrades *storetypes.StoreUpgrades) *storetypes.StoreUpgrades {
	switch {
	case upgradeInfo.Name == tv1_0_0.UpgradeName && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		// prepare store for testnet upgrade v1.0.0
		storeUpgrades = &storetypes.StoreUpgrades{
			Added:   []string{authz.ModuleName},
			Deleted: []string{"asset", "liquidity", "oracle", "vault"},
		}
	case upgradeInfo.Name == tv2_0_0.UpgradeName && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		// prepare store for testnet upgrade v2.0.0
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{
				assettypes.ModuleName,
				auctiontypes.ModuleName,
				bandoraclemoduletypes.ModuleName,
				collectortypes.ModuleName,
				esmtypes.ModuleName,
				lendtypes.ModuleName,
				liquidationtypes.ModuleName,
				liquiditytypes.ModuleName,
				lockertypes.ModuleName,
				markettypes.ModuleName,
				rewardstypes.ModuleName,
				tokenminttypes.ModuleName,
				vaulttypes.ModuleName,
				feegrant.ModuleName,
				icacontrollertypes.StoreKey,
				icahosttypes.StoreKey,
			},
		}
	case upgradeInfo.Name == tv2_0_0.UpgradeNameV2 && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		// prepare store for testnet upgrade v2.1.0
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{
				icacontrollertypes.StoreKey,
				icahosttypes.StoreKey,
			},
		}
	case upgradeInfo.Name == tv3_0_0.UpgradeName && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
	case upgradeInfo.Name == tv4_0_0.UpgradeName && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		// delete deprecated kv store instead of migrating as this is only for testnet.
		// won't be used on main net migration and to make store name similar with V1 appended for all modules
		storeUpgrades = &storetypes.StoreUpgrades{
			Deleted: []string{"vaultv1", "rewards", "collector", "locker", "lend", "auction", "liquidation", "esm"},
			Added: []string{
				vaulttypes.ModuleName, rewardstypes.ModuleName, liquidationtypes.ModuleName,
				collectortypes.ModuleName, lockertypes.ModuleName, lendtypes.ModuleName, auctiontypes.ModuleName, esmtypes.ModuleName,
			},
		}
	case upgradeInfo.Name == tv4_0_0.UpgradeNameV4_1_0 && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		storeUpgrades = &storetypes.StoreUpgrades{}
	case upgradeInfo.Name == tv4_0_0.UpgradeNameV4_2_0 && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		storeUpgrades = &storetypes.StoreUpgrades{}
	case upgradeInfo.Name == tv4_0_0.UpgradeNameV4_3_0 && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		storeUpgrades = &storetypes.StoreUpgrades{}
	case upgradeInfo.Name == tv4_0_0.UpgradeNameV4_4_0 && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		storeUpgrades = &storetypes.StoreUpgrades{}

	// prepare store for main net upgrade v5.0.0
	case upgradeInfo.Name == mv5.UpgradeName && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{
				assettypes.ModuleName,
				auctiontypes.ModuleName,
				bandoraclemoduletypes.ModuleName,
				collectortypes.ModuleName,
				esmtypes.ModuleName,
				liquidationtypes.ModuleName,
				liquiditytypes.ModuleName,
				lockertypes.ModuleName,
				markettypes.ModuleName,
				rewardstypes.ModuleName,
				tokenminttypes.ModuleName,
				vaulttypes.ModuleName,
				feegrant.ModuleName,
				icacontrollertypes.StoreKey,
				icahosttypes.StoreKey,
				authz.ModuleName,
			},
		}
	}

	return storeUpgrades
}
