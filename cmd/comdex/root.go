package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/prometheus/client_golang/prometheus"

	tmdb "github.com/cometbft/cometbft-db"
	tmcfg "github.com/cometbft/cometbft/config"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cometbft/cometbft/libs/log"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	snapshottypes "github.com/cosmos/cosmos-sdk/snapshots/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
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
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(comdex.DefaultNodeHome).
		WithViper("")

	cobra.EnableCommandSorting = false
	root := &cobra.Command{
		Use:   "comdex",
		Short: "Comdex - Decentralised Synthetic Asset Exchange",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			// bump47: recheck if customTMConfig is required, else replace with nil
			customTMConfig := initTendermintConfig()

			return server.InterceptConfigsPreRunHandler(cmd, "", nil, customTMConfig)
		},
	}

	initRootCmd(root, encodingConfig)
	return root, encodingConfig
}

func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	// these values put a higher strain on node memory
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}

func initRootCmd(rootCmd *cobra.Command, encoding comdex.EncodingConfig) {
	gentxModule := comdex.ModuleBasics[genutiltypes.ModuleName].(genutil.AppModuleBasic)
	rootCmd.AddCommand(
		genutilcli.InitCmd(comdex.ModuleBasics, comdex.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, comdex.DefaultNodeHome, gentxModule.GenTxValidator),
		genutilcli.GenTxCmd(comdex.ModuleBasics, encoding.TxConfig, banktypes.GenesisBalancesIterator{}, comdex.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(comdex.ModuleBasics),
		AddGenesisAccountCmd(comdex.DefaultNodeHome),
		// AddGenesisWasmMsgCmd(comdex.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		testnetCmd(comdex.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		MigrateStoreCmd(),
		debug.Cmd(),
		config.Cmd(),
	)

	server.AddCommands(rootCmd, comdex.DefaultNodeHome, appCreatorFunc, appExportFunc, addModuleInitFlags)
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(comdex.DefaultNodeHome),
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
		authcli.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcli.QueryTxsByEventsCmd(),
		authcli.QueryTxCmd(),
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
	)

	comdex.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func appCreatorFunc(logger log.Logger, db tmdb.DB, tracer io.Writer, options servertypes.AppOptions) servertypes.Application {
	var cache sdk.MultiStorePersistentCache
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
	snapshotDB, err := tmdb.NewDB("metadata", tmdb.GoLevelDBBackend, snapshotDir)
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
		comdex.GetWasmEnabledProposals(),
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

func appExportFunc(logger log.Logger, db tmdb.DB, tracer io.Writer, height int64,
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
		app = comdex.New(logger, db, tracer, false, map[int64]bool{}, homePath, cast.ToUint(options.Get(server.FlagInvCheckPeriod)), config, options, comdex.GetWasmEnabledProposals(), emptyWasmOpts)

		if err := app.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		app = comdex.New(logger, db, tracer, true, map[int64]bool{}, homePath, cast.ToUint(options.Get(server.FlagInvCheckPeriod)), config, options, comdex.GetWasmEnabledProposals(), emptyWasmOpts)
	}

	return app.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}
