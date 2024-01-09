package keepers

import (
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	cwasm "github.com/comdex-official/comdex/app/wasm"
	"github.com/comdex-official/comdex/x/asset"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctionkeeper "github.com/comdex-official/comdex/x/auction/keeper"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/auctionsV2"
	auctionsV2keeper "github.com/comdex-official/comdex/x/auctionsV2/keeper"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	bandoraclemodule "github.com/comdex-official/comdex/x/bandoracle"
	bandoraclemodulekeeper "github.com/comdex-official/comdex/x/bandoracle/keeper"
	bandoraclemoduletypes "github.com/comdex-official/comdex/x/bandoracle/types"
	collectorkeeper "github.com/comdex-official/comdex/x/collector/keeper"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	esmkeeper "github.com/comdex-official/comdex/x/esm/keeper"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	"github.com/comdex-official/comdex/x/lend"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	"github.com/comdex-official/comdex/x/liquidationsV2"
	liquidationsV2keeper "github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	liquidationsV2types "github.com/comdex-official/comdex/x/liquidationsV2/types"
	"github.com/comdex-official/comdex/x/liquidity"
	liquiditykeeper "github.com/comdex-official/comdex/x/liquidity/keeper"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	"github.com/comdex-official/comdex/x/market"
	marketkeeper "github.com/comdex-official/comdex/x/market/keeper"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	rewardskeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	tokenmintkeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	icq "github.com/cosmos/ibc-apps/modules/async-icq/v7"
	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v7/keeper"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v7/types"
	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7"
	ibchookskeeper "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/keeper"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	"github.com/spf13/cast"
	"path/filepath"
)

var maccPerms = map[string][]string{
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
	icatypes.ModuleName:            nil,
	ibcfeetypes.ModuleName:         nil,
	assettypes.ModuleName:          nil,
	icqtypes.ModuleName:            nil,
	nft.ModuleName:                 nil,
}

type AppKeepers struct {
	*baseapp.BaseApp

	amino *codec.LegacyAmino
	cdc   codec.Codec

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
	NFTKeeper         nftkeeper.Keeper

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
	mm *module.Manager
	// Module configurator
	configurator module.Configurator
}

