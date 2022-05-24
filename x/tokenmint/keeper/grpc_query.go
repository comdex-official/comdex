package keeper

import (
	"github.com/comdex-official/comdex/x/tokenmint/types"
)

var _ types.QueryServer = Keeper{}
