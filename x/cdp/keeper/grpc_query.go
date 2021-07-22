package keeper

import (
	"github.com/comdex-official/comdex/x/cdp/types"
)

var _ types.QueryServer = Keeper{}
