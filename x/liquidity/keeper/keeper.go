package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/liquidity/expected"
)

type Keeper struct {
	cdc             codec.BinaryCodec
	key             sdk.StoreKey
	params          paramstypes.Subspace
	liquiditykeeper expected.LiquidityKeeper
	oraclekeeper    expected.OracleKeeper
	vaultkeeper     expected.VaultKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, params paramstypes.Subspace, liquidity expected.LiquidityKeeper, oracle expected.OracleKeeper, vault expected.VaultKeeper) Keeper {
	if !params.HasKeyTable() {
		params = params.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:             cdc,
		key:             key,
		params:          params,
		liquiditykeeper: liquidity,
		oraclekeeper:    oracle,
		vaultkeeper:     vault,
	}
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
