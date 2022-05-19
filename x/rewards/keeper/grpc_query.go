package keeper

import (
	"github.com/comdex-official/comdex/x/rewards/types"
)

var _ types.QueryServer = Keeper{}
