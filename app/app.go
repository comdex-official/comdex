package app

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/comdex-official/comdex/app/keepers"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmprototypes "github.com/cometbft/cometbft/proto/tendermint/types"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"

	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/gorilla/mux"

	"github.com/rakyll/statik/fs"

	"github.com/CosmWasm/wasmd/x/wasm"
	assetclient "github.com/comdex-official/comdex/x/asset/client"
	bandoraclemoduleclient "github.com/comdex-official/comdex/x/bandoracle/client"
	lendclient "github.com/comdex-official/comdex/x/lend/client"
	liquidationsV2client "github.com/comdex-official/comdex/x/liquidationsV2/client"
	liquidationsV2types "github.com/comdex-official/comdex/x/liquidationsV2/types"
	liquidityclient "github.com/comdex-official/comdex/x/liquidity/client"
	tmdb "github.com/cometbft/cometbft-db"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"

	auctionsV2client "github.com/comdex-official/comdex/x/auctionsV2/client"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v7/types"

	mv13 "github.com/comdex-official/comdex/app/upgrades/mainnet/v13"
	tv13 "github.com/comdex-official/comdex/app/upgrades/testnet/v13"
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
		lendclient.AddAssetToPairHandler,
		lendclient.AddAssetRatesParamsHandler,
		lendclient.AddAuctionParamsHandler,
		lendclient.AddMultipleAssetToPairHandler,
		lendclient.AddMultipleLendPairsHandler,
		lendclient.AddPoolPairsHandler,
		lendclient.AddAssetRatesPoolPairsHandler,
		lendclient.AddDepreciatePoolsHandler,
		lendclient.AddEModePairsHandler,
		paramsclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
	}
	proposalHandlers = append(proposalHandlers, assetclient.AddAssetsHandler...)
	proposalHandlers = append(proposalHandlers, liquidityclient.LiquidityProposalHandler...)
	proposalHandlers = append(proposalHandlers, liquidationsV2client.LiquidationsV2Handler...)
	proposalHandlers = append(proposalHandlers, auctionsV2client.AuctionsV2Handler...)

	return proposalHandlers
}

// DefaultNodeHome default home directories for the application daemon
var (
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
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	invCheckPeriod    uint
	interfaceRegistry codectypes.InterfaceRegistry

	// keys to access the substores
	keys  map[string]*storetypes.KVStoreKey
	tkeys map[string]*storetypes.TransientStoreKey
	mKeys map[string]*storetypes.MemoryStoreKey

	AppKeepers keepers.AppKeepers

	// the module manager
	ModuleManager *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// module configurator
	configurator module.Configurator

	// custom checkTx handler
	//checkTxHandler pobabci.CheckTx
}

// New returns a reference to an initialized App.
func New(
	logger log.Logger,
	db tmdb.DB,
	traceStore io.Writer,
	loadLatest bool,
	invCheckPeriod uint,
	appOptions servertypes.AppOptions,
	wasmEnabledProposals []wasm.ProposalType,
	wasmOpts []wasm.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {

	encodingConfig := MakeEncodingConfig()
	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	baseApp := baseapp.NewBaseApp(Name, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetVersion(version.Version)
	baseApp.SetInterfaceRegistry(interfaceRegistry)

	app := &App{
		BaseApp:           baseApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		tkeys:             sdk.NewTransientStoreKeys(paramstypes.TStoreKey),
		mKeys:             sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey),
	}

	app.AppKeepers = keepers.NewAppKeepers(
		appCodec,
		baseApp,
		legacyAmino,
		keepers.GetMaccPerms(),
		appOptions,
		wasmOpts,
		"ucmdx",
	)
	skipGenesisInvariants := cast.ToBool(appOptions.Get(crisis.FlagSkipGenesisInvariants))
	app.keys = app.AppKeepers.GetKVStoreKey()
	app.ModuleManager = module.NewManager(
		appModules(app, encodingConfig, skipGenesisInvariants)...,
	)

	app.configurator = module.NewConfigurator(appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.ModuleManager.RegisterServices(app.configurator)
	app.ModuleManager.SetOrderBeginBlockers(orderBeginBlockers()...)
	app.ModuleManager.SetOrderEndBlockers(orderEndBlockers()...)
	app.ModuleManager.SetOrderInitGenesis(orderInitGenesis()...)
	app.ModuleManager.RegisterInvariants(app.AppKeepers.CrisisKeeper)
	// initialize stores
	app.MountKVStores(app.keys)
	app.MountTransientStores(app.tkeys)
	app.MountMemoryStores(app.mKeys)
	app.registerUpgradeHandlers()

	// SDK v47 - since we do not use dep inject, this gives us access to newer gRPC services.
	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.ModuleManager.Modules))
	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	wasmConfig, err := wasm.ReadWasmConfig(appOptions)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AppKeepers.AccountKeeper,
				BankKeeper:      app.AppKeepers.BankKeeper,
				FeegrantKeeper:  app.AppKeepers.FeegrantKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			GovKeeper:         app.AppKeepers.GovKeeper,
			wasmConfig:        wasmConfig,
			txCounterStoreKey: app.AppKeepers.GetKey(wasmtypes.StoreKey),
			IBCChannelKeeper:  app.AppKeepers.IbcKeeper,
			Cdc:               appCodec,
		},
	)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	if manager := app.SnapshotManager(); manager != nil {
		err = manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.AppKeepers.WasmKeeper),
		)
		if err != nil {
			panic("failed to register snapshot extension: " + err.Error())
		}
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}

		// Initialize and seal the capability keeper so all persistent capabilities
		// are loaded in-memory and prevent any further modules from creating scoped
		// sub-keepers.
		// This must be done during creation of baseapp rather than in InitChain so
		// that in-memory capabilities get regenerated on appKeepers restart.
		// Note that since this reads from the store, we can only perform it when
		// `loadLatest` is set to true.
		ctx := app.BaseApp.NewUncachedContext(true, tmprototypes.Header{})
		app.AppKeepers.CapabilityKeeper.InitMemStore(ctx)
		app.AppKeepers.CapabilityKeeper.Seal()
	}
	// set the BaseApp's parameter store
	// baseApp.SetParamStore(
	// 	app.ParamsKeeper.
	// 		Subspace(baseapp.Paramspace).
	// 		WithKeyTable(paramskeeper.ConsensusParamsKeyTable()),
	// )

	return app
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return app.ModuleManager.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block.
func (app *App) EndBlocker(ctx sdk.Context, req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return app.ModuleManager.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization.
func (app *App) InitChainer(ctx sdk.Context, req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	var state GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &state); err != nil {
		panic(err)
	}
	app.AppKeepers.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap())
	return app.ModuleManager.InitGenesis(ctx, app.AppCodec(), state)
}

