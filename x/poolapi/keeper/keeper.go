package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/poolapi/expected"
)

type Keeper struct {
	cdc             codec.BinaryCodec
	key             sdk.StoreKey
	params          paramstypes.Subspace
	liquiditykeeper expected.LiquidityKeeper
	oracle          expected.OracleKeeper
	vault           expected.VaultKeeper
	asset           expected.AssetKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, params paramstypes.Subspace, liquidity expected.LiquidityKeeper,
	oracle expected.OracleKeeper, vault expected.VaultKeeper, asset expected.AssetKeeper) Keeper {
	if !params.HasKeyTable() {
		params = params.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:             cdc,
		key:             key,
		params:          params,
		liquiditykeeper: liquidity,
		oracle:          oracle,
		vault:           vault,
		asset:           asset,
	}
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
