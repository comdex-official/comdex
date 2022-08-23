package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/comdex-official/comdex/x/asset/expected"
)

type Keeper struct {
	cdc     codec.BinaryCodec
	key     sdk.StoreKey
	params  paramstypes.Subspace
	oracle  expected.MarketKeeper
	rewards expected.RewardsKeeper
	vault   expected.VaultKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, params paramstypes.Subspace, oracle expected.MarketKeeper, rewards expected.RewardsKeeper, vault expected.VaultKeeper) Keeper {
	if !params.HasKeyTable() {
		params = params.WithKeyTable(assettypes.ParamKeyTable())
	}

	return Keeper{
		cdc:     cdc,
		key:     key,
		params:  params,
		oracle:  oracle,
		rewards: rewards,
		vault:   vault,
	}
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
