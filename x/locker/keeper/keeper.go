package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/comdex-official/comdex/x/locker/expected"
	"github.com/comdex-official/comdex/x/locker/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	key        storetypes.StoreKey
	paramstore paramtypes.Subspace
	bank       expected.BankKeeper
	asset      expected.AssetKeeper
	collector  expected.CollectorKeeper
	esm        expected.EsmKeeper
	rewards    expected.RewardsKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, ps paramtypes.Subspace, bank expected.BankKeeper, asset expected.AssetKeeper, collector expected.CollectorKeeper, esm expected.EsmKeeper, rewards expected.RewardsKeeper) Keeper {
	return Keeper{
		cdc:        cdc,
		key:        key,
		paramstore: ps,
		bank:       bank,
		asset:      asset,
		collector:  collector,
		esm:        esm,
		rewards:    rewards,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Store(ctx sdk.Context) storetypes.KVStore {
	return ctx.KVStore(k.key)
}
