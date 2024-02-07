package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/comdex-official/comdex/x/liquidity/expected"
	"github.com/comdex-official/comdex/x/liquidity/types"
)

// Keeper of the liquidity store.
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramSpace paramstypes.Subspace

	accountKeeper expected.AccountKeeper
	bankKeeper    expected.BankKeeper

	assetKeeper   expected.AssetKeeper
	marketKeeper  expected.MarketKeeper
	rewardsKeeper expected.RewardsKeeper
	tokenmint     expected.TokenMintKeeper
}

// NewKeeper creates a new liquidity Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramSpace paramstypes.Subspace,
	accountKeeper expected.AccountKeeper,
	bankKeeper expected.BankKeeper,
	assetKeeper expected.AssetKeeper,
	marketKeeper expected.MarketKeeper,
	rewardsKeeper expected.RewardsKeeper,
	tokenmint expected.TokenMintKeeper,
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
		tokenmint:     tokenmint,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
