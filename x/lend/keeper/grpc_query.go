package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
)

var (
	_ types.QueryServiceServer = (*queryServer)(nil)
)

type queryServer struct {
	Keeper
}

func NewQueryServiceServer(k Keeper) types.QueryServiceServer {
	return &queryServer{
		Keeper: k,
	}
}
