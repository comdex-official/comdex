package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
)

var _ types.QueryServer = Keeper{}
