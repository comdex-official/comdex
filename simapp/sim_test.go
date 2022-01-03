package simapp

import (
	"os"
	"testing"

	"github.com/comdex-official/comdex/app"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdkSimapp "github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	simulation2 "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Profile with:
// /usr/local/go/bin/go test -benchmem -run=^$ github.com/osmosis-labs/osmosis/simapp -bench ^BenchmarkFullAppSimulation$ -Commit=true -cpuprofile cpu.out

// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// an IAVLStore for faster simulation speed.
func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

func TestFullAppSimulation(tb *testing.T) {

	//flags
	sdkSimapp.FlagEnabledValue = true
	sdkSimapp.FlagNumBlocksValue = 20
	sdkSimapp.FlagBlockSizeValue = 25
	sdkSimapp.FlagCommitValue = true
	sdkSimapp.FlagVerboseValue = true
	sdkSimapp.FlagPeriodValue = 10
	sdkSimapp.FlagSeedValue = 10

	app.SetAccountAddressPrefixes()

	config, db, dir, logger, _, err := sdkSimapp.SetupSimulation("goleveldb-app-sim", "Simulation")
	if err != nil {
		tb.Fatalf("simulation setup failed: %s", err.Error())
	}

	defer func() {
		db.Close()
		err = os.RemoveAll(dir)
		if err != nil {
			tb.Fatal(err)
		}
	}()

	comdex := app.New(
		logger,
		db,
		nil,
		true, // load latest
		map[int64]bool{},
		app.DefaultNodeHome,
		sdkSimapp.FlagPeriodValue,
		app.MakeEncodingConfig(),
		sdkSimapp.EmptyAppOptions{},
		interBlockCacheOpt(),
		fauxMerkleModeOpt)

	// Run randomized simulation:
	_, simParams, simErr := simulation.SimulateFromSeed(
		tb,
		os.Stdout,
		comdex.BaseApp,
		AppStateFn(comdex.AppCodec(), comdex.SimulationManager()),
		simulation2.RandomAccounts,                                        // Replace with own random account function if using keys other than secp256k1
		sdkSimapp.SimulationOperations(comdex, comdex.AppCodec(), config), // Run all registered operations
		comdex.ModuleAccountAddrs(),
		config,
		comdex.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	if err = sdkSimapp.CheckExportSimulation(comdex, config, simParams); err != nil {
		tb.Fatal(err)
	}

	if simErr != nil {
		tb.Fatal(simErr)
	}

	if config.Commit {
		sdkSimapp.PrintStats(db)
	}
}

// interBlockCacheOpt returns a BaseApp option function that sets the persistent
// inter-block write-through cache.
func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}
