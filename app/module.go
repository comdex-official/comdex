package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/comdex-official/comdex/x/asset"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/auction"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/auctionsV2"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	bandoraclemodule "github.com/comdex-official/comdex/x/bandoracle"
	bandoraclemoduletypes "github.com/comdex-official/comdex/x/bandoracle/types"
	"github.com/comdex-official/comdex/x/collector"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	"github.com/comdex-official/comdex/x/esm"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	"github.com/comdex-official/comdex/x/lend"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/comdex-official/comdex/x/liquidation"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	"github.com/comdex-official/comdex/x/liquidationsV2"
	liquidationsV2types "github.com/comdex-official/comdex/x/liquidationsV2/types"
	"github.com/comdex-official/comdex/x/liquidity"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"
	"github.com/comdex-official/comdex/x/locker"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	"github.com/comdex-official/comdex/x/market"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	"github.com/comdex-official/comdex/x/rewards"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	"github.com/comdex-official/comdex/x/tokenmint"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	"github.com/comdex-official/comdex/x/vault"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	icq "github.com/cosmos/ibc-apps/modules/async-icq/v7"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v7/types"
	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibchost "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

var ModuleBasics = module.NewBasicManager(

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.

	auth.AppModuleBasic{},
	genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	gov.NewAppModuleBasic(GetGovProposalHandlers()),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	ibctm.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	ibctransfer.AppModuleBasic{},
	consensus.AppModuleBasic{},
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
	ibcfee.AppModuleBasic{},
	liquidationsV2.AppModuleBasic{},
	auctionsV2.AppModuleBasic{},
	icq.AppModuleBasic{},
	ibchooks.AppModuleBasic{},
	packetforward.AppModuleBasic{},
)

func appModules(
	app *App,
	encodingConfig EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Marshaler

	//bondDenom := app.GetChainBondDenom()

	return []module.AppModule{
		genutil.NewAppModule(
			app.AppKeepers.AccountKeeper,
			app.AppKeepers.StakingKeeper,
			app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AppKeepers.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper),
		bank.NewAppModule(appCodec, app.AppKeepers.BankKeeper, app.AppKeepers.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.AppKeepers.CapabilityKeeper, false),
		crisis.NewAppModule(app.AppKeepers.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(appCodec, &app.AppKeepers.GovKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.AppKeepers.MintKeeper, app.AppKeepers.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)), // nil -> SDK's default inflation function.
		slashing.NewAppModule(appCodec, app.AppKeepers.SlashingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.AppKeepers.DistrKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.AppKeepers.StakingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.AppKeepers.UpgradeKeeper),
		evidence.NewAppModule(app.AppKeepers.EvidenceKeeper),
		authzmodule.NewAppModule(appCodec, app.AppKeepers.AuthzKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.AppKeepers.IbcKeeper),
		ibcfee.NewAppModule(app.AppKeepers.IBCFeeKeeper),
		ica.NewAppModule(nil, &app.AppKeepers.ICAHostKeeper),
		params.NewAppModule(app.AppKeepers.ParamsKeeper),
		// app.AppKeepers.RawIcs20TransferAppModule,
		ibctransfer.NewAppModule(app.AppKeepers.IbcTransferKeeper),
		asset.NewAppModule(appCodec, app.AppKeepers.AssetKeeper),
		vault.NewAppModule(appCodec, app.AppKeepers.VaultKeeper),
		oracleModule,
		bandoracleModule,
		liquidation.NewAppModule(appCodec, app.AppKeepers.LiquidationKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper),
		locker.NewAppModule(appCodec, app.AppKeepers.LockerKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper),
		collector.NewAppModule(appCodec, app.AppKeepers.CollectorKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper),
		esm.NewAppModule(appCodec, app.AppKeepers.EsmKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.AssetKeeper),
		lend.NewAppModule(appCodec, app.AppKeepers.LendKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper),
		wasm.NewAppModule(appCodec, &app.AppKeepers.WasmKeeper, app.AppKeepers.StakingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		auction.NewAppModule(appCodec, app.AppKeepers.AuctionKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.CollectorKeeper, app.AppKeepers.AssetKeeper, app.AppKeepers.EsmKeeper),
		tokenmint.NewAppModule(appCodec, app.AppKeepers.TokenmintKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper),
		liquidity.NewAppModule(appCodec, app.AppKeepers.LiquidityKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.AssetKeeper),
		rewards.NewAppModule(appCodec, app.AppKeepers.Rewardskeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper),
		liquidationsV2.NewAppModule(appCodec, app.AppKeepers.NewliqKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper),
		auctionsV2.NewAppModule(appCodec, app.AppKeepers.NewaucKeeper, app.AppKeepers.BankKeeper),
		ibchooks.NewAppModule(app.AppKeepers.AccountKeeper),
		icq.NewAppModule(*app.AppKeepers.ICQKeeper),
		packetforward.NewAppModule(app.AppKeepers.PacketForwardKeeper),
	}
}
func orderBeginBlockers() []string {
	return []string{
		upgradetypes.ModuleName,
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
		ibchookstypes.ModuleName,
		icqtypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcfeetypes.ModuleName,
		consensusparamtypes.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
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
		ibchookstypes.ModuleName,
		icqtypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcfeetypes.ModuleName,
		consensusparamtypes.ModuleName,
	}
}

func orderInitGenesis() []string {
	return []string{
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
		ibchookstypes.ModuleName,
		icqtypes.ModuleName,
		packetforwardtypes.ModuleName,
		ibcfeetypes.ModuleName,
		consensusparamtypes.ModuleName,
	}
}
