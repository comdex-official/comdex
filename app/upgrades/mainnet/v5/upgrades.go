package v5

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

func IntializeStates(
	ctx sdk.Context,
	assetKeeper assetkeeper.Keeper,
) {
	apps := []assettypes.AppData{
		{Name: "CSWAP", ShortName: "cswap", MinGovDeposit: sdk.ZeroInt(), GovTimeInSeconds: 0, GenesisToken: []assettypes.MintGenesisToken{}},
		{Name: "HARBOR", ShortName: "hbr", MinGovDeposit: sdk.NewInt(10000000), GovTimeInSeconds: 300, GenesisToken: []assettypes.MintGenesisToken{}},
	}
	for _, app := range apps {
		err := assetKeeper.AddAppRecords(ctx, app)
		if err != nil {
			panic(err)
		}
	}

	assets := []assettypes.Asset{
		{Name: "ATOM", Denom: "uatom", Decimals: 1000000, IsOnChain: false, IsOraclePriceRequired: true},
		{Name: "CMDX", Denom: "ucmdx", Decimals: 1000000, IsOnChain: false, IsOraclePriceRequired: true},
		{Name: "CMST", Denom: "ucmst", Decimals: 1000000, IsOnChain: false, IsOraclePriceRequired: false},
		{Name: "OSMO", Denom: "uosmo", Decimals: 1000000, IsOnChain: false, IsOraclePriceRequired: true},
		{Name: "cATOM", Denom: "ucatom", Decimals: 1000000, IsOnChain: false, IsOraclePriceRequired: false},
		{Name: "cCMDX", Denom: "uccmdx", Decimals: 1000000, IsOnChain: false, IsOraclePriceRequired: false},
		{Name: "cCMST", Denom: "uccmst", Decimals: 1000000, IsOnChain: false, IsOraclePriceRequired: false},
		{Name: "cOSMO", Denom: "ucosmo", Decimals: 1000000, IsOnChain: false, IsOraclePriceRequired: false},
		{Name: "HARBOR", Denom: "uharbor", Decimals: 1000000, IsOnChain: true, IsOraclePriceRequired: false},
	}

	for _, asset := range assets {
		err := assetKeeper.AddAssetRecords(ctx, asset)
		if err != nil {
			panic(err)
		}
	}
}

// CreateUpgradeHandler creates an SDK upgrade handler for v5
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	wasmKeeper wasmkeeper.Keeper,
	assetKeeper assetkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Refs:
		// - https://docs.cosmos.network/master/building-modules/upgrade.html#registering-migrations
		// - https://docs.cosmos.network/master/migrations/chain-upgrade-guide-044.html#chain-upgrade

		// Deleting these modules from the upgrades current state
		// Add Interchain Accounts host module
		// set the ICS27 consensus version so InitGenesis is not run
		fromVM[icatypes.ModuleName] = mm.Modules[icatypes.ModuleName].ConsensusVersion()

		// create ICS27 Controller submodule params, controller module not enabled.
		controllerParams := icacontrollertypes.Params{}

		// create ICS27 Host submodule params

		// create ICS27 Host submodule params
		hostParams := icahosttypes.Params{
			HostEnabled: true,
			AllowMessages: []string{
				sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}),
				sdk.MsgTypeURL(&banktypes.MsgSend{}),
				sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}),
				sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{}),
				sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{}),
				sdk.MsgTypeURL(&stakingtypes.MsgEditValidator{}),
				sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}),
				sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}),
				sdk.MsgTypeURL(&distrtypes.MsgSetWithdrawAddress{}),
				sdk.MsgTypeURL(&distrtypes.MsgWithdrawValidatorCommission{}),
				sdk.MsgTypeURL(&distrtypes.MsgFundCommunityPool{}),
				sdk.MsgTypeURL(&govtypes.MsgVote{}),
			},
		}
		// No changes in existing module and their states,
		// This upgrades adds new modules and new states in the existing store

		icamodule, correctTypecast := mm.Modules[icatypes.ModuleName].(ica.AppModule)
		if !correctTypecast {
			panic("mm.Modules[icatypes.ModuleName] is not of type ica.AppModule")
		}
		icamodule.InitModule(ctx, controllerParams, hostParams)
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return newVM, err
		}

		IntializeStates(ctx, assetKeeper)

		// update wasm to permission
		wasmParams := wasmKeeper.GetParams(ctx)
		wasmParams.CodeUploadAccess = wasmtypes.AllowNobody
		wasmKeeper.SetParams(ctx, wasmParams)

		return newVM, err
	}
}
