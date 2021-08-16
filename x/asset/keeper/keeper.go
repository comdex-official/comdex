package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/expected"
)

type Keeper struct {
	cdc     codec.BinaryMarshaler
	key     sdk.StoreKey
	channel expected.ChannelKeeper
	scoped  expected.ScopedKeeper
}

func NewKeeper(cdc codec.BinaryMarshaler, key sdk.StoreKey) *Keeper {
	return &Keeper{
		cdc: cdc,
		key: key,
	}
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
