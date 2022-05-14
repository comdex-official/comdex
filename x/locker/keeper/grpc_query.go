package keeper

import (
	"github.com/comdex-official/comdex/x/locker/types"
)

var _ types.QueryServer = Keeper{}
