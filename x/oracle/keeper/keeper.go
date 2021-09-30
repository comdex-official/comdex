package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/comdex-official/comdex/x/oracle/expected"
)

type Keeper struct {
	cdc     codec.BinaryCodec
	key     sdk.StoreKey
	params  paramstypes.Subspace
	channel expected.ChannelKeeper
	port    expected.PortKeeper
	scoped  expected.ScopedKeeper
	asset   expected.AssetKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, params paramstypes.Subspace, channel expected.ChannelKeeper,
	port expected.PortKeeper, scoped expected.ScopedKeeper, asset expected.AssetKeeper) *Keeper {

	return &Keeper{
		cdc:     cdc,
		key:     key,
		params:  params,
		channel: channel,
		port:    port,
		scoped:  scoped,
		asset:   asset,
	}
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
