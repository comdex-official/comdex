package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetReserveFunds(_ sdk.Context, pool types.Pool) sdk.Int {
	return sdk.NewInt(int64(pool.ReserveFunds))
}