// LoadHeight loads a particular height.
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	accounts := make(map[string]bool)

	names := make([]string, 0)
	for name := range keepers.GetMaccPerms() {
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
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns App's codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Gaia's InterfaceRegistry.
func (app *App) InterfaceRegistry() codectypes.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.AppKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (a *App) RegisterAPIRoutes(server *api.Server, apiConfig serverconfig.APIConfig) {
	ctx := server.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(ctx, server.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(ctx, server.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(ctx, server.GRPCGatewayRouter)
	nodeservice.RegisterGRPCGatewayRoutes(ctx, server.GRPCGatewayRouter)

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
	tmservice.RegisterTendermintService(ctx, a.BaseApp.GRPCQueryRouter(), a.interfaceRegistry, a.Query)
}

// RegisterNodeService registers the node gRPC Query service.
func (a *App) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, a.GRPCQueryRouter())
}

func (app *App) registerUpgradeHandlers() {
	app.AppKeepers.UpgradeKeeper.SetUpgradeHandler(
		mv13.UpgradeName,
		mv13.CreateUpgradeHandlerV13(app.ModuleManager, app.configurator, app.appCodec, app.AppKeepers.ParamsKeeper,
			app.AppKeepers.ConsensusParamsKeeper, *app.AppKeepers.IbcKeeper, app.AppKeepers.ICQKeeper,
			app.AppKeepers.GovKeeper, app.AppKeepers.AssetKeeper, app.AppKeepers.LendKeeper,
			app.AppKeepers.NewliqKeeper, app.AppKeepers.NewaucKeeper),
	)
	// When a planned update height is reached, the old binary will panic
	// writing on disk the height and name of the update that triggered it
	// This will read that value, and execute the preparations for the upgrade.
	upgradeInfo, err := app.AppKeepers.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	var storeUpgrades *storetypes.StoreUpgrades

	storeUpgrades = upgradeHandlers(upgradeInfo, app, storeUpgrades)

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}

func upgradeHandlers(upgradeInfo upgradetypes.Plan, a *App, storeUpgrades *storetypes.StoreUpgrades) *storetypes.StoreUpgrades {
	switch {

	case upgradeInfo.Name == mv13.UpgradeName && !a.AppKeepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{
				icqtypes.StoreKey,
				liquidationsV2types.ModuleName,
				auctionsV2types.ModuleName,
				crisistypes.StoreKey,
				consensusparamtypes.StoreKey,
				ibcfeetypes.StoreKey,
			},
		}
	case upgradeInfo.Name == tv13.UpgradeName && !a.AppKeepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{
				crisistypes.StoreKey,
				consensusparamtypes.StoreKey,
				ibcfeetypes.StoreKey,
			},
		}
	}
	return storeUpgrades
}
