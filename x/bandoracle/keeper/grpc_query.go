package keeper

import (
	"github.com/comdex-official/comdex/x/bandoracle/types"
)

var _ types.QueryServer = Keeper{}
