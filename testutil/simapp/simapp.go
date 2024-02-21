package simapp

import (
	"github.com/comdex-official/comdex/app/params"
	"time"

	tmdb "github.com/cometbft/cometbft-db"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmprototypes "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	comdex "github.com/comdex-official/comdex/app"
)

// New creates application instance with in-memory database and disabled logging.
func New(dir string) *comdex.App {
	var (
		db       = tmdb.NewMemDB()
		logger   = log.NewNopLogger()
		encoding = params.MakeEncodingConfig()
	)

	a := comdex.New(logger, db, nil, true, map[int64]bool{}, dir, 0, encoding,
		simtestutil.EmptyAppOptions{}, comdex.GetWasmEnabledProposals(), comdex.EmptyWasmOpts)
	// InitChain updates deliverState which is required when app.NewContext is called
	a.InitChain(abcitypes.RequestInitChain{
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
