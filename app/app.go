package app

import (
	"cosmossdk.io/x/tx/signing"
	"encoding/json"
	"fmt"
	dbm "github.com/cosmos/cosmos-db"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/gogoproto/proto"

	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/gorilla/mux"
	"github.com/spf13/cast"

	"github.com/cosmos/cosmos-sdk/runtime"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"

	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8"
	ibchookskeeper "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/keeper"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/types"

	// consensus "github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	packetforward "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"

	"github.com/rakyll/statik/fs"

	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icahost "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	// ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"

	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
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
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibcfee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"

	// ibcclientclient "github.com/cosmos/ibc-go/v8/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"github.com/comdex-official/comdex/x/liquidation"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"

	"cosmossdk.io/log"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmprototypes "github.com/cometbft/cometbft/proto/tendermint/types"

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

	"github.com/comdex-official/comdex/x/common"
	commonkeeper "github.com/comdex-official/comdex/x/common/keeper"
	commontypes "github.com/comdex-official/comdex/x/common/types"

	"github.com/comdex-official/comdex/x/esm"
	esmkeeper "github.com/comdex-official/comdex/x/esm/keeper"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"

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
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"

	"github.com/comdex-official/comdex/x/liquidity"
	liquidityclient "github.com/comdex-official/comdex/x/liquidity/client"
	liquiditykeeper "github.com/comdex-official/comdex/x/liquidity/keeper"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"

	"github.com/comdex-official/comdex/x/liquidationsV2"
	liquidationsV2client "github.com/comdex-official/comdex/x/liquidationsV2/client"
	liquidationsV2keeper "github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	liquidationsV2types "github.com/comdex-official/comdex/x/liquidationsV2/types"

	"github.com/comdex-official/comdex/x/auctionsV2"
	auctionsV2client "github.com/comdex-official/comdex/x/auctionsV2/client"
	auctionsV2keeper "github.com/comdex-official/comdex/x/auctionsV2/keeper"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	icq "github.com/cosmos/ibc-apps/modules/async-icq/v8"
	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v8/keeper"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"

	"github.com/cosmos/cosmos-sdk/std"

	cwasm "github.com/comdex-official/comdex/app/wasm"
	"github.com/cosmos/cosmos-sdk/codec/address"

	mv13 "github.com/comdex-official/comdex/app/upgrades/mainnet/v13"
	tv14 "github.com/comdex-official/comdex/app/upgrades/testnet/v14"
)

