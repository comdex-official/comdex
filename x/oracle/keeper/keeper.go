package keeper

import (
	oracletypes "github.com/comdex-official/comdex/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/oracle/expected"
)

type Keeper struct {
	cdc         codec.BinaryCodec
	key         sdk.StoreKey
	params      paramstypes.Subspace
	channel     expected.ChannelKeeper
	port        expected.PortKeeper
	scoped      expected.ScopedKeeper
	assetKeeper assetkeeper.Keeper
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, params paramstypes.Subspace, channel expected.ChannelKeeper,
	port expected.PortKeeper, scoped expected.ScopedKeeper, assetKeeper assetkeeper.Keeper) *Keeper {
	if !params.HasKeyTable() {
		params = params.WithKeyTable(oracletypes.ParamKeyTable())
	}

	return &Keeper{
		cdc:         cdc,
		key:         key,
		params:      params,
		channel:     channel,
		port:        port,
		scoped:      scoped,
		assetKeeper: assetKeeper,
	}
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
