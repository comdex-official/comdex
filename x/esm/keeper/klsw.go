package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/esm/types"
)

func (k *Keeper) SetCondition(ctx sdk.Context, condition bool) {
	var (
		store = k.Store(ctx)
		key   = types.Condition
		value = k.cdc.MustMarshal(
			&protobuftypes.BoolValue{
				Value: condition,
			},
		)
	)
	fmt.Println("condition setting")
	fmt.Println(value)

	store.Set(key, value)
}

func (k *Keeper) GetCondition(ctx sdk.Context) bool {
	var (
		store = k.Store(ctx)
		key   = types.Condition
		value = store.Get(key)
	)

	var id protobuftypes.BoolValue
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}