const (
	AccountAddressPrefix = "comdex"
	Name                 = "comdex"
)

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
	}
	proposalHandlers = append(proposalHandlers, assetclient.AddAssetsHandler...)
	proposalHandlers = append(proposalHandlers, liquidityclient.LiquidityProposalHandler...)
	proposalHandlers = append(proposalHandlers, liquidationsV2client.LiquidationsV2Handler...)
	proposalHandlers = append(proposalHandlers, auctionsV2client.AuctionsV2Handler...)

	return proposalHandlers
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
	// use this for clarity in argument list
	EmptyWasmOpts []wasm.Option
	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	// ModuleBasics = module.NewBasicManager(
	// 	auth.AppModuleBasic{},
	// 	genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
	// 	bank.AppModuleBasic{},
	// 	capability.AppModuleBasic{},
	// 	staking.AppModuleBasic{},
	// 	mint.AppModuleBasic{},
	// 	distr.AppModuleBasic{},
	// 	gov.NewAppModuleBasic(GetGovProposalHandlers()),
	// 	params.AppModuleBasic{},
	// 	crisis.AppModuleBasic{},
	// 	slashing.AppModuleBasic{},
	// 	authzmodule.AppModuleBasic{},
	// 	ibc.AppModuleBasic{},
	// 	ibctm.AppModuleBasic{},
	// 	upgrade.AppModuleBasic{},
	// 	evidence.AppModuleBasic{},
	// 	ibctransfer.AppModuleBasic{},
	// 	consensus.AppModuleBasic{},
	// 	vesting.AppModuleBasic{},
	// 	vault.AppModuleBasic{},
	// 	asset.AppModuleBasic{},
	// 	esm.AppModuleBasic{},
	// 	lend.AppModuleBasic{},

	// 	market.AppModuleBasic{},
	// 	locker.AppModuleBasic{},
	// 	bandoraclemodule.AppModuleBasic{},
	// 	collector.AppModuleBasic{},
	// 	liquidation.AppModuleBasic{},
	// 	auction.AppModuleBasic{},
	// 	tokenmint.AppModuleBasic{},
	// 	wasm.AppModuleBasic{},
	// 	liquidity.AppModuleBasic{},
	// 	rewards.AppModuleBasic{},
	// 	ica.AppModuleBasic{},
	// 	ibcfee.AppModuleBasic{},
	// 	liquidationsV2.AppModuleBasic{},
	// 	auctionsV2.AppModuleBasic{},
	// 	common.AppModuleBasic{},
	// 	icq.AppModuleBasic{},
	// 	ibchooks.AppModuleBasic{},
	// 	packetforward.AppModuleBasic{},
	// )
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

	amino    *codec.LegacyAmino
	cdc      codec.Codec
	txConfig client.TxConfig

	interfaceRegistry codectypes.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*storetypes.KVStoreKey
	tkeys map[string]*storetypes.TransientStoreKey
	mkeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper     authkeeper.AccountKeeper
	FeegrantKeeper    feegrantkeeper.Keeper
	BankKeeper        bankkeeper.Keeper
	BankBaseKeeper    *bankkeeper.BaseKeeper
	AuthzKeeper       authzkeeper.Keeper
	CapabilityKeeper  *capabilitykeeper.Keeper
	StakingKeeper     *stakingkeeper.Keeper
	SlashingKeeper    slashingkeeper.Keeper
	MintKeeper        mintkeeper.Keeper
	DistrKeeper       distrkeeper.Keeper
	GovKeeper         govkeeper.Keeper
	CrisisKeeper      *crisiskeeper.Keeper
	UpgradeKeeper     *upgradekeeper.Keeper
	ParamsKeeper      paramskeeper.Keeper
	IbcKeeper         *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCFeeKeeper      ibcfeekeeper.Keeper
	IbcHooksKeeper    *ibchookskeeper.Keeper
	ICAHostKeeper     icahostkeeper.Keeper
	EvidenceKeeper    evidencekeeper.Keeper
	IbcTransferKeeper ibctransferkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper         capabilitykeeper.ScopedKeeper
	ScopedIBCTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedIBCOracleKeeper   capabilitykeeper.ScopedKeeper
	ScopedBandoracleKeeper  capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper     capabilitykeeper.ScopedKeeper
	ScopedICQKeeper         capabilitykeeper.ScopedKeeper

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
	NewliqKeeper      liquidationsV2keeper.Keeper
	NewaucKeeper      auctionsV2keeper.Keeper
	CommonKeeper      commonkeeper.Keeper

	// IBC modules
	// transfer module
	Ics20WasmHooks      *ibchooks.WasmHooks
	HooksICS4Wrapper    ibchooks.ICS4Middleware
	PacketForwardKeeper *packetforwardkeeper.Keeper
	ICQKeeper           *icqkeeper.Keeper

	ConsensusParamsKeeper consensusparamkeeper.Keeper

	WasmKeeper     wasm.Keeper
	ContractKeeper *wasmkeeper.PermissionedKeeper
	// the module manager
	mm                 *module.Manager
	BasicModuleManager module.BasicManager
	// Module configurator
	configurator module.Configurator
}

