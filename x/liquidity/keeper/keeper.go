package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/petrichormoney/petri/x/liquidity/expected"
	"github.com/petrichormoney/petri/x/liquidity/types"
)

// Keeper of the liquidity store.
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramSpace paramstypes.Subspace

	accountKeeper expected.AccountKeeper
	bankKeeper    expected.BankKeeper

	assetKeeper   expected.AssetKeeper
	marketKeeper  expected.MarketKeeper
	rewardsKeeper expected.RewardsKeeper
}

// NewKeeper creates a new liquidity Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	paramSpace paramstypes.Subspace,
	accountKeeper expected.AccountKeeper,
	bankKeeper expected.BankKeeper,
	assetKeeper expected.AssetKeeper,
	marketKeeper expected.MarketKeeper,
	rewardsKeeper expected.RewardsKeeper,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    paramSpace,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		assetKeeper:   assetKeeper,
		marketKeeper:  marketKeeper,
		rewardsKeeper: rewardsKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
