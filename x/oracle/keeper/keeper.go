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
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey,channel expected.ChannelKeeper,port expected.PortKeeper,scoped expected.ScopedKeeper) *Keeper {
	return &Keeper{
		cdc: cdc,
		key: key,
		channel: channel,
		port: port,
		scoped: scoped,

	}
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