// New returns a reference to an initialized App.
func NewComdexApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOptions servertypes.AppOptions,
	wasmOpts []wasm.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	interfaceRegistry, err := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
			},
			ValidatorAddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
			},
		},
	})
	if err != nil {
		panic(err)
	}
	appCodec := codec.NewProtoCodec(interfaceRegistry)
	legacyAmino := codec.NewLegacyAmino()
	txConfig := authtx.NewTxConfig(appCodec, authtx.DefaultSignModes)

	std.RegisterLegacyAminoCodec(legacyAmino)
	std.RegisterInterfaces(interfaceRegistry)

	var (
		tkeys = storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
		mkeys = storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
		keys  = storetypes.NewKVStoreKeys(
			authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
			minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
			govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, icahosttypes.StoreKey, upgradetypes.StoreKey,
			evidencetypes.StoreKey, ibctransfertypes.StoreKey, ibcfeetypes.StoreKey, capabilitytypes.StoreKey,
			vaulttypes.StoreKey, assettypes.StoreKey, collectortypes.StoreKey, liquidationtypes.StoreKey,
			markettypes.StoreKey, bandoraclemoduletypes.StoreKey, lockertypes.StoreKey,
			wasm.StoreKey, authzkeeper.StoreKey, auctiontypes.StoreKey, tokenminttypes.StoreKey,
			rewardstypes.StoreKey, feegrant.StoreKey, liquiditytypes.StoreKey, esmtypes.ModuleName, lendtypes.StoreKey,
			liquidationsV2types.StoreKey, auctionsV2types.StoreKey, commontypes.StoreKey, ibchookstypes.StoreKey, packetforwardtypes.StoreKey, icqtypes.StoreKey, consensusparamtypes.StoreKey, crisistypes.StoreKey,
		)
	)

	baseApp := baseapp.NewBaseApp(Name, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetVersion(version.Version)
	baseApp.SetInterfaceRegistry(interfaceRegistry)
	baseApp.SetTxEncoder(txConfig.TxEncoder())

	app := &App{
		BaseApp:           baseApp,
		amino:             legacyAmino,
		cdc:               appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
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
	// register the key tables for legacy param subspaces
	keyTable := ibcclienttypes.ParamKeyTable()
	keyTable.RegisterParamSet(&ibcconnectiontypes.Params{})

	//nolint:godox  //TODO: refactor this code
	app.ParamsKeeper.Subspace(authtypes.ModuleName).WithKeyTable(authtypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(banktypes.ModuleName).WithKeyTable(banktypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(stakingtypes.ModuleName).WithKeyTable(stakingtypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(minttypes.ModuleName).WithKeyTable(minttypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(distrtypes.ModuleName).WithKeyTable(distrtypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(slashingtypes.ModuleName).WithKeyTable(slashingtypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypesv1.ParamKeyTable())
	app.ParamsKeeper.Subspace(crisistypes.ModuleName).WithKeyTable(crisistypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(ibctransfertypes.ModuleName).WithKeyTable(ibctransfertypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(ibchost.ModuleName).WithKeyTable(keyTable)
	app.ParamsKeeper.Subspace(icahosttypes.SubModuleName).WithKeyTable(icahosttypes.ParamKeyTable())
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
	app.ParamsKeeper.Subspace(liquidationsV2types.ModuleName)
	app.ParamsKeeper.Subspace(auctionsV2types.ModuleName)
	app.ParamsKeeper.Subspace(commontypes.ModuleName)
	app.ParamsKeeper.Subspace(icqtypes.ModuleName)
	app.ParamsKeeper.Subspace(packetforwardtypes.ModuleName).WithKeyTable(packetforwardtypes.ParamKeyTable())

	// set the BaseApp's parameter store
	// baseApp.SetParamStore(
	// 	app.ParamsKeeper.
	// 		Subspace(baseapp.Paramspace).
	// 		WithKeyTable(paramskeeper.ConsensusParamsKeyTable()),
	// )
	govModAddress := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]), govModAddress, runtime.EventService{})
	baseApp.SetParamStore(app.ConsensusParamsKeeper.ParamsStore)

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
		scopedIBCOracleKeeper  = app.CapabilityKeeper.ScopeToModule(markettypes.ModuleName) // can remove it
		scopedWasmKeeper       = app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
		scopedICAHostKeeper    = app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
		scopedBandoracleKeeper = app.CapabilityKeeper.ScopeToModule(bandoraclemoduletypes.ModuleName)
		scopedICQKeeper        = app.CapabilityKeeper.ScopeToModule(icqtypes.ModuleName)
	)

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		app.cdc,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		app.ModuleAccountsPermissions(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		AccountAddressPrefix,
		govModAddress,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		app.cdc,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		nil,
		govModAddress,
		logger,
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		app.cdc,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		govModAddress,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
	app.MintKeeper = mintkeeper.NewKeeper(
		app.cdc,
		runtime.NewKVStoreService(keys[minttypes.StoreKey]),
		stakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		govModAddress,
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		app.cdc,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		govModAddress,
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		app.cdc,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		stakingKeeper,
		govModAddress,
	)
	invCheckPeriod := cast.ToUint(appOptions.Get(server.FlagInvCheckPeriod))
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.cdc,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		govModAddress,
		app.AccountKeeper.AddressCodec(),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[authzkeeper.StoreKey]),
		app.cdc,
		baseApp.MsgServiceRouter(),
		app.AccountKeeper,
	)

	// get skipUpgradeHeights from the app options
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOptions.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOptions.Get(flags.FlagHome))

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		app.cdc,
		homePath,
		app.BaseApp,
		govModAddress,
	)
	// register the staking hooks
	// NOTE: StakingKeeper above is passed by reference, so that it will contain these hooks
	// app.StakingKeeper = *stakingKeeper.SetHooks(
	// 	stakingtypes.NewMultiStakingHooks(
	// 		app.DistrKeeper.Hooks(),
	// 		app.SlashingKeeper.Hooks(),
	// 	),
	// )
	stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks()),
	)
	app.StakingKeeper = stakingKeeper

	// Create IBC Keeper
	app.IbcKeeper = ibckeeper.NewKeeper(
		app.cdc,
		app.keys[ibchost.StoreKey],
		app.GetSubspace(ibchost.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
		govModAddress,
	)

	// Configure the hooks keeper
	hooksKeeper := ibchookskeeper.NewKeeper(
		app.keys[ibchookstypes.StoreKey],
	)
	app.IbcHooksKeeper = &hooksKeeper

	cmdxPrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	wasmHooks := ibchooks.NewWasmHooks(app.IbcHooksKeeper, &app.WasmKeeper, cmdxPrefix) // The contract keeper needs to be set later
	app.Ics20WasmHooks = &wasmHooks
	app.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		app.IbcKeeper.ChannelKeeper,
		app.Ics20WasmHooks,
	)

	// Do not use this middleware for anything except x/wasm requirement.
	// The spec currently requires new channels to be created, to use it.
	// We need to wait for Channel Upgradability before we can use this for any other middleware.
	app.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec,
		app.keys[ibcfeetypes.StoreKey],
		app.HooksICS4Wrapper, // replaced with IBC middleware
		app.IbcKeeper.ChannelKeeper,
		app.IbcKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
	)

	app.PacketForwardKeeper = packetforwardkeeper.NewKeeper(
		appCodec,
		app.keys[packetforwardtypes.StoreKey],
		app.IbcTransferKeeper, // Will be zero-value here. Reference is set later on with SetTransferKeeper.
		app.IbcKeeper.ChannelKeeper,
		app.DistrKeeper,
		app.BankKeeper,
		app.IbcKeeper.ChannelKeeper,
		govModAddress,
	)

	app.IbcTransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		app.keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.PacketForwardKeeper,
		app.IbcKeeper.ChannelKeeper,
		app.IbcKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
		govModAddress,
	)

	app.PacketForwardKeeper.SetTransferKeeper(app.IbcTransferKeeper)

	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, app.keys[icahosttypes.StoreKey],
		app.GetSubspace(icahosttypes.SubModuleName),
		app.HooksICS4Wrapper,
		app.IbcKeeper.ChannelKeeper,
		app.IbcKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
		govModAddress,
	)

	app.AssetKeeper = assetkeeper.NewKeeper(
		app.cdc,
		app.keys[assettypes.StoreKey],
		app.GetSubspace(assettypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&app.Rewardskeeper,
		&app.VaultKeeper,
		&app.BandoracleKeeper,
		govModAddress,
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
		&app.LiquidationKeeper,
		&app.AuctionKeeper,
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
		govModAddress,
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
		app.IbcKeeper.PortKeeper,
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
		app.IbcKeeper.PortKeeper,
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
		&app.TokenmintKeeper,
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

	app.NewliqKeeper = liquidationsV2keeper.NewKeeper(
		app.cdc,
		app.keys[liquidationsV2types.StoreKey],
		app.keys[liquidationsV2types.MemStoreKey],
		app.GetSubspace(liquidationsV2types.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&app.AssetKeeper,
		&app.VaultKeeper,
		&app.MarketKeeper,
		&app.EsmKeeper,
		&app.Rewardskeeper,
		&app.LendKeeper,
		&app.NewaucKeeper,
		&app.CollectorKeeper,
		govModAddress,
	)

	app.NewaucKeeper = auctionsV2keeper.NewKeeper(
		app.cdc,
		app.keys[auctionsV2types.StoreKey],
		app.keys[auctionsV2types.MemStoreKey],
		app.GetSubspace(auctionsV2types.ModuleName),
		&app.NewliqKeeper,
		app.BankKeeper,
		&app.MarketKeeper,
		&app.AssetKeeper,
		&app.EsmKeeper,
		&app.VaultKeeper,
		&app.CollectorKeeper,
		&app.TokenmintKeeper,
	)

	app.CommonKeeper = commonkeeper.NewKeeper(
		app.cdc,
		app.keys[commontypes.StoreKey],
		app.keys[commontypes.MemStoreKey],
		app.GetSubspace(commontypes.ModuleName),
		&app.WasmKeeper,
		govModAddress,
	)

	// ICQ Keeper
	icqKeeper := icqkeeper.NewKeeper(
		appCodec,
		app.keys[icqtypes.StoreKey],
		app.IbcKeeper.ChannelKeeper, // may be replaced with middleware
		app.IbcKeeper.ChannelKeeper,
		app.IbcKeeper.PortKeeper,
		scopedICQKeeper,
		app.GRPCQueryRouter(),
		govModAddress,
		// NewQuerierWrapper(baseApp), // in-case of strangelove-ventures icq
	)
	app.ICQKeeper = &icqKeeper

	// Create Async ICQ module
	icqModule := icq.NewIBCModule(*app.ICQKeeper)

	// Note: the sealing is done after creating wasmd and wiring that up

	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOptions)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}
	supportedFeatures := "iterator,staking,stargate,comdex,cosmwasm_1_1,cosmwasm_1_2,cosmwasm_1_3"

	wasmOpts = append(cwasm.RegisterCustomPlugins(&app.LockerKeeper, &app.TokenmintKeeper, &app.AssetKeeper, &app.Rewardskeeper, &app.CollectorKeeper, &app.LiquidationKeeper, &app.AuctionKeeper, &app.EsmKeeper, &app.VaultKeeper, &app.LendKeeper, &app.LiquidityKeeper, &app.MarketKeeper), wasmOpts...)

	app.WasmKeeper = wasmkeeper.NewKeeper(
		app.cdc,
		runtime.NewKVStoreService(keys[wasmtypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		distrkeeper.NewQuerier(app.DistrKeeper),
		app.IBCFeeKeeper,
		app.IbcKeeper.ChannelKeeper,
		app.IbcKeeper.PortKeeper,
		scopedWasmKeeper,
		app.IbcTransferKeeper,
		baseApp.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		govModAddress,
		wasmOpts...,
	)

	// set the contract keeper for the Ics20WasmHooks
	app.ContractKeeper = wasmkeeper.NewDefaultPermissionKeeper(app.WasmKeeper)
	app.Ics20WasmHooks.ContractKeeper = &app.WasmKeeper

	// register the proposal types
	govRouter := govtypesv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypesv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(assettypes.RouterKey, asset.NewUpdateAssetProposalHandler(app.AssetKeeper)).
		AddRoute(lendtypes.RouterKey, lend.NewLendHandler(app.LendKeeper)).
		AddRoute(bandoraclemoduletypes.RouterKey, bandoraclemodule.NewFetchPriceHandler(app.BandoracleKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(app.IbcKeeper.ClientKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IbcKeeper.ClientKeeper)).
		AddRoute(liquiditytypes.RouterKey, liquidity.NewLiquidityProposalHandler(app.LiquidityKeeper)).
		AddRoute(liquidationsV2types.RouterKey, liquidationsV2.NewLiquidationsV2Handler(app.NewliqKeeper)).
		AddRoute(auctionsV2types.RouterKey, auctionsV2.NewAuctionsV2Handler(app.NewaucKeeper))

	govKeeper := govkeeper.NewKeeper(
		app.cdc, runtime.NewKVStoreService(keys[govtypes.StoreKey]), app.AccountKeeper, app.BankKeeper,
		app.StakingKeeper, app.DistrKeeper, app.MsgServiceRouter(), govtypes.DefaultConfig(), govModAddress,
	)

	govKeeper.SetLegacyRouter(govRouter)
	app.GovKeeper = *govKeeper

	// Create Transfer Stack
	var transferStack ibcporttypes.IBCModule
	transferStack = ibctransfer.NewIBCModule(app.IbcTransferKeeper)
	transferStack = ibcfee.NewIBCMiddleware(transferStack, app.IBCFeeKeeper)
	transferStack = ibchooks.NewIBCMiddleware(transferStack, &app.HooksICS4Wrapper)
	transferStack = packetforward.NewIBCMiddleware(
		transferStack,
		app.PacketForwardKeeper,
		0,
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		packetforwardkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)

	var (
		evidenceRouter      = evidencetypes.NewRouter()
		ibcRouter           = ibcporttypes.NewRouter()
		oracleModule        = market.NewAppModule(app.cdc, app.MarketKeeper, app.BandoracleKeeper, app.AssetKeeper)
		bandOracleIBCModule = bandoraclemodule.NewIBCModule(app.BandoracleKeeper)
	)

	// RecvPacket, message that originates from core IBC and goes down to app, the flow is:
	// channel.RecvPacket -> fee.OnRecvPacket -> icaHost.OnRecvPacket
	var icaHostStack ibcporttypes.IBCModule
	icaHostStack = icahost.NewIBCModule(app.ICAHostKeeper)
	icaHostStack = ibcfee.NewIBCMiddleware(icaHostStack, app.IBCFeeKeeper)

	// Create fee enabled wasm ibc Stack
	var wasmStack ibcporttypes.IBCModule
	wasmStack = wasm.NewIBCHandler(app.WasmKeeper, app.IbcKeeper.ChannelKeeper, app.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, app.IBCFeeKeeper)

	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)
	ibcRouter.AddRoute(bandoraclemoduletypes.ModuleName, bandOracleIBCModule)
	ibcRouter.AddRoute(wasm.ModuleName, wasmStack)
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostStack)
	ibcRouter.AddRoute(icqtypes.ModuleName, icqModule)
	app.IbcKeeper.SetRouter(ibcRouter)
	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	app.EvidenceKeeper = *evidencekeeper.NewKeeper(
		app.cdc,
		runtime.NewKVStoreService(keys[evidencetypes.StoreKey]),
		app.StakingKeeper,
		app.SlashingKeeper,
		app.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)
	app.EvidenceKeeper.SetRouter(evidenceRouter)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOptions.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app, txConfig),
		auth.NewAppModule(app.cdc, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(app.cdc, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(app.cdc, *app.CapabilityKeeper, false),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(app.cdc, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(app.cdc, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)), // nil -> SDK's default inflation function.
		slashing.NewAppModule(app.cdc, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		distr.NewAppModule(app.cdc, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(app.cdc, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		evidence.NewAppModule(app.EvidenceKeeper),
		authzmodule.NewAppModule(app.cdc, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IbcKeeper),
		ibcfee.NewAppModule(app.IBCFeeKeeper),
		ica.NewAppModule(nil, &app.ICAHostKeeper),
		params.NewAppModule(app.ParamsKeeper),
		// app.RawIcs20TransferAppModule,
		ibctransfer.NewAppModule(app.IbcTransferKeeper),
		asset.NewAppModule(app.cdc, app.AssetKeeper),
		vault.NewAppModule(app.cdc, app.VaultKeeper),
		oracleModule,
		bandoracleModule,
		liquidation.NewAppModule(app.cdc, app.LiquidationKeeper, app.AccountKeeper, app.BankKeeper),
		locker.NewAppModule(app.cdc, app.LockerKeeper, app.AccountKeeper, app.BankKeeper),
		collector.NewAppModule(app.cdc, app.CollectorKeeper, app.AccountKeeper, app.BankKeeper),
		esm.NewAppModule(app.cdc, app.EsmKeeper, app.AccountKeeper, app.BankKeeper, app.AssetKeeper),
		lend.NewAppModule(app.cdc, app.LendKeeper, app.AccountKeeper, app.BankKeeper),
		wasm.NewAppModule(app.cdc, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		auction.NewAppModule(app.cdc, app.AuctionKeeper, app.AccountKeeper, app.BankKeeper, app.CollectorKeeper, app.AssetKeeper, app.EsmKeeper),
		tokenmint.NewAppModule(app.cdc, app.TokenmintKeeper, app.AccountKeeper, app.BankKeeper),
		liquidity.NewAppModule(app.cdc, app.LiquidityKeeper, app.AccountKeeper, app.BankKeeper, app.AssetKeeper),
		rewards.NewAppModule(app.cdc, app.Rewardskeeper, app.AccountKeeper, app.BankKeeper),
		liquidationsV2.NewAppModule(app.cdc, app.NewliqKeeper, app.AccountKeeper, app.BankKeeper),
		common.NewAppModule(app.cdc, app.CommonKeeper, app.AccountKeeper, app.BankKeeper, app.WasmKeeper),
		auctionsV2.NewAppModule(app.cdc, app.NewaucKeeper, app.BankKeeper),
		ibchooks.NewAppModule(app.AccountKeeper),
		icq.NewAppModule(*app.ICQKeeper, app.GetSubspace(icqtypes.ModuleName)),
		packetforward.NewAppModule(app.PacketForwardKeeper, app.GetSubspace(packetforwardtypes.ModuleName)),
	)

	// BasicModuleManager defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration and genesis verification.
	// By default it is composed of all the module from the module manager.
	// Additionally, app module basics can be overwritten by passing them as argument.
	app.BasicModuleManager = module.NewBasicManagerFromManager(
		app.mm,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			govtypes.ModuleName: gov.NewAppModuleBasic(
				[]govclient.ProposalHandler{
					paramsclient.ProposalHandler,
				},
			),
		})
	app.BasicModuleManager.RegisterLegacyAminoCodec(legacyAmino)
	app.BasicModuleManager.RegisterInterfaces(interfaceRegistry)

	// NOTE: upgrade module is required to be prioritized
	app.mm.SetOrderPreBlockers(
		upgradetypes.ModuleName,
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibchost.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		feegrant.ModuleName,
		authtypes.ModuleName,
		capabilitytypes.ModuleName,
		authz.ModuleName,
		assettypes.ModuleName,
		collectortypes.ModuleName,
		vaulttypes.ModuleName,
		bandoraclemoduletypes.ModuleName,
		markettypes.ModuleName,
		lockertypes.ModuleName,
		liquidationtypes.ModuleName,
		auctiontypes.ModuleName,
		tokenminttypes.ModuleName,
		vestingtypes.ModuleName,
		paramstypes.ModuleName,
		wasmtypes.ModuleName,
		banktypes.ModuleName,
		rewardstypes.ModuleName,
		liquiditytypes.ModuleName,
		lendtypes.ModuleName,
		esmtypes.ModuleName,
		liquidationsV2types.ModuleName,
		auctionsV2types.ModuleName,
		commontypes.ModuleName,
		ibchookstypes.ModuleName,
		icqtypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcfeetypes.ModuleName,
		consensusparamtypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		genutiltypes.ModuleName,
		feegrant.ModuleName,
		vestingtypes.ModuleName,
		evidencetypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		ibctransfertypes.ModuleName,
		authtypes.ModuleName,
		slashingtypes.ModuleName,
		authz.ModuleName,
		paramstypes.ModuleName,
		capabilitytypes.ModuleName,
		upgradetypes.ModuleName,
		bandoraclemoduletypes.ModuleName,
		markettypes.ModuleName,
		lockertypes.ModuleName,
		vaulttypes.ModuleName,
		liquidationtypes.ModuleName,
		auctiontypes.ModuleName,
		tokenminttypes.ModuleName,
		wasmtypes.ModuleName,
		lendtypes.ModuleName,
		assettypes.ModuleName,
		collectortypes.ModuleName,
		banktypes.ModuleName,
		rewardstypes.ModuleName,
		liquiditytypes.ModuleName,
		esmtypes.ModuleName,
		liquidationsV2types.ModuleName,
		auctionsV2types.ModuleName,
		commontypes.ModuleName,
		ibchookstypes.ModuleName,
		icqtypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcfeetypes.ModuleName,
		consensusparamtypes.ModuleName,
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
		ibchost.ModuleName,
		icatypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		wasmtypes.ModuleName,
		authz.ModuleName,
		vestingtypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
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
		liquiditytypes.ModuleName,
		rewardstypes.ModuleName,
		crisistypes.ModuleName,
		liquidationsV2types.ModuleName,
		auctionsV2types.ModuleName,
		commontypes.ModuleName,
		ibchookstypes.ModuleName,
		icqtypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcfeetypes.ModuleName,
		consensusparamtypes.ModuleName,
	)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.cdc, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)
	app.registerUpgradeHandlers()
	// initialize stores
	app.MountKVStores(app.keys)
	app.MountTransientStores(app.tkeys)
	app.MountMemoryStores(app.mkeys)

	// SDK v47 - since we do not use dep inject, this gives us access to newer gRPC services.
	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))
	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeegrantKeeper,
				SignModeHandler: txConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			GovKeeper:             app.GovKeeper,
			wasmConfig:            wasmConfig,
			WasmKeeper:            &app.WasmKeeper,
			TXCounterStoreService: runtime.NewKVStoreService(keys[wasmtypes.StoreKey]),
			IBCChannelKeeper:      app.IbcKeeper,
			Cdc:                   appCodec,
		},
	)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	if manager := app.SnapshotManager(); manager != nil {
		err = manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
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
	app.ScopedICQKeeper = scopedICQKeeper

	app.ScopedWasmKeeper = scopedWasmKeeper
	return app
}

