package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/petrichormoney/petri/x/locker/expected"
	"github.com/petrichormoney/petri/x/locker/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	key        sdk.StoreKey
	paramstore paramtypes.Subspace
	bank       expected.BankKeeper
	asset      expected.AssetKeeper
	collector  expected.CollectorKeeper
	esm        expected.EsmKeeper
	rewards    expected.RewardsKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, ps paramtypes.Subspace, bank expected.BankKeeper, asset expected.AssetKeeper, collector expected.CollectorKeeper, esm expected.EsmKeeper, rewards expected.RewardsKeeper) Keeper {
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

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
