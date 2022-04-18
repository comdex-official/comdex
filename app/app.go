package app

import (
	"io"
	"os"
	"path/filepath"

	rewardsclient "github.com/comdex-official/comdex/x/rewards/client"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"

	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	ibcclient "github.com/cosmos/ibc-go/v2/modules/core/02-client"

	authvestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
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
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
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
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	freegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ibctransfer "github.com/cosmos/ibc-go/v2/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v2/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v2/modules/core"
	ibcporttypes "github.com/cosmos/ibc-go/v2/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v2/modules/core/24-host"
	ibctypes "github.com/cosmos/ibc-go/v2/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v2/modules/core/keeper"

	"github.com/spf13/cast"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/comdex-official/comdex/x/asset"
	assetclient "github.com/comdex-official/comdex/x/asset/client"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/auction"
	auctionkeeper "github.com/comdex-official/comdex/x/auction/keeper"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	bandoraclemodule "github.com/comdex-official/comdex/x/bandoracle"
	bandoraclemoduleclient "github.com/comdex-official/comdex/x/bandoracle/client"
	bandoraclemodulekeeper "github.com/comdex-official/comdex/x/bandoracle/keeper"
	bandoraclemoduletypes "github.com/comdex-official/comdex/x/bandoracle/types"
	"github.com/comdex-official/comdex/x/liquidation"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	"github.com/comdex-official/comdex/x/liquidity"
	liquiditykeeper "github.com/comdex-official/comdex/x/liquidity/keeper"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"
	"github.com/comdex-official/comdex/x/market"
	marketkeeper "github.com/comdex-official/comdex/x/market/keeper"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	"github.com/comdex-official/comdex/x/rewards"
	rewardskeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"

	"github.com/comdex-official/comdex/x/exchange"
	exchangekeeper "github.com/comdex-official/comdex/x/exchange/keeper"
	exchangetypes "github.com/comdex-official/comdex/x/exchange/types"

	"github.com/comdex-official/comdex/x/vault"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
)

const (
	AccountAddressPrefix = "comdex"
	Name                 = "comdex"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

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
		gov.NewAppModuleBasic(
			append(
				assetclient.AddAssetsHandler,
				bandoraclemoduleclient.AddFetchPriceHandler,
				paramsclient.ProposalHandler,
				distrclient.ProposalHandler,
				upgradeclient.ProposalHandler,
				upgradeclient.CancelProposalHandler,
				rewardsclient.AddNewMintingRewardsHandler,
				rewardsclient.DisableMintingRewardsHandler,
			)...,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		vault.AppModuleBasic{},
		asset.AppModuleBasic{},
		liquidity.AppModuleBasic{},
		asset.AppModuleBasic{},
		market.AppModuleBasic{},
		bandoraclemodule.AppModuleBasic{},
		liquidation.AppModuleBasic{},
		auction.AppModuleBasic{},
		rewards.AppModuleBasic{},
		exchange.AppModuleBasic{},
		wasm.AppModuleBasic{},
	)
)

