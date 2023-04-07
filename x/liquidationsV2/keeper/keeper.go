package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/liquidationsV2/expected"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	memKey     sdk.StoreKey
	paramstore paramtypes.Subspace
	account    expected.AccountKeeper
	bank       expected.BankKeeper
	vault      expected.VaultKeeper
	asset      expected.AssetKeeper
	market     expected.MarketKeeper
	esm        expected.EsmKeeper
	rewards    expected.RewardsKeeper
	lend       expected.LendKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	account expected.AccountKeeper,
	bank expected.BankKeeper,
	asset expected.AssetKeeper,
	vault expected.VaultKeeper,
	market expected.MarketKeeper,
	esm expected.EsmKeeper,
	rewards expected.RewardsKeeper,
	lend expected.LendKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		account:    account,
		bank:       bank,
		asset:      asset,
		vault:      vault,
		market:     market,
		esm:        esm,
		rewards:    rewards,
		lend:       lend,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}
