package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetFeeSource(ctx sdk.Context, sdkTx sdk.Tx, feePayer sdk.AccAddress, fee sdk.Coins) {
	fmt.Println(feePayer.String())
	fmt.Println(fee)
}