func NewAppKeepers(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	cdc *codec.LegacyAmino,
	maccPerms map[string][]string,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
	bondDenom string,
) AppKeepers {

	appKeepers := AppKeepers{}

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()
	keys := appKeepers.GetKVStoreKey()
	tkeys := appKeepers.GetTransientStoreKey()

	appKeepers.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, keys[consensusparamtypes.StoreKey], authtypes.NewModuleAddress(govtypes.ModuleName).String())
	bApp.SetParamStore(&appKeepers.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[capabilitytypes.StoreKey],
		appKeepers.mkeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	var (
		scopedIBCKeeper        = appKeepers.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
		scopedTransferKeeper   = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
		scopedIBCOracleKeeper  = appKeepers.CapabilityKeeper.ScopeToModule(markettypes.ModuleName) // can remove it
		scopedWasmKeeper       = appKeepers.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
		scopedICAHostKeeper    = appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
		scopedBandoracleKeeper = appKeepers.CapabilityKeeper.ScopeToModule(bandoraclemoduletypes.ModuleName)
		scopedICQKeeper        = appKeepers.CapabilityKeeper.ScopeToModule(icqtypes.ModuleName)
	)

	Bech32Prefix := "comdex"
	// add keepers
	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appKeepers.cdc,
		appKeepers.keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		maccPerms,
		Bech32Prefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appKeepers.cdc,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		nil,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[minttypes.StoreKey],
		stakingKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[distrtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appKeepers.cdc,
		cdc,
		appKeepers.keys[slashingtypes.StoreKey],
		stakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[crisistypes.StoreKey],
		appKeepers.invCheckPeriod,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appKeepers.cdc,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
	)

	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		appKeepers.keys[upgradetypes.StoreKey],
		appKeepers.cdc,
		homePath,
		appKeepers.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	// register the staking hooks
	// NOTE: StakingKeeper above is passed by reference, so that it will contain these hooks
	// appKeepers.StakingKeeper = *stakingKeeper.SetHooks(
	// 	stakingtypes.NewMultiStakingHooks(
	// 		appKeepers.DistrKeeper.Hooks(),
	// 		appKeepers.SlashingKeeper.Hooks(),
	// 	),
	// )
	stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks()),
	)
	appKeepers.StakingKeeper = stakingKeeper

	// Create IBC Keeper
	appKeepers.IbcKeeper = ibckeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[ibchost.StoreKey],
		appKeepers.GetSubspace(ibchost.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		scopedIBCKeeper,
	)

	appKeepers.NFTKeeper = nftkeeper.NewKeeper(keys[nftkeeper.StoreKey], appCodec, appKeepers.AccountKeeper, appKeepers.BankKeeper)

	// Configure the hooks keeper
	hooksKeeper := ibchookskeeper.NewKeeper(
		appKeepers.keys[ibchookstypes.StoreKey],
	)
	appKeepers.IbcHooksKeeper = &hooksKeeper

	cmdxPrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	wasmHooks := ibchooks.NewWasmHooks(appKeepers.IbcHooksKeeper, &appKeepers.WasmKeeper, cmdxPrefix) // The contract keeper needs to be set later
	appKeepers.Ics20WasmHooks = &wasmHooks
	appKeepers.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		appKeepers.IbcKeeper.ChannelKeeper,
		appKeepers.Ics20WasmHooks,
	)

	// Do not use this middleware for anything except x/wasm requirement.
	// The spec currently requires new channels to be created, to use it.
	// We need to wait for Channel Upgradability before we can use this for any other middleware.
	appKeepers.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcfeetypes.StoreKey],
		appKeepers.HooksICS4Wrapper, // replaced with IBC middleware
		appKeepers.IbcKeeper.ChannelKeeper,
		&appKeepers.IbcKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
	)

	appKeepers.PacketForwardKeeper = packetforwardkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[packetforwardtypes.StoreKey],
		appKeepers.GetSubspace(packetforwardtypes.ModuleName),
		appKeepers.IbcTransferKeeper, // Will be zero-value here. Reference is set later on with SetTransferKeeper.
		appKeepers.IbcKeeper.ChannelKeeper,
		appKeepers.DistrKeeper,
		appKeepers.BankKeeper,
		appKeepers.IbcKeeper.ChannelKeeper,
	)

	appKeepers.IbcTransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.PacketForwardKeeper,
		appKeepers.IbcKeeper.ChannelKeeper,
		&appKeepers.IbcKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		scopedTransferKeeper,
	)

	appKeepers.PacketForwardKeeper.SetTransferKeeper(appKeepers.IbcTransferKeeper)

	appKeepers.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.HooksICS4Wrapper,
		appKeepers.IbcKeeper.ChannelKeeper,
		&appKeepers.IbcKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		scopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)

	appKeepers.AssetKeeper = assetkeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[assettypes.StoreKey],
		appKeepers.GetSubspace(assettypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		&appKeepers.Rewardskeeper,
		&appKeepers.VaultKeeper,
		&appKeepers.BandoracleKeeper,
	)

	appKeepers.LendKeeper = lendkeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[lendtypes.StoreKey],
		appKeepers.keys[lendtypes.StoreKey],
		appKeepers.GetSubspace(lendtypes.ModuleName),
		appKeepers.BankKeeper,
		appKeepers.AccountKeeper,
		&appKeepers.AssetKeeper,
		&appKeepers.MarketKeeper,
		&appKeepers.EsmKeeper,
		&appKeepers.LiquidationKeeper,
		&appKeepers.AuctionKeeper,
	)

	appKeepers.EsmKeeper = esmkeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[esmtypes.StoreKey],
		appKeepers.keys[esmtypes.StoreKey],
		appKeepers.GetSubspace(esmtypes.ModuleName),
		&appKeepers.AssetKeeper,
		&appKeepers.VaultKeeper,
		appKeepers.BankKeeper,
		&appKeepers.MarketKeeper,
		&appKeepers.TokenmintKeeper,
		&appKeepers.CollectorKeeper,
	)

	appKeepers.VaultKeeper = vaultkeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[vaulttypes.StoreKey],
		appKeepers.BankKeeper,
		&appKeepers.AssetKeeper,
		&appKeepers.MarketKeeper,
		&appKeepers.CollectorKeeper,
		&appKeepers.EsmKeeper,
		&appKeepers.TokenmintKeeper,
		&appKeepers.Rewardskeeper,
	)

	appKeepers.TokenmintKeeper = tokenmintkeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[tokenminttypes.StoreKey],
		appKeepers.BankKeeper,
		&appKeepers.AssetKeeper,
	)

	appKeepers.BandoracleKeeper = bandoraclemodulekeeper.NewKeeper(
		appCodec,
		keys[bandoraclemoduletypes.StoreKey],
		keys[bandoraclemoduletypes.MemStoreKey],
		appKeepers.GetSubspace(bandoraclemoduletypes.ModuleName),
		appKeepers.IbcKeeper.ChannelKeeper,
		&appKeepers.IbcKeeper.PortKeeper,
		scopedBandoracleKeeper,
		&appKeepers.MarketKeeper,
		appKeepers.AssetKeeper,
	)
	bandoracleModule := bandoraclemodule.NewAppModule(
		appCodec,
		appKeepers.BandoracleKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedBandoracleKeeper,
		&appKeepers.IbcKeeper.PortKeeper,
		appKeepers.IbcKeeper.ChannelKeeper,
	)

	appKeepers.MarketKeeper = marketkeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[markettypes.StoreKey],
		appKeepers.GetSubspace(markettypes.ModuleName),
		scopedIBCOracleKeeper,
		appKeepers.AssetKeeper,
		&appKeepers.BandoracleKeeper,
	)

	appKeepers.LiquidationKeeper = liquidationkeeper.NewKeeper(
		appKeepers.cdc,
		keys[liquidationtypes.StoreKey],
		keys[liquidationtypes.MemStoreKey],
		appKeepers.GetSubspace(liquidationtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		&appKeepers.AssetKeeper,
		&appKeepers.VaultKeeper,
		&appKeepers.MarketKeeper,
		&appKeepers.AuctionKeeper,
		&appKeepers.EsmKeeper,
		&appKeepers.Rewardskeeper,
		&appKeepers.LendKeeper,
	)

	appKeepers.AuctionKeeper = auctionkeeper.NewKeeper(
		appKeepers.cdc,
		keys[auctiontypes.StoreKey],
		keys[auctiontypes.MemStoreKey],
		appKeepers.GetSubspace(auctiontypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		&appKeepers.MarketKeeper,
		&appKeepers.LiquidationKeeper,
		&appKeepers.AssetKeeper,
		&appKeepers.VaultKeeper,
		&appKeepers.CollectorKeeper,
		&appKeepers.TokenmintKeeper,
		&appKeepers.EsmKeeper,
		&appKeepers.LendKeeper,
	)

	appKeepers.CollectorKeeper = collectorkeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[collectortypes.StoreKey],
		appKeepers.keys[collectortypes.MemStoreKey],
		&appKeepers.AssetKeeper,
		&appKeepers.AuctionKeeper,
		&appKeepers.LockerKeeper,
		&appKeepers.Rewardskeeper,
		appKeepers.GetSubspace(collectortypes.ModuleName),
		appKeepers.BankKeeper,
	)

	appKeepers.LockerKeeper = lockerkeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[lockertypes.StoreKey],
		appKeepers.GetSubspace(lockertypes.ModuleName),
		appKeepers.BankKeeper,
		&appKeepers.AssetKeeper,
		&appKeepers.CollectorKeeper,
		&appKeepers.EsmKeeper,
		&appKeepers.Rewardskeeper,
	)

	appKeepers.LiquidityKeeper = liquiditykeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[liquiditytypes.StoreKey],
		appKeepers.GetSubspace(liquiditytypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		&appKeepers.AssetKeeper,
		&appKeepers.MarketKeeper,
		&appKeepers.Rewardskeeper,
		&appKeepers.TokenmintKeeper,
	)

	appKeepers.Rewardskeeper = rewardskeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[rewardstypes.StoreKey],
		appKeepers.keys[rewardstypes.MemStoreKey],
		appKeepers.GetSubspace(rewardstypes.ModuleName),
		&appKeepers.LockerKeeper,
		&appKeepers.CollectorKeeper,
		&appKeepers.VaultKeeper,
		&appKeepers.AssetKeeper,
		appKeepers.BankKeeper,
		appKeepers.LiquidityKeeper,
		&appKeepers.MarketKeeper,
		&appKeepers.EsmKeeper,
		&appKeepers.LendKeeper,
	)

	appKeepers.NewliqKeeper = liquidationsV2keeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[liquidationsV2types.StoreKey],
		appKeepers.keys[liquidationsV2types.MemStoreKey],
		appKeepers.GetSubspace(liquidationsV2types.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		&appKeepers.AssetKeeper,
		&appKeepers.VaultKeeper,
		&appKeepers.MarketKeeper,
		&appKeepers.EsmKeeper,
		&appKeepers.Rewardskeeper,
		&appKeepers.LendKeeper,
		&appKeepers.NewaucKeeper,
		&appKeepers.CollectorKeeper,
	)

	appKeepers.NewaucKeeper = auctionsV2keeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[auctionsV2types.StoreKey],
		appKeepers.keys[auctionsV2types.MemStoreKey],
		appKeepers.GetSubspace(auctionsV2types.ModuleName),
		&appKeepers.NewliqKeeper,
		appKeepers.BankKeeper,
		&appKeepers.MarketKeeper,
		&appKeepers.AssetKeeper,
		&appKeepers.EsmKeeper,
		&appKeepers.VaultKeeper,
		&appKeepers.CollectorKeeper,
		&appKeepers.TokenmintKeeper,
	)

	// ICQ Keeper
	icqKeeper := icqkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icqtypes.StoreKey],
		appKeepers.GetSubspace(icqtypes.ModuleName),
		appKeepers.IbcKeeper.ChannelKeeper, // may be replaced with middleware
		appKeepers.IbcKeeper.ChannelKeeper,
		&appKeepers.IbcKeeper.PortKeeper,
		scopedICQKeeper,
		appKeepers.GRPCQueryRouter(),
		// NewQuerierWrapper(baseApp), // in-case of strangelove-ventures icq
	)
	appKeepers.ICQKeeper = &icqKeeper

	// Create Async ICQ module
	icqModule := icq.NewIBCModule(*appKeepers.ICQKeeper)

	// Note: the sealing is done after creating wasmd and wiring that up

	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}
	supportedFeatures := "iterator,staking,stargate,comdex,cosmwasm_1_1,cosmwasm_1_2,cosmwasm_1_3"

	wasmOpts = append(cwasm.RegisterCustomPlugins(&appKeepers.LockerKeeper, &appKeepers.TokenmintKeeper, &appKeepers.AssetKeeper, &appKeepers.Rewardskeeper, &appKeepers.CollectorKeeper, &appKeepers.LiquidationKeeper, &appKeepers.AuctionKeeper, &appKeepers.EsmKeeper, &appKeepers.VaultKeeper, &appKeepers.LendKeeper, &appKeepers.LiquidityKeeper, &appKeepers.MarketKeeper), wasmOpts...)

	appKeepers.WasmKeeper = wasmkeeper.NewKeeper(
		appKeepers.cdc,
		keys[wasmtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		distrkeeper.NewQuerier(appKeepers.DistrKeeper),
		appKeepers.IBCFeeKeeper,
		appKeepers.IbcKeeper.ChannelKeeper,
		&appKeepers.IbcKeeper.PortKeeper,
		scopedWasmKeeper,
		appKeepers.IbcTransferKeeper,
		bApp.MsgServiceRouter(),
		appKeepers.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)

	// set the contract keeper for the Ics20WasmHooks
	appKeepers.ContractKeeper = wasmkeeper.NewDefaultPermissionKeeper(appKeepers.WasmKeeper)
	appKeepers.Ics20WasmHooks.ContractKeeper = &appKeepers.WasmKeeper

	// register the proposal types
	govRouter := govtypesv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypesv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(assettypes.RouterKey, asset.NewUpdateAssetProposalHandler(appKeepers.AssetKeeper)).
		AddRoute(lendtypes.RouterKey, lend.NewLendHandler(appKeepers.LendKeeper)).
		AddRoute(bandoraclemoduletypes.RouterKey, bandoraclemodule.NewFetchPriceHandler(appKeepers.BandoracleKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IbcKeeper.ClientKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IbcKeeper.ClientKeeper)).
		AddRoute(liquiditytypes.RouterKey, liquidity.NewLiquidityProposalHandler(appKeepers.LiquidityKeeper)).
		AddRoute(liquidationsV2types.RouterKey, liquidationsV2.NewLiquidationsV2Handler(appKeepers.NewliqKeeper)).
		AddRoute(auctionsV2types.RouterKey, auctionsV2.NewAuctionsV2Handler(appKeepers.NewaucKeeper))

	if len(wasmEnabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(appKeepers.WasmKeeper, wasmEnabledProposals))
	}

	govKeeper := govkeeper.NewKeeper(
		appKeepers.cdc, keys[govtypes.StoreKey], appKeepers.AccountKeeper, appKeepers.BankKeeper,
		appKeepers.StakingKeeper, bApp.MsgServiceRouter(), govtypes.DefaultConfig(), authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	govKeeper.SetLegacyRouter(govRouter)
	appKeepers.GovKeeper = *govKeeper

	// Create Transfer Stack
	var transferStack ibcporttypes.IBCModule
	transferStack = ibctransfer.NewIBCModule(appKeepers.IbcTransferKeeper)
	transferStack = ibcfee.NewIBCMiddleware(transferStack, appKeepers.IBCFeeKeeper)
	transferStack = ibchooks.NewIBCMiddleware(transferStack, &appKeepers.HooksICS4Wrapper)
	transferStack = packetforward.NewIBCMiddleware(
		transferStack,
		appKeepers.PacketForwardKeeper,
		0,
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		packetforwardkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)

	var (
		evidenceRouter      = evidencetypes.NewRouter()
		ibcRouter           = ibcporttypes.NewRouter()
		oracleModule        = market.NewAppModule(appKeepers.cdc, appKeepers.MarketKeeper, appKeepers.BandoracleKeeper, appKeepers.AssetKeeper)
		bandOracleIBCModule = bandoraclemodule.NewIBCModule(appKeepers.BandoracleKeeper)
	)

	// RecvPacket, message that originates from core IBC and goes down to appKeepers, the flow is:
	// channel.RecvPacket -> fee.OnRecvPacket -> icaHost.OnRecvPacket
	var icaHostStack ibcporttypes.IBCModule
	icaHostStack = icahost.NewIBCModule(appKeepers.ICAHostKeeper)
	icaHostStack = ibcfee.NewIBCMiddleware(icaHostStack, appKeepers.IBCFeeKeeper)

	// Create fee enabled wasm ibc Stack
	var wasmStack ibcporttypes.IBCModule
	wasmStack = wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IbcKeeper.ChannelKeeper, appKeepers.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, appKeepers.IBCFeeKeeper)

	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)
	ibcRouter.AddRoute(bandoraclemoduletypes.ModuleName, bandOracleIBCModule)
	ibcRouter.AddRoute(wasm.ModuleName, wasmStack)
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostStack)
	ibcRouter.AddRoute(icqtypes.ModuleName, icqModule)
	appKeepers.IbcKeeper.SetRouter(ibcRouter)
	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	appKeepers.EvidenceKeeper = *evidencekeeper.NewKeeper(
		appKeepers.cdc,
		appKeepers.keys[evidencetypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
	)
	appKeepers.EvidenceKeeper.SetRouter(evidenceRouter)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	appKeepers.ScopedIBCKeeper = scopedIBCKeeper
	appKeepers.ScopedIBCTransferKeeper = scopedTransferKeeper
	appKeepers.ScopedIBCOracleKeeper = scopedIBCOracleKeeper
	appKeepers.ScopedICAHostKeeper = scopedICAHostKeeper
	appKeepers.ScopedBandoracleKeeper = scopedBandoracleKeeper
	appKeepers.ScopedICQKeeper = scopedICQKeeper

	appKeepers.ScopedWasmKeeper = scopedWasmKeeper

	return appKeepers
}
func initParamsKeeper(cdc codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {

	paramsKeeper := paramskeeper.NewKeeper(
		cdc,
		legacyAmino,
		key,
		tkey,
	)

	paramsKeeper.Subspace(authtypes.ModuleName).WithKeyTable(authtypes.ParamKeyTable())
	paramsKeeper.Subspace(banktypes.ModuleName).WithKeyTable(banktypes.ParamKeyTable())
	paramsKeeper.Subspace(stakingtypes.ModuleName).WithKeyTable(stakingtypes.ParamKeyTable())
	paramsKeeper.Subspace(minttypes.ModuleName).WithKeyTable(minttypes.ParamKeyTable())
	paramsKeeper.Subspace(distrtypes.ModuleName).WithKeyTable(distrtypes.ParamKeyTable())
	paramsKeeper.Subspace(slashingtypes.ModuleName).WithKeyTable(slashingtypes.ParamKeyTable())
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypesv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName).WithKeyTable(crisistypes.ParamKeyTable())
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(vaulttypes.ModuleName)
	paramsKeeper.Subspace(assettypes.ModuleName)
	paramsKeeper.Subspace(collectortypes.ModuleName)
	paramsKeeper.Subspace(esmtypes.ModuleName)
	paramsKeeper.Subspace(lendtypes.ModuleName)
	paramsKeeper.Subspace(markettypes.ModuleName)
	paramsKeeper.Subspace(liquidationtypes.ModuleName)
	paramsKeeper.Subspace(lockertypes.ModuleName)
	paramsKeeper.Subspace(bandoraclemoduletypes.ModuleName)
	paramsKeeper.Subspace(wasmtypes.ModuleName).WithKeyTable(wasmtypes.ParamKeyTable())
	paramsKeeper.Subspace(auctiontypes.ModuleName)
	paramsKeeper.Subspace(tokenminttypes.ModuleName)
	paramsKeeper.Subspace(liquiditytypes.ModuleName)
	paramsKeeper.Subspace(rewardstypes.ModuleName)
	paramsKeeper.Subspace(liquidationsV2types.ModuleName)
	paramsKeeper.Subspace(auctionsV2types.ModuleName)
	paramsKeeper.Subspace(icqtypes.ModuleName)
	paramsKeeper.Subspace(packetforwardtypes.ModuleName).WithKeyTable(packetforwardtypes.ParamKeyTable())

	return paramsKeeper
}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// GetStakingKeeper implements the TestingApp interface.
func (appKeepers *AppKeepers) GetStakingKeeper() *stakingkeeper.Keeper {
	return appKeepers.StakingKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (appKeepers *AppKeepers) GetIBCKeeper() *ibckeeper.Keeper {
	return appKeepers.IbcKeeper
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (appKeepers *AppKeepers) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return appKeepers.ScopedIBCKeeper
}

// GetWasmKeeper implements the TestingApp interface.
func (appKeepers *AppKeepers) GetWasmKeeper() wasmkeeper.Keeper {
	return appKeepers.WasmKeeper
}

func GetMaccPerms() map[string][]string {
	return maccPerms
}
