package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/comdex-official/comdex/x/asset/types"
	oraclekeeper "github.com/comdex-official/comdex/x/oracle/keeper"
	"github.com/comdex-official/comdex/x/poolapi/expected"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"
)

type Keeper struct {
	cdc             codec.BinaryCodec
	key             sdk.StoreKey
	params          paramstypes.Subspace
	liquiditykeeper expected.LiquidityKeeper
	oracleKeeper    oraclekeeper.Keeper
	vaultKeeper     vaultkeeper.Keeper
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, params paramstypes.Subspace, liquidity expected.LiquidityKeeper,
	oracleKeeper oraclekeeper.Keeper, vaultKeeper vaultkeeper.Keeper) Keeper {
	if !params.HasKeyTable() {
		params = params.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:             cdc,
		key:             key,
		params:          params,
		liquiditykeeper: liquidity,
		oracleKeeper:    oracleKeeper,
		vaultKeeper:     vaultKeeper,
	}
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
