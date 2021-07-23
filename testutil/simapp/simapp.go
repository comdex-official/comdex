package simapp

import (
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	comdex "github.com/comdex-official/comdex/app"
)

// New creates application instance with in-memory database and disabled logging.
func New(dir string) *comdex.App {
	db := tmdb.NewMemDB()
	logger := log.NewNopLogger()

	encoding := comdex.MakeEncodingConfig()

	a := comdex.New(logger, db, nil, true, map[int64]bool{}, dir, 0, encoding,
		simapp.EmptyAppOptions{})
	// InitChain updates deliverState which is required when app.NewContext is called
	a.InitChain(abcitypes.RequestInitChain{
		ConsensusParams: defaultConsensusParams,
		AppStateBytes:   []byte("{}"),
	})
	return a
}

var defaultConsensusParams = &abcitypes.ConsensusParams{
	Block: &abcitypes.BlockParams{
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
