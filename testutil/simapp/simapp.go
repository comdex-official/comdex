package simapp

import (
	"time"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmprototypes "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	comdex "github.com/comdex-official/comdex/app"
)

// New creates application instance with in-memory database and disabled logging.
func New(dir string) *comdex.App {
	var (
		db       = dbm.NewMemDB()
		logger   = log.NewNopLogger()
	)

	a := comdex.NewComdexApp(logger, db, nil, true,
		simtestutil.EmptyAppOptions{}, comdex.EmptyWasmOpts)
	// InitChain updates deliverState which is required when app.NewContext is called
	a.InitChain(&abcitypes.RequestInitChain{
		ConsensusParams: defaultConsensusParams,
		AppStateBytes:   []byte("{}"),
	})
	return a
}

var defaultConsensusParams = &tmprototypes.ConsensusParams{
	Block: &tmprototypes.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmprototypes.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmprototypes.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}
