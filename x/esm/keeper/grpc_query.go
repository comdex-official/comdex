package keeper

import (
	"github.com/comdex-official/comdex/x/esm/types"
)

var _ types.QueryServer = Keeper{}
