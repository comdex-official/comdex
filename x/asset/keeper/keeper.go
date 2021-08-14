package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc codec.BinaryMarshaler
	key sdk.StoreKey
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