var (
	_ servertypes.Application = (*App)(nil)
)

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
	accountKeeper     authkeeper.AccountKeeper
	freegrantKeeper   freegrantkeeper.Keeper
	bankKeeper        bankkeeper.Keeper
	capabilityKeeper  *capabilitykeeper.Keeper
	stakingKeeper     stakingkeeper.Keeper
	slashingKeeper    slashingkeeper.Keeper
	mintKeeper        mintkeeper.Keeper
	distrKeeper       distrkeeper.Keeper
	govKeeper         govkeeper.Keeper
	crisisKeeper      crisiskeeper.Keeper
	upgradeKeeper     upgradekeeper.Keeper
	paramsKeeper      paramskeeper.Keeper
	ibcKeeper         *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	evidenceKeeper    evidencekeeper.Keeper
	ibcTransferKeeper ibctransferkeeper.Keeper

	// make scoped keepers public for test purposes
	scopedIBCKeeper         capabilitykeeper.ScopedKeeper
	scopedIBCTransferKeeper capabilitykeeper.ScopedKeeper
	scopedIBCOracleKeeper   capabilitykeeper.ScopedKeeper
	scopedBandoracleKeeper  capabilitykeeper.ScopedKeeper

	BandoracleKeeper  bandoraclemodulekeeper.Keeper
	assetKeeper       assetkeeper.Keeper
	vaultKeeper       vaultkeeper.Keeper
	liquidityKeeper   liquiditykeeper.Keeper
	marketKeeper      marketkeeper.Keeper
	liquidationKeeper liquidationkeeper.Keeper
	auctionKeeper     auctionkeeper.Keeper
	rewardsKeeper     rewardskeeper.Keeper
	exchangeKeeper    exchangekeeper.Keeper
	scopedWasmKeeper  capabilitykeeper.ScopedKeeper

	wasmKeeper wasm.Keeper
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
			govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
			evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
			vaulttypes.StoreKey, liquiditytypes.StoreKey, assettypes.StoreKey,
			markettypes.StoreKey, bandoraclemoduletypes.StoreKey, liquidationtypes.StoreKey,
			auctiontypes.StoreKey, wasm.StoreKey, rewardstypes.StoreKey, exchangetypes.StoreKey,
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

	app.paramsKeeper = paramskeeper.NewKeeper(
		app.cdc,
		app.amino,
		app.keys[paramstypes.StoreKey],
		app.tkeys[paramstypes.TStoreKey],
	)

	//TODO: refactor this code
	app.paramsKeeper.Subspace(authtypes.ModuleName)
	app.paramsKeeper.Subspace(banktypes.ModuleName)
	app.paramsKeeper.Subspace(stakingtypes.ModuleName)
	app.paramsKeeper.Subspace(minttypes.ModuleName)
	app.paramsKeeper.Subspace(distrtypes.ModuleName)
	app.paramsKeeper.Subspace(slashingtypes.ModuleName)
	app.paramsKeeper.Subspace(govtypes.ModuleName).
		WithKeyTable(govtypes.ParamKeyTable())
	app.paramsKeeper.Subspace(crisistypes.ModuleName)
	app.paramsKeeper.Subspace(liquiditytypes.ModuleName)
	app.paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	app.paramsKeeper.Subspace(ibchost.ModuleName)
	app.paramsKeeper.Subspace(vaulttypes.ModuleName)
	app.paramsKeeper.Subspace(assettypes.ModuleName)
	app.paramsKeeper.Subspace(markettypes.ModuleName)
	app.paramsKeeper.Subspace(bandoraclemoduletypes.ModuleName)
	app.paramsKeeper.Subspace(liquidationtypes.ModuleName)
	app.paramsKeeper.Subspace(auctiontypes.ModuleName)
	app.paramsKeeper.Subspace(rewardstypes.ModuleName)
	app.paramsKeeper.Subspace(exchangetypes.ModuleName)
	app.paramsKeeper.Subspace(wasmtypes.ModuleName)

	// set the BaseApp's parameter store
	baseApp.SetParamStore(
		app.paramsKeeper.
			Subspace(baseapp.Paramspace).
			WithKeyTable(paramskeeper.ConsensusParamsKeyTable()),
	)

	// add capability keeper and ScopeToModule for ibc module
	app.capabilityKeeper = capabilitykeeper.NewKeeper(
		app.cdc,
		app.keys[capabilitytypes.StoreKey],
		app.mkeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	var (
		scopedIBCKeeper       = app.capabilityKeeper.ScopeToModule(ibchost.ModuleName)
		scopedTransferKeeper  = app.capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
		scopedIBCOracleKeeper = app.capabilityKeeper.ScopeToModule(markettypes.ModuleName)
		scopedWasmKeeper      = app.capabilityKeeper.ScopeToModule(wasm.ModuleName)
	)

	// add keepers
	app.accountKeeper = authkeeper.NewAccountKeeper(
		app.cdc,
		app.keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		app.ModuleAccountsPermissions(),
	)
	app.bankKeeper = bankkeeper.NewBaseKeeper(
		app.cdc,
		app.keys[banktypes.StoreKey],
		app.accountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		app.ModuleAccountAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		app.cdc,
		app.keys[stakingtypes.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)
	app.mintKeeper = mintkeeper.NewKeeper(
		app.cdc,
		app.keys[minttypes.StoreKey],
		app.GetSubspace(minttypes.ModuleName),
		&stakingKeeper,
		app.accountKeeper,
		app.bankKeeper,
		authtypes.FeeCollectorName,
	)
	app.distrKeeper = distrkeeper.NewKeeper(
		app.cdc,
		app.keys[distrtypes.StoreKey],
		app.GetSubspace(distrtypes.ModuleName),
		app.accountKeeper,
		app.bankKeeper,
		&stakingKeeper,
		authtypes.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)
	app.slashingKeeper = slashingkeeper.NewKeeper(
		app.cdc,
		app.keys[slashingtypes.StoreKey],
		&stakingKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)
	app.crisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		app.bankKeeper,
		authtypes.FeeCollectorName,
	)
	app.upgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		app.keys[upgradetypes.StoreKey],
		app.cdc,
		homePath,
		app.BaseApp,
	)
	app.registerUpgradeHandlers()
	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.distrKeeper.Hooks(),
			app.slashingKeeper.Hooks(),
		),
	)

	// Create IBC Keeper
	app.ibcKeeper = ibckeeper.NewKeeper(
		app.cdc,
		app.keys[ibchost.StoreKey],
		app.GetSubspace(ibchost.ModuleName),
		app.stakingKeeper,
		app.upgradeKeeper,
		scopedIBCKeeper,
	)

	app.assetKeeper = assetkeeper.NewKeeper(
		app.cdc,
		app.keys[assettypes.StoreKey],
		app.GetSubspace(assettypes.ModuleName),
		&app.marketKeeper,
	)
	app.vaultKeeper = vaultkeeper.NewKeeper(
		app.cdc,
		app.keys[vaulttypes.StoreKey],
		app.bankKeeper,
		app.accountKeeper,
		&app.assetKeeper,
		&app.marketKeeper,
		&app.BandoracleKeeper,
	)

	app.liquidityKeeper = liquiditykeeper.NewKeeper(
		app.cdc,
		app.keys[liquiditytypes.StoreKey],
		app.GetSubspace(liquiditytypes.ModuleName),
		app.bankKeeper,
		app.accountKeeper,
		app.distrKeeper,
	)

	scopedBandoracleKeeper := app.capabilityKeeper.ScopeToModule(bandoraclemoduletypes.ModuleName)
	app.scopedBandoracleKeeper = scopedBandoracleKeeper

	app.BandoracleKeeper = bandoraclemodulekeeper.NewKeeper(
		appCodec,
		keys[bandoraclemoduletypes.StoreKey],
		keys[bandoraclemoduletypes.MemStoreKey],
		app.GetSubspace(bandoraclemoduletypes.ModuleName),
		app.ibcKeeper.ChannelKeeper,
		&app.ibcKeeper.PortKeeper,
		scopedBandoracleKeeper,
		&app.marketKeeper,
		app.assetKeeper,
	)
	bandoracleModule := bandoraclemodule.NewAppModule(
		appCodec,
		app.BandoracleKeeper,
		app.accountKeeper,
		app.bankKeeper,
		app.scopedBandoracleKeeper,
		&app.ibcKeeper.PortKeeper,
		app.ibcKeeper.ChannelKeeper,
	)

	app.marketKeeper = *marketkeeper.NewKeeper(
		app.cdc,
		app.keys[markettypes.StoreKey],
		app.GetSubspace(markettypes.ModuleName),
		scopedIBCOracleKeeper,
		app.assetKeeper,
		&app.BandoracleKeeper,
	)

	app.liquidationKeeper = *liquidationkeeper.NewKeeper(
		app.cdc,
		keys[liquidationtypes.StoreKey],
		keys[liquidationtypes.MemStoreKey],
		app.GetSubspace(liquidationtypes.ModuleName),
		app.accountKeeper,
		app.bankKeeper,
		&app.assetKeeper,
		&app.vaultKeeper,
		&app.marketKeeper,
		&app.auctionKeeper,
	)

	app.auctionKeeper = *auctionkeeper.NewKeeper(
		app.cdc,
		keys[auctiontypes.StoreKey],
		keys[auctiontypes.MemStoreKey],
		app.GetSubspace(auctiontypes.ModuleName),
		app.accountKeeper,
		app.bankKeeper,
		&app.assetKeeper,
		&app.vaultKeeper,
		&app.marketKeeper,
		&app.liquidationKeeper,
	)

	app.rewardsKeeper = *rewardskeeper.NewKeeper(
		app.cdc,
		keys[rewardstypes.StoreKey],
		keys[rewardstypes.MemStoreKey],
		app.GetSubspace(rewardstypes.ModuleName),
		app.bankKeeper,
		&app.assetKeeper,
		&app.vaultKeeper,
		&app.marketKeeper,
	)

	app.exchangeKeeper = *exchangekeeper.NewKeeper(
		app.cdc,
		app.keys[assettypes.StoreKey],
		app.GetSubspace(assettypes.ModuleName),
	)

	// Create Transfer Keepers
	app.ibcTransferKeeper = ibctransferkeeper.NewKeeper(
		app.cdc,
		app.keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.ibcKeeper.ChannelKeeper,
		&app.ibcKeeper.PortKeeper,
		app.accountKeeper,
		app.bankKeeper,
		scopedTransferKeeper,
	)

	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOptions)
	supportedFeatures := "iterator,staking,stargate"

	app.wasmKeeper = wasmkeeper.NewKeeper(
		app.cdc,
		keys[wasmtypes.StoreKey],
		app.GetSubspace(wasmtypes.ModuleName),
		app.accountKeeper,
		app.bankKeeper,
		app.stakingKeeper,
		app.distrKeeper,
		app.ibcKeeper.ChannelKeeper,
		&app.ibcKeeper.PortKeeper,
		scopedWasmKeeper,
		app.ibcTransferKeeper,
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
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(app.ibcKeeper.ClientKeeper)).
		AddRoute(assettypes.RouterKey, asset.NewUpdateAssetProposalHandler(app.assetKeeper)).
		AddRoute(bandoraclemoduletypes.RouterKey, bandoraclemodule.NewFetchPriceHandler(app.BandoracleKeeper)).
		AddRoute(rewardstypes.RouterKey, rewards.NewRewardsProposalHandler(app.rewardsKeeper)).
		AddRoute(exchangetypes.RouterKey, exchange.ExchnageProposalHandler(app.exchangeKeeper))

	app.govKeeper = govkeeper.NewKeeper(
		app.cdc,
		app.keys[govtypes.StoreKey],
		app.GetSubspace(govtypes.ModuleName),
		app.accountKeeper,
		app.bankKeeper,
		&stakingKeeper,
		govRouter,
	)

	var (
		evidenceRouter = evidencetypes.NewRouter()
		ibcRouter      = ibcporttypes.NewRouter()
		transferModule = ibctransfer.NewAppModule(app.ibcTransferKeeper)
		oracleModule   = market.NewAppModule(app.cdc, app.marketKeeper)
	)

	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	ibcRouter.AddRoute(bandoraclemoduletypes.ModuleName, bandoracleModule)
	app.ibcKeeper.SetRouter(ibcRouter)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	app.evidenceKeeper = *evidencekeeper.NewKeeper(
		app.cdc,
		app.keys[evidencetypes.StoreKey],
		&app.stakingKeeper,
		app.slashingKeeper,
	)
	app.evidenceKeeper.SetRouter(evidenceRouter)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOptions.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx, encoding.TxConfig),
		auth.NewAppModule(app.cdc, app.accountKeeper, nil),
		vesting.NewAppModule(app.accountKeeper, app.bankKeeper),
		bank.NewAppModule(app.cdc, app.bankKeeper, app.accountKeeper),
		capability.NewAppModule(app.cdc, *app.capabilityKeeper),
		crisis.NewAppModule(&app.crisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(app.cdc, app.govKeeper, app.accountKeeper, app.bankKeeper),
		mint.NewAppModule(app.cdc, app.mintKeeper, app.accountKeeper),
		slashing.NewAppModule(app.cdc, app.slashingKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		distr.NewAppModule(app.cdc, app.distrKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		staking.NewAppModule(app.cdc, app.stakingKeeper, app.accountKeeper, app.bankKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
		ibc.NewAppModule(app.ibcKeeper),
		params.NewAppModule(app.paramsKeeper),
		transferModule,
		asset.NewAppModule(app.cdc, app.assetKeeper),
		vault.NewAppModule(app.cdc, app.vaultKeeper),
		liquidity.NewAppModule(app.cdc, app.liquidityKeeper, app.accountKeeper, app.bankKeeper, app.distrKeeper),
		asset.NewAppModule(app.cdc, app.assetKeeper),
		oracleModule,
		bandoracleModule,
		liquidation.NewAppModule(app.cdc, app.liquidationKeeper, app.accountKeeper, app.bankKeeper),
		auction.NewAppModule(app.cdc, app.auctionKeeper, app.accountKeeper, app.bankKeeper),
		rewards.NewAppModule(app.cdc, app.rewardsKeeper, app.accountKeeper, app.bankKeeper),
		exchange.NewAppModule(app.cdc, app.exchangeKeeper),
		wasm.NewAppModule(app.cdc, &app.wasmKeeper, app.stakingKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName, minttypes.ModuleName, distrtypes.ModuleName, slashingtypes.ModuleName,
		evidencetypes.ModuleName, stakingtypes.ModuleName, liquiditytypes.ModuleName, ibchost.ModuleName,
		bandoraclemoduletypes.ModuleName, markettypes.ModuleName, liquidationtypes.ModuleName,
		auctiontypes.ModuleName, crisistypes.ModuleName, genutiltypes.ModuleName, authtypes.ModuleName,
		capabilitytypes.ModuleName, transferModule.Name(), assettypes.ModuleName, vaulttypes.ModuleName,
		vesting.AppModuleBasic{}.Name(), paramstypes.ModuleName, wasmtypes.ModuleName, banktypes.ModuleName,
		govtypes.ModuleName, rewardstypes.ModuleName, exchangetypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName, govtypes.ModuleName, stakingtypes.ModuleName, liquiditytypes.ModuleName,
		bandoraclemoduletypes.ModuleName, minttypes.ModuleName,
		distrtypes.ModuleName, genutiltypes.ModuleName, vesting.AppModuleBasic{}.Name(), evidencetypes.ModuleName, ibchost.ModuleName,
		vaulttypes.ModuleName, wasmtypes.ModuleName, authtypes.ModuleName, slashingtypes.ModuleName, paramstypes.ModuleName,
		markettypes.ModuleName, capabilitytypes.ModuleName, upgradetypes.ModuleName, transferModule.Name(),
		assettypes.ModuleName, banktypes.ModuleName, liquidationtypes.ModuleName, auctiontypes.ModuleName, rewardstypes.ModuleName,
		exchangetypes.ModuleName,
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
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		liquiditytypes.ModuleName,
		ibctransfertypes.ModuleName,
		assettypes.ModuleName,
		vaulttypes.ModuleName,
		bandoraclemoduletypes.ModuleName,
		markettypes.ModuleName,
		liquidationtypes.ModuleName,
		auctiontypes.ModuleName,
		rewardstypes.ModuleName,
		exchangetypes.ModuleName,
		paramstypes.ModuleName,
		authvestingtypes.ModuleName,
		upgradetypes.ModuleName,
		wasmtypes.ModuleName,
		vesting.AppModuleBasic{}.Name(),
		paramstypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encoding.Amino)
	app.configurator = module.NewConfigurator(app.cdc, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)
	// initialize stores
	app.MountKVStores(app.keys)
	app.MountTransientStores(app.tkeys)
	app.MountMemoryStores(app.mkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.accountKeeper,
			BankKeeper:      app.bankKeeper,
			FeegrantKeeper:  app.freegrantKeeper,
			SignModeHandler: encoding.TxConfig.SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
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
		app.capabilityKeeper.InitMemStore(ctx)
		app.capabilityKeeper.Seal()
	}

	app.scopedIBCKeeper = scopedIBCKeeper
	app.scopedIBCTransferKeeper = scopedTransferKeeper
	app.scopedIBCOracleKeeper = scopedIBCOracleKeeper

	app.scopedWasmKeeper = scopedWasmKeeper
	return app
}

// Name returns the name of the App
func (a *App) Name() string { return a.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (a *App) BeginBlocker(ctx sdk.Context, req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return a.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (a *App) EndBlocker(ctx sdk.Context, req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return a.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (a *App) InitChainer(ctx sdk.Context, req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	var state GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &state); err != nil {
		panic(err)
	}
	a.upgradeKeeper.SetModuleVersionMap(ctx, a.mm.GetVersionMap())
	return a.mm.InitGenesis(ctx, a.cdc, state)
}

// LoadHeight loads a particular height
func (a *App) LoadHeight(height int64) error {
	return a.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (a *App) ModuleAccountAddrs() map[string]bool {
	accounts := make(map[string]bool)
	for name := range a.ModuleAccountsPermissions() {
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

// InterfaceRegistry returns Gaia's InterfaceRegistry
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
	subspace, _ := a.paramsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (a *App) RegisterAPIRoutes(server *api.Server, _ serverconfig.APIConfig) {
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
		vaulttypes.ModuleName:          {authtypes.Minter, authtypes.Burner},
		liquiditytypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
		liquidationtypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		auctiontypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		rewardstypes.ModuleName:        nil,
		exchangetypes.ModuleName:       nil,
		wasm.ModuleName:                {authtypes.Burner},
	}
}
func (app *App) registerUpgradeHandlers() {
	app.upgradeKeeper.SetUpgradeHandler("v0.1.0", func(ctx sdk.Context, plan upgradetypes.Plan, _ module.VersionMap) (module.VersionMap, error) {
		// 1st-time running in-store migrations, using 1 as fromVersion to
		// avoid running InitGenesis.
		fromVM := map[string]uint64{
			authtypes.ModuleName:        auth.AppModule{}.ConsensusVersion(),
			banktypes.ModuleName:        bank.AppModule{}.ConsensusVersion(),
			capabilitytypes.ModuleName:  capability.AppModule{}.ConsensusVersion(),
			crisistypes.ModuleName:      crisis.AppModule{}.ConsensusVersion(),
			distrtypes.ModuleName:       distr.AppModule{}.ConsensusVersion(),
			evidencetypes.ModuleName:    evidence.AppModule{}.ConsensusVersion(),
			govtypes.ModuleName:         gov.AppModule{}.ConsensusVersion(),
			minttypes.ModuleName:        mint.AppModule{}.ConsensusVersion(),
			paramstypes.ModuleName:      params.AppModule{}.ConsensusVersion(),
			slashingtypes.ModuleName:    slashing.AppModule{}.ConsensusVersion(),
			stakingtypes.ModuleName:     staking.AppModule{}.ConsensusVersion(),
			upgradetypes.ModuleName:     upgrade.AppModule{}.ConsensusVersion(),
			vestingtypes.ModuleName:     vesting.AppModule{}.ConsensusVersion(),
			ibctypes.ModuleName:         ibc.AppModule{}.ConsensusVersion(),
			genutiltypes.ModuleName:     genutil.AppModule{}.ConsensusVersion(),
			ibctransfertypes.ModuleName: ibctransfer.AppModule{}.ConsensusVersion(),
			assettypes.ModuleName:       asset.AppModule{}.ConsensusVersion(),
			vaulttypes.ModuleName:       vault.AppModule{}.ConsensusVersion(),
		}
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	})

	upgradeInfo, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == "v0.1.0" && !app.upgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{wasmtypes.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