// Name returns the name of the App
func (a *App) Name() string { return a.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (a *App) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return a.mm.BeginBlock(ctx)
}

// EndBlocker application updates every end block.
func (a *App) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return a.mm.EndBlock(ctx)
}

// InitChainer application update at chain initialization.
func (a *App) InitChainer(ctx sdk.Context, req *abcitypes.RequestInitChain) (*abcitypes.ResponseInitChain, error) {
	var state GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &state); err != nil {
		panic(err)
	}
	err := a.UpgradeKeeper.SetModuleVersionMap(ctx, a.mm.GetVersionMap())
	if err != nil {
		panic(err)
	}
	response, err := a.mm.InitGenesis(ctx, a.cdc, state)
	return response, err
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
func (a *App) AppCodec() codec.Codec {
	return a.cdc
}

// InterfaceRegistry returns Gaia's InterfaceRegistry.
func (a *App) InterfaceRegistry() codectypes.InterfaceRegistry {
	return a.interfaceRegistry
}

// TxConfig returns WasmApp's TxConfig
func (a *App) TxConfig() client.TxConfig {
	return a.txConfig
}

// AutoCliOpts returns the autocli options for the app.
func (a *App) AutoCliOpts() autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range a.mm.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(a.mm.Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}

// DefaultGenesis returns a default genesis from the registered AppModuleBasic's.
func (a *App) DefaultGenesis() map[string]json.RawMessage {
	return a.BasicModuleManager.DefaultGenesis(a.cdc)
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (a *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	return a.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (a *App) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return a.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (a *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
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
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(ctx, server.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(ctx, server.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	a.BasicModuleManager.RegisterGRPCGatewayRoutes(ctx, server.GRPCGatewayRouter)
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
	cmtservice.RegisterTendermintService(ctx, a.BaseApp.GRPCQueryRouter(), a.interfaceRegistry, a.Query)
}

// RegisterNodeService registers the node gRPC Query service.
func (a *App) RegisterNodeService(clientCtx client.Context, cfg serverconfig.Config) {
	nodeservice.RegisterNodeService(clientCtx, a.GRPCQueryRouter(), cfg)
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
		lendtypes.ModuleAcc4:           {authtypes.Minter, authtypes.Burner},
		lendtypes.ModuleAcc5:           {authtypes.Minter, authtypes.Burner},
		lendtypes.ModuleAcc6:           {authtypes.Minter, authtypes.Burner},
		lendtypes.ModuleAcc7:           {authtypes.Minter, authtypes.Burner},
		lendtypes.ModuleAcc8:           {authtypes.Minter, authtypes.Burner},
		lendtypes.ModuleAcc9:           {authtypes.Minter, authtypes.Burner},
		lendtypes.ModuleAcc10:          {authtypes.Minter, authtypes.Burner},
		lendtypes.ModuleAcc11:          {authtypes.Minter, authtypes.Burner},
		liquidationtypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		auctiontypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		lockertypes.ModuleName:         {authtypes.Minter, authtypes.Burner},
		esmtypes.ModuleName:            {authtypes.Burner},
		wasm.ModuleName:                {authtypes.Burner},
		liquiditytypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
		rewardstypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		liquidationsV2types.ModuleName: {authtypes.Minter, authtypes.Burner},
		auctionsV2types.ModuleName:     {authtypes.Minter, authtypes.Burner},
		commontypes.ModuleName:         nil,
		icatypes.ModuleName:            nil,
		ibcfeetypes.ModuleName:         nil,
		assettypes.ModuleName:          nil,
		icqtypes.ModuleName:            nil,
	}
}

func (app *App) PreBlocker(ctx sdk.Context, _ *abcitypes.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.mm.PreBlock(ctx)
}

func (a *App) registerUpgradeHandlers() {
	a.UpgradeKeeper.SetUpgradeHandler(
		tv14.UpgradeName,
		tv14.CreateUpgradeHandlerV14(a.mm, a.configurator, a.CommonKeeper),
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

func upgradeHandlers(upgradeInfo upgradetypes.Plan, a *App, storeUpgrades *storetypes.StoreUpgrades) *storetypes.StoreUpgrades {
	switch {

	case upgradeInfo.Name == mv13.UpgradeName && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
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
	case upgradeInfo.Name == tv14.UpgradeName && !a.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height):
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{
				commontypes.StoreKey,
			},
		}
	}
	return storeUpgrades
}
