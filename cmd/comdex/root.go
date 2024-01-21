package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	storetypes "cosmossdk.io/store/types"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/prometheus/client_golang/prometheus"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/snapshots"
	snapshottypes "cosmossdk.io/store/snapshots/types"
	cmtcfg "github.com/cometbft/cometbft/config"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	comdex "github.com/comdex-official/comdex/app"
)

func NewRootCmd() (*cobra.Command, comdex.EncodingConfig) {
	encodingConfig := comdex.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.FlagBroadcastMode).
		WithHomeDir(comdex.DefaultNodeHome).
		WithViper("")

	cobra.EnableCommandSorting = false
	root := &cobra.Command{
		Use:           "comdex",
		Short:         "Comdex - DeFi Infrastructure layer for the Cosmos ecosystem",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx = initClientCtx.WithCmdContext(cmd.Context())
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}
			// This needs to go after ReadFromClientConfig, as that function
			// sets the RPC client needed for SIGN_MODE_TEXTUAL. This sign mode
			// is only available if the client is online.
			if !initClientCtx.Offline {
				enabledSignModes := append(tx.DefaultSignModes, signing.SignMode_SIGN_MODE_TEXTUAL)
				txConfigOpts := tx.ConfigOptions{
					EnabledSignModes:           enabledSignModes,
					TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(initClientCtx),
				}
				txConfig, err := tx.NewTxConfigWithOptions(
					initClientCtx.Codec,
					txConfigOpts,
				)
				if err != nil {
					return err
				}

				initClientCtx = initClientCtx.WithTxConfig(txConfig)
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}
			// 2 seconds + 1 second tendermint = 3 second blocks
			timeoutCommit := 2 * time.Second

			customAppTemplate, customAppConfig := initAppConfig()
			customCMTConfig := initCometBFTConfig(timeoutCommit)

			// Force faster block times
			os.Setenv("COMDEX_CONSENSUS_TIMEOUT_COMMIT", cast.ToString(timeoutCommit))

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig)
		},
	}

	initRootCmd(root, encodingConfig, initClientCtx)
	return root, encodingConfig
}

func initCometBFTConfig(timeoutCommit time.Duration) *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	// these values put a higher strain on node memory
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

	// While this is set, it only applies to new configs.
	cfg.Consensus.TimeoutCommit = timeoutCommit

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	type CustomAppConfig struct {
		serverconfig.Config

		Wasm wasmtypes.WasmConfig `mapstructure:"wasm"`
	}

	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()

	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
		Wasm:   wasmtypes.DefaultWasmConfig(),
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate + wasmtypes.DefaultConfigTemplate()

	return customAppTemplate, customAppConfig
}

func initRootCmd(rootCmd *cobra.Command, encoding comdex.EncodingConfig, clientCtx client.Context) {
	cfg := sdk.GetConfig()
	cfg.Seal()

	gentxModule := comdex.ModuleBasics[genutiltypes.ModuleName].(genutil.AppModuleBasic)
	rootCmd.AddCommand(
		genutilcli.InitCmd(comdex.ModuleBasics, comdex.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, comdex.DefaultNodeHome, gentxModule.GenTxValidator, clientCtx.TxConfig.SigningContext().ValidatorAddressCodec()),
		genutilcli.GenTxCmd(comdex.ModuleBasics, encoding.TxConfig, banktypes.GenesisBalancesIterator{}, comdex.DefaultNodeHome, clientCtx.TxConfig.SigningContext().ValidatorAddressCodec()),
		genutilcli.ValidateGenesisCmd(comdex.ModuleBasics),
		AddGenesisAccountCmd(comdex.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		testnetCmd(comdex.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		MigrateStoreCmd(),
		debug.Cmd(),
		confixcmd.ConfigCommand(),
	)

	server.AddCommands(rootCmd, comdex.DefaultNodeHome, appCreatorFunc, appExportFunc, addModuleInitFlags)
	rootCmd.AddCommand(
		server.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(),
	)
}

func addModuleInitFlags(cmd *cobra.Command) {
	crisis.AddModuleInitFlags(cmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.ValidatorCommand(),
		rpc.QueryEventForTxCmd(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockResultsCmd(),
	)

	comdex.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcli.GetSignCommand(),
		authcli.GetSignBatchCommand(),
		authcli.GetMultiSignCommand(),
		authcli.GetMultiSignBatchCmd(),
		authcli.GetValidateSignaturesCommand(),
		authcli.GetBroadcastCommand(),
		authcli.GetEncodeCommand(),
		authcli.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
	)

	comdex.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func appCreatorFunc(logger log.Logger, db dbm.DB, tracer io.Writer, options servertypes.AppOptions) servertypes.Application {
	var cache storetypes.MultiStorePersistentCache
	if cast.ToBool(options.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, height := range cast.ToIntSlice(options.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(height)] = true
	}

	pruningOptions, err := server.GetPruningOptionsFromFlags(options)
	if err != nil {
		panic(err)
	}

	snapshotDir := filepath.Join(cast.ToString(options.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := dbm.NewDB("metadata", dbm.GoLevelDBBackend, snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}
	var wasmOpts []wasm.Option
	if cast.ToBool(options.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}
	snapshotOptions := snapshottypes.NewSnapshotOptions(
		cast.ToUint64(options.Get(server.FlagStateSyncSnapshotInterval)),
		cast.ToUint32(options.Get(server.FlagStateSyncSnapshotKeepRecent)),
	)

	homeDir := cast.ToString(options.Get(flags.FlagHome))
	chainID := cast.ToString(options.Get(flags.FlagChainID))
	if chainID == "" {
		// fallback to genesis chain-id
		appGenesis, err := tmtypes.GenesisDocFromFile(filepath.Join(homeDir, "config", "genesis.json"))
		if err != nil {
			panic(err)
		}

		chainID = appGenesis.ChainID
	}

	return comdex.New(
		logger, db, tracer, true, skipUpgradeHeights,
		cast.ToString(options.Get(flags.FlagHome)),
		cast.ToUint(options.Get(server.FlagInvCheckPeriod)),
		comdex.MakeEncodingConfig(),
		options,
		wasmOpts,
		baseapp.SetPruning(pruningOptions),
		baseapp.SetMinGasPrices(cast.ToString(options.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(options.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(options.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(options.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(options.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(options.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshot(snapshotStore, snapshotOptions),
		baseapp.SetChainID(chainID),
	)
}

func appExportFunc(logger log.Logger, db dbm.DB, tracer io.Writer, height int64,
	forZeroHeight bool, jailAllowedAddrs []string, options servertypes.AppOptions, modulesToExport []string,
) (servertypes.ExportedApp, error) {
	config := comdex.MakeEncodingConfig()
	config.Marshaler = codec.NewProtoCodec(config.InterfaceRegistry)
	homePath, ok := options.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home is not set")
	}
	var emptyWasmOpts []wasm.Option
	var app *comdex.App
	if height != -1 {
		app = comdex.New(logger, db, tracer, false, map[int64]bool{}, homePath, cast.ToUint(options.Get(server.FlagInvCheckPeriod)), config, options, emptyWasmOpts)

		if err := app.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		app = comdex.New(logger, db, tracer, true, map[int64]bool{}, homePath, cast.ToUint(options.Get(server.FlagInvCheckPeriod)), config, options, emptyWasmOpts)
	}

	return app.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
