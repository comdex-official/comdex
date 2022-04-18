package keeper

import (
	exchangetypes "github.com/comdex-official/comdex/x/exchange/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	cdc    codec.BinaryCodec
	key    sdk.StoreKey
	params paramstypes.Subspace
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, params paramstypes.Subspace) *Keeper {
	if !params.HasKeyTable() {
		params = params.WithKeyTable(exchangetypes.ParamKeyTable())
	}

	return &Keeper{
		cdc:    cdc,
		key:    key,
		params: params,
	}
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
