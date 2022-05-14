package keeper

import (
	"github.com/comdex-official/comdex/x/collector/types"
)

var _ types.QueryServer = Keeper{}
