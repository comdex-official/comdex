package keeper

import (
	"github.com/comdex-official/comdex/x/auctionsV2/types"
)

var _ types.QueryServer = Keeper{}

type QueryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &QueryServer{
		Keeper: k,
	}
}
